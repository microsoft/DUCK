package ducklib

import (
	"reflect"
	"testing"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
)

func TestNewDatabase(t *testing.T) {

	tests := []struct {
		name   string
		config structs.DBConf
		want   *database
	}{
		{"Just setting Config to a DBConf ", structs.DBConf{Name: "Testname"}, &database{Config: structs.DBConf{Name: "Testname"}}},
	}
	for _, tt := range tests {
		if got := NewDatabase(tt.config); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewDatabase() = %v, want %v", tt.name, got, tt.want)
		}
	}
}



func Test_database_Init(t *testing.T) {

	tests := []struct {
		name     string
		database *database
		wantErr  bool
	}{
		{"Empty config", NewDatabase(structs.DBConf{}), false},
		{"INVALIDDBNAME", NewDatabase(structs.DBConf{Name: "INVALIDDBNAME"}), true},
	}
	for _, tt := range tests {
		if err := tt.database.Init(); (err != nil) != tt.wantErr {
			t.Errorf("%q. database.Init() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_database_GetLogin(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name     string
		database *database
		args     args
		wantID   string
		wantPw   string
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotID, gotPw, err := tt.database.GetLogin(tt.args.email)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. database.GetLogin() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if gotID != tt.wantID {
			t.Errorf("%q. database.GetLogin() gotID = %v, want %v", tt.name, gotID, tt.wantID)
		}
		if gotPw != tt.wantPw {
			t.Errorf("%q. database.GetLogin() gotPw = %v, want %v", tt.name, gotPw, tt.wantPw)
		}
	}
}

func Test_database_GetUser(t *testing.T) {
	type args struct {
		userid string
	}
	tests := []struct {
		name     string
		database *database
		args     args
		want     structs.User
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.database.GetUser(tt.args.userid)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. database.GetUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. database.GetUser() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_database_DeleteUser(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		database *database
		args     args
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.database.DeleteUser(tt.args.id); (err != nil) != tt.wantErr {
			t.Errorf("%q. database.DeleteUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_database_PutUser(t *testing.T) {
	type args struct {
		user structs.User
	}
	tests := []struct {
		name     string
		database *database
		args     args
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.database.PutUser(tt.args.user); (err != nil) != tt.wantErr {
			t.Errorf("%q. database.PutUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_database_PostUser(t *testing.T) {
	type args struct {
		user structs.User
	}
	tests := []struct {
		name     string
		database *database
		args     args
		wantID   string
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotID, err := tt.database.PostUser(tt.args.user)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. database.PostUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if gotID != tt.wantID {
			t.Errorf("%q. database.PostUser() = %v, want %v", tt.name, gotID, tt.wantID)
		}
	}
}

func Test_database_GetUserDict(t *testing.T) {
	type args struct {
		userid string
	}
	tests := []struct {
		name     string
		database *database
		args     args
		want     structs.Dictionary
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.database.GetUserDict(tt.args.userid)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. database.GetUserDict() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. database.GetUserDict() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_database_PutUserDict(t *testing.T) {
	type args struct {
		dict   structs.Dictionary
		userID string
	}
	tests := []struct {
		name     string
		database *database
		args     args
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.database.PutUserDict(tt.args.dict, tt.args.userID); (err != nil) != tt.wantErr {
			t.Errorf("%q. database.PutUserDict() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_database_GetDocument(t *testing.T) {
	type args struct {
		documentid string
	}
	tests := []struct {
		name     string
		database *database
		args     args
		want     structs.Document
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.database.GetDocument(tt.args.documentid)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. database.GetDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. database.GetDocument() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_database_GetDocumentSummariesForUser(t *testing.T) {
	type args struct {
		userid string
	}
	tests := []struct {
		name     string
		database *database
		args     args
		want     []structs.Document
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.database.GetDocumentSummariesForUser(tt.args.userid)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. database.GetDocumentSummariesForUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. database.GetDocumentSummariesForUser() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_database_DeleteDocument(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		database *database
		args     args
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.database.DeleteDocument(tt.args.id); (err != nil) != tt.wantErr {
			t.Errorf("%q. database.DeleteDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_database_PutDocument(t *testing.T) {
	type args struct {
		doc structs.Document
	}
	tests := []struct {
		name     string
		database *database
		args     args
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.database.PutDocument(tt.args.doc); (err != nil) != tt.wantErr {
			t.Errorf("%q. database.PutDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_database_PostDocument(t *testing.T) {
	type args struct {
		doc structs.Document
	}
	tests := []struct {
		name     string
		database *database
		args     args
		wantID   string
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotID, err := tt.database.PostDocument(tt.args.doc)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. database.PostDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if gotID != tt.wantID {
			t.Errorf("%q. database.PostDocument() = %v, want %v", tt.name, gotID, tt.wantID)
		}
	}
}
