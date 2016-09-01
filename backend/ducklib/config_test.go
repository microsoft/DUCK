package ducklib

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	wrongPath, correctPath, testPath string
)

func TestConfig(t *testing.T) {

	wrongPath = filepath.Join(goPath, "/src/github.com/Microsoft/DUCK/backend/nofile")
	correctPath = filepath.Join(goPath, "/src/github.com/Microsoft/DUCK/backend/configuration.json")
	testPath = "structs/testdata/configuration.json"
	//get env vars if they are set
	dbLocation := os.Getenv("DUCK_DATABASE.LOCATION")
	dbPort := os.Getenv("DUCK_DATABASE.PORT")
	dbName := os.Getenv("DUCK_DATABASE.NAME")
	dbUsername := os.Getenv("DUCK_DATABASE.USERNAME")
	dbPassword := os.Getenv("DUCK_DATABASE.PASSWORD")

	jwtkey := os.Getenv("DUCK_JWTKEY")
	webdir := os.Getenv("DUCK_WEBDIR")
	rbdir := os.Getenv("DUCK_RULEBASEDIR")
	gpr := os.Getenv("DUCK_GOPATHRELATIVE")
	loadtd := os.Getenv("DUCK_LOADTESTDATA")

	//set them all to zero	//set env to prior values
	os.Setenv("DUCK_DATABASE.LOCATION", "")
	os.Setenv("DUCK_DATABASE.PORT", "")
	os.Setenv("DUCK_DATABASE.NAME", "")
	os.Setenv("DUCK_DATABASE.USERNAME", "")
	os.Setenv("DUCK_DATABASE.PASSWORD", "")

	os.Setenv("DUCK_JWTKEY", "")
	os.Setenv("DUCK_WEBDIR", "")
	os.Setenv("DUCK_RULEBASEDIR", "")
	os.Setenv("DUCK_GOPATHRELATIVE", "")
	os.Setenv("DUCK_LOADTESTDATA", "")

	//t.Error("AHHHHHH")
	t.Run("File=1", testNoFile)

	//set env to prior values
	os.Setenv("DUCK_DATABASE.LOCATION", dbLocation)
	os.Setenv("DUCK_DATABASE.PORT", dbPort)
	os.Setenv("DUCK_DATABASE.NAME", dbName)
	os.Setenv("DUCK_DATABASE.USERNAME", dbUsername)
	os.Setenv("DUCK_DATABASE.PASSWORD", dbPassword)

	os.Setenv("DUCK_JWTKEY", jwtkey)
	os.Setenv("DUCK_WEBDIR", webdir)
	os.Setenv("DUCK_RULEBASEDIR", rbdir)
	os.Setenv("DUCK_GOPATHRELATIVE", gpr)
	os.Setenv("DUCK_LOADTESTDATA", loadtd)

}

func testNoFile(t *testing.T) {
	c := NewConfiguration("")

	//should be default
	if c.DBConfig != nil {
		t.Errorf("Configuration with no File: Database Object should be nil, is %+v", c.DBConfig)
	}

}
func testWrongFile(t *testing.T) {
	c := NewConfiguration(wrongPath)

	//should be default
	if c.DBConfig != nil {
		t.Errorf("Configuration with no File: Database Object should be nil, is %+v", c.DBConfig)
	}

}

func testEnvGopath(t *testing.T) {
	c := NewConfiguration(correctPath)

	type teststruct struct {
		envar  string
		setval string
		hasval string
	}
	//TestTable: map[EnvVar_description] [envar, setval, hasval]
	testtable := map[string]teststruct{
		"DUCK_JWTKEY":      {envar: "DUCK_JWTKEY", setval: "abcde", hasval: "abcde"},
		"DUCK_WEBDIR":      {envar: "DUCK_WEBDIR", setval: "abcde", hasval: "abcde"},
		"DUCK_RULEBASEDIR": {envar: "DUCK_RULEBASEDIR", setval: "abcde", hasval: "abcde"},
		"location":         {envar: "DUCK_DATABASE.LOCATION", setval: "abcde", hasval: "abcde"},
		"port":             {envar: "DUCK_DATABASE.PORT", setval: "1234", hasval: "1234"},
		"port_wrong":       {envar: "DUCK_DATABASE.PORT", setval: "abcde", hasval: "5984"},
		"name":             {envar: "DUCK_DATABASE.NAME", setval: "abcde", hasval: "abcde"},
		"username":         {envar: "DUCK_DATABASE.USERNAME", setval: "abcde", hasval: "abcde"},
		"Password":         {envar: "DUCK_DATABASE.PASSWORD", setval: "abcde", hasval: "abcde"},
		"9":                {envar: "DUCK_JWTKEY", setval: "abcde", hasval: "abcde"},
		"11":               {envar: "DUCK_JWTKEY", setval: "abcde", hasval: "abcde"},
		"12":               {envar: "DUCK_JWTKEY", setval: "abcde", hasval: "abcde"},
		"14":               {envar: "DUCK_JWTKEY", setval: "abcde", hasval: "abcde"},
	}
	for key, val := range testtable {
		t.Errorf("%s%s", key, val)
	}
	//should be default
	if c.DBConfig != nil {
		t.Errorf("Configuration with no File: Database Object should be nil, is %+v", c.DBConfig)
	}

}
