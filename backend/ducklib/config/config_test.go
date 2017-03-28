package config

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var (
	wrongPath, correctPath, testPath string
)

func TestConfig(t *testing.T) {

	wrongPath = filepath.Join(os.Getenv("GOPATH"), "/src/github.com/Microsoft/DUCK/backend/nofile")
	correctPath = filepath.Join(os.Getenv("GOPATH"), "/src/github.com/Microsoft/DUCK/backend/configuration.json")
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

	//set them all to zero	//set env to prior values
	os.Setenv("DUCK_DATABASE.LOCATION", "")
	os.Setenv("DUCK_DATABASE.PORT", "")
	os.Setenv("DUCK_DATABASE.NAME", "")
	os.Setenv("DUCK_DATABASE.USERNAME", "")
	os.Setenv("DUCK_DATABASE.PASSWORD", "")

	os.Setenv("DUCK_JWTKEY", "")
	os.Setenv("DUCK_WEBDIR", "")
	os.Setenv("DUCK_RULEBASEDIR", "")

	//t.Error("AHHHHHH")
	t.Run("File=1", testNoFile)
	t.Run("File=2", testWrongFile)
	t.Run("File=3", testCorrectFile)
	t.Run("Env=1", testEnvGopath)

	//set env to prior values
	os.Setenv("DUCK_DATABASE.LOCATION", dbLocation)
	os.Setenv("DUCK_DATABASE.PORT", dbPort)
	os.Setenv("DUCK_DATABASE.NAME", dbName)
	os.Setenv("DUCK_DATABASE.USERNAME", dbUsername)
	os.Setenv("DUCK_DATABASE.PASSWORD", dbPassword)

	os.Setenv("DUCK_JWTKEY", jwtkey)
	os.Setenv("DUCK_WEBDIR", webdir)
	os.Setenv("DUCK_RULEBASEDIR", rbdir)

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
func testCorrectFile(t *testing.T) {

	c := NewConfiguration(correctPath)
	goPath := os.Getenv("GOPATH")

	dat, err := ioutil.ReadFile(correctPath)
	if err != nil {
		t.Error("Skipped")
		t.Skip("Could not load testing config file")

	}
	var i map[string]interface{}
	err = json.Unmarshal(dat, &i)
	if err != nil {
		t.Error("Skipped")
		t.Skip("Could not load testing config file")

	}
	str, ok := i["jwtkey"].(string)
	if !ok {
		t.Error("Skipped")
		t.Skip("Could not load jwtkey from testing config file")

	}
	j, err := base64.StdEncoding.DecodeString(str)
	if !ok {
		t.Error("Skipped")
		t.Skip("Could not decode jwtkey in byte")

	}

	if len(c.JwtKey) != len(j) {
		t.Errorf("Configuration with correct File: JWT Key does not equal correct key got: %s, want: %s", c.JwtKey, j)
	}

	for l := range j {
		if j[l] != c.JwtKey[l] {
			t.Errorf("Configuration with correct File: JWT Key does not equal correct key got: %s, want: %s", c.JwtKey, j)
		}

	}
	rb := i["rulebasedir"].(string)
	wb := i["webdir"].(string)

	if goPath != "" {
		if !filepath.IsAbs(rb) {
			rb = filepath.Join(goPath, rb)
		}

		if !filepath.IsAbs(wb) {
			wb = filepath.Join(goPath, wb)
		}
	}

	if c.RulebaseDir != rb {
		t.Errorf("Configuration with correct File: RulebaseDir does not equal correct Dir. got: %s, want: %s", c.RulebaseDir, rb)
	}
	if c.WebDir != wb {
		t.Errorf("Configuration with correct File: WebDir does not equal correct Dir. got: %s, want: %s", c.WebDir, wb)
	}

	if c.DBConfig == nil {
		t.Errorf("Configuration with correct File: DBConfig is nil, should not be nil")
	}

}

func testEnvGopath(t *testing.T) {

	type teststruct struct {
		envar   string
		setval  string
		wantval string
	}
	abs, err := filepath.Abs("abcde")
	if err != nil {
		t.Errorf("Configuration with env Var, error getting absolute filepath: %s", err)
	}
	goabs := filepath.Join(os.Getenv("GOPATH"), "abcde")
	//TestTable: map[EnvVar_description] [envar, setval, hasval]
	testtable := map[string]teststruct{
		"DUCK_JWTKEY":          {envar: "DUCK_JWTKEY", setval: "abcde", wantval: "abcde"},
		"DUCK_WEBDIR":          {envar: "DUCK_WEBDIR", setval: "abcde", wantval: goabs},
		"DUCK_WEBDIR_ABS":      {envar: "DUCK_WEBDIR", setval: abs, wantval: abs},
		"DUCK_RULEBASEDIR":     {envar: "DUCK_RULEBASEDIR", setval: "abcde", wantval: goabs},
		"DUCK_RULEBASEDIR_ABS": {envar: "DUCK_RULEBASEDIR", setval: abs, wantval: abs},
		"location":             {envar: "DUCK_DATABASE.LOCATION", setval: "abcde", wantval: "abcde"},
		"port":                 {envar: "DUCK_DATABASE.PORT", setval: "1234", wantval: "1234"},
		"port_wrong":           {envar: "DUCK_DATABASE.PORT", setval: "abcde", wantval: "5984"},
		"name":                 {envar: "DUCK_DATABASE.NAME", setval: "abcde", wantval: "abcde"},
		"username":             {envar: "DUCK_DATABASE.USERNAME", setval: "abcde", wantval: "abcde"},
		"Password":             {envar: "DUCK_DATABASE.PASSWORD", setval: "abcde", wantval: "abcde"},
	}
	for key, val := range testtable {

		os.Setenv(val.envar, val.setval)

		c := NewConfiguration(correctPath)

		switch val.envar {
		case "DUCK_JWTKEY":
			if string(c.JwtKey) != val.wantval {
				t.Errorf("Testing environment Variable setting. Key: %s.  Wanted %s, got %s", key, val.wantval, c.JwtKey)
			}
		case "DUCK_WEBDIR":
			if c.WebDir != val.wantval {
				t.Errorf("Testing environment Variable setting. Key: %s. Wanted %s, got %s", key, val.wantval, c.WebDir)
			}
		case "DUCK_RULEBASEDIR":
			if c.RulebaseDir != val.wantval {
				t.Errorf("Testing environment Variable setting. Key: %s. Wanted %s, got %s", key, val.wantval, c.RulebaseDir)
			}
		case "DUCK_DATABASE.LOCATION":
			if c.DBConfig.Location != val.wantval {
				t.Errorf("Testing environment Variable setting. Key: %s. Wanted %s, got %s", key, val.wantval, c.DBConfig.Location)
			}
		case "DUCK_DATABASE.PORT":
			if strconv.Itoa(c.DBConfig.Port) != val.wantval {
				t.Errorf("Testing environment Variable setting. Key: %s. Wanted %s, got %d", key, val.wantval, c.DBConfig.Port)
			}
		case "DUCK_DATABASE.NAME":
			if c.DBConfig.Name != val.wantval {
				t.Errorf("Testing environment Variable setting. Key: %s. Wanted %s, got %s", key, val.wantval, c.DBConfig.Name)
			}
		case "DUCK_DATABASE.USERNAME":
			if c.DBConfig.Username != val.wantval {
				t.Errorf("Testing environment Variable setting. Key: %s. Wanted %s, got %s", key, val.wantval, c.DBConfig.Username)
			}
		case "DUCK_DATABASE.PASSWORD":
			if c.DBConfig.Password != val.wantval {
				t.Errorf("Testing environment Variable setting. Key: %s. Wanted %s, got %s", key, val.wantval, c.DBConfig.Password)
			}
		}

		os.Setenv(val.envar, "")

	}

}
