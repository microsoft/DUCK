package db

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"fmt"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	_ "github.com/Microsoft/DUCK/backend/plugins/mockdb"
)

var testDB *Database
var users map[string]structs.User
var documents map[string]struct {
	Pass     bool
	Document structs.Document
}
var dictionaries map[string]struct {
	Error bool
	Put   structs.Dictionary
	Want  structs.Dictionary
}

func TestDatabase(t *testing.T) {

	t.Run("NewDatabase", testNewDatabase)

	tDB, err := NewDatabase(structs.DBConf{Name: "Testname"})
	testDB = tDB
	if err != nil {
		t.Error("Cannot initialize test database")
	}
	users = make(map[string]structs.User)

	users = map[string]structs.User{
		"user a": structs.User{Email: "duckTEST@example.com", Password: "duck", Firstname: "TESTBDudley", Lastname: "Duck", Locale: "en"},
		"user b": structs.User{Email: "dä@example.com", Password: "<'//''$=äÜ", Firstname: "Sören", Lastname: "Duck", Locale: "de"},
		"user c": structs.User{Email: "dTEST@example.com", Password: "duck"},

		"user d": structs.User{Password: "duck", Firstname: "François", Lastname: "Duck", Locale: "fr"},
		"user e": structs.User{Email: "TEST@example.com", Firstname: "TEST", Lastname: "Duck", Locale: "en"},
		"user f": structs.User{},
		"user g": structs.User{Email: "duckTEST@example.com", Password: "otherpassword", Firstname: "TESTBDudley", Lastname: "Duck", Locale: "en"},
	}

	t.Run("PostUser", testDatabase_PostUser)
	t.Run("GetLogin", testDatabase_GetLogin)
	t.Run("GetUser", testDatabase_GetUser)
	t.Run("PutUser", testDatabase_PutUser)
	t.Run("DeleteUser", testDatabase_DeleteUser)

	t.Run("UserDicts", testDatabase_DICTS)

	if err := loadDocs(); err != nil {
		t.Error(err.Error())
		t.Skip("No testfixtures no Documenttests")
	}

	t.Run("PostDocument", testDatabase_PostDocument)
	t.Run("GetDocument", testDatabase_GetDocument)
	t.Run("GetDocumentSummariesForUser", testDatabase_GetDocumentSummariesForUser)
	t.Run("PutDocument", testDatabase_PutDocument)
	t.Run("DeleteDocument", testDatabase_DeleteDocument)

}

func loadDocs() error {
	dat, err := ioutil.ReadFile("../handlers/documents/testdata/document.json")
	if err = json.Unmarshal(dat, &documents); err != nil {
		return fmt.Errorf("Testfixture Documents not correctly loading")
	}

	return nil
}

