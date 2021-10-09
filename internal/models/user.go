package models

type User struct {
	Email    	 string `json:"email,omitempty"`
	Password 	 string `json:"password,omitempty"`
	Nickname     string `json:"nickname,omitempty"`
}
