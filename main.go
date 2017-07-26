package main

import (
	"flag"
	"fmt"

	"os"

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
	iTunes bool
}

func wrapEmoji(e string) string {
	if e == "" {
		return e
	}
	return ":" + e + ":"
}

func main() {
	s := internal.Settings

	// Parse arguments
	var f Flags
	flag.BoolVar(&f.iTunes, "i", false, "Append information of the music playing on iTunes")
	flag.Parse()
	id := flag.Arg(0)
	if id == "" && !f.iTunes {
		usage()
	}

	var t, e string
	if id != "" {
		tmpl, ok := s.Templates[id]
		if !ok {
			fmt.Fprintln(os.Stderr, `Template "`+id+`" is not defined in settings.yml`)
			fmt.Fprintln(os.Stderr, "")
			usage()
		}
		t = tmpl.Text
		e = wrapEmoji(tmpl.Emoji)
	}

	if f.iTunes {
		st := internal.GetITunesStatus()
		if st.Valid {
			msg := "PLAYING " + st.Artist + " - " + st.Name
			if e == "" {
				e = ":musical_note:"
			} else {
				t += " :musical_note: "
			}
			t += msg
		}
	}

	internal.SetSlackUserStatus(s.Token, t, e)
	emoji.Println(e + " " + t)
}
