// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package carneades

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"path/filepath"

	"github.com/Microsoft/DUCK/backend/ducklib/db"
	"github.com/Microsoft/DUCK/backend/ducklib/structs"
)

type normalizer struct {
	original    structs.Document
	normalized  *NormalizedDocument
	docTaxonomy structs.Taxonomy

	//database     *Database
	//categoryDict map[string]map[string]*structs.DictionaryEntry
	//codeDict     map[string]map[string]*structs.DictionaryEntry
	//  [azure]-> DictionaryEntry
	GlobalDict structs.Dictionary
}

//NormalizedDocument wraps structs.Document and adds an extra field 'Parts'.
//The Parts field maps a slice of codes to a more specific code
//
//  parts:
//   part1:
//     - c1
//     - c2
//   part2:
//     - part1
//IsA translates a custom code into a standard one
//relationship is as follows: KEY is a VALUE
//eg. ThingA is a capability, ThingB is a third_party_services
type NormalizedDocument struct {
	structs.Document
	Statements []NormalizedStatement
	IsA        map[string]string
	Facts      []string
}

type NormalizedStatement struct {
	structs.Statement
	UseScopeLocation    string
	SourceScopeLocation string
	ResultScopeLocation string
	PlaceInStruct       int
}

//NewNormalizer returns a new initialized normalizer
func NewNormalizer(doc structs.Document, db *db.Database, webdir string) (*normalizer, error) {
	//norm := Normalizer{original: doc, database: db}
	norm := normalizer{original: doc}

	//put everything into new

	user, err := db.GetUser(doc.Owner)
	if err != nil {
		return &norm, err
	}

	// set dictionary
	//
	//DictionaryEntry for MIcrosoft Azure
	//("microsoft_azure", {
	//	value : "Microsoft Azure",
	//	type : "scope",
	//	code : "microsoft_azure",
	//	category : "2",
	//	dictionaryType : "global"
	//})
	norm.GlobalDict = user.GlobalDictionary

	//Taxonomy

	docTaxPath := fmt.Sprintf("/assets/config/taxonomy-%s.json", doc.Locale)
	docPath := filepath.Join(webdir, docTaxPath)
	dat, err := ioutil.ReadFile(docPath)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(dat, &norm.docTaxonomy); err != nil {
		return nil, err
	}
	return &norm, nil
}

func (n *normalizer) GetNormalized() (*NormalizedDocument, error) {

	n.normalized = new(NormalizedDocument)
	//put all the other fields from the original into the normalized struct
	n.normalized.ID = n.original.ID
	n.normalized.Locale = n.original.Locale
	n.normalized.Name = n.original.Name
	n.normalized.Owner = n.original.Owner
	n.normalized.Revision = n.original.Revision

	//creates dict and moves statements into norm dict
	if err := n.CreateDict(); err != nil {
		return n.normalized, err
	}

	if err := n.SetLocation(); err != nil {
		return n.normalized, err
	}

	if err := n.Straighten(); err != nil {
		return n.normalized, err
	}

	return n.normalized, nil
}

//Straighten moves the except and and clauses int their own statements
func (n *normalizer) Straighten() error {
	/*for _, stmt := range n.original.Statements {
		for i, dcat := range stmt.DataCategories {
			if dcat.Op == structs.AND {
				statement := new(structs.Statement)

			}
		}
	}*/
	return nil
}

//GetLocation sets the Loaction fields in the Normalized Document
func (n *normalizer) SetLocation() error {
	fmt.Println("SETLOCATION")
	for i, stmt := range n.normalized.Statements {
		if stmt.UseScopeLocation == "" {
			n.normalized.Statements[i].UseScopeLocation = "null"
			fmt.Println("null")
		} else {
			fmt.Println("full")
		}
		if stmt.SourceScopeLocation == "" {
			n.normalized.Statements[i].SourceScopeLocation = "null"
		}
		if stmt.ResultScopeLocation == "" {
			n.normalized.Statements[i].ResultScopeLocation = "null"
		}
	}
	return nil
}

