package internal

import (
	"fmt"
	"regexp"
	"strings"

	"os"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

func WrapEmoji(e string) string {
	if e == "" {
		return e
	}
	return ":" + e + ":"
}

func LimitStringByLength(str string, maxlen int) string {
	r := []rune(str)
	if len(r) <= maxlen {
		return str
	}
	return string(r[:maxlen-1]) + "…"
}

var splitEmojiRegexp = regexp.MustCompile(`^:[^: ]+: *`)

func SplitEmoji(text string) (string, string) {
	text = strings.Trim(text, " ")
	e := ""
	m := splitEmojiRegexp.FindString(text)
	if m != "" {
		e = strings.Trim(m, " ")
		text = text[len(m):]
	}
	return e, text
}

var bggray = color.New(color.BgHiBlack)

func PrintStatus(e, t string) {
	if e != "" {
		bggray.Print(emoji.Sprint(e))
		fmt.Print(" ")
	}
	emoji.Println(t)
}

var yellow = color.New(color.FgYellow)

func Warn(msgs ...string) {
	for _, msg := range msgs {
		lines := strings.Split(msg, "\n")
		for _, line := range lines {
			yellow.Fprintln(os.Stderr, `[warning] `+line)
		}
	}
	fmt.Fprintln(os.Stderr, "")
}
