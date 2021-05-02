package config

import (
	"io/ioutil"
	"net/http"
	"net/url"

	auth "github.com/abbot/go-http-auth"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
	"gopkg.in/yaml.v2"
)

type matchers []*labels.Matcher
type upstream struct {
	*url.URL
}

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
	Upstream upstream          `yaml:"upstream"`
	Matchers matchers          `yaml:"matcher,omitempty"`
	Headers  map[string]string `yaml:"headers,omitempty"`
}

// VirtualHost is a basic configuration unit
type VirtualHost struct {
	Username   string     `yaml:"username"`
	Prometheus Prometheus `yaml:"prometheus,omitempty"`
}

func (p *Prometheus) URL() *url.URL {
	return p.Upstream.URL
}

// AddHeaders adds custom headers to request
func (p *Prometheus) UpdateRequest(req *http.Request) {
	req.Host = p.Upstream.Host
	for i, v := range p.Headers {
		req.Header.Set(i, v)
	}
}

// Config represents the base configuration file
type Config struct {
	VirtualHosts []VirtualHost `yaml:"hosts"`
	Htpasswd     string        `yaml:"htpasswd"`
	original     string
}

func (c *Config) NewBasicAuthenticator() *auth.BasicAuth {
	htpasswd := auth.HtpasswdFileProvider(c.Htpasswd)
	return auth.NewBasicAuthenticator("Basic Realm", htpasswd)
}

// Find particular VirtualHost configuration
func (c *Config) Find(name string) (*VirtualHost, error) {
	for _, element := range c.VirtualHosts {
		if element.Username == name {
			return &element, nil
		}
	}
	return nil, errors.New("Unable to find virtual host")
}

// UnmarshalYAML a regular string type to Prometheus matcher type
func (m *matchers) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buf string
	err := unmarshal(&buf)
	if err != nil {
		return errors.New("Unable to unmarshal string")
	}
	expr, err := promql.ParseExpr(buf)
	if err != nil {
		return errors.New("Unable to parse PromQL")
	}

	switch n := expr.(type) {
	case *promql.VectorSelector:
		*m = n.LabelMatchers
		return nil
	default:
		return errors.New("Invalid matcher declaration")
	}
}

func (u *upstream) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	parsed, err := url.Parse(s)
	if err != nil {
		return err
	}
	u.URL = parsed
	return nil
}
