package models

type Playlist struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	User int64  `json:"user"`
}
