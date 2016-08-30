package ducklib

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	wrongPath, correctPath string
)

func TestMain(m *testing.M) {

	wrongPath = filepath.Join(goPath, "/src/github.com/Microsoft/DUCK/backend/nofile")
	correctPath = filepath.Join(goPath, "/src/github.com/Microsoft/DUCK/backend/configuration.json")
	conf = NewConfiguration()

	rrun := m.run()

	os.Exit(rrun)
}
