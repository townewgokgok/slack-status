package internal

import (
	"io/ioutil"

	"path/filepath"
	"runtime"

	"time"

	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type StatusTemplate struct {
	Text  string `yaml:"text,omitempty"`
	Emoji string `yaml:"emoji,omitempty"`
}

type ITunesSettings struct {
	WatchIntervalSec time.Duration `yaml:"watch_interval_sec"`
}

type LastFMSettings struct {
	WatchIntervalSec time.Duration `yaml:"watch_interval_sec"`
	UserName         string        `yaml:"user_name"`
	APIKey           string        `yaml:"api_key"`
	Secret           string        `yaml:"secret"`
}

type SettingsRoot struct {
	Token        string                    `yaml:"token"`
	Templates    map[string]StatusTemplate `yaml:"templates"`
	PlayingEmoji string                    `yaml:"playing_emoji"`
	PlayingText  string                    `yaml:"playing_text"`
	ITunes       ITunesSettings            `yaml:"itunes,omitempty"`
	LastFM       LastFMSettings            `yaml:"lastfm,omitempty"`
}

var Settings SettingsRoot
var projectDir, settingsPath, settingsExamplePath string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to detect settings file location")
	}
	projectDir = filepath.Dir(filepath.Dir(filename))
	settingsPath = filepath.Join(projectDir, "settings.yml")
	settingsExamplePath = filepath.Join(projectDir, "settings.sample.yml")

	_, err := os.Stat(settingsPath)
	if err != nil {
		b, err := ioutil.ReadFile(settingsExamplePath)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(settingsPath, b, 0600)
		if err != nil {
			panic(err)
		}
	}

	data, err := ioutil.ReadFile(settingsPath)
	if err != nil {
		panic("Failed to load settings: " + err.Error())
	}
	err = yaml.Unmarshal(data, &Settings)
	if err != nil {
		panic("Failed to unmarshall settings: " + err.Error())
	}
}

func Edit() {
	cmd := exec.Command("vi", settingsPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
