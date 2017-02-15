package dbplugin

type MyDatabase struct {
}

func (mdb *MyDatabase) Init(url string) error {
	//fmt.Println("Testextension initialized")
	return nil
}

func init() {
	//db := &MyDatabase{}
	//pluginregistry.RegisterDatabase(db)
}
