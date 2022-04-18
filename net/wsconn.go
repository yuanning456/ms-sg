package net

type ReqBody struct {
	Seq     int64		`json:"seq"`
	Name 	string 		`json:"name"`
	Msg		interface{}	`json:"msg"`
	Proxy	string		`json:"proxy"`
}

type RspBody struct {
	Seq     int64		`json:"seq"`
	Name 	string 		`json:"name"`
	Code	int			`json:"code"`
	Msg		interface{}	`json:"msg"`
}

//类似http的一个请求标准
type WsMsgReq struct {
	Body	*ReqBody
	Conn	WSConn
}

type WsMsgRsp struct {
	Body*	RspBody
}
//理解为 request请求 请求会有参数 请求中放参数 取参数
type WSConn interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}

type Handshake struct {
	Key string `json:"key"`
}
type Heartbeat struct {
	CTime int64	`json:"ctime"`
	STime int64	`json:"stime"`
}
