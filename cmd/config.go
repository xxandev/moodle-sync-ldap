package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"moodle-sync-ldap/internal/ldap"
	"moodle-sync-ldap/internal/moodle"
	"moodle-sync-ldap/internal/utils"
	"os"

	"gopkg.in/yaml.v3"
)

type FormatConfig int8

const (
	JSON FormatConfig = iota
	XML
	YAML
)

type Config struct {
	Moodle   moodle.Config `json:"moodle,omitempty" xml:"moodle,omitempty" yaml:"moodle,omitempty"`
	LDAP     ldap.Config   `json:"ldif,omitempty" xml:"ldif,omitempty" yaml:"ldif,omitempty"`
	fileLDIF string
}

func (c *Config) Init() (bool, error) {
	switch {
	case utils.IsStatFile("send-gmail.json"):
		content, err := os.ReadFile("send-gmail.json")
		if err != nil {
			return true, err
		}
		if err := c.Unmarshal(JSON, content); err != nil {
			return true, err
		}
	case utils.IsStatFile("send-gmail.xml"):
		content, err := os.ReadFile("send-gmail.xml")
		if err != nil {
			return true, err
		}
		if err := c.Unmarshal(XML, content); err != nil {
			return true, err
		}
	case utils.IsStatFile("send-gmail.yaml"):
		content, err := os.ReadFile("send-gmail.yaml")
		if err != nil {
			return true, err
		}
		if err := c.Unmarshal(YAML, content); err != nil {
			return true, err
		}
	}
	return false, errors.New("config file not found")
}

func (c *Config) Marshal(format FormatConfig) ([]byte, error) {
	switch format {
	case JSON:
		return json.MarshalIndent(c, "", "\t")
	case XML:
		return xml.MarshalIndent(c, "", "\t")
	case YAML:
		return yaml.Marshal(c)
	}
	return nil, errors.New("unknown configuration format")
}

func (c *Config) Unmarshal(format FormatConfig, data []byte) error {
	switch format {
	case JSON:
		return json.Unmarshal(data, c)
	case XML:
		return xml.Unmarshal(data, c)
	case YAML:
		return yaml.Unmarshal(data, c)
	}
	return errors.New("unknown configuration format")
}

func (c *Config) Check() error {
	if len(c.Moodle.URL) < 8 {
		return errors.New("invalid moodle url")
	}
	if len(c.Moodle.Token) < 8 {
		return errors.New("invalid moodle token")
	}
	if len(c.LDAP.DN) < 4 {
		return errors.New("invalid ldap dn")
	}
	return nil
}
