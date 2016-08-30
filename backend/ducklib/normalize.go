package ducklib

/*
// Normalizer ...
type Normalizer struct {
	original     structs.Document
	normalized   structs.Document
	taxonomy     structs.Taxonomy
	dictionary   structs.Dictionary
	database     *Database
	categoryDict map[string]map[string]*structs.DictionaryEntry
	codeDict     map[string]map[string]*structs.DictionaryEntry
}

//NewNormalizer returns a new initialized Normalizer
func NewNormalizer(doc structs.Document, db *Database) (*Normalizer, error) {
	norm := Normalizer{original: doc, database: db}
	// set dictionary
	user, err := norm.database.GetUser(doc.Owner)
	if err != nil {
		return nil, err
	}
	norm.dictionary = user.Dictionary
	for _, entry := range norm.dictionary {
		// for better searchability save pointer to dict entry in map
		// entries in categoryDict are ordered by Type (e.g. "scope" or "action" etc)
		// and category (e.g. 2).
		// entries in codeDict are ordered by Type (e.g. "scope" or "action" etc)
		// and code (e.g. "account_data" or "linked_data" etc.).
		norm.categoryDict[entry.Type][entry.Value] = &entry
		norm.codeDict[entry.Type][entry.Code] = &entry
	}

	//Taxonomy
	goPath := os.Getenv("GOPATH")
	taxPath := fmt.Sprintf("/src/github.com/Microsoft/DUCK/frontend/src/assets/config/taxonomy-%s.json", user.Locale)
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
func (n *Normalizer) Normalize() *structs.Document {
	for _, statement := range n.original.Statements {
		normStmt := structs.Statement{}
		normStmt.ActionCode = n.getCode("action", statement.ActionCode)
		normStmt.DataCategoryCode = n.getCode("dataCategory", statement.DataCategoryCode)
		normStmt.Passive = statement.Passive
		normStmt.QualifierCode = n.getCode("qualifier", statement.QualifierCode)
		normStmt.ResultScopeCode = n.getCode("scope", statement.ResultScopeCode)
		normStmt.SourceScopeCode = n.getCode("scope", statement.SourceScopeCode)
		normStmt.TrackingID = statement.TrackingID
		normStmt.UseScopeCode = n.getCode("scope", statement.UseScopeCode)
		n.normalized.Statements = append(n.normalized.Statements, normStmt)
	}
	n.normalized.ID = n.original.ID
	n.normalized.Locale = n.original.Locale
	n.normalized.Name = n.original.Name
	n.normalized.Owner = n.original.Owner
	n.normalized.Revision = n.original.Revision

	return &n.normalized
}

// get Code from taxonomy. For this a dictionary entry is retrieved from the codeDict
//in the taxonomy is then looked for the category of the dictionary entry since this
//should be the same regardless of the code value the corresponding code in the
//taxonomy is then returned if one is found
func (n *Normalizer) getCode(Type string, Code string) string {

	dict, prs := n.codeDict[Type][Code]
	if prs {
		switch Type {
		case "action":
			for _, typ := range n.taxonomy.Action {
				if dict.Category == typ.Category {
					return typ.Code
				}
			}
		case "dataCategory":
			for _, typ := range n.taxonomy.DataCategory {
				if dict.Category == typ.Category {
					return typ.Code
				}
			}
		case "qualifier":
			for _, typ := range n.taxonomy.Qualifier {
				if dict.Category == typ.Category {
					return typ.Code
				}
			}
		case "scope":
			for _, typ := range n.taxonomy.Scope {
				if dict.Category == typ.Category {
					return typ.Code
				}
			}

		}

	}
	return Code

}

//Denormalize denormalizes a Document after validation
func (n *Normalizer) Denormalize() *structs.Document {
	return &n.original
}

//DenormalizeVariants denormalises valid variants of a document
func (n *Normalizer) DenormalizeVariants() []structs.Document {

	return nil
}
*/
