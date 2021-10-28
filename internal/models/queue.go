package models

type Current struct {
	Index int     `json:"index"`
	Time  float64 `json:"time"`
}

type TrackData struct {
	Artist    string `json:"artist"`
	Cover     string `json:"cover"`
	TotalTime string `json:"total_time"`
	Track     string `json:"track"`
	Url       string `json:"url"`
}

type Queue struct {
	Volume  float64     `json:"volume"`
	Mute    bool        `json:"mute"`
	Shuffle bool        `json:"shuffle"`
	Current Current     `json:"current"`
	Queue   []TrackData `json:"queue"`
}