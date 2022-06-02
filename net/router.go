package net

import (
	"log"
	"strings"
)


type group struct {
	prefix string
	handlerMap map[string]HandlerFunc
}

func (g *group) AddRouter(name string, h HandlerFunc) {
	g.handlerMap[name] = h
}

func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix: prefix,
		handlerMap: make(map[string]HandlerFunc),
	}
	r.group = append(r.group, g)
	return g
}
type HandlerFunc func(req *WsMsgReq, rsq *WsMsgRsp)

type Router struct {
	group []*group
}

func (r Router) Run(req *WsMsgReq, rsq *WsMsgRsp) {
	// req.Body.Name 路径 登陆业务account.login(account组表示) login路由表示
	strs := strings.Split(req.Body.Name, ".")
	prefix 	:= ""
	name 	:= ""
	if len(strs) == 2 {
		prefix	= strs[0]
		name 	= strs[1]
	}
	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, rsq)
		} else if g.prefix == "*" {
			g.exec(name, req, rsq)
		} 
	}  
}

func (g *group) exec(name string, req *WsMsgReq, rsq *WsMsgRsp) () {
	h := g.handlerMap[name]
	if h!= nil {
		h(req, rsq)
	} else {
		h := g.handlerMap["*"]
		if h!= nil {
			h(req, rsq)
		} else {
			log.Println("找不到路由")
		}
	}
}