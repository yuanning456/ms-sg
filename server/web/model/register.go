package model


type RegisterReq struct {
	Username	string `form:"username" json:"username"`
	Password 	string `form:"passrod" json:"passrod"`
	Ip 			string `form:"ip" json:"ip"`
	Hardware 	string	`form:"hardware" json:"hardware"`
}