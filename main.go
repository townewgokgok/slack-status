package main

import (
	"flag"
	"fmt"

	"os"

	"time"

	"strings"

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
	flag.BoolVar(&f.watch, "w", false, "Watch changes (with -i or -l)")
	flag.Parse()
	id := flag.Arg(0)

	if f.edit {
		internal.Edit()
	}
	if s.Token == "" || strings.ContainsRune(s.Token, '.') {
		fmt.Fprintln(os.Stderr, `settings.yml seems to be not customized. Try "slack-status -e" to edit it.`)
		fmt.Fprintln(os.Stderr, "")
		usage()
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

func appendInfo(emoji, text, emojiToAppend, headerToAppend, textToAppend string) (string, string) {
	if emojiToAppend != "" {
		if emoji == "" {
			emoji = wrapEmoji(emojiToAppend)
		} else {
			text += " " + wrapEmoji(emojiToAppend)
		}
	}
	if headerToAppend != "" {
		text += " " + headerToAppend
	}
	text += " " + textToAppend
	if text[0] == ' ' {
		text = text[1:]
	}
	return emoji, text
}

var lastText string
var lastEmoji string
var updatedCount int

func update(f *Flags, e, t string) {
	if f.iTunes {
		i := internal.GetITunesStatus()
		if i.Valid {
			e, t = appendInfo(e, t, s.PlayingEmoji, s.PlayingText, i.Artist+" - "+i.Name)
		}
	}

	if f.lastFM {
		l := internal.GetLastFMStatus()
		if l.Valid {
			e, t = appendInfo(e, t, s.PlayingEmoji, s.PlayingText, l.Artist+" - "+l.Name)
		}
	}

	changed := updatedCount == 0 || !(t == lastText && e == lastEmoji)

	if changed {
		if !f.dryRun {
			internal.SetSlackUserStatus(s.Token, t, e)
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
