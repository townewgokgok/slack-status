package internal

import (
	"io/ioutil"

	"path/filepath"
	"runtime"

	"time"

	"os"
	"os/exec"

	"bytes"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type StatusTemplate struct {
	Emoji string `yaml:"emoji,omitempty"`
	Text  string `yaml:"text,omitempty"`
}

type MusicSettings struct {
	WatchIntervalSec time.Duration `yaml:"watch_interval_sec"`
	Emoji            string        `yaml:"emoji"`
	Format           string        `yaml:"format"`
}

var Settings struct {
	Slack struct {
		Token string `yaml:"token"`
	} `yaml:"slack"`
	ITunes struct {
		MusicSettings `yaml:",inline"`
	} `yaml:"itunes,omitempty"`
	LastFM struct {
		MusicSettings `yaml:",inline"`
		UserName      string `yaml:"user_name"`
		APIKey        string `yaml:"api_key"`
		Secret        string `yaml:"secret"`
	} `yaml:"lastfm,omitempty"`
	Templates map[string]StatusTemplate `yaml:"templates"`
}

var projectDir, settingsExamplePath string
var SettingsPath string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to detect settings file location")
	}
	projectDir = filepath.Dir(filepath.Dir(filename))
	homeDir, err := homedir.Dir()
	if err != nil {
		panic("Failed to detect your home directory")
	}
	settingsExamplePath = filepath.Join(projectDir, "settings.sample.yml")
	SettingsPath = filepath.Join(homeDir, ".slack-status.settings.yml")

	_, err = os.Stat(SettingsPath)
	if err != nil {
		b, err := ioutil.ReadFile(settingsExamplePath)
		if err != nil {
			panic(err)
		}
		if runtime.GOOS == "windows" {
			b = bytes.Replace(b, []byte("\x0A"), []byte("\x0D\x0A"), -1)
		}
		err = ioutil.WriteFile(SettingsPath, b, 0600)
		if err != nil {
			panic(err)
		}
	}

	data, err := ioutil.ReadFile(SettingsPath)
	if err != nil {
		panic("Failed to load settings: " + err.Error())
	}
	err = yaml.Unmarshal(data, &Settings)
	if err != nil {
		panic("Failed to unmarshall settings: " + err.Error())
	}
}

func Edit() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "start", "notepad", SettingsPath)
	} else {
		cmd = exec.Command("vi", SettingsPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
