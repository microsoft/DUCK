package rulebases

import (
	"testing"

	"github.com/labstack/echo"
)

func TestHandler_CheckDoc(t *testing.T) {
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
	t.Error("rulebases.Handler.CheckDoc() tests not implemented")
	for _, tt := range tests {
		if err := tt.h.CheckDoc(tt.args.c); (err != nil) != tt.wantErr {
			t.Errorf("%q. Handler.CheckDoc() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestHandler_CheckDocID(t *testing.T) {
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
	t.Error("rulebases.Handler.CheckDocID() tests not implemented")
	for _, tt := range tests {
		if err := tt.h.CheckDocID(tt.args.c); (err != nil) != tt.wantErr {
			t.Errorf("%q. Handler.CheckDocID() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestHandler_GetRulebases(t *testing.T) {
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
	t.Error("rulebases.Handler.GetRulebases() tests not implemented")
	for _, tt := range tests {
		if err := tt.h.GetRulebases(tt.args.c); (err != nil) != tt.wantErr {
			t.Errorf("%q. Handler.GetRulebases() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
