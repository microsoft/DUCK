// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package couchdb

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

create Rulebase		✓
get Rulebase		✓
update Rulebase		✓
delete Rulebase		✓


*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"io"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/Microsoft/DUCK/backend/pluginregistry"
)

//not using default Client: https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
var netClient = &http.Client{Timeout: time.Second * 10}

//Couchbase implements the pluginregistry.DBPlugin interface for CouchDB
type Couchbase struct {
	url      string
	database string
	user     string
	password string
	cookie   *http.Cookie
	auth     bool //Do we need submit auth info to cb?
}

// GetLogin returns ID and Password for the matching username from the couchbase Database
func (cb *Couchbase) GetLogin(email string) (id string, pw string, err error) {
	url := fmt.Sprintf("%s/%s/_design/app/_view/user_login?key=\"%s\"", cb.url, cb.database, email)
	//cb.url + "/" + cb.database + "/_design/app/_view/user?key='" + username + "'"

	bdy, err := cb.doGet(url)
	if err != nil {
		return "", "", err
	}
	rows, err := getRows(bdy)
	if err != nil {
		//Wrap error? http://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully

		if err.Error() == "No Data returned" {
			return "", "", structs.NewHTTPError("User not found", 404)
		}
		return "", "", err
	}
	if len(rows) > 1 {
		return "", "", structs.NewHTTPError("User not unique", 409)
	}

	row := rows[0].(map[string]interface{})

	pw, prs := row["value"].(string)
	if !prs || len(pw) <= 0 {
		return "", "", structs.NewHTTPError("Password not found", 401)
	}

	id, prs = row["id"].(string)
	if !prs || len(id) <= 0 {
		return "", "", structs.NewHTTPError("ID not found", 404)
	}

	//log.Printf("id: %s, pw: %s", id, pw)
	return

}

//GetUserDict returs the Global Dictionary of the specified user
func (cb *Couchbase) GetUserDict(id string) (structs.Dictionary, error) {

	mp, err := cb.getCouchbaseDocument(id)
	if err != nil {
		return nil, err
	}

	u := userFromValueMap(mp)
	//fmt.Printf("%+v\n", mp)
	return u.GlobalDictionary, nil

}

//GetUser returns a User with the sppecified ID from the Couchbase Database
func (cb *Couchbase) GetUser(id string) (structs.User, error) {

	var u structs.User
	mp, err := cb.getCouchbaseDocument(id)
	if err != nil {
		return u, err
	}

	u = userFromValueMap(mp)

	return u, nil

}

//GetDocument returns a Data use document with the sppecified ID from the Couchbase Database
func (cb *Couchbase) GetDocument(id string) (structs.Document, error) {
	var doc structs.Document
	mp, err := cb.getCouchbaseDocument(id)
	if err != nil {
		return doc, err
	}
	//fmt.Printf("%+v\n", mp)
	doc = docFromValueMap(mp)

	return doc, nil

}

/*
//GetRulebase returns a Rulebase with the sppecified ID from the Couchbase Database
func (cb *Couchbase) GetRulebase(id string) (document map[string]interface{}, err error) {

	return cb.getCouchbaseDocument(id)
}*/

func (cb *Couchbase) getCouchbaseDocument(cbDocID string) (document map[string]interface{}, err error) {

	url := fmt.Sprintf("%s/%s/%s", cb.url, cb.database, cbDocID)

	cbDoc, err := cb.doGet(url)

	if _, prs := cbDoc["_id"]; !prs {
		return nil, structs.NewHTTPError("No Data", 409)
	}
	return cbDoc, nil

}

