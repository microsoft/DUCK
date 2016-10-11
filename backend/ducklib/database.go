package ducklib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/Microsoft/DUCK/backend/pluginregistry"
	"github.com/twinj/uuid"
)

type database struct {
	Config structs.DBConf
}

var db pluginregistry.DBPlugin

//NewDatabase returns an intialized database struct
func NewDatabase(config structs.DBConf) *database {
	return &database{Config: config}
}

//Put this into plugin
func FillTestdata() error {

	var listOfData []interface{}
	testDataFile := "src/github.com/Microsoft/DUCK/testdata.json"
	goPath := os.Getenv("GOPATH")
	if goPath != "" && !filepath.IsAbs(testDataFile) {
		testDataFile = filepath.Join(goPath, testDataFile)
	}

	if _, err := os.Stat(testDataFile); os.IsNotExist(err) {
		testDataFile = "testdata.json"
		if _, err := os.Stat(testDataFile); os.IsNotExist(err) {
			return structs.NewHTTPError("testdata.json not found", 404)
		}
	}

	dat, err := ioutil.ReadFile(testDataFile)
	if err != nil {
		log.Printf("Error in FillTestdata while trying to read from the file: %s", err)
		return structs.WrapErrWith(err, structs.NewHTTPError(err.Error(), 500))
	}

	if err := json.Unmarshal(dat, &listOfData); err != nil {
		return err
	}

	for _, l := range listOfData {

		mp := l.(map[string]interface{})
		entryType := mp["type"].(string)

		if entryType == "user" {
			var u structs.User
			u.FromValueMap(mp)

			docs, err := db.GetDocumentSummariesForUser(u.ID)
			if err == nil {
				for _, val := range docs {
					db.DeleteDocument(val.ID)
				}

			}

		}

	}

	for _, l := range listOfData {

		mp := l.(map[string]interface{})

		entryType := mp["type"].(string)

		switch entryType {
		case "document":
			var d structs.Document
			d.FromValueMap(mp)

			if err := db.NewDocument(d); err != nil {
				return structs.WrapErrWith(err, structs.NewHTTPError(err.Error(), 500))
			}
		case "user":
			var u structs.User
			u.FromValueMap(mp)
			if err := db.NewUser(u); err != nil {
				log.Printf("tried to put user %s with ID %s into database but ID already exists", u.Email, u.ID)

			}

		}

	}

	return nil
}

//Init initializes the database and checks for connection errors
func (database *database) Init() error {

	db = pluginregistry.DatabasePlugin

	return db.Init(database.Config)

}

/*
User DB operations
*/

//GetLogin returns id and password for username
func (database *database) GetLogin(email string) (id string, pw string, err error) {
	return db.GetLogin(email)
}

func (database *database) GetUser(userid string) (structs.User, error) {

	return db.GetUser(userid)

}

func (database *database) DeleteUser(id string) error {

	return db.DeleteUser(id)

}

func (database *database) PutUser(user structs.User) error {

	return db.UpdateUser(user)

}

func (database *database) PostUser(user structs.User) (ID string, err error) {
	//check for duplicate
	if user.Email == "" {
		return "", structs.NewHTTPError("No email submitted", 400)
	}
	if user.Password == "" {
		return "", structs.NewHTTPError("No password submitted", 400)
	}

	_, _, err = db.GetLogin(user.Email)

	// if user is not in database we can create a new one
	//TODO: What if another Database Plugin returns another Error when getting an nonexistant User?
	if err != nil && err.Error() == "User not found" {
		u := uuid.NewV4()
		uuid := uuid.Formatter(u, uuid.Clean)
		user.ID = uuid
		return uuid, db.NewUser(user)

	}

	//We don't know if the user exists because we got an error checking this
	if err != nil {
		return "", err
	}

	return "", structs.NewHTTPError("User already exists", 409)
}

func (database *database) GetUserDict(userid string) (structs.Dictionary, error) {

	return db.GetUserDict(userid)

}
func (database *database) PutUserDict(dict structs.Dictionary, userID string) error {

	return db.UpdateUserDict(dict, userID)

}

/*
Document DB operations

*/
func (database *database) GetDocument(documentid string) (structs.Document, error) {

	return db.GetDocument(documentid)

}
func (database *database) GetDocumentSummariesForUser(userid string) ([]structs.Document, error) {

	return db.GetDocumentSummariesForUser(userid)

}

func (database *database) DeleteDocument(id string) error {

	return db.DeleteDocument(id)

}

func (database *database) PutDocument(doc structs.Document) error {

	return db.UpdateDocument(doc)

}

func (database *database) PostDocument(doc structs.Document) (ID string, err error) {
	if doc.Name == "" {
		return "", structs.NewHTTPError("No Document Name submitted", 400)
	}

	if doc.Owner == "" {
		return "", structs.NewHTTPError("No Document Owner submitted", 400)
	}

	smap := make(map[string]bool)
	for _, s := range doc.Statements {
		if _, prs := smap[s.TrackingID]; prs {
			return "", structs.NewHTTPError("Document contains two Statements with the same statement ID", 400)
		}
		smap[s.TrackingID] = true
	}

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
