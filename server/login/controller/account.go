package controller

import (
	"ms-sg/net"
 	"ms-sg/server/login/proto"
)

var DefaultAccount = &Account{}
type Account struct {

}
func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsq *net.WsMsgRsp) {
	rsq.Body.Code = 0
	loginRes := &proto.LoginRsq{}
	loginRes.UId = 1
	loginRes.UserName = "admin"
	loginRes.Password = "admin"
	rsq.Body.Msg = loginRes
	
}