//GetDocumentSummariesForUser returns a list all data use documents a user owns
//Summaries only include the documents name and ID
func (cb *Couchbase) GetDocumentSummariesForUser(userid string) ([]structs.Document, error) {
	url := fmt.Sprintf("%s/%s/_design/app/_view/documents_by_user?startkey=[\"%s\",\"\"]&endkey=[\"%s\",{}]",
		cb.url, cb.database, userid, userid)

	bdy, err := cb.doGet(url)
	if err != nil {
		return nil, err
	}

	rows, err := getRows(bdy)
	if err != nil {
		return nil, structs.WrapErrWith(err, structs.NewHTTPError(err.Error(), http.StatusNotFound))
	}

	var documents []structs.Document

	for _, intf := range rows {

		doc := intf.(map[string]interface{})
		var document structs.Document
		if id, ok := doc["id"]; ok {
			document.ID = id.(string)
		}

		if name, ok := doc["value"]; ok {
			document.Name = name.(string)
		}

		documents = append(documents, document)

	}

	return documents, nil

}

// DeleteDocument deletes a Data Use Document from the Couchbase Database
func (cb *Couchbase) DeleteDocument(id string) error {

	doc, err := cb.getCouchbaseDocument(id)
	if err != nil {

		return err
	}
	if rev, prs := doc["_rev"]; prs {
		err := cb.deleteCbDocument(id, rev.(string))
		if err != nil {
			return err
		}
		return nil
	}

	return structs.NewHTTPError("Could not delete Entry", 409)

}

// DeleteUser deletes a User from the Couchbase Database
func (cb *Couchbase) DeleteUser(id string) error {
	doc, err := cb.getCouchbaseDocument(id)
	if err != nil {

		return err
	}
	if rev, prs := doc["_rev"]; prs {
		err := cb.deleteCbDocument(id, rev.(string))
		if err != nil {
			return err
		}
		return nil
	}

	return structs.NewHTTPError("Could not delete Entry", 409)
}

/*
// DeleteRulebase deletes a User from the Couchbase Database
func (cb *Couchbase) DeleteRulebase(id string, rev string) error {
	return cb.deleteCbDocument(id, rev)
}*/

func (cb *Couchbase) deleteCbDocument(id string, rev string) error {
	url := fmt.Sprintf("%s/%s/%s?rev=%s", cb.url, cb.database, id, rev)

	jsonbody, err := cb.doRequest(http.MethodDelete, url, nil, false)
	if err != nil {
		return err
	}
	if _, prs := jsonbody["ok"]; prs {
		return nil
	}

	if err, prs := jsonbody["error"]; prs {
		reason := jsonbody["reason"].(string)
		e := "Error:" + err.(string) + ", Reason: " + reason
		return structs.NewHTTPError(e, 400)
	}

	return structs.NewHTTPError("Could not decrypt Couchbase response", 500)

}

// NewUser creates a new user in the couchbase Database
func (cb *Couchbase) NewUser(user structs.User) error {
	return cb.putUser(user)
}

//NewDocument creates a new Data use document in the couchbase Database
func (cb *Couchbase) NewDocument(doc structs.Document) error {
	return cb.putDocument(doc)
}

/*
//NewRulebase creates a new Rulebase in the couchbase Database
func (cb *Couchbase) NewRulebase(id string, entry string) error {
	return cb.putEntry(id, entry, "rulebase")
}*/

//UpdateUserDict updates an existing UserDict in the Couchbase database
func (cb *Couchbase) UpdateUserDict(dict structs.Dictionary, userID string) error {
	user, err := cb.GetUser(userID)
	if err != nil {
		return err
	}

	user.GlobalDictionary = dict

	//fmt.Printf("%+v", user)
	return cb.putUser(user)
}

//UpdateUser replaces an existing User in the Couchbase database
func (cb *Couchbase) UpdateUser(user structs.User) error {
	return cb.putUser(user)
}

//UpdateDocument replaces an existing Data Use Document in the Couchbase database
func (cb *Couchbase) UpdateDocument(doc structs.Document) error {
	return cb.putDocument(doc)
}

