/*
create DB


create User
get user
update user
delete user

create document
get document
update document
delete document

create RuleSet
get RuleSet
update RuleSet
delete RuleSet


*/
package couchbase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Metaform/duck/backend/pluginregistry"
)

type Couchbase struct {
	url      string
	database string
}

func getMap(resp io.Reader) (map[string]interface{}, error) {

	content, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	//TODO: remove this when not used anymore
	fmt.Println(string(content))

	var data map[string]interface{}

	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}
	return data, nil
}

//Print prints sth
func (cb *Couchbase) Print() {
	fmt.Println("Testextension Printed sth")
}

//Save saves sth
func (cb *Couchbase) Save() {
	fmt.Println("Testextension saved sth")
}

func (cb *Couchbase) GetLogin(username string) (string, error) {
	url := fmt.Sprintf("%s/%s/_design/app/_view/user?key=\"%s\"", cb.url, cb.database, username)
	//cb.url + "/" + cb.database + "/_design/app/_view/user?key='" + username + "'"
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	jsonbody, err := getMap(resp.Body)
	if err != nil {
		return "", err
	}
	rows, prs := jsonbody["rows"].([]interface{})

	if !prs || len(rows) != 1 {
		return "", errors.New("User not found")
	}
	row := rows[0].(map[string]interface{})

	pw, prs := row["value"].(string)
	if !prs || len(pw) <= 0 {
		return "", errors.New("Password not found")
	}
	return pw, nil

}

//Init initializes the Couchbase DB & tests for connection errors
func (cb *Couchbase) Init(url string, database string) error {
	fmt.Println("Couchase initialization")
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

	fmt.Println("Testextension initialized")
	return nil
}

func init() {
	db := &Couchbase{}
	pluginregistry.RegisterDatabase(db)
}
