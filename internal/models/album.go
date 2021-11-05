package models

type Album struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	Year           int64  `json:"year,omitempty"`
	Artist         string `json:"artist,omitempty"`
	Artwork        string `json:"artwork,omitempty"`
	TracksCount    int64  `json:"tracksCount,omitempty"`
	TracksDuration int64  `json:"tracksDuration,omitempty"`
	ArtWorkColor   string `json:"artwork_color,omitempty"`
}
