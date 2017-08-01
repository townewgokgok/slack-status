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

var version = "1.1.0"

var cyan = color.New(color.FgCyan)
var yellow = color.New(color.FgYellow)
var red = color.New(color.FgRed)
var bggray = color.New(color.BgHiBlack)

func warn(msgs ...string) {
	for _, msg := range msgs {
		lines := strings.Split(msg, "\n")
		for _, line := range lines {
			yellow.Fprintln(os.Stderr, `[warning] `+line)
		}
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
	watch  bool
}

var flags Flags

func cliError(msgs ...string) *cli.ExitError {
	return cli.NewExitError(red.Sprint(strings.Join(msgs, "\n")), 1)
}

func listTemplates(indent string) string {
	maxlen := 0
	ids := []string{}
	for id := range internal.Settings.Templates {
		if maxlen < len(id) {
			maxlen = len(id)
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)
	result := ""
	for _, id := range ids {
		tmpl := internal.Settings.Templates[id]
		str := fmt.Sprintf("%s%-"+strconv.Itoa(maxlen)+"s = %s\n", indent, id, tmpl)
		result += emoji.Sprint(str)
	}
	return result
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
				e, t, err := internal.GetSlackUserStatus()
				if err != nil {
					return cli.NewExitError(err, 1)
				}
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
				fmt.Print(listTemplates(""))
				return nil
			},
		},
		{
			Name:      "example",
			Aliases:   []string{"x"},
			Usage:     "Shows an example settings schema",
			ArgsUsage: " ",
			Action: func(ctx *cli.Context) error {
				fmt.Println(internal.SettingsExample)
				return nil
			},
		},
		{
			Name:    "set",
			Aliases: []string{"s"},
			Usage:   "Updates your status",
			Description: `Template IDs:` + "\n" +
				listTemplates("     ") +
				"   \n" +
				`   Special template IDs:` + "\n" +
				`     ` + settings.ITunes.TemplateID + ` = appends information about the music playing on iTunes` + "\n" +
				`     ` + settings.LastFM.TemplateID + ` = appends information about the music scrobbled to last.fm`,
			ArgsUsage: "[<template ID> ...]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dryrun, d",
					Usage: "just print the composed status text (your status will be not changed)",
				},
				cli.BoolFlag{
					Name:  "watch, w",
					Usage: `watch changes (with "` + settings.ITunes.TemplateID + `" or "` + settings.LastFM.TemplateID + `" template)`,
				},
			},
			Action: func(ctx *cli.Context) error {
				flags.dryRun = ctx.Bool("dryrun")
				flags.watch = ctx.Bool("watch")
				ids := []string(ctx.Args())

				withInfo := 0
				interval := time.Duration(30)
				for _, id := range ids {
					switch id {
					case settings.ITunes.TemplateID:
						withInfo++
						interval = settings.ITunes.WatchIntervalSec
						if interval < 1 {
							warn(`itunes.watch_interval_sec must be >= 1`)
							interval = 5
						}
					case settings.LastFM.TemplateID:
						withInfo++
						interval = settings.LastFM.WatchIntervalSec
						if interval < 15 {
							warn(`lastfm.watch_interval_sec must be >= 15`)
							interval = 15
						}
					}
				}
				if 1 < withInfo {
					return cliError(fmt.Sprintf(
						`Both "%s" and "%s" cannot be specified at the same time.`,
						settings.ITunes.TemplateID,
						settings.LastFM.TemplateID,
					))
				}
				if withInfo == 0 {
					flags.watch = false
				}
				if len(ids) == 0 && withInfo == 0 {
					cli.ShowCommandHelp(ctx, "set")
					os.Exit(1)
				}

				err := update(&flags, ids)
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				for flags.watch {
					time.Sleep(interval * time.Second)
					err = update(&flags, ids)
					if err != nil {
						warn(err.Error())
					}
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

	if 0 < len(internal.SettingsWarnings) {
		w := append([]string{
			`Your setting file seems to be corrupted.`,
			`Try "slack-status example" to show an example settings schema.`,
			``,
		}, internal.SettingsWarnings...)
		warn(w...)
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

func appendMusicInfo(settings *internal.MusicSettings, status *internal.MusicStatus) string {
	r := regexp.MustCompile(`%\w`)
	return r.ReplaceAllStringFunc(settings.Format, func(m string) string {
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

func update(flags *Flags, templateIDs []string) error {
	now := time.Now().Format("[15:04:05] ")
	notice := []string{}

	tmpls := []string{}
	for _, id := range templateIDs {
		tmpl, ok := settings.Templates[id]
		if ok {
			tmpls = append(tmpls, tmpl)
			continue
		}
		switch id {
		case settings.ITunes.TemplateID:
			status := &internal.GetITunesStatus().MusicStatus
			if status.Ok {
				tmpls = append(tmpls, appendMusicInfo(&settings.ITunes.MusicSettings, status))
				continue
			}
			n := "iTunes seems to be stopped"
			if status.Err != "" {
				n += " : " + status.Err
			}
			notice = append(notice, n)
		case settings.LastFM.TemplateID:
			status := &internal.GetLastFMStatus().MusicStatus
			if status.Ok {
				tmpls = append(tmpls, appendMusicInfo(&settings.LastFM.MusicSettings, status))
				continue
			}
			n := "Failed to fetch music information from last.fm"
			if status.Err != "" {
				n += " : " + status.Err
			}
			notice = append(notice, n)
		default:
			return cliError(
				`Template "`+id+`" is not defined in settings file.`,
				`Try "slack-status list" to list your templates.`,
			)
		}
	}

	t := strings.Join(tmpls, " ")
	e, t := splitEmoji(t)
	t = limitStringByLength(t, internal.SlackUserStatusMaxLength)
	changed := updatedCount == 0 || !(t == lastText && e == lastEmoji)

	if changed {
		var err error
		if !flags.dryRun {
			err = internal.SetSlackUserStatus(t, e)
		}
		if 0 < len(notice) {
			if flags.watch {
				red.Fprint(os.Stderr, now)
			}
			red.Fprintln(os.Stderr, strings.Join(notice, "\n"))
		}
		if flags.watch {
			cyan.Print(now)
		}
		printStatus(e, t)
		if err != nil {
			return err
		}
	}

	lastText = t
	lastEmoji = e
	updatedCount++

	return nil
}
