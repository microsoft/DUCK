package ducklib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	_ "github.com/Microsoft/DUCK/backend/plugins/mockdb"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

/*
TODO

checkDocHandler(c echo.Context)
checkDocIDHandler(c echo.Context)
getRulebasesHandler(c echo.Context)


*/
var (
	conf structs.Configuration
	e    *echo.Echo

	users map[string]struct {
		Pass bool         `json:"pass"`
		User structs.User `json:"user"`
	}
	//eg document["document_a"]
	documents map[string]struct {
		Pass     bool             `json:"pass"`
		Document structs.Document `json:"document"`
	}
	dicts struct {
		User         structs.User
		Dictionaries map[string]struct {
			Pass     bool               `json:"pass"`
			EntryIn  structs.Dictionary `json:"entryIn"`
			EntryOut structs.Dictionary `json:"entryOut"`
		} `json:"dictionaries"`
		Entries map[string]struct {
			EntryIn  structs.DictionaryEntry `json:"entryIn"`
			EntryOut structs.DictionaryEntry `json:"entryOut"`
		} `json:"entries"`
	}

	//eg userIDs["user_a"]="a structs.User.ID"
	userIDs     = make(map[string]string)
	documentIDs = make(map[string]string)
	//documents owners and # of documents they own
	owners = make(map[string]int)
)

//Move this into main.go
func TestMain(m *testing.M) {

	conf = NewConfiguration(filepath.Join(os.Getenv("GOPATH"), "/src/github.com/Microsoft/DUCK/backend/configuration.json"))

	e = GetServer(conf)
	if e != nil {

		os.Exit(m.Run())
	}

}


func TestUserHandler(t *testing.T) {

	dat, err := ioutil.ReadFile("structs/testdata/user.json")

	if err = json.Unmarshal(dat, &users); err != nil {
		t.Error("Testfixture User not correctly loading")
		t.Skip("No testfixtures no usertests")
	}
	//t.Logf("User %+v\n", users)
	//t.Error("AHHHHHH")
	t.Run("PostUser=1", testPostUserHandler)
	t.Run("PostUser=2", testPostUserHandlerAgain)
	t.Run("Login=1", testLoginHandler)
	t.Run("login=2", testWrongLogin)
	t.Run("PutUser=1", testPutUserHandler)
	t.Run("DeleteUser=1", testDeleteUserHandler)

}

func TestDocumentHandler(t *testing.T) {

	dat, err := ioutil.ReadFile("structs/testdata/document.json")

	if err = json.Unmarshal(dat, &documents); err != nil {
		t.Error("Testfixture Document not correctly loading")
		t.Skip("No testfixtures no documenttests")
	}
	//set documentowners

	for _, value := range documents {
		if o := value.Document.Owner; o != "" && value.Pass {
			owners[o]++
		}
	}

	//t.Logf("User %+v\n", documents)
	t.Run("PostDocument=1", testpostDocHandler)
	t.Run("GetDocument=1", testGetDocHandler)
	t.Run("Summaries=1", testGetDocSummaries)
	t.Run("CopyDocument=1", testCopyDocHandler)
	t.Run("PutDocument=1", testPutDocHandler)

	t.Run("DeleteDocument=1", testDeleteDocHandler)
	//t.Error("AHHHHHH")

}

func TestDictionaryHandler(t *testing.T) {

	dat, err := ioutil.ReadFile("structs/testdata/dictionary.json")

	if err = json.Unmarshal(dat, &dicts); err != nil {
		t.Error("Testfixture Dictionary not correctly loading")
		t.Skip("No testfixtures no dictionary tests")
	}
	id, err := datab.PostUser(dicts.User)
	if err != nil {
		t.Skip("Not able to save user to mockdb in dictionary test. Skipping tests..")
	}
	defer func() {
		err := datab.DeleteUser(id)
		if err != nil {
			t.Log("Could not delete user from mockDB in dictionary test this can interfere with other tests")
		}
	}()
	user, err := datab.GetUser(id)
	if err != nil {
		t.Skip("Not able to read user from mockdb in dictionary test. Skipping tests ..")
	}
	dicts.User = user

	//t.Logf("User %+v\n", documents)
	/*	t.Run("PostDocument=1", testpostDocHandler)
		t.Run("GetDocument=1", testGetDocHandler)
		t.Run("Summaries=1", testGetDocSummaries)
		t.Run("CopyDocument=1", testCopyDocHandler)
		t.Run("PutDocument=1", testPutDocHandler)

		t.Run("DeleteDocument=1", testDeleteDocHandler)*/

}

