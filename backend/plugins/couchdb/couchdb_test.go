package couchdb

import (
	"io"
	"reflect"
	"testing"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
)

func Test_getMap(t *testing.T) {
	type args struct {
		resp io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := getMap(tt.args.resp)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. getMap() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getMap() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_getRows(t *testing.T) {
	type args struct {
		resp io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := getRows(tt.args.resp)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. getRows() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getRows() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCouchbase_GetLogin(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantId  string
		wantPw  string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotId, gotPw, err := tt.cb.GetLogin(tt.args.email)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.GetLogin() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if gotId != tt.wantId {
			t.Errorf("%q. Couchbase.GetLogin() gotId = %v, want %v", tt.name, gotId, tt.wantId)
		}
		if gotPw != tt.wantPw {
			t.Errorf("%q. Couchbase.GetLogin() gotPw = %v, want %v", tt.name, gotPw, tt.wantPw)
		}
	}
}

func TestCouchbase_GetUserDict(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		want    structs.Dictionary
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.cb.GetUserDict(tt.args.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.GetUserDict() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Couchbase.GetUserDict() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCouchbase_GetUser(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		want    structs.User
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.cb.GetUser(tt.args.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.GetUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Couchbase.GetUser() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCouchbase_GetDocument(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		want    structs.Document
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.cb.GetDocument(tt.args.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.GetDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Couchbase.GetDocument() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCouchbase_getCouchbaseDocument(t *testing.T) {
	type args struct {
		cbDocID string
	}
	tests := []struct {
		name         string
		cb           *Couchbase
		args         args
		wantDocument map[string]interface{}
		wantErr      bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotDocument, err := tt.cb.getCouchbaseDocument(tt.args.cbDocID)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.getCouchbaseDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotDocument, tt.wantDocument) {
			t.Errorf("%q. Couchbase.getCouchbaseDocument() = %v, want %v", tt.name, gotDocument, tt.wantDocument)
		}
	}
}

func TestCouchbase_GetDocumentSummariesForUser(t *testing.T) {
	type args struct {
		userid string
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		want    []structs.Document
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.cb.GetDocumentSummariesForUser(tt.args.userid)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.GetDocumentSummariesForUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Couchbase.GetDocumentSummariesForUser() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCouchbase_DeleteDocument(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.DeleteDocument(tt.args.id); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.DeleteDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_DeleteUser(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.DeleteUser(tt.args.id); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.DeleteUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_deleteCbDocument(t *testing.T) {
	type args struct {
		id  string
		rev string
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.deleteCbDocument(tt.args.id, tt.args.rev); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.deleteCbDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_NewUser(t *testing.T) {
	type args struct {
		user structs.User
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.NewUser(tt.args.user); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.NewUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_NewDocument(t *testing.T) {
	type args struct {
		doc structs.Document
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.NewDocument(tt.args.doc); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.NewDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_UpdateUserDict(t *testing.T) {
	type args struct {
		dict   structs.Dictionary
		userID string
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.UpdateUserDict(tt.args.dict, tt.args.userID); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.UpdateUserDict() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_UpdateUser(t *testing.T) {
	type args struct {
		user structs.User
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.UpdateUser(tt.args.user); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.UpdateUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_UpdateDocument(t *testing.T) {
	type args struct {
		doc structs.Document
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.UpdateDocument(tt.args.doc); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.UpdateDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_putUser(t *testing.T) {
	type args struct {
		u structs.User
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.putUser(tt.args.u); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.putUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_putDocument(t *testing.T) {
	type args struct {
		d structs.Document
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.putDocument(tt.args.d); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.putDocument() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_putEntry(t *testing.T) {
	type args struct {
		entry      map[string]interface{}
		designfile bool
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.putEntry(tt.args.entry, tt.args.designfile); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.putEntry() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_Init(t *testing.T) {
	type args struct {
		config structs.DBConf
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.Init(tt.args.config); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.Init() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_testFileExists(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		cb      *Couchbase
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.cb.testFileExists(tt.args.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.testFileExists() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. Couchbase.testFileExists() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCouchbase_createDatabase(t *testing.T) {
	tests := []struct {
		name    string
		cb      *Couchbase
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.cb.createDatabase(); (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.createDatabase() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestCouchbase_testDBExists(t *testing.T) {
	tests := []struct {
		name    string
		cb      *Couchbase
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.cb.testDBExists()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Couchbase.testDBExists() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. Couchbase.testDBExists() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_init(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for range tests {
		init()
	}
}
