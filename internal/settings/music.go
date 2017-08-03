package settings

import (
	"time"
)

type MusicSettings struct {
	TemplateID       string        `yaml:"template_id"`
	WatchIntervalSec time.Duration `yaml:"watch_interval_sec"`
	Format           string        `yaml:"format"`
}

type ITunesSettings struct {
	MusicSettings `yaml:",inline"`
}

type LastFMSettings struct {
	MusicSettings `yaml:",inline"`
	UserName      string `yaml:"user_name"`
	APIKey        string `yaml:"api_key"`
}
