package main

import (
	"ms-sg/config"
	"ms-sg/net"
	"ms-sg/server/gate"
)

/**
1.登陆功能account.login 需要通过网关转发到登陆服务器
2.网关(websocket  客户端)和登陆服务器(websocket 服务端)交互
3.网关和游戏的客户端交互 这时候网关它是websocket的服务端
4.websocket服务端 已经实现了
5.websocket客户端 需要实现
6.网关代理服务器 代理请求（代理地址 多个；代理连接通道）  客户端连接websocket的连接
7.路由：接受所有的请求(网关websocket服务端的功能)
8.握手协议 检测第一次连接合法
*/

func main() {
	host := config.File.MustValue("gate_server", "host", "127.0.0.1")
	port := config.File.MustValue("gate_server", "port", "9992")

	s := net.NewServer(host + ":" + port)
	s.NeedSecret(true)
	gate.Init()
	s.Router(gate.Router)
	s.Start()
}