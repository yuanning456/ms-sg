package net

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"ms-sg/constant"
	"ms-sg/utils"
	"sync"
	"time"

	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type syncCtx struct {
	//Goroutine 的上下文，包含 Goroutine 的运行状态、环境、现场等信息
	ctx context.Context
	cancel context.CancelFunc
	outChan chan *RspBody
}

func NewSyncCtx() *syncCtx {
	ctx, cancal := context.WithTimeout(context.Background(), 15 * time.Second)
	return &syncCtx{
		ctx: ctx,
		cancel: cancal,
		outChan: make(chan *RspBody),
	}
}

func (s *syncCtx) wait() *RspBody {
	select {
	case msg := <- s.outChan:
		return msg
	case  <- s.ctx.Done():
		log.Println("代理服务响应消息超时")
		return nil 
	}
}

type ClientConn struct {
	wsConn *websocket.Conn
	isClosed bool
	property 		map[string]interface{} //设置一些属性
	propertyLock  	sync.RWMutex
	Seq				int64
	handshake bool //握手的状态
	handshakeChan chan bool //握手消息成功 通知通道 
	onPush func(conn *ClientConn,body *RspBody) // 给代理服务器发送消息的
	onClose    		func(conn*ClientConn) //关闭的时候需要处理的 一些操作
	syncCtxMap map[int64]*syncCtx //int64客户端的连接 *syncCtx是我们接受到的数据
	syncCtxLock sync.RWMutex
}

func NewClientConn(wsConn *websocket.Conn) *ClientConn{
	return &ClientConn{
		wsConn: wsConn,
		handshakeChan: make(chan bool),
		Seq: 0,
		isClosed: false,
		property: make(map[string]interface{}),
		syncCtxMap: map[int64]*syncCtx{},
	}
}

func (c *ClientConn) Start() bool {
	//一直不停的接收消息
	//等待握手的消息返回
	c.handshake = false
	go c.wsReadLoop()
	return c.waitHandShake()
}

func (c *ClientConn) Close() {
	_ = c.wsConn.Close()
}

func (c *ClientConn)wsReadLoop() {
	// for {
	// 	_, data, err := c.wsConn.ReadMessage()
	// 	fmt.Println(data)
	// 	fmt.Println(err)
	// 	//读取消息  会收到 握手 心跳 请求消息

	// 	c.handshake = true
	// 	c.handshakeChan <- true
	// 	//收到握手消息了

	// }

	defer func() {
		if err := recover(); err != nil {
			log.Println("捕捉到异常", err)
			c.Close()
		}
	}()
	for {
		//这个data 要向 reqBody转换
		_, data, err := c.wsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误", err)
			break
		}
		//业务处理
		//1. data解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("数据格式不对", err)
			continue
		}
		//2. 前端消息是加密消息 需要解密
		secretKey, err := c.GetProperty("secretKey")
		if err != nil {
			log.Println("客户端secretKey获取失败", err)
			continue
		}
		if err == nil {
			key := secretKey.(string)
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("解密失败", err)
				// c.Handshake()
			} else {
				data = d
			}
		}
		//3. 转json 转为body
		body := &RspBody{}
		err = json.Unmarshal(data, body)
		fmt.Println("body数据：", body)

		if err != nil {
			log.Println("读取数据un json失败", err)
			continue
		} else {
			//握手 还是比的请求
			if body.Seq == 0 {
				if body.Name == HandshakeMsg {
					//获取密钥
					hs := &Handshake{}
					mapstructure.Decode(body.Msg, hs)
					if hs.Key != "" {
						c.SetProperty("secretKey", hs.Key)
					} else {
						c.RemoveProperty("secretKey")
					}
					c.handshake = true
					c.handshakeChan <- true
				} else {
					if c.onPush != nil {
						c.onPush(c, body)
					}
				}
			} else {
				c.syncCtxLock.RLock()
				ctx, ok := c.syncCtxMap[body.Seq]
				c.syncCtxLock.RUnlock()
				if ok {
					ctx.outChan <- body
				} else {
					log.Println("no seq syncCtx find")
				}
			}
			
		}
	}
	c.Close()
}

func (c ClientConn)waitHandShake() bool {
	//等待握手消息 成功
	//程序问题超时了  需要设置一个超时时间
	//context包维护了一个时间
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	select { 
	case _ = <- c.handshakeChan:
		log.Println("握手成功")
		return true
	case _ = <- ctx.Done():
		log.Println("握手超时....")
		return false
	}
}


func (c *ClientConn) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *ClientConn) GetProperty(key string)  (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("key 不存在")
	}
}
func (c *ClientConn) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
func (c *ClientConn) Addr() string {
	return c.wsConn.RemoteAddr().String()
}
func (c *ClientConn) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	// w.outChan <- rsp
	fmt.Println(rsp)
	c.write(rsp.Body)
}
func (c *ClientConn) write(body interface{}) error {
	log.Println("服务端写入数据：", body)
	data, err := json.Marshal(body)
	if err != nil {
		log.Println("写入数据json失败", err)
		return err
	}
	// secretKey, err := c.GetProperty("secretKey")
	// if err == nil {
	// 	key := secretKey.(string)
	// 	//加密
	// 	data, err = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	// 	if err != nil {
	// 		return err
	// 	}
	// } 
	//压缩
	if data, err := utils.Zip(data); err == nil {
		err = c.wsConn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil { 
			log.Println("写数据失败", err)
			return err
		}
	} else {
		log.Println("压缩数据失败", err)
		return err
	}
	return nil 
}


func (c *ClientConn) SetOnPush(hook func(conn *ClientConn,body *RspBody)) {
	c.onPush = hook
}
func (c *ClientConn) Send(name string, msg interface{}) *RspBody {
	//把请求发送给代理服务器 登陆服务器 等待返回
	c.Seq += 1
	seq := c.Seq
	sc := NewSyncCtx()
	c.syncCtxLock.Lock()
	c.syncCtxMap[seq] = sc
	c.syncCtxLock.Unlock()

	//构建一个req请求 
	req := &ReqBody{
		Seq: seq,
		Name: name,
		Msg: msg,
	}
	rsp := &RspBody{
		Name: name,
		Seq: seq,
		Code: constant.OK,
	}
	err := c.write(req)
	if err != nil {
		sc.cancel()
	} else {
		r := sc.wait()
		if r == nil {
			rsp.Code = constant.ProxyConnectError
		} else {
			rsp = r
		}
	} 
	c.syncCtxLock.Lock()
	delete(c.syncCtxMap, seq)
	c.syncCtxLock.Unlock()
	return rsp 
}