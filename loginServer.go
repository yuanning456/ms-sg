package main

import (
	"ms-sg/config"
	"ms-sg/net"
	"ms-sg/server/login"
)

// http://localhost:8080/api/login
// websocket ws://localhost:8080 发消息的 封装为路由



func main() {
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "9991")

	s := net.NewServer(host + ":" + port)
	s.NeedSecret(false)
	login.Init()
	s.Router(login.Router)
	s.Start()

}

