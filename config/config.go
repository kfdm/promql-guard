package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/prometheus/prometheus/pkg/labels"
)

// Load parses the YAML input s into a Config.
func Load(s string) (*Config, error) {
	cfg := &Config{}

	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	cfg.original = s
	return cfg, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, errors.Wrapf(err, "parsing YAML file %s", filename)
	}
	return cfg, nil
}

// Prometheus guard configuration
type Prometheus struct {
	Upstream string            `yaml:"upstream"`
	Labels   map[string]string `yaml:"labels"`
}

// VirtualHost is a basic configuration unit
type VirtualHost struct {
	Hostname   string     `yaml:"hostname"`
	Prometheus Prometheus `yaml:"prometheus,omitempty"`
}

// Config represents the base configuration file
type Config struct {
	VirtualHosts []VirtualHost `yaml:"hosts"`
	original     string
}

// Find particular VirtualHost configuration
func (c *Config) Find(name string) (*VirtualHost, error) {
	for _, element := range c.VirtualHosts {
		if element.Hostname == name {
			return &element, nil
		}
	}
	return nil, errors.New("Unable to find virtual host")
}

// Matchers from Prometheus Config
func (vh VirtualHost) Matchers() ([]*labels.Matcher, error) {
	res := make([]*labels.Matcher, 0, len(vh.Prometheus.Labels))
	for name, value := range vh.Prometheus.Labels {
		res = append(res, &labels.Matcher{
			Name:  name,
			Value: value,
			Type:  labels.MatchEqual,
		})
	}
	return res, nil
}
