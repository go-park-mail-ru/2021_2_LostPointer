package models

type User struct {
	Email    	 string `json:"email" form:"email" query:"email"`
	Password 	 string `json:"password" form:"password" query:"password"`
	Nickname     string `json:"nickname" form:"nickname" query:"nickname"`
}
