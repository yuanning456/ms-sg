package logic

import (
	"log"
	"ms-sg/constant"
	"ms-sg/db"
	"ms-sg/server/common"
	"ms-sg/server/models"
	"ms-sg/server/web/model"
	"ms-sg/utils"
	"time"
)



var DefaultAccountLogic = &AccountLogic{}
type AccountLogic struct {

}

func (a AccountLogic) Register(req *model.RegisterReq) error {
	user := &models.User{}
	ok, err := db.Engine.Table(user).Where("username=?", req.Username).Get(user)
	if err != nil {
		log.Println("注册查询失败",err)
		return common.New(constant.DBError,"数据库异常")
	}
	if ok {
		return common.New(constant.UserExist,"该用户名称已经被注册")
	} else {
		user.Username	= req.Username
		user.Passcode 	= utils.RandSeq(6)
		user.Passwd  	= utils.Password(req.Password, user.Passcode)
		user.Ctime		= time.Now()
		user.Mtime		= time.Now()
		user.Hardware 	= req.Hardware
		_, err := db.Engine.Table(user).Insert(user)
		if err != nil {
			log.Println("数据库插入失败", err)
			return common.New(constant.DBError,"数据库异常")
		}
	}
	return nil
}