//userhandlertests
func testPostUserHandler(t *testing.T) {

	/*e := GetServer(conf, goPath)
	if e == nil {
		t.Fatal("Get Server Failed")
	}*/
	for key, value := range users {

		userJSON, err := json.Marshal(value.User)
		if err != nil {
			t.Errorf("Test with %s: user Post into json Marshalling not functioning", key)

		}

		req, err := http.NewRequest(echo.POST, "/users", bytes.NewReader(userJSON))
		if err != nil {
			t.Errorf("Test with %s: Error creating User: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
		err = postUserHandler(c)
		if err != nil {
			t.Errorf("Test with %s: Error creating User during post:%s", key, err)
		}

		if value.Pass {

			if rec.Code != http.StatusCreated {
				t.Errorf("Test with %s: user creation does not return HTTP code %d but %d.", key, http.StatusCreated, rec.Code)
			} else {

				// compare with user fields since some fields are unique

				var res structs.User

				if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
					t.Errorf("Test with %s: user Creation does not return valid User struct", key)
				}

				if res.Email != value.User.Email {
					t.Errorf("Test with %s: User creation returns User email %s, wants %s", key, res.Email, value.User.Email)
				}
				if res.Firstname != value.User.Firstname {
					t.Errorf("Test with %s: User creation returns User Firstname %s, wants %s", key, res.Firstname, value.User.Firstname)
				}
				if res.Lastname != value.User.Lastname {
					t.Errorf("Test with %s: User creation returns User Lastname %s, wants %s", key, res.Lastname, value.User.Lastname)
				}
				//Passwords are not retuned anymore
				/*if res.Password != value.User.Password {
					t.Errorf("Test with %s: User creation returns User Password %s, wants %s", key, res.Password, value.User.Password)
				}*/
				userIDs[key] = res.ID
				value.User.ID = res.ID
			}

		} else { //  test missing values

			if rec.Code < 400 {
				t.Errorf("Test with %s: user creation does not return a HTTP error code (>=400) but %d.", key, rec.Code)
			}
		}

	}

}

func testPostUserHandlerAgain(t *testing.T) {

	//test if already existing user is not saved again
	for key, value := range users {
		userJSON, err := json.Marshal(value.User)
		if err != nil {
			t.Errorf("Test with %s: second user Post  into json Marshalling not functioning", key)

		}
		req, err := http.NewRequest(echo.POST, "/users", bytes.NewReader(userJSON))
		if err != nil {
			t.Errorf("Test with %s: Error creating User: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
		err = postUserHandler(c)
		if err != nil {
			t.Errorf("Test with %s: Error creating User during post:%s", key, err)
		}

		if rec.Code < 400 {
			t.Errorf("Test with %s: user creation does not return a HTTP error code (>=400) but %d.", key, rec.Code)
		}

	}

}

func testLoginHandler(t *testing.T) {
	for key, value := range users {

		value.User.ID = userIDs[key]

		login := structs.Login{Email: value.User.Email, Password: value.User.Password}

		userJSON, err := json.Marshal(login)
		if err != nil {
			t.Errorf("Test with %s: user login json Marshal not functioning", key)

		}

		req, err := http.NewRequest(echo.POST, "/login", bytes.NewReader(userJSON))

		if err != nil {
			t.Errorf("Test with %s: Error creating User: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		err = loginHandler(c)

		if value.Pass {
			if err != nil {
				t.Errorf("Test with %s: Error logging in: %s", key, err)
			}
			if rec.Code != http.StatusOK {
				t.Errorf("Test with %s: user Login does not return HTTP code %d but %d.", key, http.StatusOK, rec.Code)
			} else {

				var dat map[string]interface{}

				//log.Println(rec.Body.String())

				if err := json.Unmarshal(rec.Body.Bytes(), &dat); err != nil {
					t.Errorf("Test with %s: user login does not return valid JSON", key)
				}

				if _, prs := dat["token"]; !prs {
					t.Errorf("Test with %s: user login does not return token", key)

				}

				if s, prs := dat["firstName"]; prs {
					if value.User.Firstname != s.(string) {
						t.Errorf("Test with %s: User login returns User Firstname %s, wants %s", key, s.(string), value.User.Firstname)
					}
				} else {
					t.Errorf("Test with %s: User login does not return Firstname", key)

				}

				if s, prs := dat["lastName"]; prs {
					if value.User.Lastname != s.(string) {
						t.Errorf("Test with %s: User login returns User Lastname %s, wants %s", key, s.(string), value.User.Lastname)
					}
				} else {
					t.Errorf("Test with %s: User login does not return Lastname", key)

				}

				if s, prs := dat["id"]; prs {
					if value.User.ID != s.(string) {
						t.Errorf("Test with %s: User login returns User id %s, wants %s", key, s.(string), value.User.ID)
					}
				} else {
					t.Errorf("Test with %s: User login does not return id", key)

				}

			}
		} else {
			if rec.Code != echo.ErrUnauthorized.Code {
				t.Errorf("Test with %s: user login does not return HTTP code %d but %d.", key, echo.ErrUnauthorized.Code, rec.Code)
			}
		}

	}
}

func testWrongLogin(t *testing.T) {
	key := "user_a"
	value := users[key]
	value.User.Password += "WrongPassword"
	logins := make([]interface{}, 4)

	logins[0] = structs.Login{Email: value.User.Email, Password: value.User.Password}
	logins[1] = structs.Login{Email: value.User.Email}
	docs := make([]structs.Document, 2)
	logins[2] = structs.Response{Documents: &docs}
	logins[3] = "teststring"

	for i, login := range logins {
		userJSON, err := json.Marshal(login)
		if err != nil {
			t.Errorf("Test with login %d: user login with wrong password json Marshal not functioning", i)

		}

		req, err := http.NewRequest(echo.POST, "/login", bytes.NewReader(userJSON))

		if err != nil {
			t.Errorf("Test with login %d: Error login with wrong password: %s", i, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		err = loginHandler(c)

		if rec.Code != echo.ErrUnauthorized.Code && rec.Code != echo.ErrNotFound.Code {
			t.Errorf("Test with login %d: wrong user login test does not return HTTP code %d or %d but %d.", i, echo.ErrUnauthorized.Code, echo.ErrNotFound.Code, rec.Code)
		}
	}
}

func testPutUserHandler(t *testing.T) {
	for key, value := range users {
		if !value.Pass {
			continue
		}
		value.User.ID = userIDs[key]

		value.User.Firstname = fmt.Sprintf("xx%s~", value.User.Firstname)

		userJSON, err := json.Marshal(value.User)
		if err != nil {
			t.Errorf("Test with %s: user update json Marshal not functioning", key)

		}

		req, err := http.NewRequest(echo.PUT, "/users/:id", bytes.NewReader(userJSON))

		if err != nil {
			t.Errorf("Test with %s: Error creating User: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
		c.SetParamNames("id")
		c.SetParamValues(value.User.ID)
		err = putUserHandler(c)
		if err != nil {
			t.Errorf("Test with %s: Error creating User during post:%s", key, err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Test with %s: user update does not return HTTP code %d but %d.", key, http.StatusOK, rec.Code)
		} else {

			var res structs.User

			if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
				t.Errorf("Test with %s: user update does not return valid User struct", key)
			}

			if res.Email != value.User.Email {
				t.Errorf("Test with %s: User update returns User email %s, wants %s", key, res.Email, value.User.Email)
			}
			if res.Firstname != value.User.Firstname {
				t.Errorf("Test with %s: User update returns User Firstname %s, wants %s", key, res.Firstname, value.User.Firstname)
			}
			if res.Lastname != value.User.Lastname {
				t.Errorf("Test with %s: User update returns User Lastname %s, wants %s", key, res.Lastname, value.User.Lastname)
			}
			if res.Password != value.User.Password {
				t.Errorf("Test with %s: User update returns User Password %s, wants %s", key, res.Password, value.User.Password)
			}
		}

	}
}

func testDeleteUserHandler(t *testing.T) {
	for key, value := range users {

		req, err := http.NewRequest(echo.DELETE, "/users/:id", nil)
		if err != nil {
			t.Errorf("Test with %s: Error deleting User: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		c.SetParamNames("id")
		c.SetParamValues(userIDs[key])

		err = deleteUserHandler(c)
		if err != nil {
			t.Errorf("Test with %s: Error deleting User during post:%s", key, err)
		}

		if value.Pass {

			if rec.Code != http.StatusOK {
				t.Errorf("Test with %s: user deletion does not return HTTP code %d but %d.", key, http.StatusOK, rec.Code)
			}

		} else { //  test missing user

			if rec.Code != http.StatusNotFound {
				t.Errorf("Test with %s: user deletion does not return HTTP code %d but %d.", key, http.StatusNotFound, rec.Code)
			}
		}

	}
}

func testpostDocHandler(t *testing.T) {

	for key, value := range documents {

		userJSON, err := json.Marshal(value.Document)
		if err != nil {
			t.Errorf("Test with %s: Document Post into json Marshalling not functioning", key)

		}

		req, err := http.NewRequest(echo.POST, "/documents", bytes.NewReader(userJSON))
		if err != nil {
			t.Errorf("Test with %s: Error creating Document: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
		err = postDocHandler(c)
		if err != nil {
			t.Errorf("Test with %s: Error creating Document during post:%s", key, err)
		}

		if value.Pass {

			if rec.Code != http.StatusCreated {
				t.Errorf("Test with %s: Document creation does not return HTTP code %d but %d.", key, http.StatusCreated, rec.Code)
			} else {

				// compare with user fields since some fields are unique

				var res structs.Document

				if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
					t.Errorf("Test with %s: Document Creation does not return valid Document struct", key)
				}

				if res.Name != value.Document.Name {
					t.Errorf("Test with %s: Document creation returns Document Name %s, wants %s", key, res.Name, value.Document.Name)
				}
				if res.Owner != value.Document.Owner {
					t.Errorf("Test with %s: Document creation returns Document Owner %s, wants %s", key, res.Owner, value.Document.Owner)
				}
				if len(res.Statements) != len(value.Document.Statements) {
					t.Errorf("Test with %s: Document creation returns  Document with %d Statements, wants %d", key, len(res.Statements), len(value.Document.Statements))
				}
				if res.Locale != value.Document.Locale {
					t.Errorf("Test with %s: Document creation returns Document Locale %s, wants %s", key, res.Locale, value.Document.Locale)
				}
				documentIDs[key] = res.ID
				value.Document.ID = res.ID
			}

		} else { //  test missing values

			//this might not always be 500
			if rec.Code != http.StatusBadRequest {
				t.Errorf("Test with %s: Document creation does not return HTTP code %d but %d.", key, http.StatusBadRequest, rec.Code)
			}
		}

	}

}

func testPutDocHandler(t *testing.T) {

	for key, value := range documents {
		if !value.Pass {
			continue
		}
		value.Document.ID = documentIDs[key]

		value.Document.Name = fmt.Sprintf("xx%s~", value.Document.Name)

		docJSON, err := json.Marshal(value.Document)
		if err != nil {
			t.Errorf("Test with %s: user update json Marshal not functioning", key)

		}

		req, err := http.NewRequest(echo.PUT, "/documents/", bytes.NewReader(docJSON))

		if err != nil {
			t.Errorf("Test with %s: Error creating User: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
		err = putDocHandler(c)
		if err != nil {
			t.Errorf("Test with %s: Error creating User during post:%s", key, err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Test with %s: document update does not return HTTP code %d but %d.", key, http.StatusOK, rec.Code)
		} else {

			var res structs.Document

			if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
				t.Errorf("Test with %s: document update does not return valid User struct", key)
			}

			if res.Name != value.Document.Name {
				t.Errorf("Test with %s: Document creation returns Document Name %s, wants %s", key, res.Name, value.Document.Name)
			}
			if res.Owner != value.Document.Owner {
				t.Errorf("Test with %s: Document creation returns Document Owner %s, wants %s", key, res.Owner, value.Document.Owner)
			}
			if len(res.Statements) != len(value.Document.Statements) {
				t.Errorf("Test with %s: Document creation returns  Document with %d Statements, wants %d", key, len(res.Statements), len(value.Document.Statements))
			}
			if res.Locale != value.Document.Locale {
				t.Errorf("Test with %s: Document creation returns Document Locale %s, wants %s", key, res.Locale, value.Document.Locale)
			}
		}

	}
}

func testGetDocHandler(t *testing.T) {
	for key, value := range documents {

		value.Document.ID = documentIDs[key]

		req, err := http.NewRequest(echo.GET, "/documents/:docid", nil)

		if err != nil {
			t.Errorf("Test with %s: Error getting Document: %s", key, err)
			continue
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		c.SetParamNames("docid")
		c.SetParamValues(documentIDs[key])
		err = getDocHandler(c)
		if err != nil {
			t.Errorf("Test with %s: Error getting Document during get:%s", key, err)

			continue
		}

		if value.Pass {

			if rec.Code != http.StatusOK {
				t.Errorf("Test with %s: Document get does not return HTTP code %d but %d.", key, http.StatusOK, rec.Code)
			} else {

				var res structs.Document

				if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
					t.Errorf("Test with %s:  Document get does not return valid Document struct", key)
				}

				if res.Name != value.Document.Name {
					t.Errorf("Test with %s:  Document get returns Document Name %s, wants %s", key, res.Name, value.Document.Name)
				}
				if res.Owner != value.Document.Owner {
					t.Errorf("Test with %s:  Document get returns Document Owner %s, wants %s", key, res.Owner, value.Document.Owner)
				}
				if len(res.Statements) != len(value.Document.Statements) {
					t.Errorf("Test with %s:  Document get returns  Document with %d Statements, wants %d", key, len(res.Statements), len(value.Document.Statements))
				}
				if res.Locale != value.Document.Locale {
					t.Errorf("Test with %s:  Document get returns Document Locale %s, wants %s", key, res.Locale, value.Document.Locale)
				}

			}

		} else { //  test missing values

			//this might not always be 500
			if rec.Code != http.StatusNotFound {
				t.Errorf("Test with %s:  Document get does not return HTTP code %d but %d.", key, http.StatusNotFound, rec.Code)
			}
		}

	}
}

func testCopyDocHandler(t *testing.T) {

	var copys []structs.Document

	for key, value := range documents {

		/*
			copyDocHandler is called with docID as param and a Document as POST load.
			This POSTed document has a (new) locale, name and owner.
			The statements have to be copied.

			So we need just the statements from the original documents and have to invent the rest
		*/

		value.Document.ID = documentIDs[key]

		cp := structs.Document{Name: fmt.Sprintf("NEW%s", value.Document.Name),
			Locale: "de", Owner: value.Document.Owner}

		docJSON, err := json.Marshal(cp)
		if err != nil {
			t.Errorf("Test with %s: user copy json Marshal not functioning", key)

		}
		req, err := http.NewRequest(echo.POST, "/documents/copy/:docid", bytes.NewReader(docJSON))

		if err != nil {
			t.Errorf("Test with %s: Error posting Document: %s", key, err)
			continue
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		c.SetParamNames("docid")
		c.SetParamValues(documentIDs[key])
		err = copyStatementsHandler(c)
		if err != nil {
			t.Errorf("Test with %s: Error posting Document during copy:%s", key, err)

			continue
		}

		if value.Pass {

			if rec.Code != http.StatusOK {
				t.Errorf("Test with %s: Document copy does not return HTTP code %d but %d.", key, http.StatusOK, rec.Code)
			} else {

				var res structs.Document

				if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
					t.Errorf("Test with %s:  Document get does not return valid Document struct", key)

				}

				if res.ID == documentIDs[key] {
					t.Errorf("Test with %s:  Document copy returns Document ID %s which is the same as %s, wants a different one", key, res.Name, documentIDs[key])
				}

				if res.Name != cp.Name {
					t.Errorf("Test with %s:  Document copy returns Document Name %s, wants %s", key, res.Name, value.Document.Name)
				}
				if res.Owner != value.Document.Owner {
					t.Errorf("Test with %s:  Document copy returns Document Owner %s, wants %s", key, res.Owner, value.Document.Owner)
				}
				if len(res.Statements) != len(value.Document.Statements) {
					t.Errorf("Test with %s:  Document copy returns  Document with %d Statements, wants %d", key, len(res.Statements), len(value.Document.Statements))
				}
				if res.Locale != cp.Locale {
					t.Errorf("Test with %s:  Document copy returns Document Locale %s, wants %s", key, res.Locale, value.Document.Locale)
				}
				//add to documents
				copys = append(copys, res)

			}

		} else { //  test missing values

			//this might not always be 500
			if rec.Code != http.StatusNotFound {
				t.Errorf("Test with %s:  Document get does not return HTTP code %d but %d.", key, http.StatusNotFound, rec.Code)
			}
		}

	}
	count := 1
	for _, val := range copys {
		name := fmt.Sprintf("document_copy_%i", count)
		cp := documents[name]
		cp.Pass = true
		cp.Document = val
		documents[name] = cp
		documentIDs[name] = val.ID
		count++
	}

}

func testGetDocSummaries(t *testing.T) {
	for owner, count := range owners {

		req, err := http.NewRequest(echo.GET, "/documents/:userid/summary", nil)

		if err != nil {
			t.Errorf("Test with %s: Error creating request for summary: %s", owner, err)
			continue
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
		c.SetParamNames("userid")
		c.SetParamValues(owner)
		err = getDocSummaries(c)
		if err != nil {
			t.Errorf("Test with %s: Error getting summary during HTTP GET:%s", owner, err)
			continue
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Test with %s: Documents summary does not return HTTP code %d but %d.", owner, http.StatusOK, rec.Code)
		} else {

			var res []structs.Document

			if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
				t.Errorf("Test with %s: Documents summary does not return valid Documents list", owner)
			}

			if len(res) != count {
				t.Errorf("Test with %s: Document creation returns  Document with %d Statements, wants %d", owner, len(res), count)
			}

		}

	}
}

func testDeleteDocHandler(t *testing.T) {
	for key, value := range documents {

		req, err := http.NewRequest(echo.DELETE, "/documents/:docid", nil)
		if err != nil {
			t.Errorf("Test with %s: Error deleting Document: %s", key, err)
			continue
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		c.SetParamNames("docid")
		c.SetParamValues(documentIDs[key])

		err = deleteDocHandler(c)
		if err != nil {
			t.Errorf("Test with %s: Error deleting Document during post:%s", key, err)
			continue
		}

		if value.Pass {

			if rec.Code != http.StatusOK {
				t.Errorf("Test with %s: Document deletion does not return HTTP code %d but %d.", key, http.StatusOK, rec.Code)
				continue
			}

		} else { //  test missing Document

			if rec.Code != http.StatusNotFound {
				t.Errorf("Test with %s: Document deletion does not return HTTP code %d but %d.", key, http.StatusNotFound, rec.Code)
				continue
			}
		}

	}
}
