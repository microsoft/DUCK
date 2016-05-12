package pluginregistry

type DBPlugin interface {
	Save()
	Print()
}

var DatabasePlugin DBPlugin

// RegisterDatabase registers a database
func RegisterDatabase(db DBPlugin) {
	DatabasePlugin = db
}
