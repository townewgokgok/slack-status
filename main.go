package main

import (
	"flag"
	"fmt"

	"os"

	"time"

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
	iTunes bool
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
	flag.BoolVar(&f.iTunes, "i", false, "Append information of the music playing on iTunes")
	flag.BoolVar(&f.watch, "w", false, "Watch changes (with -i)")
	flag.Parse()
	id := flag.Arg(0)
	if !f.iTunes {
		f.watch = false
	}

	if id == "" && !f.iTunes {
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

	update(&f, t0, e0)

	interval := s.WatchIntervalSec
	if interval < 1 {
		interval = 1
	}
	for f.watch {
		time.Sleep(interval * time.Second)
		update(&f, t0, e0)
	}
}

var lastText string
var lastEmoji string
var updatedCount int

func update(f *Flags, t, e string) {
	if f.iTunes {
		i := internal.GetITunesStatus()
		if i.Valid {
			if s.PlayingEmoji != "" {
				if e == "" {
					e = wrapEmoji(s.PlayingEmoji)
				} else {
					t += " " + wrapEmoji(s.PlayingEmoji)
				}
			}
			if s.PlayingText != "" {
				t += " " + s.PlayingText
			}
			t += " " + i.Artist + " - " + i.Name
			if t[0] == ' ' {
				t = t[1:]
			}
		}
	}

	if updatedCount == 0 || !(t == lastText && e == lastEmoji) {
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
