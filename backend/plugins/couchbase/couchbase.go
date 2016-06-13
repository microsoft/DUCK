package couchbase

/*
implement following functions:

create DB


create User 		✓
get user 			✓
update user			✓
delete user			✓

create document		✓
get document		✓
update document		✓
delete document		✓
Document overview	✓

create RuleSet		✓
get RuleSet			✓
update RuleSet		✓
delete RuleSet		✓


*/

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Microsoft/DUCK/backend/ducklib"
	"github.com/Microsoft/DUCK/backend/pluginregistry"
)

var designDoc = `{"_id":"_design/app","views":{"foo":{"map":"function(doc){ emit(doc._id, doc._rev)}"},` +
	`"by_date":{"map":"function(doc) { if(doc.date && doc.title) {   emit(doc.date, doc.title);  }}"},` +
	`"user_login":{"map":"function(doc) { if(doc.type =='user') {   emit(doc.email,  doc.password);  }}"},` +
	`"user":{"map":"function(doc) { if(doc.type =='user') {   emit(doc._id, doc);  }}"},` +
	`"documents":{"map":"function(doc) { if(doc.type =='document') {   emit(doc._id, doc);  }}"},` +
	`"rulebases":{"map":"function(doc) { if(doc.type =='rulebase') {   emit(doc._id, doc._rev);  }}"},` +
	`"documents_by_user":{"map":"function(doc) { if(doc.type =='document') {   emit([doc.owner, doc._id], doc.name);  }}"}},` +
	`"language":"javascript"}`

//Couchbase implements the pluginregistry.DBPlugin interface for the Couchbase Dabase
type Couchbase struct {
	url      string
	database string
}

