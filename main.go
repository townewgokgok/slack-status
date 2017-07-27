package main

import (
	"flag"
	"fmt"

	"os"

	"time"

	"strings"

	"regexp"

	"sort"
	"strconv"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"github.com/townewgokgok/slack-status/internal"
)

var cyan = color.New(color.FgCyan)
var red = color.New(color.FgRed)

func warn(msg string) {
	red.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, "")
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: slack-status [options..] [<template ID>]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Templates:")
	maxlen := 0
	ids := []string{}
	for id := range internal.Settings.Templates {
		if maxlen < len(id) {
			maxlen = len(id)
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		tmpl := internal.Settings.Templates[id]
		str := fmt.Sprintf("  %-"+strconv.Itoa(maxlen)+"s = %s %s", id, wrapEmoji(tmpl.Emoji), tmpl.Text)
		emoji.Fprintln(os.Stderr, str)
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
		warn(`settings.yml seems to be not customized. Try "slack-status -e" to edit it.`)
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
		warn(`Both -i and -l cannot be specified at the same time`)
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
			warn(`Template "` + id + `" is not defined in settings.yml`)
			usage()
		}
		t0 = tmpl.Text
		e0 = wrapEmoji(tmpl.Emoji)
	}

	update(&f, e0, t0, f.watch)

	for f.watch {
		time.Sleep(interval * time.Second)
		update(&f, e0, t0, f.watch)
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

func update(f *Flags, e, t string, printTime bool) {
	now := time.Now().Format("[15:04:05] ")
	notice := ""

	if f.iTunes {
		status := &internal.GetITunesStatus().MusicStatus
		if status.Ok {
			e, t = appendMusicInfo(e, t, &s.ITunes.MusicSettings, status)
		} else {
			notice = "iTunes seems to be stopped"
			if status.Err != "" {
				notice += " : " + status.Err
			}
		}
	}

	if f.lastFM {
		status := &internal.GetLastFMStatus().MusicStatus
		if status.Ok {
			e, t = appendMusicInfo(e, t, &s.LastFM.MusicSettings, status)
		} else {
			notice = "Failed to fetch music information from last.fm"
			if status.Err != "" {
				notice += " : " + status.Err
			}
		}
	}

	t = limitStringByLength(t, internal.SlackUserStatusMaxLength)

	changed := updatedCount == 0 || !(t == lastText && e == lastEmoji)

	if changed {
		if !f.dryRun {
			internal.SetSlackUserStatus(t, e)
		}
		if notice != "" {
			if printTime {
				red.Fprint(os.Stderr, now)
			}
			red.Fprintln(os.Stderr, notice)
		}
		if printTime {
			cyan.Print(now)
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