func (cb *Couchbase) putUser(u structs.User) error {
	entryMap := make(map[string]interface{})
	entryMap["type"] = "user"
	entryMap["_id"] = u.ID
	entryMap["email"] = u.Email
	entryMap["password"] = u.Password
	entryMap["firstname"] = u.Firstname
	entryMap["lastname"] = u.Lastname
	entryMap["locale"] = u.Locale
	entryMap["assumptionSet"] = u.AssumptionSet
	if u.Revision != "" {
		entryMap["_rev"] = u.Revision
	}
	dict := make(map[string]map[string]string)
	for key, val := range u.GlobalDictionary {
		dc := make(map[string]string)
		dc["value"] = val.Value
		dc["case_1"] = val.Case_1
		dc["case_2"] = val.Case_2
		dc["location"] = val.Location
		dc["type"] = val.Type
		dc["code"] = val.Code
		dc["category"] = val.Category
		dc["dictionaryType"] = val.DictionaryType

		dict[key] = dc
	}
	entryMap["dictionary"] = dict
	return cb.putEntry(entryMap, false)

}

func (cb *Couchbase) putDocument(d structs.Document) error {
	entryMap := make(map[string]interface{})
	entryMap["type"] = "document"
	entryMap["_id"] = d.ID
	entryMap["name"] = d.Name
	entryMap["owner"] = d.Owner
	entryMap["locale"] = d.Locale
	entryMap["assumptionSet"] = d.AssumptionSet
	entryMap["description"] = d.Description
	if d.Revision != "" {
		entryMap["_rev"] = d.Revision
	}
	dict := make(map[string]map[string]string)
	for key, val := range d.Dictionary {
		dc := make(map[string]string)
		dc["value"] = val.Value
		dc["case_1"] = val.Case_1
		dc["case_2"] = val.Case_2
		dc["location"] = val.Location
		dc["type"] = val.Type
		dc["code"] = val.Code
		dc["category"] = val.Category
		dc["dictionaryType"] = val.DictionaryType

		dict[key] = dc
	}
	entryMap["dictionary"] = dict

	var stmts []map[string]interface{}
	for _, statement := range d.Statements {
		stmt := make(map[string]interface{})
		stmt["useScopeCode"] = statement.UseScopeCode
		stmt["qualifierCode"] = statement.QualifierCode
		stmt["dataCategoryCode"] = statement.DataCategoryCode
		stmt["sourceScopeCode"] = statement.SourceScopeCode
		stmt["actionCode"] = statement.ActionCode
		stmt["resultScopeCode"] = statement.ResultScopeCode
		stmt["trackingId"] = statement.TrackingID
		if statement.Tag != nil {
			stmt["tag"] = *statement.Tag
		}

		if statement.DataCategories != nil {
			var dcats []map[string]interface{}
			for _, cat := range statement.DataCategories {
				dcat := make(map[string]interface{})
				dcat["qualifierCode"] = cat.QualifierCode
				dcat["dataCategoryCode"] = cat.DataCategoryCode
				dcat["operator"] = cat.Op

				dcats = append(dcats, dcat)
			}
			stmt["dataCategories"] = dcats
		}

		if statement.Passive {
			stmt["passive"] = "true"
		} else {
			stmt["passive"] = "false"
		}

		stmts = append(stmts, stmt)
	}
	if stmts != nil {
		entryMap["statements"] = stmts
	}

	return cb.putEntry(entryMap, false)
}

