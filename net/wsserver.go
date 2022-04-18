package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"ms-sg/utils"
	"sync"

	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
)

//websocket服务
type wsServer struct {
	wsConn *websocket.Conn
	router *Router
	outChan chan *WsMsgRsp //写队列
	Seq int64
	property map[string]interface{}
	propertyLock sync.RWMutex
}


func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn: wsConn,
		outChan: make(chan *WsMsgRsp, 1000),
		property:  make(map[string]interface{}),
		Seq: 0,
	}
}

func (w *wsServer) Router(router *Router) {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}
func (w *wsServer) GetProperty(key string)  (value interface{}, err error) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	if value, ok := w.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("key 不存在")
	}
}
func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}
func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}
func (w *wsServer) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	w.outChan <- rsp
}
//
func (w *wsServer) Start() {
	//启动读写数据处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) readMsgLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			w.Close()
		}
	}()

	for {
		//这个data 要向 reqBody转换
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误", err)
			break
		}
		// w.outChan <- 
		log.Println("husju", data)
		//业务处理
		//1. data解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("数据格式不对", err)
			continue
		}
		
		//2. 前端消息是加密消息 需要解密
		secretKey, err := w.GetProperty("secretKey")
		if err != nil {
			log.Println("secretKey获取失败", err)
			continue
		}
		if err == nil {
			key := secretKey.(string)
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("解密失败", err)
				w.Handshake()
				// w.
			} else {
				data = d
			}
		}
		//3. 转json 转为body
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("读取数据un json失败", err)
			continue
		} else {
			//获取到前端传递的数据 处理业务
			req := &WsMsgReq{Conn: w, Body: body}
			rsq := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}
			w.router.Run(req, rsq)
			w.outChan <- rsq
		}
	}
	w.Close()
}
func(w *wsServer) writeMsgLoop() {
	//读消息 处理 回复消息
	for {
		select {
		case msg := <- w.outChan:
			fmt.Println(msg)
			w.Write(msg)
		
		}
	}
}

func (w *wsServer) Close() {
	w.wsConn.Close()
}

func (w *wsServer) Write(msg *WsMsgRsp) {
	data, err := json.Marshal(msg.Body)
	if err != nil {
		log.Println("写入数据json失败", err)
		
	}
	secretKey, err := w.GetProperty("secretKey")
	if err == nil {
		key := secretKey.(string)
		//加密
		data, _ := utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		//压缩
		if data, err := utils.Zip(data); err == nil {
			w.wsConn.WriteMessage(websocket.BinaryMessage, data)
		}
		d, _ := json.Marshal(msg)
		log.Println("服务端写入数据：", d)
	}
}

const HandshakeMsg = "handshare"
// 握手协议 当游戏客户端发送请求的时候 
// 后端会发送加密key给客户端
// 客户端在发送数据的时候就会使用此key进行加密处理

func (w *wsServer) Handshake() {
	secretKey := ""
	key, err := w.GetProperty("secretKey")
	if err == nil {
		secretKey = key.(string)
	} else {
		secretKey = utils.RandSeq(16)
	}
	handshake := &Handshake{Key: secretKey}
	body := &RspBody{Name: HandshakeMsg, Msg: handshake}
	if data, err := json.Marshal(body); err == nil {
		if secretKey != "" {
			w.SetProperty("secretKey", secretKey)
		} else {
			w.RemoveProperty("secretKey")
		}
		if data, err := utils.Zip(data); err == nil {
			w.wsConn.WriteMessage(websocket.BinaryMessage, data)
		}
	}

}





