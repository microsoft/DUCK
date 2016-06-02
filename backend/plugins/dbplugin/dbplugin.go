package dbplugin

import (
	"fmt"

	"github.com/Metaform/duck/backend/pluginregistry"
)

type MyDatabase struct {
}

func (mdb *MyDatabase) Init(url string) error {
	fmt.Println("Testextension initialized")
}

func init() {
	db := &MyDatabase{}
	pluginregistry.RegisterDatabase(db)
}
