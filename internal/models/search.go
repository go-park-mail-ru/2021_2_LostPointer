package models

type ArtistShort struct {
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

type SearchResponseData struct {
	Tracks  []Track       `json:"tracks,omitempty"`
	Artists []ArtistShort `json:"artists,omitempty"`
}
