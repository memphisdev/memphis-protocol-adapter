package adapter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/g41797/sputnik"
	"github.com/tkanos/gonfig"
)

// ConfigFactory returns implementation of sputnik.ConfFactory based on github.com/tkanos/gonfig
// - JSON format of config files
// - Automatic matching of environment variables
// - Env. variable for configuration "example" and key "kname" should be set in environment as "EXAMPLE_KNAME"
// - Value in environment automatically overrides value from the file
// - Temporary implementation github.com/tkanos/gonfig  will be replaced with github.com/g41797/gonfig
func ConfigFactory(cfPath string) sputnik.ConfFactory {
	cnf := newConfig(cfPath)
	return cnf.unmarshal
}

type config struct {
	confPath string
}

func newConfig(cfPath string) *config {
	return &config{confPath: cfPath}
}

func (conf *config) unmarshal(confName string, result any) error {
	fPath := filepath.Join(conf.confPath, strings.ToLower(confName))
	fPath += ".json"
	_, err := os.Open(fPath)
	if err != nil {
		return err
	}
	err = gonfig.GetConf(fPath, result)
	return err
}
