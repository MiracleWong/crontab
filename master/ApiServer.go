package master

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

// 结构体：任务的HTTP接口
type ApiServer struct {
	httpServer *http.Server
}


// 一个单例对象
var (
	G_apiServer *ApiServer
)

// 保存任务的接口
func saveHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
}
// 初始化任务
func InitApiServer() (err error) {
	// 定义变量
	var (
		mux *http.ServeMux
		listener net.Listener
		httpServer *http.Server
	)

	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save",saveHandle)


	// 启动TCP监听
	if listener, err = net.Listen("tcp",":" + strconv.Itoa(G_config.ApiPort)); err != nil{
		return
	}

	// 创建一个HTTP 的服务
	httpServer = &http.Server{
		ReadTimeout: time.Duration(G_config.ApiReadTimeout) * time.Second,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Second,
		Handler: mux,
	}

	// 赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 启动服务端
	go httpServer.Serve(listener)

	return err
}
