package settings

import (
	"io/ioutil"

	"path/filepath"
	"runtime"

	"os"
	"os/exec"

	"strings"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

var Settings struct {
	Slack     SlackSettings    `yaml:"slack"`
	Templates TemplateSettings `yaml:"templates"`
	ITunes    ITunesSettings   `yaml:"itunes,omitempty"`
	LastFM    LastFMSettings   `yaml:"lastfm,omitempty"`
}

var SettingsPath string
var SettingsWarnings []string

func init() {
	SettingsWarnings = []string{}

	homeDir, err := homedir.Dir()
	if err != nil {
		SettingsWarnings = append(SettingsWarnings, "Failed to detect your home directory")
	}
	SettingsPath = filepath.Join(homeDir, ".slack-status.yml")

	_, err = os.Stat(SettingsPath)
	if err != nil {
		example := SettingsExample
		if runtime.GOOS == "windows" {
			example = strings.Replace(example, "\x0A", "\x0D\x0A", -1)
		}
		err = ioutil.WriteFile(SettingsPath, []byte(example), 0600)
		if err != nil {
			SettingsWarnings = append(SettingsWarnings, err.Error())
		}
	}

	data, err := ioutil.ReadFile(SettingsPath)
	if err != nil {
		SettingsWarnings = append(SettingsWarnings, "Failed to load settings: "+err.Error())
	}
	err = yaml.Unmarshal(data, &Settings)
	if err != nil {
		SettingsWarnings = append(SettingsWarnings, "Failed to unmarshall settings: "+err.Error())
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
