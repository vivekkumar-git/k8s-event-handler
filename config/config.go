package config

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	WatchMode RunMode = "watch"
)

type RunMode string

type Config struct {
	Mode       RunMode
	Resources  []Resource
	Namespaces []string
	Notifier   Notifier
}

type Resource struct {
	Kind string
}

type Notifier struct {
	Slack Slack
}

// Slack contains slack configuration
type Slack struct {
	Enabled bool
	Token   string
	Channel string
	Title   string
}

func (c *Config) init() {
	if c.Mode == "" {
		c.Mode = WatchMode
	}

	if len(c.Namespaces) == 0 {
		c.Namespaces = append(c.Namespaces, "all")
	}
}

func (c *Config) validate() error {
	for _, ns := range c.Namespaces {
		if ns == "all" {
			if len(c.Namespaces) > 1 {
				return errors.New("cannot specify a namespace after selecting all")
			}
		}
	}
	return nil
}

// New returns new Config
func GetConfig(filepath string) (*Config, error) {
	c := &Config{}
	config, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer config.Close()

	b, err := ioutil.ReadAll(config)
	if err != nil {
		return nil, err
	}

	if len(b) != 0 {
		yaml.Unmarshal(b, c)
	}

	c.init()
	err = c.validate()
	if err != nil {
		return nil, err
	}

	return c, nil
}
