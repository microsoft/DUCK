package dbplugin

import (
	"fmt"

	"github.com/Metaform/duck/backend/pluginregistry"
)

type MyDatabase struct {
}

func (mdb *MyDatabase) Print() {
	fmt.Println("Extension Printed sth")
}

func (mdb *MyDatabase) Save() {
	fmt.Println("Extension saved sth")
}

func init() {
	db := &MyDatabase{}
	pluginregistry.RegisterDatabase(db)
}
