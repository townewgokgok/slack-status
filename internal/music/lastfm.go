package music

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

func GetLastFMStatus(apikey, username string) *LastFMStatus {
	spn := spinner.New(spinner.CharSets[14], time.Second/30)
	spn.Start()
	defer spn.Stop()
	api := lastfm.New(apikey, "")
	tracks, err := api.User.GetRecentTracks(lastfm.P{
		"user": username,
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
