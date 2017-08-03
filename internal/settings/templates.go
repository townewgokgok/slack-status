package settings

import (
	"fmt"
	"sort"
	"strconv"

	"strings"

	"github.com/kyokomi/emoji"
	"github.com/townewgokgok/slack-status/internal/generator"
	"github.com/townewgokgok/slack-status/internal/helper"
	"github.com/townewgokgok/slack-status/internal/music"
)

type TemplateSettings map[string]string

func (s TemplateSettings) Dump(indent string) string {
	maxlen := 0
	ids := []string{}
	for id := range s {
		if maxlen < len(id) {
			maxlen = len(id)
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)
	result := ""
	for _, id := range ids {
		tmpl := s[id]
		str := fmt.Sprintf("%s%-"+strconv.Itoa(maxlen)+"s = %s\n", indent, id, tmpl)
		result += emoji.Sprint(str)
	}
	return result
}

func generateByTemplate(tmplID string) (string, string) {
	tmpl, ok := Settings.Templates[tmplID]
	gen := generator.NewGenerator()

	if !ok {
		switch tmplID {
		case Settings.ITunes.TemplateID:
			status := &music.GetITunesStatus().MusicStatus
			if !status.Ok {
				n := "iTunes seems to be stopped"
				if status.Err != "" {
					n += " : " + status.Err
				}
				return "", n
			}
			gen.AddReplacer(status.Replacer)
			tmpl = Settings.ITunes.MusicSettings.Format
		case Settings.LastFM.TemplateID:
			status := &music.GetLastFMStatus(Settings.LastFM.APIKey, Settings.LastFM.UserName).MusicStatus
			if !status.Ok {
				n := "Failed to fetch music information from last.fm"
				if status.Err != "" {
					n += " : " + status.Err
				}
				return "", n
			}
			gen.AddReplacer(status.Replacer)
			tmpl = Settings.LastFM.MusicSettings.Format
		default:
			return "", ""
		}
	}

	return gen.Execute(tmpl), ""
}

func Generate(templateIDs []string) (string, string, []string, error) {
	notice := []string{}
	texts := []string{}

	for _, id := range templateIDs {
		t, n := generateByTemplate(id)
		if n != "" {
			notice = append(notice, n)
		}
		if t == "" {
			return "", "", notice, fmt.Errorf(
				`Template "%s" is not defined in settings file.`+"\n"+
					`Try "slack-status list" to list your templates.`,
				id,
			)
		}
		texts = append(texts, t)
	}

	emj, txt := helper.SplitEmoji(strings.Join(texts, " "))
	return emj, txt, notice, nil
}
