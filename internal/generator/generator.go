package generator

import (
	"fmt"
	"strings"

	"github.com/townewgokgok/slack-status/internal/helper"
	"github.com/townewgokgok/slack-status/internal/music"
	s "github.com/townewgokgok/slack-status/internal/settings"
)

func generateByTemplate(tmplID string) (string, string, bool) {
	tmpl, ok := s.Settings.Templates[tmplID]
	rep := newReplacerChain()

	if !ok {
		switch tmplID {
		case s.Settings.ITunes.TemplateID:
			status := &music.GetITunesStatus().MusicStatus
			if !status.Ok {
				n := "iTunes seems to be stopped"
				if status.Err != "" {
					n += " : " + status.Err
				}
				return "", n, true
			}
			rep.AddReplacer(status.Replacer)
			tmpl = s.Settings.ITunes.MusicSettings.Format
		case s.Settings.LastFM.TemplateID:
			status := &music.GetLastFMStatus(s.Settings.LastFM.APIKey, s.Settings.LastFM.UserName).MusicStatus
			if !status.Ok {
				n := "Failed to fetch music information from last.fm"
				if status.Err != "" {
					n += " : " + status.Err
				}
				return "", n, true
			}
			rep.AddReplacer(status.Replacer)
			tmpl = s.Settings.LastFM.MusicSettings.Format
		default:
			return "", "", false
		}
	}

	return rep.execute(tmpl), "", true
}

func Generate(templateIDs []string) (string, string, []string, error) {
	notice := []string{}
	texts := []string{}

	for _, id := range templateIDs {
		t, n, ok := generateByTemplate(id)
		if n != "" {
			notice = append(notice, n)
		}
		if !ok {
			return "", "", notice, fmt.Errorf(
				`Template "%s" is not defined in settings file.`+"\n"+
					`Try "slack-status list" to list your templates.`,
				id,
			)
		}
		if t != "" {
			texts = append(texts, t)
		}
	}

	emj, txt := helper.SplitEmoji(strings.Join(texts, " "))
	return emj, txt, notice, nil
}
