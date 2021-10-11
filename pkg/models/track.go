package models

type Track struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Explicit    bool   `json:"explicit"`
	Genre       string `json:"genre"`
	Number      int64  `json:"number"`
	ListenCount int64  `json:"listenCount"`
	Duration    int64  `json:"duration"`
	Cover       string `json:"cover"`
}