package master

import (
	"context"
	"encoding/json"
	"github.com/MiracleWong/crontab/master/common"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
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
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	// 初始化配置
	config = clientv3.Config{
		//Endpoints:   []string{"127.0.0.1:2379"},
		//DialTimeout: 5 * time.Second,
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Microsecond,
	}
	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}
	// 获取kv 和 lease
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	// 赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

// 保存任务
func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	// 把任务保存到/cron/jobs/任务名 -> json
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj common.Job
	)
	// etcd 保存key
	jobKey = common.JOB_SAVE_DIR + job.Name

	// 任务信息json
	if jobValue, err = json.Marshal(job); err != nil {

		return
	}
	//println(jobValue)

	// 保存到etcd
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	if putResp.PrevKv != nil {
		//
		err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj)
		if err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

// 删除任务
func (jobMgr *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
	// 删除任务
	var (
		jobKey    string
		delResp   *clientv3.DeleteResponse
		oldJobObj common.Job
	)
	// etcd 删除的key
	jobKey = common.JOB_SAVE_DIR + name
	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	if len(delResp.PrevKvs) != 0 {
		err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj)
		if err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}
