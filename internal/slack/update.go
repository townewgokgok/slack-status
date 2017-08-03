package slack

import (
	"os"
	"strings"
	"time"

	"github.com/townewgokgok/slack-status/internal/generator"
	"github.com/townewgokgok/slack-status/internal/helper"
)

var lastText string
var lastEmoji string
var updatedCount int

func Update(watch, dryRun bool, templateIDs []string) error {
	now := time.Now().Format("[15:04:05] ")

	emj, txt, notice, err := generator.Generate(templateIDs)
	if err != nil {
		return err
	}

	txt = helper.LimitStringByLength(txt, SlackUserStatusMaxLength)
	changed := updatedCount == 0 || !(txt == lastText && emj == lastEmoji)

	if changed {
		var err error
		if !dryRun {
			err = SetSlackUserStatus(txt, emj)
		}
		if 0 < len(notice) {
			if watch {
				helper.Red.Fprint(os.Stderr, now)
			}
			helper.Red.Fprintln(os.Stderr, strings.Join(notice, "\n"))
		}
		if watch {
			helper.Cyan.Print(now)
		}
		helper.PrintStatus(emj, txt)
		if err != nil {
			return err
		}
	}

	lastText = txt
	lastEmoji = emj
	updatedCount++

	return nil
}