func (cb *Couchbase) putEntry(entry map[string]interface{}, designfile bool) error {

	var entryReader io.Reader
	var url string
	//set url which is different when we want to create a new designfile
	if !designfile {
		var entryBytes []byte

		entryBytes, err := json.Marshal(entry)
		if err != nil {
			return structs.WrapErrWith(err, structs.NewHTTPError(err.Error(), 500))
		}
		entryReader = bytes.NewReader(entryBytes)
		url = fmt.Sprintf("%s/%s/%s", cb.url, cb.database, entry["_id"])

	} else {

		//if we want to create a new design file, the whole file is in one entry in the entry map
		//the key is entry.
		//This is a hack which might have to be fixed somewhen
		if _, ok := entry["entry"]; !ok {
			return structs.NewHTTPError("Expected a map with one entry: \"entry\" but did not get it.", 409)
		}
		entryReader = strings.NewReader(entry["entry"].(string))
		url = fmt.Sprintf("%s/%s/_design/app", cb.url, cb.database)
	}
	//submit data &
	//check if we succeeded, couchdb answers with a map containing the field "ok"
	jsonbody, err := cb.doRequest(http.MethodPut, url, entryReader, false)
	if err != nil {
		return err
	}

	if _, prs := jsonbody["ok"]; prs {
		return nil
	}

	if err, prs := jsonbody["error"]; prs {
		reason := jsonbody["reason"].(string)

		e := "Error:" + err.(string) + ", Reason: " + reason
		return structs.NewHTTPError(e, 400)
	}

	return structs.NewHTTPError("Could not understand Couchbase response", 502)
}

//Init initializes the Couchbase DB & tests for connection errors
func (cb *Couchbase) Init(config structs.DBConf) error {
	log.Println("Couchbase initialization")

	//this is an design document with th id "_design/app" and some view functions that are used to access the user and documents
	//the view functions are used in this source code file and should be therefore present
	designDoc := `{"_id":"_design/app","views":{"foo":{"map":"function(doc){ emit(doc._id, doc._rev)}"},` +
		`"by_date":{"map":"function(doc) { if(doc.date && doc.title) {   emit(doc.date, doc.title);  }}"},` +
		`"user_login":{"map":"function(doc) { if(doc.type =='user') {   emit(doc.email,  doc.password);  }}"},` +
		`"user":{"map":"function(doc) { if(doc.type =='user') {   emit(doc._id, doc);  }}"},` +
		`"documents":{"map":"function(doc) { if(doc.type =='document') {   emit(doc._id, doc);  }}"},` +
		`"rulebases":{"map":"function(doc) { if(doc.type =='rulebase') {   emit(doc._id, doc._rev);  }}"},` +
		`"documents_by_user":{"map":"function(doc) { if(doc.type =='document') {   emit([doc.owner, doc._id], doc.name);  }}"}},` +
		`"language":"javascript"}`

	designMap := map[string]interface{}{"entry": designDoc}
	port := strconv.Itoa(config.Port)
	if port == "" {
		return structs.NewHTTPError("couchDB needs an port entry in config", 400)
	}
	if config.Location == "" {
		return structs.NewHTTPError("couchDB needs an url entry in config", 400)
	}

	if config.Username != "" && config.Password != "" {
		cb.user = config.Username
		cb.password = config.Password
		cb.auth = true
		log.Println("Database Username and Password set.")
	} else {
		cb.auth = false
		log.Println("Warning: Username or password missing in couchdb config. Assuming no auth needed. This is *not* recommended.")
	}

	url := config.Location + ":" + port

	jsonbody, err := cb.doGet(url)
	if err != nil {
		return structs.WrapErrWith(err, structs.NewHTTPError("Could not connect to CouchDB, please check if CouchDB is available", 500))
	}

	cdb, prs := jsonbody["couchdb"].(string)
	if !prs || cdb != "Welcome" {

		return structs.NewHTTPError("Connection to couchdb failed", 400)
	}
	cb.url = url

	if config.Name == "" {
		return structs.NewHTTPError("Couchdb needs a database entry in the config file to know the name of the database", 400)
	}
	cb.database = config.Name

	log.Printf("couchdb url is: %s; database name is: %s.", cb.url, cb.database)
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

		err := cb.putEntry(designMap, true)
		if err != nil {
			log.Printf("ERROR: %#+v\n", err)
		}
	}

	log.Println("Testextension initialized")
	return nil
}

