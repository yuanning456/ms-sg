package controller

import (
	"log"
	"ms-sg/config"
	"ms-sg/constant"
	"ms-sg/net"
	"strings"
	"sync"
)

var GateHandler = &Handler{
	proxyMap: make(map[string]map[int64]*net.ProxyClient),
}

type Handler struct {
	proxyMutex sync.Mutex
	//代理地址 -> 客户端连接 (游戏ID -> 连接)
	proxyMap map[string]map[int64]*net.ProxyClient // map[string]是代理的地址  map[int64]是每一个客户端的连接 
	loginProxy string
	gameProxy string
}

func (h *Handler) Router(r *net.Router) {
	h.loginProxy = config.File.MustValue("gate_server", "login_proxy", "127.0.0.1:9993")
	h.gameProxy = config.File.MustValue("gate_server", "game_proxy", "127.0.0.1:9994")
	g := r.Group("*")
	g.AddRouter("*", h.all)
}


func (h *Handler) all(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	log.Println("网关的处理器")
	//account  转发
	name := req.Body.Name
	proxyStr := ""
	if isAccount(name) {
		proxyStr = h.loginProxy
	}
	if proxyStr == "" {
		rsp.Body.Code = constant.ProxyNotInConnect
		return 
	}
	h.proxyMutex.Lock()
	_, ok := h.proxyMap[proxyStr]
	if !ok {
		h.proxyMap[proxyStr] = make(map[int64]*net.ProxyClient)
	}
	h.proxyMutex.Unlock() 
	//客户端ID
	c, err := req.Conn.GetProperty("cid")
	if err != nil {
		log.Println("cid  没取到")
		rsp.Body.Code = constant.InvalidParam
		return
	}  
	cid := c.(int64)
	proxy := h.proxyMap[proxyStr][cid]
	if proxy == nil {
		proxy = net.NewProxyClient(proxyStr)
		err := proxy.Connet() 
		if err != nil {
			h.proxyMutex.Lock()
			delete(h.proxyMap[proxyStr], cid)
			h.proxyMutex.Unlock()
			rsp.Body.Code = constant.ProxyNotInConnect
			return  
		}
		h.proxyMap[proxyStr][cid] = proxy
		proxy.SetProperty("cid", cid)
		proxy.SetProperty("proxy", proxyStr)
		proxy.SetProperty("gateConn", req.Conn)
		proxy.SetOnPush(h.onPush)
	}
	
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	r, err := proxy.Send(req.Body.Name, req.Body.Msg)
	if r != nil {
		rsp.Body.Code  = r.Code
		rsp.Body.Msg  = r.Msg
	} else {
		rsp.Body.Code = constant.ProxyNotInConnect
		return 
	}

}

func (h *Handler) onPush(conn *net.ClientConn,body *net.RspBody) {
	gc, err := conn.GetProperty("gateConn")
	if err != nil {
		log.Println("onpush getConn", err)
	}
	getConn := gc.(net.WSConn)
	getConn.Push(body.Name, body.Msg)
}

func isAccount(name string) bool {
	return strings.HasPrefix(name, "account.")
}