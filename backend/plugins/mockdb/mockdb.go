package mockdb

import (
	"errors"
	"log"
	"net/http"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/Microsoft/DUCK/backend/pluginregistry"
)

// Mock is a Mock database for testing purposes
type Mock struct {
	DataUseDocuments map[string]structs.Document
	User             map[string]structs.User
}

//Init initializes the Mock
// throws error if name equals INVALIDDBNAME
func (m *Mock) Init(dbconf structs.DBConf) error {

	if dbconf.Name == "INVALIDDBNAME" {
		return errors.New("Error initializing MockDB: Invalid Database Name")
	}

	m.User = make(map[string]structs.User)
	m.DataUseDocuments = make(map[string]structs.Document)
	_, ok := pluginregistry.DatabasePlugin.(*Mock)

	if ok {
		log.Println("MockDB already registered")

		return nil
	}
	if pluginregistry.DatabasePlugin != nil {
		log.Println("Already a registered Database. Registered anyway.")

	} else {
		log.Println("Registered first")
	}
	pluginregistry.RegisterDatabase(m)

	return nil
}

func init() {

	m := &Mock{}

	_, ok := pluginregistry.DatabasePlugin.(*Mock)

	if ok {
		log.Println("MockDB already registered")

	} else if pluginregistry.DatabasePlugin != nil {
		log.Println("Already a registered Database. Registered anyway.")

	} else {
		log.Println("Mockdb Registered first")
	}
	pluginregistry.RegisterDatabase(m)
}

//GetLogin returns ID and Password for the matching username
func (m *Mock) GetLogin(userMail string) (id string, pw string, err error) {
	for _, u := range m.User {
		if userMail == u.Email {
			id = u.ID
			pw = u.Password
			return
		}
	}
	err = errors.New("User not found")
	return
}

//GetUser returns a user struct
func (m *Mock) GetUser(id string) (structs.User, error) {

	if u, prs := m.User[id]; prs {
		return u, nil
	}
	return structs.User{}, errors.New("User not found")
}

//GetUserDict returns a user dictionary
func (m *Mock) GetUserDict(id string) (structs.Dictionary, error) {

	if u, prs := m.User[id]; prs {
		return u.GlobalDictionary, nil
	}
	return structs.Dictionary{}, errors.New("User not found")
}

//UpdateUserDict updates a user dictionary
func (m *Mock) UpdateUserDict(dict structs.Dictionary, userID string) error {

	if u, prs := m.User[userID]; prs {
		u.GlobalDictionary = dict
		m.User[userID] = u
		return nil
	}
	return errors.New("Could not update Dictionary")
}

//DeleteUser deletes a user
func (m *Mock) DeleteUser(id string) error {
	if _, prs := m.User[id]; prs {
		delete(m.User, id)
		return nil
	}
	return errors.New("Cannot delete user: User not found")
}

// NewUser creates a new User
func (m *Mock) NewUser(user structs.User) error {

	if _, prs := m.User[user.ID]; !prs {
		m.User[user.ID] = user
		return nil
	}
	return errors.New("Cannot create user: User already exists")
}

// UpdateUser updates an existing User
func (m *Mock) UpdateUser(user structs.User) error {

	if _, prs := m.User[user.ID]; prs {
		m.User[user.ID] = user
		return nil
	}
	return errors.New("Cannot Update user: User not found")
}

//GetDocumentSummariesForUser returns all documents for a user
//A summary consists only of Document ID and Name
func (m *Mock) GetDocumentSummariesForUser(userid string) ([]structs.Document, error) {
	var l []structs.Document

	if len(m.DataUseDocuments) == 0 {
		return l, errors.New("No Documents found")
	}

	for id, doc := range m.DataUseDocuments {

		if userid == doc.Owner {

			var d structs.Document
			d.ID = id
			d.Name = doc.Name
			l = append(l, d)

		}
	}
	if len(l) == 0 {
		return nil, structs.NewHTTPError("No Data returned", http.StatusNotFound)
	}
	return l, nil
}

//GetDocument returns a Document
func (m *Mock) GetDocument(id string) (structs.Document, error) {
	if d, prs := m.DataUseDocuments[id]; prs {
		return d, nil
	}
	return structs.Document{}, errors.New("Document not found")
}

//NewDocument creates a new document
func (m *Mock) NewDocument(doc structs.Document) error {

	if _, prs := m.DataUseDocuments[doc.ID]; !prs {
		m.DataUseDocuments[doc.ID] = doc
		return nil
	}
	return errors.New("Cannot create Document: Document already exists")
}

//UpdateDocument updates a Document
func (m *Mock) UpdateDocument(doc structs.Document) error {
	if _, prs := m.DataUseDocuments[doc.ID]; prs {
		m.DataUseDocuments[doc.ID] = doc
		return nil
	}
	return errors.New("Cannot Update Document: Document not found")
}

//DeleteDocument deletes a document
func (m *Mock) DeleteDocument(id string) error {
	if _, prs := m.DataUseDocuments[id]; prs {
		delete(m.DataUseDocuments, id)
		return nil
	}
	return errors.New("Cannot delete Document: Document not found")
}

/*

	//GetStatement(id string) (document map[string]interface{}, err error)
*/
