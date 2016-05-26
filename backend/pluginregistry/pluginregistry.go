package pluginregistry

// DBPlugin is the interface the Databse Plugin has to satisfy
type DBPlugin interface {
	Save()
	Print()
	Init(url string, databasename string) error
	GetLogin(username string) (id string, pw string, err error)
	GetEntry(id string) (user map[string]interface{}, err error)
	GetDocumentSummariesForUser(userid string) (documents []map[string]string, err error)
	PutEntry(id string, entry string) error
	DeleteEntry(id string, rev string) error
}

// DatabasePlugin is the Plugin
var DatabasePlugin DBPlugin

// RegisterDatabase registers a database
func RegisterDatabase(db DBPlugin) {
	DatabasePlugin = db
}
