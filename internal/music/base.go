package music

type MusicStatus struct {
	Ok     bool
	Err    string
	Artist string
	Album  string
	Title  string
}

func (s *MusicStatus) Replacer(m string) (string, bool) {
	switch m {
	case "%A":
		return s.Artist, true
	case "%a":
		return s.Album, true
	case "%t":
		return s.Title, true
	}
	return "", false
}
