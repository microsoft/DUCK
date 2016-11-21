package carneades

import (
	"io"
	"reflect"
	"testing"
)

func TestMakeComplianceCheckerPlugin(t *testing.T) {
	type args struct {
		ruleBaseDir string
	}
	tests := []struct {
		name    string
		args    args
		want    *ComplianceCheckerPlugin
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := MakeComplianceCheckerPlugin(tt.args.ruleBaseDir)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. MakeComplianceCheckerPlugin() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. MakeComplianceCheckerPlugin() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestComplianceCheckerPlugin_Intialize(t *testing.T) {
	tests := []struct {
		name    string
		c       *ComplianceCheckerPlugin
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.c.Intialize(); (err != nil) != tt.wantErr {
			t.Errorf("%q. ComplianceCheckerPlugin.Intialize() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestComplianceCheckerPlugin_Shutdown(t *testing.T) {
	tests := []struct {
		name string
		c    ComplianceCheckerPlugin
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.c.Shutdown()
	}
}

func TestComplianceCheckerPlugin_ruleBaseReader(t *testing.T) {
	type args struct {
		ruleBaseID string
	}
	tests := []struct {
		name string
		c    *ComplianceCheckerPlugin
		args args
		want io.Reader
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := tt.c.ruleBaseReader(tt.args.ruleBaseID); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ComplianceCheckerPlugin.ruleBaseReader() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestComplianceCheckerPlugin_IsCompliant(t *testing.T) {
	type args struct {
		ruleBaseID string
		document   *NormalizedDocument
	}
	tests := []struct {
		name    string
		c       *ComplianceCheckerPlugin
		args    args
		want    bool
		want1   Explanation
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, got1, err := tt.c.IsCompliant(tt.args.ruleBaseID, tt.args.document)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. ComplianceCheckerPlugin.IsCompliant() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. ComplianceCheckerPlugin.IsCompliant() got = %v, want %v", tt.name, got, tt.want)
		}
		if !reflect.DeepEqual(got1, tt.want1) {
			t.Errorf("%q. ComplianceCheckerPlugin.IsCompliant() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}

func TestComplianceCheckerPlugin_CompliantDocuments(t *testing.T) {
	type args struct {
		ruleBaseID string
		document   *NormalizedDocument
		maxResults int
		offset     int
	}
	tests := []struct {
		name    string
		c       *ComplianceCheckerPlugin
		args    args
		want    bool
		want1   []*NormalizedDocument
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, got1, err := tt.c.CompliantDocuments(tt.args.ruleBaseID, tt.args.document, tt.args.maxResults, tt.args.offset)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. ComplianceCheckerPlugin.CompliantDocuments() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. ComplianceCheckerPlugin.CompliantDocuments() got = %v, want %v", tt.name, got, tt.want)
		}
		if !reflect.DeepEqual(got1, tt.want1) {
			t.Errorf("%q. ComplianceCheckerPlugin.CompliantDocuments() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
