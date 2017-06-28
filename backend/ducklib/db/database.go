// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package db

import (
	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/Microsoft/DUCK/backend/pluginregistry"
	"github.com/twinj/uuid"
)

//Database is a Wrapper around the DBPlugin and handles the database communication
type Database struct {
	Config structs.DBConf
	db     pluginregistry.DBPlugin
}

var db pluginregistry.DBPlugin

//NewDatabase returns an intialized database struct;
//returns an error if connection problems were detected
func NewDatabase(config structs.DBConf) (*Database, error) {

	database := &Database{Config: config}

	database.db = pluginregistry.DatabasePlugin
	err := database.db.Init(database.Config)
	if err != nil {
		return &Database{}, err
	}
	return database, nil

}

/*
User DB operations
*/

//GetLogin returns id and password for username.
func (database *Database) GetLogin(email string) (id string, pw string, err error) {
	return database.db.GetLogin(email)
}

//GetUser returns a User from the plugged in Database.
func (database *Database) GetUser(userid string) (structs.User, error) {

	return database.db.GetUser(userid)

}

//DeleteUser deletes a structs.User from the plugged in Database.
func (database *Database) DeleteUser(id string) error {

	return database.db.DeleteUser(id)

}

//PutUser updates an existing User in the plugged in Database.
func (database *Database) PutUser(user structs.User) error {

	return database.db.UpdateUser(user)

}

//PostUser creates a new User in the plugged in Database and returns its ID, if no other user has this users mail address.
func (database *Database) PostUser(user structs.User) (ID string, err error) {
	//check for duplicate
	if user.Email == "" {
		return "", structs.NewHTTPError("No email submitted", 400)
	}
	if user.Password == "" {
		return "", structs.NewHTTPError("No password submitted", 400)
	}

	_, _, err = database.db.GetLogin(user.Email)

	// if user is not in database we can create a new one
	//TODO: What if another Database Plugin returns another Error when getting an nonexistant User?
	if err != nil && err.Error() == "User not found" {
		u := uuid.NewV4()
		uuid := uuid.Formatter(u, uuid.Clean)
		user.ID = uuid
		return uuid, database.db.NewUser(user)

	}

	//We don't know if the user exists because we got an error checking this
	if err != nil {
		return "", err
	}

	return "", structs.NewHTTPError("User already exists", 409)
}

//GetUserDict returns the dictionary struct of the user.
func (database *Database) GetUserDict(userid string) (structs.Dictionary, error) {

	return database.db.GetUserDict(userid)

}

//PutUserDict sets tthe users dictionary to the specified ditcionary struct.
func (database *Database) PutUserDict(dict structs.Dictionary, userID string) error {

	return database.db.UpdateUserDict(dict, userID)

}

/*
Document DB operations

*/

//GetDocument returns the Document with the specified id.
func (database *Database) GetDocument(documentid string) (structs.Document, error) {

	return database.db.GetDocument(documentid)

}

//GetDocumentSummariesForUser returns a list all data use documents a user owns.
//Summaries only include the documents name and ID.
func (database *Database) GetDocumentSummariesForUser(userid string) ([]structs.Document, error) {

	return database.db.GetDocumentSummariesForUser(userid)

}

//DeleteDocument deletes the Document with the specified id.
func (database *Database) DeleteDocument(id string) error {

	return database.db.DeleteDocument(id)

}

//PutDocument updates the given Document with a document in te database with the same ID.
func (database *Database) PutDocument(doc structs.Document) error {

	return database.db.UpdateDocument(doc)

}

//PostDocument creates a new document in the database.
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

	return uuid, database.db.NewDocument(doc)

}
