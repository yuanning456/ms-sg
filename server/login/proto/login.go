package proto



type LoginRsq struct {
	UserName	string		`json:"username"`
	Password 	string		`json:"password"`
	Session 	string		`json:"session"`
	UId 		int			`json:"uid"`
}

type LoginReq struct {
	UserName	string		`json:"username"`
	Password 	string		`json:"password"`
	Ip 			string		`json:"ip"`
	Hardware 	string		`json:"hardware"`
}


