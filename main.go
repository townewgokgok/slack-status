package main

import (
	"fmt"

	"github.com/urfave/cli"

	"os"

	"time"

	"strings"

	"sort"

	"github.com/townewgokgok/slack-status/internal/helper"
	s "github.com/townewgokgok/slack-status/internal/settings"
	"github.com/townewgokgok/slack-status/internal/slack"
)

var version = "1.2.0"

func cliError(msgs ...string) *cli.ExitError {
	return cli.NewExitError(helper.Red.Sprint(strings.Join(msgs, "\n")), 1)
}

func warn(msgs ...string) {
	for _, msg := range msgs {
		lines := strings.Split(msg, "\n")
		for _, line := range lines {
			helper.Yellow.Fprintln(os.Stderr, `[warning] `+line)
		}
	}
	fmt.Fprintln(os.Stderr, "")
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
				e, t, err := slack.GetSlackUserStatus()
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				helper.PrintStatus(e, t)
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
				dryRun := ctx.Bool("dryrun")
				watch := ctx.Bool("watch")
				ids := []string(ctx.Args())

				withInfo := 0
				interval := time.Duration(30)
				for _, id := range ids {
					switch id {
					case s.Settings.ITunes.TemplateID:
						withInfo++
						interval = s.Settings.ITunes.WatchIntervalSec
						if interval < 1 {
							warn(`itunes.watch_interval_sec must be >= 1`)
							interval = 5
						}
					case s.Settings.LastFM.TemplateID:
						withInfo++
						interval = s.Settings.LastFM.WatchIntervalSec
						if interval < 15 {
							warn(`lastfm.watch_interval_sec must be >= 15`)
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
					watch = false
				}
				if len(ids) == 0 && withInfo == 0 {
					cli.ShowCommandHelp(ctx, "set")
					os.Exit(1)
				}

				err := slack.Update(watch, dryRun, ids)
				if err != nil {
					return cliError(err.Error())
				}

				for watch {
					time.Sleep(interval * time.Second)
					err = slack.Update(watch, dryRun, ids)
					if err != nil {
						warn(err.Error())
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
		warn(w...)
	}

	app.Run(os.Args)
}
