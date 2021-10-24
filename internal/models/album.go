package models

type Album struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	Year           int64  `json:"year"`
	Artist         string `json:"artist,omitempty"`
	Artwork        string `json:"artwork"`
	TracksCount    int64  `json:"tracksCount"`
	TracksDuration int64  `json:"tracksDuration"`
}
