// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package dictionaries

import (
	"encoding/json"
	"io/ioutil"
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
	conf config.Configuration
	e    *echo.Echo

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
	dih Handler
)

func TestDictionaryHandler(t *testing.T) {

	conf = config.NewConfiguration(filepath.Join(os.Getenv("GOPATH"), "/src/github.com/Microsoft/DUCK/backend/configuration.json"))
	e = echo.New()

	dih = Handler{}
	dab, err := db.NewDatabase(*conf.DBConfig)
	if err != nil {
		t.Skip("User Handler test failed; was not able to datab.Init()")
	}
	dih.Db = dab
	dat, err := ioutil.ReadFile("testdata/dictionary.json")

	if err = json.Unmarshal(dat, &dicts); err != nil {
		t.Error("Testfixture Dictionary not correctly loading")
		t.Skip("No testfixtures no dictionary tests")
	}
	id, err := dih.Db.PostUser(dicts.User)
	if err != nil {
		t.Skip("Not able to save user to mockdb in dictionary test. Skipping tests..")
	}
	defer func() {
		err := dih.Db.DeleteUser(id)
		if err != nil {
			t.Log("Could not delete user from mockDB in dictionary test this can interfere with other tests")
		}
	}()
	user, err := dih.Db.GetUser(id)
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

func TestHandler_GetUserDict(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		h       *Handler
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.h.GetUserDict(tt.args.c); (err != nil) != tt.wantErr {
			t.Errorf("%q. Handler.GetUserDict() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestHandler_GetDictItem(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		h       *Handler
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.h.GetDictItem(tt.args.c); (err != nil) != tt.wantErr {
			t.Errorf("%q. Handler.GetDictItem() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestHandler_DeleteDictItem(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		h       *Handler
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.h.DeleteDictItem(tt.args.c); (err != nil) != tt.wantErr {
			t.Errorf("%q. Handler.DeleteDictItem() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestHandler_PutDictItem(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		h       *Handler
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.h.PutDictItem(tt.args.c); (err != nil) != tt.wantErr {
			t.Errorf("%q. Handler.PutDictItem() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestHandler_PutUserDict(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		h       *Handler
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.h.PutUserDict(tt.args.c); (err != nil) != tt.wantErr {
			t.Errorf("%q. Handler.PutUserDict() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
