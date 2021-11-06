package models

type AuthData struct {
	Email    string `json:"email,omitempty" form:"email" query:"email"`
	Password string `json:"password,omitempty" form:"password" query:"password"`
}
