package internal

import (
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/shkh/lastfm-go/lastfm"
)

type LastFMStatus struct {
	MusicStatus
	Date time.Time
}

func GetLastFMStatus() *LastFMStatus {
	spn := spinner.New(spinner.CharSets[14], time.Second/30)
	spn.Start()
	defer spn.Stop()
	s := Settings.LastFM
	api := lastfm.New(s.APIKey, "")
	tracks, err := api.User.GetRecentTracks(lastfm.P{
		"user": s.UserName,
	})
	if err != nil {
		return &LastFMStatus{
			MusicStatus: MusicStatus{Err: err.Error()},
		}
	}
	if len(tracks.Tracks) == 0 {
		return &LastFMStatus{}
	}
	tr := tracks.Tracks[0]
	uts, _ := strconv.ParseInt(tr.Date.Uts, 10, 64)
	return &LastFMStatus{
		MusicStatus: MusicStatus{
			Ok:     true,
			Artist: tr.Artist.Name,
			Album:  tr.Album.Name,
			Title:  tr.Name,
		},
		Date: time.Unix(uts, 0),
	}
}
