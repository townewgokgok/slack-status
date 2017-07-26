package internal

import (
	"io/ioutil"

	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

type StatusTemplate struct {
	Text  string `yaml:"text,omitempty"`
	Emoji string `yaml:"emoji,omitempty"`
}

type SettingsRoot struct {
	Token     string                    `yaml:"token"`
	Templates map[string]StatusTemplate `yaml:"templates"`
	PlayingEmoji string `yaml:"playing_emoji"`
	PlayingText string `yaml:"playing_text"`
}

var Settings SettingsRoot

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to detect settings file location")
	}
	path := filepath.Join(filepath.Dir(filepath.Dir(filename)), "settings.yml")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic("Failed to load settings: " + err.Error())
	}
	err = yaml.Unmarshal(data, &Settings)
	if err != nil {
		panic("Failed to unmarshall settings: " + err.Error())
	}
}
