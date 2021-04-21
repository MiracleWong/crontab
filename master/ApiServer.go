package master

import (
	"encoding/json"
	"github.com/MiracleWong/crontab/master/common"
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
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)

	// 1. 解析POST表单
	err = req.ParseForm()
	if err != nil {
		goto ERR
	}

	// 2. 取出Job的值
	postJob = req.PostForm.Get("job")

	// 3. 反序列化
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

	// 4. 保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	// 5. 应答
	if bytes, err = common.BuildResponse(0,"success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 6. 错误应答
	if bytes, err = common.BuildResponse(-1, err.Error(), oldJob); err == nil {
		resp.Write(bytes)
	}
}

// 初始化任务
func InitApiServer() (err error) {
	// 定义变量
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)

	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	// 启动TCP监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	// 创建一个HTTP 的服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Second,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Second,
		Handler:      mux,
	}

	// 赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 启动服务端
	go httpServer.Serve(listener)

	return err
}
