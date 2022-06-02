package net

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

type ProxyClient struct {
	proxy string //代理的地址
	conn *ClientConn
}

func (c *ProxyClient) SetProperty(key string, data interface{}) {
	if c.conn != nil {
		c.conn.SetProperty(key, data)
	}
}
func (c *ProxyClient) SetOnPush(hook func(conn *ClientConn,body *RspBody)) {
	if c.conn != nil {
		c.conn.SetOnPush(hook)
	}
}
func (c *ProxyClient) Send(name string, msg interface{}) (*RspBody, error) {
	if c.conn != nil {
		return c.conn.Send(name, msg), nil
	}
	return nil, errors.New("未找到连接")
}
func NewProxyClient(proxy string) *ProxyClient {
	return &ProxyClient{
		proxy: proxy,
	}
}

//固定的 抄下来就行
func (c *ProxyClient) Connet() error {
	//去连接websocket服务
	//通过Dialer连接websocket服务器
	var dialer = websocket.Dialer{
		Subprotocols:     []string{"p1", "p2"},
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 30 * time.Second,
	}
	ws, _, err := dialer.Dial(c.proxy, nil)
	if err == nil{
		c.conn = NewClientConn(ws)
		if !c.conn.Start(){
			return errors.New("握手失败")
		}
	}
	return err
}