package mockdb

// Mock is a Mock database for testing purposes
type Mock struct {
	DataUseDocuments []map[string]string
	User             []map[string]string
}

//Init initializes the Mock
func (m *Mock) Init(url string, databasename string) error {
	return nil
}

//GetLogin returns ID and Password for the matching username
func (m *Mock) GetLogin(username string) (id string, pw string, err error) {
	return
}

//GetUser returns a user map
func (m *Mock) GetUser(id string) (user map[string]interface{}, err error) {
	return
}

//DeleteUser deletes a user
func (m *Mock) DeleteUser(id string, rev string) error {
	return nil
}

// NewUser creates a new User
func (m *Mock) NewUser(id string, entry string) (eid string, err error) {
	return
}

//GetDocumentSummariesForUser returns all documents for a user
func (m *Mock) GetDocumentSummariesForUser(userid string) (documents []map[string]string, err error) {
	return
}

//GetDocument returns a Document
func (m *Mock) GetDocument(id string) (document map[string]interface{}, err error) {
	return
}

//NewDocument creates a new document
func (m *Mock) NewDocument(id string, entry string) (eid string, err error) {
	return
}

//UpdateDocument updates a Document
func (m *Mock) UpdateDocument(id string, entry string) (eid string, err error) {
	return
}

//DeleteDocument deletes a document
func (m *Mock) DeleteDocument(id string, rev string) error {
	return nil
}

/*

	//GetStatement(id string) (document map[string]interface{}, err error)
*/
