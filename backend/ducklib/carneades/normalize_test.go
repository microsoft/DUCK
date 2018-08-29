// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package carneades

import (
	"reflect"
	"testing"

	"github.com/Microsoft/DUCK/backend/ducklib/db"
	"github.com/Microsoft/DUCK/backend/ducklib/structs"
)

func TestNewNormalizer(t *testing.T) {
	type args struct {
		doc    structs.Document
		db     *db.Database
		webdir string
	}
	tests := []struct {
		name    string
		args    args
		want    *Normalizer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	t.Errorf("Implement Normalize tests")
	for _, tt := range tests {
		got, err := NewNormalizer(tt.args.doc, tt.args.db, tt.args.webdir)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. NewNormalizer() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewNormalizer() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_normalizer_CreateDict(t *testing.T) {
	tests := []struct {
		name    string
		n       *Normalizer
		want    *NormalizedDocument
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		err := tt.n.createDict()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. normalizer.CreateDict() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(tt.n.normalized.IsA, tt.want) {
			t.Errorf("%q. normalizer.CreateDict() = %v, want %v", tt.name, tt.n.normalized.IsA, tt.want)
		}
	}
}

func Test_normalizer_getCode(t *testing.T) {
	type args struct {
		Type string
		Code string
	}
	tests := []struct {
		name string
		n    *Normalizer
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := tt.n.getCode(tt.args.Type, tt.args.Code); got != tt.want {
			t.Errorf("%q. normalizer.getCode() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_normalizer_Denormalize(t *testing.T) {
	tests := []struct {
		name string
		n    *Normalizer
		want *structs.Document
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := tt.n.Denormalize(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. normalizer.Denormalize() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
