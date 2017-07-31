package main

import (
	"fmt"

	"github.com/urfave/cli"

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

var version = "1.0.1"

var cyan = color.New(color.FgCyan)
var yellow = color.New(color.FgYellow)
var red = color.New(color.FgRed)
var bggray = color.New(color.BgHiBlack)

func warn(msgs ...string) {
	for _, msg := range msgs {
		yellow.Fprintln(os.Stderr, `[warning] `+msg)
	}
	fmt.Fprintln(os.Stderr, "")
}

func wrapEmoji(e string) string {
	if e == "" {
		return e
	}
	return ":" + e + ":"
}

var settings = internal.Settings

type Flags struct {
	dryRun bool
	iTunes bool
	lastFM bool
	watch  bool
}

var flags Flags

func cliError(msgs ...string) *cli.ExitError {
	return cli.NewExitError(red.Sprint(strings.Join(msgs, "\n")), 1)
}

func main() {
	// Parse arguments
	app := cli.NewApp()
	app.Name = "slack-status"
	app.Version = version
	app.Usage = "Updates your Slack user status from CLI"
	//app.Authors = []cli.Author{
	//	{
	//		Name:  "townewgokgok",
	//		Email: "townewgokgok@gmail.com",
	//	},
	//}
	app.HelpName = "slack-status"

	app.Commands = []cli.Command{
		{
			Name:      "edit",
			Aliases:   []string{"e"},
			Usage:     "Opens your settings file in the editor",
			ArgsUsage: " ",
			Action: func(ctx *cli.Context) error {
				internal.Edit()
				return nil
			},
		},
		{
			Name:      "get",
			Aliases:   []string{"g"},
			Usage:     "Shows your current status",
			ArgsUsage: " ",
			Action: func(ctx *cli.Context) error {
				e, t := internal.GetSlackUserStatus()
				printStatus(e, t)
				return nil
			},
		},
		{
			Name:      "list",
			Aliases:   []string{"l"},
			Usage:     "Lists your templates",
			ArgsUsage: " ",
			Action: func(ctx *cli.Context) error {
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
					str := fmt.Sprintf("%-"+strconv.Itoa(maxlen)+"s = %s", id, tmpl)
					emoji.Fprintln(os.Stderr, str)
				}
				return nil
			},
		},
		{
			Name:      "set",
			Aliases:   []string{"s"},
			Usage:     "Updates your status",
			ArgsUsage: "[<template ID>]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dryrun, d",
					Usage: "just print the composed status text (your status will be not changed)",
				},
				cli.BoolFlag{
					Name:  "itunes, i",
					Usage: "append information about the music playing on iTunes",
				},
				cli.BoolFlag{
					Name:  "lastfm, l",
					Usage: "append information about the music playing on last.fm",
				},
				cli.BoolFlag{
					Name:  "watch, w",
					Usage: "watch changes (with -i or -l)",
				},
			},
			Action: func(ctx *cli.Context) error {
				flags.dryRun = ctx.Bool("dryrun")
				flags.iTunes = ctx.Bool("itunes")
				flags.lastFM = ctx.Bool("lastfm")
				flags.watch = ctx.Bool("watch")
				id := ctx.Args().Get(0)

				withInfo := 0
				interval := time.Duration(1)
				if flags.iTunes {
					withInfo++
					interval = settings.ITunes.WatchIntervalSec
					if interval < 1 {
						warn(`itunes.watch_interval_sec must be >= 1`)
						interval = 5
					}
				}
				if flags.lastFM {
					withInfo++
					interval = settings.LastFM.WatchIntervalSec
					if interval < 15 {
						warn(`lastfm.watch_interval_sec must be >= 15`)
						interval = 15
					}
				}
				if 1 < withInfo {
					return cliError(`Both -i and -l cannot be specified at the same time.`)
				}
				if withInfo == 0 {
					flags.watch = false
				}
				if id == "" && withInfo == 0 {
					cli.ShowCommandHelp(ctx, "set")
					os.Exit(1)
				}

				var t0 string
				if id != "" {
					tmpl, ok := settings.Templates[id]
					if !ok {
						return cliError(
							`Template "`+id+`" is not defined in settings file.`,
							`Try "slack-status list" to list your templates.`,
						)
					}
					t0 = tmpl
				}

				update(&flags, t0)

				for flags.watch {
					time.Sleep(interval * time.Second)
					update(&flags, t0)
				}

				return nil
			},
		},
	}

	app.Action = func(ctx *cli.Context) error {
		if settings.Slack.Token == "" || strings.ContainsRune(settings.Slack.Token, '.') {
			return cliError(
				`Your settings file seems to be not configured correctly.`,
				`The example settings file has been created at `+internal.SettingsPath,
				`Try "slack-status edit" to edit it.`,
			)
		}
		cli.ShowAppHelp(ctx)
		return nil
	}

	app.Run(os.Args)
}

func appendInfo(text, textToAppend string) string {
	if text != "" {
		text += " "
	}
	text += textToAppend
	return text
}

func appendMusicInfo(text string, settings *internal.MusicSettings, status *internal.MusicStatus) string {
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
	return appendInfo(text, info)
}

func limitStringByLength(str string, maxlen int) string {
	r := []rune(str)
	if len(r) <= maxlen {
		return str
	}
	return string(r[:maxlen-1]) + "…"
}

var splitEmojiRegexp = regexp.MustCompile(`^:[^: ]+: *`)

func splitEmoji(text string) (string, string) {
	text = strings.Trim(text, " ")
	e := ""
	m := splitEmojiRegexp.FindString(text)
	if m != "" {
		e = strings.Trim(m, " ")
		text = text[len(m):]
	}
	return e, text
}

func printStatus(e, t string) {
	if e != "" {
		bggray.Print(emoji.Sprint(e))
		fmt.Print(" ")
	}
	emoji.Println(t)
}

var lastText string
var lastEmoji string
var updatedCount int

func update(flags *Flags, t string) {
	now := time.Now().Format("[15:04:05] ")
	notice := ""

	if flags.iTunes {
		status := &internal.GetITunesStatus().MusicStatus
		if status.Ok {
			t = appendMusicInfo(t, &settings.ITunes.MusicSettings, status)
		} else {
			notice = "iTunes seems to be stopped"
			if status.Err != "" {
				notice += " : " + status.Err
			}
		}
	}

	if flags.lastFM {
		status := &internal.GetLastFMStatus().MusicStatus
		if status.Ok {
			t = appendMusicInfo(t, &settings.LastFM.MusicSettings, status)
		} else {
			notice = "Failed to fetch music information from last.fm"
			if status.Err != "" {
				notice += " : " + status.Err
			}
		}
	}

	e, t := splitEmoji(t)

	t = limitStringByLength(t, internal.SlackUserStatusMaxLength)

	changed := updatedCount == 0 || !(t == lastText && e == lastEmoji)

	if changed {
		if !flags.dryRun {
			internal.SetSlackUserStatus(t, e)
		}
		if notice != "" {
			if flags.watch {
				red.Fprint(os.Stderr, now)
			}
			red.Fprintln(os.Stderr, notice)
		}
		if flags.watch {
			cyan.Print(now)
		}
		printStatus(e, t)
	}

	lastText = t
	lastEmoji = e
	updatedCount++
}
