package structs

import (
	"reflect"
	"testing"
)

func TestNewHTTPError(t *testing.T) {
	type args struct {
		err  string
		code int
	}
	tests := []struct {
		name string
		args args
		want HTTPError
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := NewHTTPError(tt.args.err, tt.args.code); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewHTTPError() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestWrapErrWith(t *testing.T) {
	type args struct {
		err  error
		herr HTTPError
	}
	tests := []struct {
		name string
		args args
		want HTTPError
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := WrapErrWith(tt.args.err, tt.args.herr); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. WrapErrWith() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestHTTPError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    HTTPError
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := tt.e.Error(); got != tt.want {
			t.Errorf("%q. HTTPError.Error() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
