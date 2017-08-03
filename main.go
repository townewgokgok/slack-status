package main

import (
	"fmt"

	"github.com/urfave/cli"

	"os"

	"time"

	"strings"

	"sort"

	"github.com/fatih/color"
	"github.com/townewgokgok/slack-status/internal"
	"github.com/townewgokgok/slack-status/internal/music"
	s "github.com/townewgokgok/slack-status/internal/settings"
)

var version = "1.2.0"

var cyan = color.New(color.FgCyan)
var red = color.New(color.FgRed)

type Flags struct {
	dryRun bool
	watch  bool
}

var flags Flags

func cliError(msgs ...string) *cli.ExitError {
	return cli.NewExitError(red.Sprint(strings.Join(msgs, "\n")), 1)
}

func main() {
	// Parse arguments
	app := cli.NewApp()
	app.EnableBashCompletion = true
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
				s.Edit()
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
				internal.PrintStatus(e, t)
				return nil
			},
		},
		{
			Name:      "list",
			Aliases:   []string{"l"},
			Usage:     "Lists your templates",
			ArgsUsage: " ",
			Action: func(ctx *cli.Context) error {
				fmt.Print(s.Settings.Templates.Dump(""))
				return nil
			},
		},
		{
			Name:      "example",
			Aliases:   []string{"x"},
			Usage:     "Shows an example settings schema",
			ArgsUsage: " ",
			Action: func(ctx *cli.Context) error {
				fmt.Println(s.SettingsExample)
				return nil
			},
		},
		{
			Name:    "set",
			Aliases: []string{"s"},
			Usage:   "Updates your status",
			Description: `Template IDs:` + "\n" +
				s.Settings.Templates.Dump("     ") +
				"   \n" +
				`   Special template IDs:` + "\n" +
				`     ` + s.Settings.ITunes.TemplateID + ` = appends information about the music playing on iTunes` + "\n" +
				`     ` + s.Settings.LastFM.TemplateID + ` = appends information about the music scrobbled to last.fm`,
			ArgsUsage: "[<template ID> ...]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dryrun, d",
					Usage: "just print the composed status text (your status will be not changed)",
				},
				cli.BoolFlag{
					Name:  "watch, w",
					Usage: `watch changes (with "` + s.Settings.ITunes.TemplateID + `" or "` + s.Settings.LastFM.TemplateID + `" template)`,
				},
			},
			BashComplete: func(c *cli.Context) {
				ids := []string{
					s.Settings.ITunes.TemplateID,
					s.Settings.LastFM.TemplateID,
				}
				for id := range s.Settings.Templates {
					ids = append(ids, id)
				}
				sort.Strings(ids)
				fmt.Println(strings.Join(ids, "\n"))
			},
			Action: func(ctx *cli.Context) error {
				flags.dryRun = ctx.Bool("dryrun")
				flags.watch = ctx.Bool("watch")
				ids := []string(ctx.Args())

				withInfo := 0
				interval := time.Duration(30)
				for _, id := range ids {
					switch id {
					case s.Settings.ITunes.TemplateID:
						withInfo++
						interval = s.Settings.ITunes.WatchIntervalSec
						if interval < 1 {
							internal.Warn(`itunes.watch_interval_sec must be >= 1`)
							interval = 5
						}
					case s.Settings.LastFM.TemplateID:
						withInfo++
						interval = s.Settings.LastFM.WatchIntervalSec
						if interval < 15 {
							internal.Warn(`lastfm.watch_interval_sec must be >= 15`)
							interval = 15
						}
					}
				}
				if 1 < withInfo {
					return cliError(fmt.Sprintf(
						`Both "%s" and "%s" cannot be specified at the same time.`,
						s.Settings.ITunes.TemplateID,
						s.Settings.LastFM.TemplateID,
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
						internal.Warn(err.Error())
					}
				}

				return nil
			},
		},
	}

	app.Action = func(ctx *cli.Context) error {
		if s.Settings.Slack.Token == "" || strings.ContainsRune(s.Settings.Slack.Token, '.') {
			return cliError(
				`Your settings file seems to be not configured correctly.`,
				`The example settings file has been created at `+s.SettingsPath,
				`Try "slack-status edit" to edit it.`,
			)
		}
		cli.ShowAppHelp(ctx)
		return nil
	}

	if 0 < len(s.SettingsWarnings) {
		w := append([]string{
			`Your setting file seems to be corrupted.`,
			`Try "slack-status example" to show an example settings schema.`,
			``,
		}, s.SettingsWarnings...)
		internal.Warn(w...)
	}

	app.Run(os.Args)
}

var lastText string
var lastEmoji string
var updatedCount int

func update(flags *Flags, templateIDs []string) error {
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
			return cliError(
				`Template "`+id+`" is not defined in settings file.`,
				`Try "slack-status list" to list your templates.`,
			)
		}
	}

	t := strings.Join(tmpls, " ")
	e, t := internal.SplitEmoji(t)
	t = internal.LimitStringByLength(t, internal.SlackUserStatusMaxLength)
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
		internal.PrintStatus(e, t)
		if err != nil {
			return err
		}
	}

	lastText = t
	lastEmoji = e
	updatedCount++

	return nil
}
