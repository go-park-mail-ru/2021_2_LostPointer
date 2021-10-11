package models

type Album struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	Year           int64  `json:"year"`
	Artist         string `json:"artist"`
	ArtWork        string `json:"artWork"`
	TracksCount    int64  `json:"tracksCount"`
	TracksDuration int64  `json:"tracksDuration"`
}
