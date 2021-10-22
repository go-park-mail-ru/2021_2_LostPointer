package models

type Csrf struct {
	Token string `json:"token,omitempty"`
}
