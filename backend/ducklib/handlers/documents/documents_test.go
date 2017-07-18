// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package documents

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

	"github.com/Microsoft/DUCK/backend/ducklib/config"
	"github.com/Microsoft/DUCK/backend/ducklib/db"
	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	_ "github.com/Microsoft/DUCK/backend/plugins/mockdb"
	"github.com/labstack/echo"
)

var (
	doh   Handler
	datab *db.Database
	conf  config.Configuration
	e     *echo.Echo

	users map[string]struct {
		Pass bool         `json:"pass"`
		User structs.User `json:"user"`
	}
	//eg document["document_a"]
	documents map[string]struct {
		Pass     bool             `json:"pass"`
		Document structs.Document `json:"document"`
	}

	//eg userIDs["user_a"]="a structs.User.ID"
	userIDs     = make(map[string]string)
	documentIDs = make(map[string]string)
	//documents owners and # of documents they own
	owners = make(map[string]int)
)

func TestDocumentHandler(t *testing.T) {
	conf = config.NewConfiguration(filepath.Join(os.Getenv("GOPATH"), "/src/github.com/Microsoft/DUCK/backend/configuration.json"))
	e = echo.New()

	doh = Handler{}
	dab, err := db.NewDatabase(*conf.DBConfig)
	if err != nil {
		t.Skip("User Handler test failed; was not able to datab.Init()")
	}
	doh.Db = dab

	dat, err := ioutil.ReadFile("testdata/document.json")

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
	t.Run("PostDocument=1", testpostDoc)
	t.Run("GetDocument=1", testGetDoc)
	t.Run("Summaries=1", testGetDocSummaries)
	t.Run("CopyDocument=1", testCopyDoc)
	t.Run("PutDocument=1", testPutDoc)

	t.Run("DeleteDocument=1", testDeleteDoc)
	//t.Error("AHHHHHH")

}

func testpostDoc(t *testing.T) {

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

		c := e.NewContext(req, rec)
		err = doh.PostDoc(c)
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

func testPutDoc(t *testing.T) {

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

		c := e.NewContext(req, rec)
		err = doh.PutDoc(c)
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

func testGetDoc(t *testing.T) {
	for key, value := range documents {

		value.Document.ID = documentIDs[key]

		req, err := http.NewRequest(echo.GET, "/documents/:docid", nil)

		if err != nil {
			t.Errorf("Test with %s: Error getting Document: %s", key, err)
			continue
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.SetParamNames("docid")
		c.SetParamValues(documentIDs[key])
		err = doh.GetDoc(c)
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

func testCopyDoc(t *testing.T) {

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

		c := e.NewContext(req, rec)

		c.SetParamNames("docid")
		c.SetParamValues(documentIDs[key])
		err = doh.CopyStatements(c)
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

		c := e.NewContext(req, rec)
		c.SetParamNames("userid")
		c.SetParamValues(owner)
		err = doh.GetDocSummaries(c)
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

func testDeleteDoc(t *testing.T) {
	for key, value := range documents {

		req, err := http.NewRequest(echo.DELETE, "/documents/:docid", nil)
		if err != nil {
			t.Errorf("Test with %s: Error deleting Document: %s", key, err)
			continue
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("docid")
		c.SetParamValues(documentIDs[key])

		err = doh.DeleteDoc(c)
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
