package models

type AuthData struct {
	Email    string `json:"email,omitempty" form:"email" query:"email"`
	Password string `json:"password,omitempty" form:"password" query:"password"`
}

type RegisterData struct {
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
	Nickname string `json:"nickname" form:"nickname" query:"nickname"`
}
