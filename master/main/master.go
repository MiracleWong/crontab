package main

import (
	"flag"
	"fmt"
	"github.com/MiracleWong/crontab/master"
	"runtime"
	"time"
)

var (
	confFile string // 配置文件路径
)

// 初始化线程数量
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// 初始化命令行参数
func initArgs() {
	// master -config ./master.json
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
	flag.Parse() // 解析所有的命令行参数
}

func main() {
	var (
		err error
	)
	// 初始化线程
	initEnv()
	// 初始化命令行参数
	initArgs()
	// 加载配置
	if err = master.IninConfig(confFile); err != nil {
		//return err
		goto ERR
	}

	// 任务管理器
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	// 启动API HTTP 服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}

	// 正常退出
	return

ERR:
	fmt.Println(err)
}
