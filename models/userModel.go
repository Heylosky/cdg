package models

type User struct {
	Name      string   `form:"name" json:"name" xml:"name"`
	Password  string   `form:"password" json:"password" xml:"password"`
	Email     string   `form:"email" json:"email" xml:"email"`
	TokenList []string `form:"token_list" json:"token_list" xml:"token_list"`
}
