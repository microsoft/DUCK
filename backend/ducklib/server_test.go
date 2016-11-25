package ducklib

import (
	"reflect"
	"testing"

	"github.com/Microsoft/DUCK/backend/ducklib/config"
	"github.com/labstack/echo"
)

//three tests:
//valid server start
//config with dbconfig name = INVALIDDBNAME
//valid with loadtestdata = true

func TestGetServer(t *testing.T) {
	type args struct {
		conf config.Configuration
	}
	tests := []struct {
		name string
		args args
		want *echo.Echo
	}{
	// TODO: Add test cases.
	}
	t.Errorf("GetServer() tests not implemented")
	for _, tt := range tests {
		if got := GetServer(tt.args.conf); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetServer() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
