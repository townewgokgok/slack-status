package internal

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type StatusTemplate struct {
	Text  string `yaml:"text,omitempty"`
	Emoji string `yaml:"emoji,omitempty"`
}

type SettingsRoot struct {
	Token     string                    `yaml:"token"`
	Templates map[string]StatusTemplate `yaml:"templates"`
}

var Settings SettingsRoot

func init() {
	data, err := ioutil.ReadFile("settings.yml")
	if err != nil {
		panic("Failed to load settings: " + err.Error())
	}
	err = yaml.Unmarshal(data, &Settings)
	if err != nil {
		panic("Failed to unmarshall settings: " + err.Error())
	}
}
