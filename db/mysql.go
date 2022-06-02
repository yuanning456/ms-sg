package db

import (
	"fmt"
	"log"
	"ms-sg/config"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)


var Engine *xorm.Engine
func Init() {
	// import "gorm.io/driver/mysql"
	// refer: https://gorm.io/docs/connecting_to_the_database.html#MySQL
	
	mysqlConfig, err := config.File.GetSection("mysql")
	if err != nil {
		fmt.Println("数据库配置错误", err)
		panic(err)
	}
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", 
		mysqlConfig["user"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"],
		mysqlConfig["charset"],
		)

	Engine, err = xorm.NewEngine("mysql", conn)
	if err != nil {
		fmt.Println("数据库连接失败", err)
		panic(err)
	}
	err = Engine.Ping()
	if err != nil {
		fmt.Println("数据库ping不通", err)
		panic(err)
	}

	//设置一些配置
	maxIdle := config.File.MustInt("mysql", "max_idle", 2)
	maxConn := config.File.MustInt("mysql", "max_conn", 2)
	Engine.SetMaxIdleConns(maxIdle)
	Engine.SetMaxOpenConns(maxConn)
	Engine.ShowSQL(true)
	log.Println("数据库初始化完成...")
} 