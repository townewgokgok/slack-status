package music

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/townewgokgok/slack-status/internal/helper"
)

type ITunesStatus struct {
	MusicStatus
	Position float64
	Duration float64
	Start    float64
	Finish   float64
}

func GetITunesStatus() *ITunesStatus {
	var tsv string
	if runtime.GOOS == "windows" {
		tsv = getITunesStatusForWindows()
	} else {
		tsv = getITunesStatusForMac()
	}
	if tsv == "" {
		return &ITunesStatus{}
	}
	values := strings.Split(tsv, "\t")
	if values[1] == "" {
		return &ITunesStatus{}
	}
	return &ITunesStatus{
		MusicStatus: MusicStatus{
			Ok:     true,
			Artist: values[4],
			Album:  values[5],
			Title:  values[6],
		},
		Position: helper.ParseFloat64(values[0]),
		Duration: helper.ParseFloat64(values[1]),
		Start:    helper.ParseFloat64(values[2]),
		Finish:   helper.ParseFloat64(values[3]),
	}
}

func getITunesStatusForWindows() string {
	d := `+"` + "`t" + `"+`
	scpt := `$i = New-Object -Com "iTunes.Application";` +
		`$t = $i.CurrentTrack;` +
		`""+` +
		`$i.PlayerPosition` + d +
		`$t.Duration` + d +
		`$t.Start` + d +
		`$t.Finish` + d +
		`$t.Artist` + d +
		`$t.Album` + d +
		`$t.Name` + d +
		`""`
	sjis, err := exec.Command("powershell", "-noprofile", "-noninteractive", "-command", scpt).Output()
	if err != nil {
		return ""
	}
	tsv, err := helper.ConvertFromShiftJIS(sjis)
	if err != nil {
		return ""
	}
	return string(tsv)
}

func getITunesStatusForMac() string {
	scpt := `tell app "itunes" to ""` +
		` & player position` +
		` & tab & duration of current track` +
		` & tab & start of current track` +
		` & tab & finish of current track` +
		` & tab & artist of current track` +
		` & tab & album of current track` +
		` & tab & name of current track` +
		` & tab`
	tsv, err := exec.Command("osascript", "-e", scpt).Output()
	if err != nil {
		return ""
	}
	return string(tsv)
}
