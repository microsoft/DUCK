package ducklib

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
)

var cfg structs.Configuration

func init() {
	flag.StringVar(&cfg.WebDir, "webdir", "", "The root directory for serving web content")
	flag.StringVar(&cfg.JwtKey, "jwtkey", "", "The secret used to sign the JWT")
	flag.StringVar(&cfg.RulebaseDir, "rulebasedir", "", "The Directory to the Rulebases")
	flag.BoolVar(&cfg.Gopathrelative, "gopathrelative", true, "Defines if webdir and rulebasedir are relative to the GOPATH")
	flag.BoolVar(&cfg.Loadtestdata, "loadtestdata", false, "If this is true, testdata will be loaded into the database")
	flag.Parse()
}

//NewConfiguration is the Constructor for a new structs.Configuration struct.
//it uses information from a cofiguration file, command flags, environment Variables and its own defaults to
//decide initial values for the configuration
func NewConfiguration(confpath string) structs.Configuration {

	c := structs.Configuration{}

	//setting defaults
	c.JwtKey = "secret"
	c.WebDir = "/src/github.com/Microsoft/DUCK/frontend/dist"
	c.RulebaseDir = "/src/github.com/Microsoft/DUCK/RuleBases"
	c.Gopathrelative = true
	c.Loadtestdata = false

	//overwrite defaults with information from config file
	if err := getFileConfig(&c, confpath); err != nil {
		log.Printf("Could not load configuration file: %s", err)

	}
	//overwrite with information from environment
	getEnv(&c)
	//overwrite with information from flags
	getFlags(&c)
	return c
}

//getFlags populates the configuration struct with data from the flags this program is called with
func getFlags(config *structs.Configuration) {

	if cfg.JwtKey != "" {
		config.JwtKey = cfg.JwtKey
	}
	if cfg.WebDir != "" {
		config.WebDir = cfg.WebDir
	}
	if cfg.RulebaseDir != "" {
		config.RulebaseDir = cfg.RulebaseDir
	}

	if cfg.Gopathrelative != true {
		config.Gopathrelative = cfg.Gopathrelative
	}

	if cfg.Loadtestdata != false {
		config.Loadtestdata = cfg.Loadtestdata
	}

}

//getFileConfig reads the config from a JSON formatted config file and
// sets the read values as configuration values
func getFileConfig(config *structs.Configuration, confpath string) error {
	dat, err := ioutil.ReadFile(confpath)
	if err != nil {
		return err

	}
	err = json.Unmarshal(dat, &config)
	return err

}

func getEnv(c *structs.Configuration) {
	//Get Environment Variables

	env := os.Getenv("DUCK_JWTKEY")
	if env != "" {
		c.JwtKey = env
	}

	env = os.Getenv("DUCK_WEBDIR")
	if env != "" {
		c.WebDir = env
	}

	env = os.Getenv("DUCK_RULEBASEDIR")
	if env != "" {
		c.RulebaseDir = env
	}
	//has to be not empty and also something like a boolean to be set
	env = os.Getenv("DUCK_GOPATHRELATIVE")
	if env != "" {
		if gpr, err := strconv.ParseBool(env); err == nil {
			c.Gopathrelative = gpr
		} else {
			log.Printf("Could not read value for GOPATHRELATIVE: %s", err)
		}
	}
	env = os.Getenv("DUCK_LOADTESTDATA")
	if env != "" {
		if ldt, err := strconv.ParseBool(env); err == nil {
			c.Loadtestdata = ldt
		} else {
			log.Printf("Could not read value for LOADTESTDATA: %s", err)
		}
	}
	env = os.Getenv("DUCK_DATABASE.LOCATION")
	if env != "" {
		c.DBConfig.Location = env
	}
	env = os.Getenv("DUCK_DATABASE.PORT")
	if env != "" {
		if p, err := strconv.Atoi(env); err == nil {
			c.DBConfig.Port = p
		} else {
			log.Printf("Could not read value for PORT: %s", err)
		}
	}
	env = os.Getenv("DUCK_DATABASE.USERNAME")
	if env != "" {
		c.DBConfig.Username = env
	}
	env = os.Getenv("DUCK_DATABASE.PASSWORD")
	if env != "" {
		c.DBConfig.Password = env
	}
	env = os.Getenv("DUCK_DATABASE.NAME")
	if env != "" {
		c.DBConfig.Name = env
	}
}
