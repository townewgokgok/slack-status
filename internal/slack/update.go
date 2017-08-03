package slack

import (
	"os"
	"strings"
	"time"

	"fmt"

	"github.com/townewgokgok/slack-status/internal/helper"
	"github.com/townewgokgok/slack-status/internal/music"
	s "github.com/townewgokgok/slack-status/internal/settings"
)

var lastText string
var lastEmoji string
var updatedCount int

func Update(watch, dryRun bool, templateIDs []string) error {
	now := time.Now().Format("[15:04:05] ")
	notice := []string{}

	tmpls := []string{}
	for _, id := range templateIDs {
		tmpl, ok := s.Settings.Templates[id]
		if ok {
			tmpls = append(tmpls, tmpl)
			continue
		}
		switch id {
		case s.Settings.ITunes.TemplateID:
			status := &music.GetITunesStatus().MusicStatus
			if status.Ok {
				tmpls = append(tmpls, s.Settings.ITunes.MusicSettings.ReplacePlaceholder(status))
				continue
			}
			n := "iTunes seems to be stopped"
			if status.Err != "" {
				n += " : " + status.Err
			}
			notice = append(notice, n)
		case s.Settings.LastFM.TemplateID:
			status := &music.GetLastFMStatus(s.Settings.LastFM.APIKey, s.Settings.LastFM.UserName).MusicStatus
			if status.Ok {
				tmpls = append(tmpls, s.Settings.LastFM.MusicSettings.ReplacePlaceholder(status))
				continue
			}
			n := "Failed to fetch music information from last.fm"
			if status.Err != "" {
				n += " : " + status.Err
			}
			notice = append(notice, n)
		default:
			return fmt.Errorf(
				`Template "%s" is not defined in settings file.`+"\n"+
					`Try "slack-status list" to list your templates.`,
				id,
			)
		}
	}

	t := strings.Join(tmpls, " ")
	e, t := helper.SplitEmoji(t)
	t = helper.LimitStringByLength(t, SlackUserStatusMaxLength)
	changed := updatedCount == 0 || !(t == lastText && e == lastEmoji)

	if changed {
		var err error
		if !dryRun {
			err = SetSlackUserStatus(t, e)
		}
		if 0 < len(notice) {
			if watch {
				helper.Red.Fprint(os.Stderr, now)
			}
			helper.Red.Fprintln(os.Stderr, strings.Join(notice, "\n"))
		}
		if watch {
			helper.Cyan.Print(now)
		}
		helper.PrintStatus(e, t)
		if err != nil {
			return err
		}
	}

	lastText = t
	lastEmoji = e
	updatedCount++

	return nil
}
