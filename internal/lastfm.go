package internal

import (
	"strconv"
	"time"

	"github.com/shkh/lastfm-go/lastfm"
)

type LastFMStatus struct {
	MusicStatus
	Date time.Time
}

func GetLastFMStatus() *LastFMStatus {
	s := Settings.LastFM
	api := lastfm.New(s.APIKey, s.Secret)
	tracks, err := api.User.GetRecentTracks(lastfm.P{
		"user": s.UserName,
	})
	if err != nil {
		panic(err)
	}
	if len(tracks.Tracks) == 0 {
		return &LastFMStatus{}
	}
	tr := tracks.Tracks[0]
	uts, _ := strconv.ParseInt(tr.Date.Uts, 10, 64)
	return &LastFMStatus{
		MusicStatus: MusicStatus{
			Valid:  true,
			Artist: tr.Artist.Name,
			Album:  tr.Album.Name,
			Title:  tr.Name,
		},
		Date: time.Unix(uts, 0),
	}
}
