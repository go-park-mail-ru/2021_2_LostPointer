package models

type Track struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Explicit    bool   `json:"explicit"`
	Genre       string `json:"genre"`
	Number      int64  `json:"number"`
	File        string `json:"file"`
	ListenCount int64  `json:"listenCount"`
	Duration    int64  `json:"duration"`
	Cover       string `json:"cover"`
}

type Album struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	Year           int64  `json:"year"`
	Artist         string `json:"artist"`
	ArtWork        string `json:"artWork"`
	TracksCount     int64  `json:"tracksCount"`
	TracksDuration int64  `json:"tracksDuration"`
}

type Artist struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
}

type Playlist struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	User int64  `json:"user"`
}

type SelectionFroHomePage struct {
	Tracks    []Track    `json:"tracks"`
	Albums    []Album    `json:"albums"`
	Playlists []Playlist `json:"playlists"`
	Artists   []Artist   `json:"artists"`
}
