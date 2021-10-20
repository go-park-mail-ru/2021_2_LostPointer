package models

type Artist struct {
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Tracks []Track `json:"tracks,omitempty"`
	Albums []Album `json:"albums,omitempty"`
}
