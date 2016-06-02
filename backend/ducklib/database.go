package ducklib

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Metaform/duck/backend/pluginregistry"
	"github.com/twinj/uuid"
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

func FillTestdata(data []byte) error {

	var listOfData []interface{}

	if err := json.Unmarshal(data, &listOfData); err != nil {
		return err
	}
	for _, l := range listOfData {

		mp := l.(map[string]interface{})
		id := mp["_id"].(string)

		entry, err := json.Marshal(l)
		if err != nil {
			return err
		}
		if _, err := db.PutEntry(id, string(entry)); err != nil {
			return err
		}

	}

	return nil
}

//Init initializes the database and checks for connection errors
func (database *Database) Init() {
	err := db.Init(database.url, database.databasename)
	if err != nil {
		fmt.Println(err)
	}
}

//GetLogin returns id and password for username
func (database *Database) GetLogin(username string) (id string, pw string, err error) {
	return db.GetLogin(username)
}

func (database *Database) GetUser(userid string) (User, error) {
	var u User
	mp, err := db.GetEntry(userid)
	if err != nil {
		return u, err
	}

	u.fromValueMap(mp)

	return u, err
}

func (database *Database) GetDocument(documentid string) (Document, error) {
	var doc Document
	mp, err := db.GetEntry(documentid)
	if err != nil {
		return doc, err
	}

	doc.fromValueMap(mp)

	return doc, err
}
func (database *Database) GetDocumentSummariesForUser(userid string) ([]Document, error) {
	var docs []Document
	list, err := db.GetDocumentSummariesForUser(userid)
	if err != nil {
		fmt.Println(err.Error())
		return docs, err
	}

	for _, item := range list {
		docs = append(docs, Document{Name: item["name"], ID: item["id"]})
	}

	return docs, nil

}

func (database *Database) Delete(id string) error {

	doc, err := db.GetEntry(id)
	if err != nil {

		return err
	}
	if rev, prs := doc["_rev"]; prs {
		err := db.DeleteEntry(id, rev.(string))
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Could not delete Entry")

}

func (database *Database) PutEntry(id string, content []byte) (eid string, err error) {

	eid, err = db.PutEntry(id, string(content))

	return

}

func (database *Database) PostEntry(content []byte) (string, error) {
	u := uuid.NewV4()

	return database.PutEntry(uuid.Formatter(u, uuid.Clean), content)

}
