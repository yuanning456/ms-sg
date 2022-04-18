package net

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type server struct {
	addr string
	router *Router
}

func NewServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

func (s *server) Router(router *Router) {
	s.router = router
}

//启动服务   
func (s *server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
}

//http升级websocket配置
var wsUpGrader = websocket.Upgrader{
	//允许所有cors跨域请求
	CheckOrigin : func(r *http.Request) bool {
		return true
	},
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	//websocket 
	//1.http协议升级为websocket
	wsConn, err := wsUpGrader.Upgrade(w, r, nil)
	if err != nil {
		//打印日志同时还会推出应用程序
		log.Fatal("websocket服务链接出错")
		//不退出程序
		// log.Println("websocket服务链接出错")
	}

	//websocket 通道建立后 客户端和服务端都可以收发消息
	//发消息的时候 把消息当成路由来处理 消息是有格式的，先定义消息的格式
	// wsConn.WriteMessage()
	// wsConn.ReadMessage()
	//客户端发消息的时候  按照{Name:"account.login"} 解析路由
	// err = wsConn.WriteMessage(websocket.BinaryMessage, []byte("hello"))
	// fmt.Println(err)
	// for {
	// 	i, p, _ := wsConn.ReadMessage()
	// 	fmt.Println(i, string(p))
	// }
	wsServer := NewWsServer(wsConn)
	wsServer.router = s.router
	wsServer.Start()
	wsServer.Handshake()
}





