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

func NewConfiguration(confpath string) structs.Configuration {

	c := structs.Configuration{}

	//setting defaults
	c.JwtKey = "secret"
	c.WebDir = "/src/github.com/Microsoft/DUCK/frontend/dist"
	c.RulebaseDir = "/src/github.com/Microsoft/DUCK/RuleBases"
	c.Gopathrelative = true
	c.Loadtestdata = false

	if err := getFileConfig(&c, confpath); err != nil {
		log.Printf("Could not load configuration file: %s", err)

	}

	getEnv(&c)

	getFlags(&c)

	log.Printf("Config: %+v", c)
	log.Printf("Datab: %+v", c.DBConfig)
	return c
}

func getFlags(config *structs.Configuration) {

	flag.StringVar(&config.WebDir, "webdir", config.WebDir, "The root directory for serving web content")
	flag.StringVar(&config.JwtKey, "jwtkey", config.JwtKey, "The secret used to sign the JWT")
	flag.StringVar(&config.RulebaseDir, "rulebasedir", config.RulebaseDir, "The Directory to the Rulebases")
	flag.BoolVar(&config.Gopathrelative, "gopathrelative", config.Gopathrelative, "Defines if webdir and rulebasedir are relative to the GOPATH")
	flag.BoolVar(&config.Loadtestdata, "loadtestdata", config.Loadtestdata, "If this is true, testdata will be loaded into the database")

	flag.Parse()

}
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
