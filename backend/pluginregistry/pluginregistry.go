package pluginregistry

import "github.com/Microsoft/DUCK/backend/ducklib/structs"

// DBPlugin is the interface the Databse Plugin has to satisfy
type DBPlugin interface {
	Init(url string, databasename string) error
	GetLogin(username string) (id string, pw string, err error)

	GetUser(id string) (structs.User, error)
	DeleteUser(id string) error
	NewUser(user structs.User) error
	UpdateUser(user structs.User) error

	GetDocumentSummariesForUser(userid string) ([]structs.Document, error)

	GetDocument(id string) (structs.Document, error)
	NewDocument(doc structs.Document) error
	UpdateDocument(doc structs.Document) error
	DeleteDocument(id string) error

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
