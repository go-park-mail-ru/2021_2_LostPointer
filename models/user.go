package models

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string
	Name 	 string `json:"name"`
}
