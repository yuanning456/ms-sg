package config

import (
	"errors"
	"log"
	"os"

	"github.com/Unknwon/goconfig"
)

const configFile = "/conf/conf.ini"
var File *goconfig.ConfigFile
func init() {

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + configFile
	// os.Args 获取执行命令 里面的参数
	// len := len(os.Args)
	// if len > 1 {
	// 	dir := os.Args[1]
	// 	if dir != "" {
	// 		configPath = dir + configFile
	// 	}
	// }
	if !fileExists(configPath) {
		panic(errors.New("配置文件不存在"))
	}
	//文件系统的读取
	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("读取文件出错：", err)
		panic(err)
	}

}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
