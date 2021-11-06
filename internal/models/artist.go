package models

type Artist struct {
	Id     int64   `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Avatar string  `json:"avatar,omitempty"`
	Video  string  `json:"video,omitempty"`
}
