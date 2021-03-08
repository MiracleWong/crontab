package master

import (
	"fmt"
	clientv3 "github.com/coreos/etcd/client/v3"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

// 单例
var (
	G_jobMgr *JobMgr
)


// 初始化设备管理
func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)
	config = clientv3.Config{
		//Endpoints:   []string{"127.0.0.1:2379"},
		//DialTimeout: 5 * time.Second,
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Microsecond,
	}
	if client, err = clientv3.New(config); err != nil{
		fmt.Println(err)
		return
	}
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)


	G_jobMgr = &JobMgr{
		client:client,
		kv:kv,
		lease:lease,
	}
	return
}