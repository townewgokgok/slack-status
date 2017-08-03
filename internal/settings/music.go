package settings

import (
	"regexp"
	"time"

	"github.com/townewgokgok/slack-status/internal/music"
)

type MusicSettings struct {
	TemplateID       string        `yaml:"template_id"`
	WatchIntervalSec time.Duration `yaml:"watch_interval_sec"`
	Format           string        `yaml:"format"`
}

func (s *MusicSettings) ReplacePlaceholder(status *music.MusicStatus) string {
	r := regexp.MustCompile(`%\w`)
	return r.ReplaceAllStringFunc(s.Format, func(m string) string {
		switch m[1] {
		case 'A':
			return status.Artist
		case 'a':
			return status.Album
		case 't':
			return status.Title
		case '%':
			return "%"
		}
		return m
	})
}

type ITunesSettings struct {
	MusicSettings `yaml:",inline"`
}

type LastFMSettings struct {
	MusicSettings `yaml:",inline"`
	UserName      string `yaml:"user_name"`
	APIKey        string `yaml:"api_key"`
}
