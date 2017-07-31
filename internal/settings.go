package internal

import (
	"io/ioutil"

	"path/filepath"
	"runtime"

	"time"

	"os"
	"os/exec"

	"strings"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type MusicSettings struct {
	TemplateID       string        `yaml:"template_id"`
	WatchIntervalSec time.Duration `yaml:"watch_interval_sec"`
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
	Templates map[string]string `yaml:"templates"`
}

var SettingsPath string

func init() {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic("Failed to detect your home directory")
	}
	SettingsPath = filepath.Join(homeDir, ".slack-status.yml")

	_, err = os.Stat(SettingsPath)
	if err != nil {
		sample := SettingsSample
		if runtime.GOOS == "windows" {
			sample = strings.Replace(sample, "\x0A", "\x0D\x0A", -1)
		}
		err = ioutil.WriteFile(SettingsPath, []byte(sample), 0600)
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

	if Settings.ITunes.TemplateID == "" {
		Settings.ITunes.TemplateID = "itunes"
	}
	if Settings.LastFM.TemplateID == "" {
		Settings.LastFM.TemplateID = "lastfm"
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
