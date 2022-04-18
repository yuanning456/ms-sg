package login

import (
	"ms-sg/net"
	"ms-sg/server/login/controller"
)


var Router = &net.Router{}

func Init() {
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}