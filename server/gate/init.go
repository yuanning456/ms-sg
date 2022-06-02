package gate

import (
	"ms-sg/net"
	"ms-sg/server/gate/controller"
)


var Router = &net.Router{}

func Init() {
	 initRouter()
}

func initRouter() {
	controller.GateHandler.Router(Router)
}