// getMap returns a map[string]interface{} containing the unmarshaled JSON from the io.Reader
func getMap(resp io.Reader) (map[string]interface{}, error) {

	content, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	//TODO: remove this when not used anymore
	//fmt.Println(string(content))

	var data map[string]interface{}

	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func getRows(resp io.Reader) ([]interface{}, error) {

	jsonbody, err := getMap(resp)
	if err != nil {
		return nil, err
	}
	rows, prs := jsonbody["rows"].([]interface{})

	if !prs || len(rows) < 1 {
		return nil, errors.New("No Data returned")
	}
	return rows, nil
}

// GetLogin returns ID and Password for the matching username from the couchbase Database
func (cb *Couchbase) GetLogin(username string) (id string, pw string, err error) {
	url := fmt.Sprintf("%s/%s/_design/app/_view/user_login?key=\"%s\"", cb.url, cb.database, username)
	//cb.url + "/" + cb.database + "/_design/app/_view/user?key='" + username + "'"

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	rows, err := getRows(resp.Body)
	if err != nil {
		return "", "", err
	}
	if len(rows) > 1 {
		return "", "", errors.New("User not unique")
	}

	row := rows[0].(map[string]interface{})

	pw, prs := row["value"].(string)
	if !prs || len(pw) <= 0 {
		return "", "", errors.New("Password not found")
	}

	id, prs = row["id"].(string)
	if !prs || len(id) <= 0 {
		return "", "", errors.New("ID not found")
	}

	return

}

//GetUser returns a User with the sppecified ID from the Couchbase Database
func (cb *Couchbase) GetUser(id string) (user map[string]interface{}, err error) {

	return cb.getCouchbaseDocument(id)
}

//GetDocument returns a Data use document with the sppecified ID from the Couchbase Database
func (cb *Couchbase) GetDocument(id string) (document map[string]interface{}, err error) {

	return cb.getCouchbaseDocument(id)
}

//GetRuleset returns a ruleset with the sppecified ID from the Couchbase Database
func (cb *Couchbase) GetRuleset(id string) (document map[string]interface{}, err error) {

	return cb.getCouchbaseDocument(id)
}

func (cb *Couchbase) getCouchbaseDocument(cbDocID string) (document map[string]interface{}, err error) {

	url := fmt.Sprintf("%s/%s/%s", cb.url, cb.database, cbDocID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cbDoc, err := getMap(resp.Body)

	if _, prs := cbDoc["_id"]; prs {
		return cbDoc, nil
	}

	return nil, errors.New("No Data")
}

//GetDocumentSummariesForUser returns a list all data use documents a user owns
func (cb *Couchbase) GetDocumentSummariesForUser(userid string) (documents []map[string]string, err error) {
	url := fmt.Sprintf("%s/%s/_design/app/_view/documents_by_user?startkey=[\"%s\",\"\"]&endkey=[\"%s\",{}]",
		cb.url, cb.database, userid, userid)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rows, err := getRows(resp.Body)
	if err != nil {
		return nil, err
	}

	documents = make([]map[string]string, len(rows))

	for row, intf := range rows {

		doc := intf.(map[string]interface{})
		document := make(map[string]string)

		document["name"] = doc["value"].(string)
		document["id"] = doc["id"].(string)
		documents[row] = document

	}

	return

}

//GetRulebases returns a list of all rulebases in the DB
func (cb *Couchbase) GetRulebases() ([]ducklib.Rulebase, error) {
	url := fmt.Sprintf("%s/%s/_design/app/_view/rulebases", cb.url, cb.database)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rows, err := getRows(resp.Body)
	if err != nil {
		return nil, err
	}

	documents := make([]ducklib.Rulebase, len(rows))

	for row, intf := range rows {

		doc := intf.(map[string]interface{})
		var document ducklib.Rulebase

		document.Name = doc["value"].(string)
		document.ID = doc["id"].(string)
		documents[row] = document

	}

	return documents, nil

}

// DeleteDocument deletes a Data Use Document from the Couchbase Database
func (cb *Couchbase) DeleteDocument(id string, rev string) error {
	return cb.deleteCbDocument(id, rev)
}

// DeleteUser deletes a User from the Couchbase Database
func (cb *Couchbase) DeleteUser(id string, rev string) error {
	return cb.deleteCbDocument(id, rev)
}

// DeleteRuleset deletes a User from the Couchbase Database
func (cb *Couchbase) DeleteRuleset(id string, rev string) error {
	return cb.deleteCbDocument(id, rev)
}

func (cb *Couchbase) deleteCbDocument(id string, rev string) error {
	url := fmt.Sprintf("%s/%s/%s?rev=%s", cb.url, cb.database, id, rev)

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(""))
	//request.SetBasicAuth("admin", "admin")
	//request.ContentLength = 0
	resp, err := client.Do(request)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonbody, err := getMap(resp.Body)
	if err != nil {
		return err
	}
	if _, prs := jsonbody["ok"]; prs {
		return nil
	}

	if err, prs := jsonbody["error"]; prs {
		reason := jsonbody["reason"].(string)

		return errors.New("Error:" + err.(string) + ", Reason: " + reason)
	}

	return errors.New("Could not decrypt Couchbase response")

}

// NewUser creates a new user in the couchbase Databse
func (cb *Couchbase) NewUser(id string, entry string) error {
	return cb.putEntry(id, entry, "user")
}

//NewDocument creates a new Data use document in the couchbase Database
func (cb *Couchbase) NewDocument(id string, entry string) error {
	return cb.putEntry(id, entry, "document")
}

//NewRuleset creates a new Ruleset in the couchbase Database
func (cb *Couchbase) NewRuleset(id string, entry string) error {
	return cb.putEntry(id, entry, "ruleset")
}

//UpdateUser replaces an existing User in the Couchbase database
func (cb *Couchbase) UpdateUser(id string, entry string) error {
	return cb.putEntry(id, entry, "user")
}

//UpdateDocument replaces an existing Data Use Document in the Couchbase database
func (cb *Couchbase) UpdateDocument(id string, entry string) error {
	return cb.putEntry(id, entry, "document")
}

//UpdateRuleset replaces an existing Ruleset in the Couchbase database
func (cb *Couchbase) UpdateRuleset(id string, entry string) error {
	return cb.putEntry(id, entry, "ruleset")
}

func (cb *Couchbase) putEntry(id, entry, entryType string) error {
	//check type of entry (document/user/ruleset)

	entryBytes := []byte(entry)
	var entryMap map[string]interface{}
	if err := json.Unmarshal(entryBytes, &entryMap); err != nil {
		return err
	}
	fieldType, prs := entryMap["type"]
	if !prs {
		entryMap["type"] = entryType
	}

	if prs && fieldType != entryType {
		err := fmt.Errorf("Couchbase Document type mismatch. Want %s, got %s", entryType, fieldType)
		
		return err
	}

	entryBytes, err := json.Marshal(entryMap)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s/%s", cb.url, cb.database, id)

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(entryBytes))
	//request.SetBasicAuth("admin", "admin")
	//request.ContentLength = 0
	resp, err := client.Do(request)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonbody, err := getMap(resp.Body)
	if err != nil {
		return err
	}
	if _, prs := jsonbody["ok"]; prs {
		return nil
	}

	if _, prs := jsonbody["error"]; prs {
		reason := jsonbody["reason"].(string)

		return errors.New(reason)
	}

	return errors.New("Could not undeerstand Couchbase response")
}

