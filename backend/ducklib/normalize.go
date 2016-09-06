package ducklib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
)

// Normalizer ...
type Normalizer struct {
	original   structs.Document
	normalized *NormalizedDocument
	taxonomy   structs.Taxonomy

	//database     *Database
	//categoryDict map[string]map[string]*structs.DictionaryEntry
	//codeDict     map[string]map[string]*structs.DictionaryEntry
	//  [azure]-> DictionaryEntry
	GlobalDict *structs.Dictionary
}

type NormalizedDocument struct {
	structs.Document
	Parts map[string][]string
}

/*
parts:
  azure:
    - c1
    - c2
  p1:
    - azure

*/

//NewNormalizer returns a new initialized Normalizer
func NewNormalizer(doc structs.Document, db *Database) (*Normalizer, error) {
	//norm := Normalizer{original: doc, database: db}
	norm := Normalizer{original: doc}
	// set dictionary
	/*
		for _, entry := range doc.Dictionary {
			// for better searchability save pointer to dict entry in map
			// entries in categoryDict are ordered by Type (e.g. "scope" or "action" etc)
			// and category (e.g. 2).
			// entries in codeDict are ordered by Type (e.g. "scope" or "action" etc)
			// and code (e.g. "account_data" or "linked_data" etc.).
			//norm.categoryDict[entry.Type][entry.Value] = &entry
			//norm.codeDict[entry.Type][entry.Code] = &entry

			//[microsoft_azure]-> {DictionaryEntry}
			norm.codeDict[entry.Code] = &entry

		}*/
	//DictionaryEntry for MIcrosoft Azure
	//("microsoft_azure", {
	//	value : "Microsoft Azure",
	//	type : "scope",
	//	code : "microsoft_azure",
	//	category : "2",
	//	dictionaryType : "global"
	//})

	//Taxonomy
	goPath := os.Getenv("GOPATH")
	taxPath := fmt.Sprintf("/src/github.com/Microsoft/DUCK/frontend/src/assets/config/taxonomy-%s.json", doc.Locale)
	path := filepath.Join(goPath, taxPath)

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(dat, norm.taxonomy); err != nil {
		return nil, err
	}

	return &norm, nil
}

//Normalize normalizes a Document for further validation
func (n *Normalizer) CreateDict() *NormalizedDocument {
	n.normalized = new(NormalizedDocument)

	n.normalized.Statements = n.original.Statements

	//make sure we have every part only once for each code
	parts := make(map[string]map[string]struct{})
	for _, statement := range n.original.Statements {

		if returnCode := n.getCode("action", statement.ActionCode); returnCode != "" {
			parts[statement.ActionCode][returnCode] = struct{}{}
		}
		if returnCode := n.getCode("qualifier", statement.QualifierCode); returnCode != "" {
			parts[statement.QualifierCode][returnCode] = struct{}{}
		}
		if returnCode := n.getCode("dataUseCategory", statement.DataCategoryCode); returnCode != "" {
			parts[statement.DataCategoryCode][returnCode] = struct{}{}
		}
		if returnCode := n.getCode("scope", statement.UseScopeCode); returnCode != "" {
			parts[statement.UseScopeCode][returnCode] = struct{}{}
		}
		if returnCode := n.getCode("scope", statement.ResultScopeCode); returnCode != "" {
			parts[statement.ResultScopeCode][returnCode] = struct{}{}
		}
		if returnCode := n.getCode("scope", statement.SourceScopeCode); returnCode != "" {
			parts[statement.SourceScopeCode][returnCode] = struct{}{}
		}

	}
	//put codes into list
	for key, value := range parts {
		for code := range value {
			n.normalized.Parts[key] = append(n.normalized.Parts[key], code)
		}
	}

	n.normalized.ID = n.original.ID
	n.normalized.Locale = n.original.Locale
	n.normalized.Name = n.original.Name
	n.normalized.Owner = n.original.Owner
	n.normalized.Revision = n.original.Revision

	return n.normalized
}

// get Code from taxonomy. For this a dictionary entry is retrieved from the codeDict
//in the taxonomy is then looked for the category of the dictionary entry since this
//should be the same regardless of the code value the corresponding code in the
//taxonomy is then returned if one is found
func (n *Normalizer) getCode(Type string, Code string) string {

	dict, prs := n.original.Dictionary[Code]
	if !prs {
		return ""
	}

	tax, prs := n.taxonomy[Type]
	if !prs {
		return ""
	}

	for _, typ := range tax {
		if dict.Category == typ.Category {
			return typ.Code
		}
	}
	return ""

}

//Denormalize denormalizes a Document after validation
func (n *Normalizer) Denormalize() *structs.Document {
	return &n.original
}

//DenormalizeVariants denormalises valid variants of a document
func (n *Normalizer) DenormalizeVariants() []structs.Document {

	return nil
}
