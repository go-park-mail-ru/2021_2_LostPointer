package models

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type AvatarResponse struct {
	Status int    `json:"status"`
	Avatar string `json:"avatar"`
}
