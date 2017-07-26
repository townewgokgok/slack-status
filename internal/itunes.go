package internal

import (
	"os/exec"
	"strconv"
	"strings"
)

type ITunesStatus struct {
	Valid    bool
	Position float64
	Duration float64
	Start    float64
	Finish   float64
	Artist   string
	Name     string
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
		` & tab & name of current track` +
		` & tab`
	tsv, err := exec.Command("osascript", "-e", scpt).Output()
	if err != nil {
		return &ITunesStatus{Valid: false}
	}
	values := strings.Split(string(tsv), "\t")
	return &ITunesStatus{
		Valid:    true,
		Position: parseFloat64(values[0]),
		Duration: parseFloat64(values[1]),
		Start:    parseFloat64(values[2]),
		Finish:   parseFloat64(values[3]),
		Artist:   values[4],
		Name:     values[5],
	}
}