func (cb *Couchbase) testFileExists(id string) (bool, error) {
	url := fmt.Sprintf("%s/%s/%s", cb.url, cb.database, id)

	jsonbody, err := cb.doGet(url)
	if err != nil {
		return false, err
	}

	if _, prs := jsonbody["_id"]; !prs {

		cdb, prs := jsonbody["error"]
		if !prs {
			return false, structs.NewHTTPError("Could not decrypt Couchbase response", 502)
		}
		e := cdb.(string)
		if e == "not_found" {
			return false, nil
		}
		return false, structs.NewHTTPError(e, 502)
	}
	return true, nil

}

func (cb *Couchbase) createDatabase() error {
	url := fmt.Sprintf("%s/%s", cb.url, cb.database)

	jsonbody, err := cb.doRequest(http.MethodPut, url, nil, false)
	if err != nil {
		return err
	}
	if _, prs := jsonbody["ok"]; prs {
		return nil
	}

	if _, prs := jsonbody["error"]; prs {
		reason := jsonbody["reason"].(string)

		return structs.NewHTTPError(reason, 502)
	}

	return structs.NewHTTPError("Could not decrypt CouchDB response", 502)
}

func (cb *Couchbase) testDBExists() (bool, error) {
	url := fmt.Sprintf("%s/%s", cb.url, cb.database)

	jsonbody, err := cb.doGet(url)
	if err != nil {

		return false, structs.WrapErrWith(err, structs.NewHTTPError("Could not connect to CouchDB, please check if CouchDB is available", 500))
	}

	cdb, prs := jsonbody["db_name"]
	if !prs {

		cdb, prs = jsonbody["error"]
		if !prs {
			return false, structs.NewHTTPError("Could not decrypt CouchDB response", 502)
		}
		e := cdb.(string)
		if e == "not_found" {
			return false, nil
		}
		return false, structs.NewHTTPError(e, 502)
	}
	return true, nil
}

func (cb *Couchbase) login() error {

	fmt.Println("login()")
	data := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{
		cb.user,
		cb.password,
	}
	databytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(databytes)

	request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:5984/_session", bodyReader)
	if err != nil {
		return structs.WrapErrWith(err, structs.NewHTTPError(err.Error(), 502))
	}
	request.SetBasicAuth(cb.user, cb.password)
	request.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	jsonbody, err := getMap(resp.Body)

	if err != nil {
		return err
	}

	if success, ok := jsonbody["ok"]; !ok || success != true {
		return structs.NewHTTPError("Login to couchdb failed", 400)
	}

	//fmt.Printf("%#v\n", resp.Cookies()[0])

	cb.cookie = resp.Cookies()[0]
	return nil
}

func (cb *Couchbase) doRequest(method, url string, body io.Reader, dologin bool) (map[string]interface{}, error) {

	if dologin {
		if err := cb.login(); err != nil {
			return nil, structs.NewHTTPError("Login to couchdb failed", 400)
		}
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, structs.WrapErrWith(err, structs.NewHTTPError(err.Error(), 502))
	}
	if cb.auth && cb.cookie != nil {
		//fmt.Printf("%s=%s counter: %d url: %s\n", cb.cookie.Name, cb.cookie.Value, counter, url)

		request.Header.Add("Cookie", fmt.Sprintf("%s=%s", cb.cookie.Name, cb.cookie.Value))

	}

	resp, err := netClient.Do(request)

	if err != nil {
		return nil, structs.WrapErrWith(err, structs.NewHTTPError(err.Error(), 502))
	}
	defer resp.Body.Close()
	//check if we succeeded, couchdb answers with a map containing the field "ok"
	jsonbody, err := getMap(resp.Body)
	if err != nil {
		return nil, err
	}
	if cb.auth && !dologin {
		if autherr, prs := jsonbody["reason"]; prs && autherr == "Authentication required." {
			return cb.doRequest(method, url, body, true)
		}
	}
	return jsonbody, nil
}

func (cb *Couchbase) doGet(url string) (map[string]interface{}, error) {
	return cb.doRequest(http.MethodGet, url, nil, false)
}

func init() {

	db := &Couchbase{}
	pluginregistry.RegisterDatabase(db)
	log.Println("Couchdb registered")
}