//Init initializes the Couchbase DB & tests for connection errors
func (cb *Couchbase) Init(url string, database string) error {
	log.Println("Couchase initialization")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonbody, err := getMap(resp.Body)
	if err != nil {
		return err
	}

	cdb, prs := jsonbody["couchdb"].(string)
	if !prs || cdb != "Welcome" {
		return errors.New("Connection to couchdb not successfull.")
	}
	cb.url = url
	cb.database = database
	ok, err := cb.testDBExists()
	if err != nil {
		return err
	}
	if !ok {
		log.Printf("Database %s does not exist. Creating database", cb.database)
		cb.createDatabase()
	}
	ok, err = cb.testFileExists("_design/app")
	if err != nil {
		return err
	}
	if !ok {
		log.Println("Designfile does not exist. Creating now")

		err := cb.putEntry("_design/app", designDoc, "design")
		if err != nil {
			log.Printf("ERROR: %#+v\n", err)
		}
	}

	log.Println("Testextension initialized")
	return nil
}

func (cb *Couchbase) testFileExists(id string) (bool, error) {
	url := fmt.Sprintf("%s/%s/%s", cb.url, cb.database, id)

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	jsonbody, err := getMap(resp.Body)
	if err != nil {
		return false, err
	}

	cdb, prs := jsonbody["_id"]
	if !prs {

		cdb, prs = jsonbody["error"]
		if !prs {
			return false, errors.New("Could not decrypt Couchbase response")
		}
		e := cdb.(string)
		if e == "not_found" {
			return false, nil
		}
		return false, errors.New(e)
	}
	return true, nil

}

func (cb *Couchbase) createDatabase() error {
	url := fmt.Sprintf("%s/%s", cb.url, cb.database)

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, url, strings.NewReader(""))
	//request.SetBasicAuth("admin", "admin")
	//request.ContentLength = 0
	resp, err := client.Do(request)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonbody, err := getMap(resp.Body)
	if err != nil {
		return err
	}
	if _, prs := jsonbody["ok"]; prs {
		return nil
	}

	if _, prs := jsonbody["error"]; prs {
		reason := jsonbody["reason"].(string)

		return errors.New(reason)
	}

	return errors.New("Could not decrypt Couchbase response")
}

func (cb *Couchbase) testDBExists() (bool, error) {
	url := fmt.Sprintf("%s/%s", cb.url, cb.database)

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	jsonbody, err := getMap(resp.Body)
	if err != nil {
		return false, err
	}

	cdb, prs := jsonbody["db_name"]
	if !prs {

		cdb, prs = jsonbody["error"]
		if !prs {
			return false, errors.New("Could not decrypt Couchbase response")
		}
		e := cdb.(string)
		if e == "not_found" {
			return false, nil
		}
		return false, errors.New(e)
	}
	return true, nil
}

func init() {

	db := &Couchbase{}
	pluginregistry.RegisterDatabase(db)
}
