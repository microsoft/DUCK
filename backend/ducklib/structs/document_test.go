package structs

import "testing"

func TestDocument_FromValueMap(t *testing.T) {
	type args struct {
		mp map[string]interface{}
	}
	tests := []struct {
		name string
		d    *Document
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.d.FromValueMap(tt.args.mp)
	}
}

func TestStatement_FromInterfaceMap(t *testing.T) {
	type args struct {
		mp map[string]interface{}
	}
	tests := []struct {
		name string
		s    *Statement
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.s.FromInterfaceMap(tt.args.mp)
	}
}

func Test_getFieldValue(t *testing.T) {
	type args struct {
		mp    map[string]interface{}
		field string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := getFieldValue(tt.args.mp, tt.args.field); got != tt.want {
			t.Errorf("%q. getFieldValue() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_getFieldBooleanValue(t *testing.T) {
	type args struct {
		mp    map[string]interface{}
		field string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := getFieldBooleanValue(tt.args.mp, tt.args.field); got != tt.want {
			t.Errorf("%q. getFieldBooleanValue() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
