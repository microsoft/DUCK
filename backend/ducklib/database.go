package ducklib

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Microsoft/DUCK/backend/pluginregistry"
	"github.com/twinj/uuid"
)

type Database struct {
	url          string
	username     string
	password     string
	databasename string
}

var db pluginregistry.DBPlugin

func NewDatabase() *Database {
	return &Database{databasename: "duck", url: "http://127.0.0.1:5984"}
}

//Put this into plugin
func FillTestdata(data []byte) error {

	var listOfData []interface{}

	if err := json.Unmarshal(data, &listOfData); err != nil {
		return err
	}
	for _, l := range listOfData {

		mp := l.(map[string]interface{})
		id := mp["_id"].(string)
		entryType := mp["type"].(string)
		entry, err := json.Marshal(l)
		if err != nil {
			return err
		}
		switch entryType {
		case "document":
			if err := db.NewDocument(id, string(entry)); err != nil {
				return err
			}
		case "user":
			if err := db.NewUser(id, string(entry)); err != nil {
				return err
			}

		}

	}

	return nil
}

//Init initializes the database and checks for connection errors
func (database *Database) Init() {

	db = pluginregistry.DatabasePlugin

	err := db.Init(database.url, database.databasename)
	if err != nil {
		fmt.Println(err)
	}
}

/*
User DB operations
*/

//GetLogin returns id and password for username
func (database *Database) GetLogin(username string) (id string, pw string, err error) {
	return db.GetLogin(username)
}

func (database *Database) GetUser(userid string) (User, error) {
	var u User
	mp, err := db.GetUser(userid)
	if err != nil {
		return u, err
	}

	u.fromValueMap(mp)

	return u, err
}

func (database *Database) DeleteUser(id string) error {

	doc, err := db.GetDocument(id)
	if err != nil {

		return err
	}
	if rev, prs := doc["_rev"]; prs {
		err := db.DeleteDocument(id, rev.(string))
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Could not delete Entry")

}

func (database *Database) PutUser(id string, content []byte) error {

	return db.UpdateDocument(id, string(content))

}

func (database *Database) PostUser(content []byte) (string, error) {
	u := uuid.NewV4()
	uuid := uuid.Formatter(u, uuid.Clean)

	return uuid, db.NewDocument(uuid, string(content))

}

/*
Document DB operations

*/
func (database *Database) GetDocument(documentid string) (Document, error) {
	var doc Document
	mp, err := db.GetDocument(documentid)
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

func (database *Database) DeleteDocument(id string) error {

	doc, err := db.GetDocument(id)
	if err != nil {

		return err
	}
	if rev, prs := doc["_rev"]; prs {
		err := db.DeleteDocument(id, rev.(string))
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Could not delete Entry")

}

func (database *Database) PutDocument(id string, content []byte) error {

	return db.UpdateDocument(id, string(content))

}

func (database *Database) PostDocument(content []byte) (string, error) {
	u := uuid.NewV4()
	uuid := uuid.Formatter(u, uuid.Clean)

	return uuid, db.NewDocument(uuid, string(content))

}

/*
Ruleset DB operations
*/
/*
func (database *Database) GetRuleset(id string) (User, error) {
	var u User
	mp, err := db.GetRuleset(id)
	if err != nil {
		return u, err
	}

	u.fromValueMap(mp)

	return u, err
}

func (database *Database) DeleteRuleset(id string) error {

	doc, err := db.GetRuleset(id)
	if err != nil {

		return err
	}
	if rev, prs := doc["_rev"]; prs {
		err := db.DeleteRuleset(id, rev.(string))
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Could not delete Entry")

}

func (database *Database) PutRuleset(id string, content []byte) error {

	return db.UpdateRuleset(id, string(content))

}

func (database *Database) PostRuleset(content []byte) (string, error) {
	u := uuid.NewV4()
	uuid := uuid.Formatter(u, uuid.Clean)

	return uuid, db.NewRuleset(uuid, string(content))

}
*/
