package models

type ArtistShort struct {
	ID     int64  `json:"id,omitempty"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

type SearchPayload struct {
	Tracks  []Track       `json:"tracks,omitempty"`
	Artists []ArtistShort `json:"artists,omitempty"`
}
