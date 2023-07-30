package adapter

import (
	"github.com/g41797/sputnik"
	"github.com/spf13/viper"
)

// ConfigFactory returns implementation of sputnik.ConfFactory based on github.com/spf13/viper
// - JSON format of config files
// - Automatic matching of environment variables
// - Env. variable for configuration "example" should have prefix "EXAMPLE_"
func ConfigFactory(cfPath string) sputnik.ConfFactory {
	cnf := newConfig(cfPath)
	return cnf.unmarshal
}

type config struct {
	v *viper.Viper
}

func newConfig(cfPath string) *config {
	v := viper.New()
	v.AddConfigPath(cfPath)
	v.AutomaticEnv()
	v.SetConfigType("json")
	return &config{v: v}
}

func (conf *config) unmarshal(confName string, result any) error {
	conf.v.SetConfigName(confName)
	conf.v.SetEnvPrefix(confName)
	err := conf.v.ReadInConfig()
	if err == nil {
		err = conf.v.Unmarshal(result)
	}
	return err
}
