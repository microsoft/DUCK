package ducklib

import (
	"encoding/json"
	"errors"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/Microsoft/DUCK/backend/pluginregistry"
	"github.com/twinj/uuid"
)

type Database struct {
	Config structs.DBConf
}

var db pluginregistry.DBPlugin

func NewDatabase(config structs.DBConf) *Database {
	return &Database{Config: config}
}

//Put this into plugin
func FillTestdata(data []byte) error {

	var listOfData []interface{}

	if err := json.Unmarshal(data, &listOfData); err != nil {
		return err
	}
	for _, l := range listOfData {

		mp := l.(map[string]interface{})

		entryType := mp["type"].(string)

		switch entryType {
		case "document":
			var d structs.Document
			d.FromValueMap(mp)

			if err := db.NewDocument(d); err != nil {
				return err
			}
		case "user":
			var u structs.User
			u.FromValueMap(mp)
			if err := db.NewUser(u); err != nil {
				return err
			}

		}

	}

	return nil
}

//Init initializes the database and checks for connection errors
func (database *Database) Init() error {

	db = pluginregistry.DatabasePlugin

	err := db.Init(database.Config)
	if err != nil {
		return err
	}
	return nil
}

/*
User DB operations
*/

//GetLogin returns id and password for username
func (database *Database) GetLogin(email string) (id string, pw string, err error) {
	return db.GetLogin(email)
}

func (database *Database) GetUser(userid string) (structs.User, error) {

	return db.GetUser(userid)

}

func (database *Database) DeleteUser(id string) error {

	return db.DeleteUser(id)

}

func (database *Database) PutUser(user structs.User) error {

	return db.UpdateUser(user)

}

func (database *Database) PostUser(user structs.User) (ID string, err error) {
	//check for duplicate
	_, _, err = db.GetLogin(user.Email)

	if err == nil || err.Error() != "No Data returned" {

		return "", errors.New("User already exists")
	}

	u := uuid.NewV4()
	uuid := uuid.Formatter(u, uuid.Clean)
	user.ID = uuid
	return uuid, db.NewUser(user)

}

/*
Document DB operations

*/
func (database *Database) GetDocument(documentid string) (structs.Document, error) {

	return db.GetDocument(documentid)

}
func (database *Database) GetDocumentSummariesForUser(userid string) ([]structs.Document, error) {

	return db.GetDocumentSummariesForUser(userid)

}

func (database *Database) DeleteDocument(id string) error {

	return db.DeleteDocument(id)

}

func (database *Database) PutDocument(doc structs.Document) error {

	return db.UpdateDocument(doc)

}

func (database *Database) PostDocument(doc structs.Document) (ID string, err error) {
	u := uuid.NewV4()
	uuid := uuid.Formatter(u, uuid.Clean)
	doc.ID = uuid
	return uuid, db.NewDocument(doc)

}

/*
Rulebase DB operations
*/
/*
func (database *Database) GetRulebase(id string) (User, error) {
	var u User
	mp, err := db.GetRulebase(id)
	if err != nil {
		return u, err
	}

	u.fromValueMap(mp)

	return u, err
}

func (database *Database) DeleteRulebase(id string) error {

	doc, err := db.GetRulebase(id)
	if err != nil {

		return err
	}
	if rev, prs := doc["_rev"]; prs {
		err := db.DeleteRulebase(id, rev.(string))
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Could not delete Entry")

}

func (database *Database) PutRulebase(id string, content []byte) error {

	return db.UpdateRulebase(id, string(content))

}

func (database *Database) PostRulebase(content []byte) (string, error) {
	u := uuid.NewV4()
	uuid := uuid.Formatter(u, uuid.Clean)

	return uuid, db.NewRulebase(uuid, string(content))

}
*/
