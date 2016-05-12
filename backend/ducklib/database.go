package ducklib

import "github.com/Metaform/duck/backend/pluginregistry"

type Database struct {
	url          string
	username     string
	password     string
	databasename string
}

var db = pluginregistry.DatabasePlugin

func NewDatabase() *Database {
	return &Database{databasename: "DUCK"}
}

func TestDB() {

	db.Print()
	db.Save()
}

func (self *Database) Init() {

}
