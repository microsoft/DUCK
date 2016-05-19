package pluginregistry

// DBPlugin is the interface the Databse Plugin has to satisfy
type DBPlugin interface {
	Save()
	Print()
	Init(url string, databasename string) error
	GetLogin(username string) (string, error)
}

// DatabasePlugin is the Plugin
var DatabasePlugin DBPlugin

// RegisterDatabase registers a database
func RegisterDatabase(db DBPlugin) {
	DatabasePlugin = db
}
