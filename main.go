package main

import (
	"flag"
	"fmt"

	"os"

	"time"

	"strings"

	"regexp"

	"github.com/kyokomi/emoji"
	"github.com/townewgokgok/slack-status/internal"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: slack-status [options..] <template ID>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Templates:")
	for id, tmpl := range internal.Settings.Templates {
		emoji.Fprintln(os.Stderr, "  - "+id+" : "+wrapEmoji(tmpl.Emoji)+" "+tmpl.Text)
	}
	os.Exit(1)
}

type Flags struct {
	dryRun bool
	edit   bool
	iTunes bool
	lastFM bool
	view   bool
	watch  bool
}

func wrapEmoji(e string) string {
	if e == "" {
		return e
	}
	return ":" + e + ":"
}

var s = internal.Settings

func main() {
	// Parse arguments
	var f Flags
	flag.BoolVar(&f.dryRun, "d", false, "Dry run")
	flag.BoolVar(&f.edit, "e", false, "Edit settings")
	flag.BoolVar(&f.iTunes, "i", false, "Append information of the music playing on iTunes")
	flag.BoolVar(&f.lastFM, "l", false, "Append information of the music playing on last.fm")
	flag.BoolVar(&f.view, "v", false, "View current status")
	flag.BoolVar(&f.watch, "w", false, "Watch changes (with -i or -l)")
	flag.Parse()
	id := flag.Arg(0)

	if f.edit {
		internal.Edit()
	}
	if s.Slack.Token == "" || strings.ContainsRune(s.Slack.Token, '.') {
		fmt.Fprintln(os.Stderr, `settings.yml seems to be not customized. Try "slack-status -e" to edit it.`)
		fmt.Fprintln(os.Stderr, "")
		usage()
	}

	if f.view {
		emoji.Println(internal.GetSlackUserStatus())
		os.Exit(0)
	}

	withInfo := 0
	interval := time.Duration(1)
	if f.iTunes {
		withInfo++
		interval = s.ITunes.WatchIntervalSec
	}
	if f.lastFM {
		withInfo++
		interval = s.LastFM.WatchIntervalSec
	}
	if interval < 1 {
		interval = 1
	}
	if 1 < withInfo {
		fmt.Fprintln(os.Stderr, `Both -i and -l cannot be specified at the same time`)
		fmt.Fprintln(os.Stderr, "")
		usage()
	}
	if withInfo == 0 {
		f.watch = false
	}
	if id == "" && withInfo == 0 {
		emoji.Fprintln(os.Stderr, "Current status: " + internal.GetSlackUserStatus())
		fmt.Fprintln(os.Stderr, "")
		usage()
	}

	var t0, e0 string
	if id != "" {
		tmpl, ok := s.Templates[id]
		if !ok {
			fmt.Fprintln(os.Stderr, `Template "`+id+`" is not defined in settings.yml`)
			fmt.Fprintln(os.Stderr, "")
			usage()
		}
		t0 = tmpl.Text
		e0 = wrapEmoji(tmpl.Emoji)
	}

	update(&f, e0, t0)

	for f.watch {
		time.Sleep(interval * time.Second)
		update(&f, e0, t0)
	}
}

func appendInfo(emoji, text, emojiToAppend, textToAppend string) (string, string) {
	if emojiToAppend != "" {
		if emoji == "" {
			emoji = wrapEmoji(emojiToAppend)
		} else {
			if text != "" {
				text += " "
			}
			text += wrapEmoji(emojiToAppend)
		}
	}
	if text != "" {
		text += " "
	}
	text += textToAppend
	return emoji, text
}

func appendMusicInfo(emoji, text string, settings *internal.MusicSettings, status *internal.MusicStatus) (string, string) {
	if !status.Valid {
		return emoji, text
	}
	r := regexp.MustCompile(`%\w`)
	info := r.ReplaceAllStringFunc(settings.Format, func(m string) string {
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
	return appendInfo(emoji, text, settings.Emoji, info)
}

func limitStringByLength(str string, maxlen int) string {
	r := []rune(str)
	if len(r) <= maxlen {
		return str
	}
	return string(r[:maxlen-1]) + "…"
}

var lastText string
var lastEmoji string
var updatedCount int

func update(f *Flags, e, t string) {
	if f.iTunes {
		e, t = appendMusicInfo(e, t, &s.ITunes.MusicSettings, &internal.GetITunesStatus().MusicStatus)
	}

	if f.lastFM {
		e, t = appendMusicInfo(e, t, &s.LastFM.MusicSettings, &internal.GetLastFMStatus().MusicStatus)
	}

	t = limitStringByLength(t, internal.SlackUserStatusMaxLength)

	changed := updatedCount == 0 || !(t == lastText && e == lastEmoji)

	if changed {
		if !f.dryRun {
			internal.SetSlackUserStatus(t, e)
		}
		if e != "" {
			emoji.Print(e + " ")
		}
		emoji.Println(t)
	}

	lastText = t
	lastEmoji = e
	updatedCount++
}
