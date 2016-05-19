package ducklib

import (
	"fmt"

	"github.com/Metaform/duck/backend/pluginregistry"
)

type Database struct {
	url          string
	username     string
	password     string
	databasename string
}

var db = pluginregistry.DatabasePlugin

func NewDatabase() *Database {
	return &Database{databasename: "duck", url: "http://127.0.0.1:5984"}
}

func TestDB() {

	db.Print()
	db.Save()
}

//Init initializes the database and checks for connection errors
func (database *Database) Init() {
	err := db.Init(database.url, database.databasename)
	if err != nil {
		fmt.Println(err)
	}
}

//Init initializes the database and checks for connection errors
func (database *Database) GetLogin(username string) (string, error) {
	return db.GetLogin(username)
}
