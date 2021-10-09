package models

type User struct {
	Email    	 string `json:"email,omitempty"`
	Password 	 string `json:"password,omitempty"`
	NickName     string `json:"nickname,omitempty"`
}
