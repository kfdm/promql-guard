package config

import (
	"io/ioutil"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/yaml.v2"
)

// Config Basic config struct
type Config []struct {
	Hostname string `yaml:"hostname"`
}

// New Load new configuration file
func New(path string, logger log.Logger) Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		level.Error(logger).Log("msg", "Error loading file")
		panic(err)
	}

	config := Config{}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		level.Error(logger).Log("msg", "Error loading file", "err", err)
		panic(err)
	}
	return config
}