func testNewDatabase(t *testing.T) {
	tests := []struct {
		name    string
		config  structs.DBConf
		want    *Database
		wantErr bool
	}{
		{"Just setting Config to a DBConf ", structs.DBConf{Name: "Testname"}, &Database{Config: structs.DBConf{Name: "Testname"}}, false},
		{"Empty config", structs.DBConf{}, &Database{}, false},
		{"INVALIDDBNAME", structs.DBConf{Name: "INVALIDDBNAME"}, &Database{}, true},
	}
	for _, tt := range tests {
		got, err := NewDatabase(tt.config)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. NewDatabase() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got.Config, tt.want.Config) || (got.db == nil && tt.want.db != nil) {
			t.Errorf("%q. NewDatabase() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func testDatabase_GetLogin(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name    string
		email   string
		wantID  string
		wantPw  string
		wantErr bool
	}{
		{"user a", users["user a"].Email, users["user a"].ID, users["user a"].Password, false},
		{"user b", users["user b"].Email, users["user b"].ID, users["user b"].Password, false},
		{"user c", users["user c"].Email, users["user c"].ID, users["user c"].Password, false},
		{"user d", users["user d"].Email, users["user d"].ID, users["user d"].Password, true},
		{"user e", users["user e"].Email, users["user e"].ID, users["user e"].Password, true},
		{"user f", users["user f"].Email, users["user f"].ID, users["user f"].Password, true},
		{"user g", users["user g"].Email, users["user a"].ID, users["user a"].Password, false}, //same id as user a since it is the same user
	}
	for _, tt := range tests {
		gotID, gotPw, err := testDB.GetLogin(tt.email)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Database.GetLogin() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if err == nil {
			if gotID != tt.wantID {
				t.Errorf("%q. database.GetLogin() gotID = %v, want %v", tt.name, gotID, tt.wantID)
			}
			if gotPw != tt.wantPw {
				t.Errorf("%q. database.GetLogin() gotPw = %v, want %v", tt.name, gotPw, tt.wantPw)
			}
		}
	}
}

func testDatabase_GetUser(t *testing.T) {
	tests := []struct {
		name    string
		userid  string
		want    structs.User
		wantErr bool
	}{
		{"user a", users["user a"].ID, users["user a"], false},
		{"user b", users["user b"].ID, users["user b"], false},
		{"user c", users["user c"].ID, users["user c"], false},
		{"user d", users["user d"].ID, users["user d"], true},
		{"user e", users["user e"].ID, users["user e"], true},
		{"user f", users["user f"].ID, users["user f"], true},
		{"user g", users["user g"].ID, users["user a"], true},
	}
	for _, tt := range tests {
		got, err := testDB.GetUser(tt.userid)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Database.GetUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if ((err == nil) || !tt.wantErr) && !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Database.GetUser() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func testDatabase_DeleteUser(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"user a", users["user a"].ID, false},
		{"user b", users["user b"].ID, false},
		{"user c", users["user c"].ID, false},
		{"user d", users["user d"].ID, true},
		{"user e", users["user e"].ID, true},
		{"user f", users["user f"].ID, true},
		{"user g", users["user g"].ID, true},
	}
	for _, tt := range tests {
		if err := testDB.DeleteUser(tt.id); (err != nil) != tt.wantErr {
			t.Errorf("%q. Database.DeleteUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func testDatabase_PutUser(t *testing.T) {

	for s, u := range users {
		u.Firstname = s
		users[s] = u
	}
	tests := []struct {
		name    string
		user    structs.User
		wantErr bool
	}{
		{"user a", users["user a"], false},
		{"user b", users["user b"], false},
		{"user c", users["user c"], false},
		{"user d", users["user d"], true},
		{"user e", users["user e"], true},
		{"user f", users["user f"], true},
		{"user g", users["user g"], true},
	}
	for _, tt := range tests {
		err := testDB.PutUser(tt.user)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Database.PutUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)

		}
		if err != nil {
			continue
		}
		got, err := testDB.GetUser(tt.user.ID)

		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Database.PutUser() verifying with GetUser(): error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.user) {
			t.Errorf("%q. Database.PutUser() verifying: database.GetUserDict() = %#v, want %#v", tt.name, got, tt.user)
		}
	}
}

func testDatabase_PostUser(t *testing.T) {
	tests := []struct {
		name    string
		user    structs.User
		wantErr bool
	}{
		{"user a", users["user a"], false},
		{"user b", users["user b"], false},
		{"user c", users["user c"], false},
		{"user d", users["user d"], true},

		{"user e", users["user e"], true},
		{"user f", users["user f"], true},
		{"user g", users["user g"], true},
	}
	for _, tt := range tests {
		gotID, err := testDB.PostUser(tt.user) //if we have an error we don't care, else id is returned
		u := tt.user
		u.ID = gotID
		users[tt.name] = u

		if (err != nil) != tt.wantErr {
			t.Errorf("%q. database.PostUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}

	}
}

func testDatabase_DICTS(t *testing.T) {

	ids := map[string]string{
		"a": "",
		"b": "",
		"c": "",
		"d": "",
	}
	for user := range ids {
		mail := user + "@example.com"
		id, _ := testDB.PostUser(structs.User{Email: mail, Password: "password", Lastname: user})
		ids[user] = id

	}

	defer func() {
		for _, id := range ids {
			if id == "" {
				continue
			}
			err := testDB.DeleteUser(id)
			if err != nil {
				t.Fatalf("Dictionary test cleanup failed: %s", err)
			}
		}
	}()

	dat, err := ioutil.ReadFile("testdata/dictionary.json")
	if err = json.Unmarshal(dat, &dictionaries); err != nil {
		t.Error("Testfixture User not correctly loading")
		t.Skip("No testfixtures no usertests")
	}

	test_database_PutUserDict := putUserDictClosure(ids)
	test_database_GetUserDict := getUserDictClosure(ids)

	t.Run("PutUserDict", test_database_PutUserDict)
	t.Run("GetUserDict", test_database_GetUserDict)

}

func getUserDictClosure(ids map[string]string) func(t *testing.T) {

	return func(t *testing.T) {

		tests := []struct {
			name    string
			userID  string
			want    structs.Dictionary
			wantErr bool
		}{
			{"Dict a", ids["a"], dictionaries["dict a"].Want, dictionaries["dict a"].Error},
			{"Dict b", ids["b"], dictionaries["dict b"].Want, dictionaries["dict b"].Error},
			{"Dict c", ids["c"], dictionaries["dict c"].Want, dictionaries["dict c"].Error},
			{"Dict d", ids["d"], structs.Dictionary(nil), false},
		}

		for _, tt := range tests {
			got, err := testDB.GetUserDict(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("%q. database.GetUserDict() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				continue
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q. database.GetUserDict() = %#v, want %#v", tt.name, got, tt.want)
			}
		}
	}
}

func putUserDictClosure(ids map[string]string) func(t *testing.T) {
	return func(t *testing.T) {

		tests := []struct {
			name    string
			userID  string
			dict    structs.Dictionary
			wantErr bool
		}{
			{"Dict a", ids["a"], dictionaries["dict a"].Put, false},
			{"Dict b", ids["b"], dictionaries["dict b"].Put, false},
			{"Dict c", ids["c"], dictionaries["dict c"].Put, false},
			{"Dict d", "-5", structs.Dictionary{}, true},
		}

		for _, tt := range tests {
			if err := testDB.PutUserDict(tt.dict, tt.userID); (err != nil) != tt.wantErr {
				t.Errorf("%q. database.PutUserDict() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		}
	}
}

func testDatabase_GetDocument(t *testing.T) {
	tests := []struct {
		name       string
		documentid string
		want       structs.Document
		wantErr    bool
	}{
		{"document_a", documents["document_a"].Document.ID, documents["document_a"].Document, (!documents["document_a"].Pass && documents["document_a"].Document.ID != "")},
		{"document_b", documents["document_b"].Document.ID, documents["document_b"].Document, (!documents["document_b"].Pass && documents["document_b"].Document.ID != "")},
		{"document_c", documents["document_c"].Document.ID, documents["document_c"].Document, (!documents["document_c"].Pass && documents["document_c"].Document.ID != "")},
		{"document_d", documents["document_d"].Document.ID, documents["document_d"].Document, (!documents["document_d"].Pass && documents["document_d"].Document.ID != "")},
		{"document_e", documents["document_e"].Document.ID, documents["document_e"].Document, (!documents["document_e"].Pass && documents["document_e"].Document.ID != "")},
		{"document_f", documents["document_f"].Document.ID, documents["document_f"].Document, (!documents["document_f"].Pass && documents["document_f"].Document.ID != "")},
		{"document_g", documents["document_g"].Document.ID, documents["document_g"].Document, (!documents["document_g"].Pass && documents["document_g"].Document.ID != "")},
		{"document_h", documents["document_h"].Document.ID, documents["document_h"].Document, (!documents["document_h"].Pass && documents["document_h"].Document.ID != "")},
		{"document_i", documents["document_i"].Document.ID, documents["document_i"].Document, (!documents["document_i"].Pass && documents["document_i"].Document.ID != "")},
	}
	for _, tt := range tests {
		got, err := testDB.GetDocument(tt.documentid)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Database.GetDocument() error = %v, wantErr %v, ID %v", tt.name, err, tt.wantErr, documents[tt.name].Document.ID)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Database.GetDocument() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func testDatabase_GetDocumentSummariesForUser(t *testing.T) {
	type args struct {
		userid string
	}
	tests := []struct {
		name     string
		database *Database
		args     args
		want     []structs.Document
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.database.GetDocumentSummariesForUser(tt.args.userid)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Database.GetDocumentSummariesForUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Database.GetDocumentSummariesForUser() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func testDatabase_DeleteDocument(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		database *Database
		args     args
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.database.DeleteDocument(tt.args.id); (err != nil) != tt.wantErr {
			t.Errorf("%q. Database.DeleteDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func testDatabase_PutDocument(t *testing.T) {
	type args struct {
		doc structs.Document
	}
	t.Error("E2")

	tests := []struct {
		name     string
		database *Database
		args     args
		wantErr  bool
	}{}
	for _, tt := range tests {
		if err := tt.database.PutDocument(tt.args.doc); (err != nil) != tt.wantErr {
			t.Errorf("%q. Database.PutDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func testDatabase_PostDocument(t *testing.T) {
	tests := []struct {
		name    string
		doc     structs.Document
		wantErr bool
	}{
		{"document_a", documents["document_a"].Document, !documents["document_a"].Pass},
		{"document_b", documents["document_b"].Document, !documents["document_b"].Pass},
		{"document_c", documents["document_c"].Document, !documents["document_c"].Pass},
		{"document_d", documents["document_d"].Document, !documents["document_d"].Pass},
		{"document_e", documents["document_e"].Document, !documents["document_e"].Pass},
		{"document_f", documents["document_f"].Document, !documents["document_f"].Pass},
		{"document_g", documents["document_g"].Document, !documents["document_g"].Pass},
		{"document_h", documents["document_h"].Document, !documents["document_h"].Pass},
		{"document_i", documents["document_i"].Document, !documents["document_i"].Pass},
	}

	for _, tt := range tests {
		gotID, err := testDB.PostDocument(tt.doc)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Database.PostDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if gotID != "" {
			d := documents[tt.name]
			d.Document.ID = gotID
			documents[tt.name] = d
		}
	}
}
