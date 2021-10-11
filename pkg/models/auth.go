package models

type Auth struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
