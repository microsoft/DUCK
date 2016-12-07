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
type NormalizedDocument struct {
	structs.Document
	Parts map[string][]string
}

//NewNormalizer returns a new initialized normalizer
func NewNormalizer(doc structs.Document, db *db.Database, webdir string) (*normalizer, error) {
	//norm := Normalizer{original: doc, database: db}
	norm := normalizer{original: doc}
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

//Normalize normalizes a Document for further validation
func (n *normalizer) CreateDict() (*NormalizedDocument, error) {
	n.normalized = new(NormalizedDocument)

	//make sure we have every part only once for each code
	//for this we make a map for every code which we will later transform into a list
	parts := make(map[string]map[string]struct{})

	//get original taxonomy code for each field in each statement
	//and save it into parts map
	//after that we check if we have mising fields in these Statements
	for _, statement := range n.original.Statements {

		if returnCode := n.getCode("action", statement.ActionCode); returnCode != "" {
			if parts[statement.ActionCode] == nil {
				parts[statement.ActionCode] = make(map[string]struct{})
			}
			parts[statement.ActionCode][returnCode] = struct{}{}
		}

		if returnCode := n.getCode("qualifier", statement.QualifierCode); returnCode != "" {
			if parts[statement.QualifierCode] == nil {
				parts[statement.QualifierCode] = make(map[string]struct{})
			}
			parts[statement.QualifierCode][returnCode] = struct{}{}
		}

		if returnCode := n.getCode("dataUseCategory", statement.DataCategoryCode); returnCode != "" {
			if parts[statement.DataCategoryCode] == nil {
				parts[statement.DataCategoryCode] = make(map[string]struct{})
			}
			parts[statement.DataCategoryCode][returnCode] = struct{}{}
		}

		if returnCode := n.getCode("scope", statement.UseScopeCode); returnCode != "" {
			if parts[statement.UseScopeCode] == nil {
				parts[statement.UseScopeCode] = make(map[string]struct{})
			}
			parts[statement.UseScopeCode][returnCode] = struct{}{}
		}

		if returnCode := n.getCode("scope", statement.ResultScopeCode); returnCode != "" {
			if parts[statement.ResultScopeCode] == nil {
				parts[statement.ResultScopeCode] = make(map[string]struct{})
			}
			parts[statement.ResultScopeCode][returnCode] = struct{}{}
		}

		if returnCode := n.getCode("scope", statement.SourceScopeCode); returnCode != "" {
			if parts[statement.SourceScopeCode] == nil {
				parts[statement.SourceScopeCode] = make(map[string]struct{})
			}
			parts[statement.SourceScopeCode][returnCode] = struct{}{}
		}

		//if data use category, data category or all scopes are missing, we cant work with this statement
		if statement.UseScopeCode == "" && statement.ResultScopeCode == "" && statement.SourceScopeCode == "" {
			return n.normalized, fmt.Errorf("statement is missing all scope fields: %s", statement.TrackingID)
		}
		if statement.ActionCode == "" {
			return n.normalized, fmt.Errorf("statement is missing data use field: %s", statement.TrackingID)
		}

		if statement.DataCategoryCode == "" {
			return n.normalized, fmt.Errorf("statement is missing data category field: %s", statement.TrackingID)
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
		n.normalized.Statements = append(n.normalized.Statements, statement)
	}

	//put codes into list
	if n.normalized.Parts == nil {
		n.normalized.Parts = make(map[string][]string)
	}
	for key, value := range parts {
		for code := range value {
			n.normalized.Parts[key] = append(n.normalized.Parts[key], code)
		}
	}
	//put all the other fields from the original into the normalized struct
	n.normalized.ID = n.original.ID
	n.normalized.Locale = n.original.Locale
	n.normalized.Name = n.original.Name
	n.normalized.Owner = n.original.Owner
	n.normalized.Revision = n.original.Revision

	return n.normalized, nil
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
	return ""

}

//Denormalize denormalizes a Document after validation
func (n *normalizer) Denormalize() *structs.Document {
	// we have the original, so why Denormalize?
	return &n.original
}
