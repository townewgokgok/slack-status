package internal

import (
	"os/exec"
	"strconv"
	"strings"
)

type ITunesStatus struct {
	MusicStatus
	Position float64
	Duration float64
	Start    float64
	Finish   float64
}

func parseFloat64(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return .0
	}
	return val
}

func GetITunesStatus() *ITunesStatus {
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
		return &ITunesStatus{}
	}
	values := strings.Split(string(tsv), "\t")
	return &ITunesStatus{
		MusicStatus: MusicStatus{
			Valid:  true,
			Artist: values[4],
			Album:  values[5],
			Title:  values[6],
		},
		Position: parseFloat64(values[0]),
		Duration: parseFloat64(values[1]),
		Start:    parseFloat64(values[2]),
		Finish:   parseFloat64(values[3]),
	}
}