//Normalize normalizes a Document for further validation
func (n *normalizer) CreateDict() error {

	//make sure we have every part only once for each code
	//for this we make a map for every code which we will later transform into a list
	isA := make(map[string]string)

	// we check if we have missing fields in a Statements
	//if not we get original taxonomy code for each field in each statement
	//and save it into parts map
	for _, statement := range n.original.Statements {

		//if data use category, data category or all scopes are missing, we cant work with this statement
		if statement.UseScopeCode == "" && statement.ResultScopeCode == "" && statement.SourceScopeCode == "" {
			return fmt.Errorf("statement is missing all scope fields: %s", statement.TrackingID)
		}
		if statement.ActionCode == "" {
			return fmt.Errorf("statement is missing data use field: %s", statement.TrackingID)
		}

		if statement.DataCategoryCode == "" {
			return fmt.Errorf("statement is missing data category field: %s", statement.TrackingID)
		}

		//find original code for the one used

		if returnCode := n.getCode("action", statement.ActionCode); returnCode != "" {

			if _, prs := isA[statement.ActionCode]; prs == false {
				isA[statement.ActionCode] = returnCode
			} else if prs == true && isA[statement.ActionCode] != returnCode {
				return fmt.Errorf("The following custom code can be two or more things, which should not be possible: %s", statement.ActionCode)
			}
		}

		if returnCode := n.getCode("qualifier", statement.QualifierCode); returnCode != "" {
			if _, prs := isA[statement.QualifierCode]; prs == false {
				isA[statement.QualifierCode] = returnCode
			} else if prs == true && isA[statement.QualifierCode] != returnCode {
				return fmt.Errorf("The following custom code can be two or more things, which should not be possible: %s", statement.QualifierCode)
			}
		}

		if returnCode := n.getCode("dataUseCategory", statement.DataCategoryCode); returnCode != "" {
			if _, prs := isA[statement.DataCategoryCode]; prs == false {
				isA[statement.DataCategoryCode] = returnCode
			} else if prs == true && isA[statement.DataCategoryCode] != returnCode {
				return fmt.Errorf("The following custom code can be two or more things, which should not be possible: %s", statement.DataCategoryCode)
			}
		}

		if returnCode := n.getCode("scope", statement.UseScopeCode); returnCode != "" {
			if _, prs := isA[statement.UseScopeCode]; prs == false {
				isA[statement.UseScopeCode] = returnCode
			} else if prs == true && isA[statement.UseScopeCode] != returnCode {
				return fmt.Errorf("The following custom code can be two or more things, which should not be possible: %s", statement.UseScopeCode)
			}
		}

		if returnCode := n.getCode("scope", statement.ResultScopeCode); returnCode != "" {
			if _, prs := isA[statement.ResultScopeCode]; prs == false {
				isA[statement.ResultScopeCode] = returnCode
			} else if prs == true && isA[statement.ResultScopeCode] != returnCode {
				return fmt.Errorf("The following custom code can be two or more things, which should not be possible: %s", statement.ResultScopeCode)
			}
		}

		if returnCode := n.getCode("scope", statement.SourceScopeCode); returnCode != "" {
			if _, prs := isA[statement.SourceScopeCode]; prs == false {
				isA[statement.SourceScopeCode] = returnCode
			} else if prs == true && isA[statement.SourceScopeCode] != returnCode {
				return fmt.Errorf("The following custom code can be two or more things, which should not be possible: %s", statement.SourceScopeCode)
			}
		}

		// if qualifier is missing that means the qualifier is "unqualified"
		if statement.QualifierCode == "" {
			statement.QualifierCode = "unqualified"
		}

		//if we have at least one scope we can fill the other two (19944  10.2.2.1)

		if statement.UseScopeCode != "" {
			if statement.SourceScopeCode == "" {
				statement.SourceScopeCode = statement.UseScopeCode
			}
			if statement.ResultScopeCode == "" {
				statement.ResultScopeCode = statement.UseScopeCode
			}
		}

		if statement.SourceScopeCode != "" {
			if statement.UseScopeCode == "" {

				statement.UseScopeCode = statement.SourceScopeCode
			}
			if statement.ResultScopeCode == "" {
				statement.ResultScopeCode = statement.SourceScopeCode
			}
		}
		if statement.ResultScopeCode != "" {
			if statement.UseScopeCode == "" {
				statement.UseScopeCode = statement.ResultScopeCode
			}
			if statement.SourceScopeCode == "" {
				statement.SourceScopeCode = statement.ResultScopeCode
			}
		}
		//add statement to normalized Document
		normstmt := NormalizedStatement{}
		normstmt.ActionCode = statement.ActionCode
		normstmt.DataCategories = make([]structs.DataCategories, len(statement.DataCategories))
		copy(normstmt.DataCategories, statement.DataCategories)
		normstmt.DataCategoryCode = statement.DataCategoryCode
		normstmt.Passive = statement.Passive
		normstmt.QualifierCode = statement.QualifierCode
		normstmt.ResultScopeCode = statement.ResultScopeCode
		normstmt.SourceScopeCode = statement.SourceScopeCode
		normstmt.Tag = statement.Tag
		normstmt.TrackingID = statement.TrackingID
		normstmt.UseScopeCode = statement.UseScopeCode

		n.normalized.Statements = append(n.normalized.Statements, normstmt)
	}

	//write partsOf and isA map into Facts
	n.getFacts()

	return nil
}

//getFacts transforms the IsA and Parts maps into a list of CHR facts
//in the form of "IsA(Thing, capability)." and "PartOf(Thing, OtherThing)."
func (n *normalizer) getFacts() {
	for k, v := range n.normalized.IsA {
		n.normalized.Facts = append(n.normalized.Facts, fmt.Sprintf("isA(%s,%s).", k, v))
	}
	//TODO: also add parts to this list: partOf(A,B)

}

// get Code from taxonomy. For this a dictionary entry is retrieved from the codeDict
//in the taxonomy is then looked for the category of the dictionary entry since this
//should be the same regardless of the code value the corresponding code in the
//taxonomy is then returned if one is found
func (n *normalizer) getCode(Type string, Code string) string {

	//if code is empty return
	if Code == "" || Type == "" {
		return ""
	}

	dicto, prso := n.original.Dictionary[Code]
	dictg, prsg := n.GlobalDict[Code]

	if !prso && !prsg {
		return ""
	}
	// document dictionary takes precendence
	if prso {
		tax, prs := n.docTaxonomy[Type]
		if !prs {
			return ""
		}

		for _, typ := range tax {
			if dicto.Category == typ.Category {
				return typ.Code
			}
		}
	}
	//if we found a code in the document dict and were able to match it to a code in the taxonomy
	//we have already returned, if we failed we will try to look for a code from the user/global dict
	if prsg {
		tax, prs := n.docTaxonomy[Type]
		if !prs {
			return ""
		}

		for _, typ := range tax {
			if dictg.Category == typ.Category {
				return typ.Code
			}
		}
	}
	// if this also failed we return nothing
	//that means we also return nothing if the code is a standard code.
	//since if that is the case we don't have to translate it.
	return ""

}

//Denormalize denormalizes a Document after validation
func (n *normalizer) Denormalize() *structs.Document {
	// we have the original
	return &n.original
}
