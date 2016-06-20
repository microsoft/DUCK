package pluginregistry

// DBPlugin is the interface the Databse Plugin has to satisfy
type DBPlugin interface {
	Init(url string, databasename string) error
	GetLogin(username string) (id string, pw string, err error)

	GetUser(id string) (user map[string]interface{}, err error)
	DeleteUser(id string, rev string) error
	NewUser(id string, entry string) error
	UpdateUser(id string, entry string) error

	GetDocumentSummariesForUser(userid string) (documents []map[string]string, err error)

	GetDocument(id string) (document map[string]interface{}, err error)
	NewDocument(id string, entry string) error
	UpdateDocument(id string, entry string) error
	DeleteDocument(id string, rev string) error

//	GetRuleset(id string) (document map[string]interface{}, err error)
//	NewRuleset(id string, entry string) error
//	UpdateRuleset(id string, entry string) error
//	DeleteRuleset(id string, rev string) error

//	GetStatement(id string) (document map[string]interface{}, err error)
}

// DatabasePlugin is the Plugin
var DatabasePlugin DBPlugin

// RegisterDatabase registers a database
func RegisterDatabase(db DBPlugin) {
	DatabasePlugin = db
}
