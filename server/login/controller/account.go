package controller

import (
	"fmt"
	"log"
	"ms-sg/constant"
	"ms-sg/db"
	"ms-sg/net"
	"ms-sg/server/login/model"
	"ms-sg/server/login/proto"
	"ms-sg/utils"
	"time"

	"github.com/mitchellh/mapstructure"
)

var DefaultAccount = &Account{}
type Account struct {

}
func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRouter("login", a.login)
}

/**
	1.用户名 密码 硬件信息 
	2.根据用户名 查询用户信息
	3.进行密码对比 成功登陆
	4.保存用户的登陆记录 
	5.保存用户最后一次的登陆信息
	6.客户端需要一个session jwt生成一个加密字符串
	7.客户端登陆的时候根据这个加密字符串判断这个用户是否合法  
**/
func (a *Account) login(req *net.WsMsgReq, rsq *net.WsMsgRsp) {
	fmt.Println("login start")
	loginRes := &proto.LoginRsq{}
	loginReq := &proto.LoginReq{}
	mapstructure.Decode(req.Body.Msg, loginReq)
	username := loginReq.UserName
	fmt.Println("loginReq:", loginReq)

	user := &model.User{}
	db.Engine.TableName(user.TableName())
	ok, err := db.Engine.Where(fmt.Sprintf("username = '%s'", username)).Get(user)
	if err != nil {
		log.Println("数据库查询出错", err)
		return
	} 
	if !ok {
		rsq.Body.Code = constant.UserNotExist
		return
	}

	pwd := utils.Password(loginReq.Password, user.Passcode)
	fmt.Println("pwd:", pwd)

	//暂时密码不正确
	// if user.Passwd != pwd {
	// 	rsq.Body.Code = constant.PwdIncorrect
	// 	return
	// }

	//jwt A.B.C 三部分 A定义加密算法 B定义放入的数据 C部分 根据秘钥+A和B生成加密字符串
	token, err := utils.Award(user.UId)
	if err != nil {
		log.Println("生成token失败")
		return
	}
	loginRes.Session = token
	loginRes.UId = user.UId
	loginRes.UserName = user.Username
	loginRes.Password = ""
	 
	rsq.Body.Msg = loginRes
	rsq.Body.Code = constant.OK

	
	//保存用户登陆信息
	lh := &model.LoginHistory{
		UId: user.UId, 
		CTime: time.Now(), 
		Ip: loginReq.Ip,
		Hardware: loginReq.Hardware, 
		State: model.Login,
	}
	db.Engine.Table(lh).Insert(lh)

	//最后一次登录的状态记录
	ll := &model.LoginLast{}
	ok ,_ = db.Engine.Table(ll).Where("uid=?", user.UId).Get(ll)
	if ok {
		//有数据 更新
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		db.Engine.Table(ll).Update(ll)
	}else{
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		ll.UId = user.UId
		_, err := db.Engine.Table(ll).Insert(ll)
		if err != nil {
			log.Println(err)
		}
	}
	//缓存一下 此用户和当前的ws连接
	// wsmgr := &net.WsMgr{}
	// wsmgr.UserLogin(req.Conn, user.UId, token)
	net.Mgr.UserLogin(req.Conn, user.UId, token)
	// net.Mgr.UserLogin(req.Conn,user.UId,token)
}


