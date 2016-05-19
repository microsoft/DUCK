package dbplugin

import (
	"fmt"

	"github.com/Metaform/duck/backend/pluginregistry"
)

type MyDatabase struct {
}

func (mdb *MyDatabase) Print() {
	fmt.Println("Testextension Printed sth")
}

func (mdb *MyDatabase) Save() {
	fmt.Println("Testextension saved sth")
}
func (mdb *MyDatabase) Init(url string) error {
	fmt.Println("Testextension initialized")
}

func init() {
	db := &MyDatabase{}
	pluginregistry.RegisterDatabase(db)
}
