package helper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kyokomi/emoji"
	"golang.org/x/text/encoding/japanese"

	"bytes"
	"io/ioutil"
	"strconv"

	"golang.org/x/text/transform"
)

func ParseFloat64(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return .0
	}
	return val
}

func ConvertFromShiftJIS(sjis []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(sjis), japanese.ShiftJIS.NewDecoder())
	utf8, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(utf8), nil
}

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

func PrintStatus(e, t string) {
	if e != "" {
		BgHiBlack.Print(emoji.Sprint(e))
		fmt.Print(" ")
	}
	emoji.Println(t)
}
