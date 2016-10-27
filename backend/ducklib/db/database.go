package db

import (
	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/Microsoft/DUCK/backend/pluginregistry"
	"github.com/twinj/uuid"
)

//Database handles the database communication
type Database struct {
	Config structs.DBConf
}

var db pluginregistry.DBPlugin

//NewDatabase returns an intialized database struct
func NewDatabase(config structs.DBConf) *Database {
	return &Database{Config: config}
}

//Init initializes the database and checks for connection errors
func (database *Database) Init() error {

	db = pluginregistry.DatabasePlugin

	return db.Init(database.Config)

}

/*
User DB operations
*/

//GetLogin returns id and password for username
func (database *Database) GetLogin(email string) (id string, pw string, err error) {
	return db.GetLogin(email)
}

//GetUser ...
func (database *Database) GetUser(userid string) (structs.User, error) {

	return db.GetUser(userid)

}

//DeleteUser ..
func (database *Database) DeleteUser(id string) error {

	return db.DeleteUser(id)

}

//PutUser ..
func (database *Database) PutUser(user structs.User) error {

	return db.UpdateUser(user)

}

//PostUser ..
func (database *Database) PostUser(user structs.User) (ID string, err error) {
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

//GetUserDict ..
func (database *Database) GetUserDict(userid string) (structs.Dictionary, error) {

	return db.GetUserDict(userid)

}

//PutUserDict ..
func (database *Database) PutUserDict(dict structs.Dictionary, userID string) error {

	return db.UpdateUserDict(dict, userID)

}

/*
Document DB operations

*/

//GetDocument ..
func (database *Database) GetDocument(documentid string) (structs.Document, error) {

	return db.GetDocument(documentid)

}

//GetDocumentSummariesForUser ..
func (database *Database) GetDocumentSummariesForUser(userid string) ([]structs.Document, error) {

	return db.GetDocumentSummariesForUser(userid)

}

//DeleteDocument ..
func (database *Database) DeleteDocument(id string) error {

	return db.DeleteDocument(id)

}

//PutDocument ..
func (database *Database) PutDocument(doc structs.Document) error {

	return db.UpdateDocument(doc)

}

//PostDocument ..
func (database *Database) PostDocument(doc structs.Document) (ID string, err error) {
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
