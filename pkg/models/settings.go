package models

type Settings struct {
	Email    string `json:"email,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}
