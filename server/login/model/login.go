package model

import "time"

const (
	Login = iota
	LoginOut
)

type LoginHistory struct {
	Id       int       `xorm:"id pk autoincr"`
	UId      int       `xorm:"uid"`
	CTime    time.Time `xorm:"ctime"`
	Ip       string    `xorm:"ip"`
	State    int8      `xorm:"state"`
	Hardware string    `xorm:"hardware"`
}

type LoginLast struct {
	Id         int       `xorm:"id pk autoincr"`
	UId        int       `xorm:"uid"`
	LoginTime  time.Time `xorm:"login_time"`
	LogoutTime time.Time `xorm:"logout_time"`
	Ip         string    `xorm:"ip"`
	Session    string    `xorm:"session"`
	IsLogout   int8      `xorm:"is_logout"`
	Hardware   string    `xorm:"hardware"`
}

func (*LoginHistory) TableName() string {
	return "login_history"
}

func (*LoginLast) TableName() string {
	return "login_last" 
}