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
	"sort"
	"fmt"
	"strconv"
	"github.com/kyokomi/emoji"
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
	} `yaml:"lastfm,omitempty"`
	Templates map[string]string `yaml:"templates"`
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

func ListTemplates(indent string) string {
	maxlen := 0
	ids := []string{}
	for id := range Settings.Templates {
		if maxlen < len(id) {
			maxlen = len(id)
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)
	result := ""
	for _, id := range ids {
		tmpl := Settings.Templates[id]
		str := fmt.Sprintf("%s%-"+strconv.Itoa(maxlen)+"s = %s\n", indent, id, tmpl)
		result += emoji.Sprint(str)
	}
	return result
}
