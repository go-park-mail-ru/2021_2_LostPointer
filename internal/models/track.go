package models

type Track struct {
	Id     int64  `json:"id"`
	Title  string `json:"title"`
	Explicit    bool   `json:"explicit"`
	Genre       string `json:"genre,omitempty"`
	Number      int64  `json:"number,omitempty"`
	File        string `json:"file"`
	ListenCount int64  `json:"listenCount"`
	Duration    int64  `json:"duration"`
	Lossless    bool   `json:"lossless"`
	Album  Album `json:"album"`
	Artist Artist `json:"artist"`
}
