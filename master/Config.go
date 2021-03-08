package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// 程序配置
type Config struct {
	ApiPort int `json:"apiPort"`
	ApiReadTimeout int `json:"apiReadTimeout"`
	ApiWriteTimeout int `json:"apiWriteTimeout"`
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EtcdDialTimeout int `json:"etcdDialTimeout"`
}

var (
	// 单例
	G_config *Config
)

// 加载配置
func IninConfig(filename string) (err error) {
	var (
		content []byte
		conf Config
	)
	// 1. 把配置文件读出来
	if content, err = ioutil.ReadFile(filename); err != nil{
		return
	}
	
	// 2. 做JSON的反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	// 3. 赋值单例
	G_config = &conf

	fmt.Println(conf)
	return
}