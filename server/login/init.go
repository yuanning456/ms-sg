package login

import (
	"ms-sg/db"
	"ms-sg/net"
	"ms-sg/server/login/controller"
)


var Router = &net.Router{}

func Init() {
	initRouter()
	db.Init()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}