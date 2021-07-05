package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"zinx/ziface"
)

// 存储zinx框架的全局参数，这些参数可以通过zinx.json 由用户来配置

type GlobalObj struct {
	// Server 相关配置
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器监听的端口号
	Name      string         // 当前服务器的名称

	// Zinx 相关配置
	Version        string // 当前Zinx的版本号
	MaxConn        int    // 当前服务器主机允许的最大连接数
	MaxPackageSize uint32 // 当前Zinx框架数据包的最大值
	WorkerPoolSize uint32 // 当前业务workerPoll中的Goroutine数量
	MaxWorkerTaskLen uint32  // 每个worker对应的消息队列的任务的数量最大值
}

var GlobalObject *GlobalObj

func (c *GlobalObj) LoadConfig() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		fmt.Println("ioutil read file error: ", err)
		return 
	}

	// 将json数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 用来初始化全局对象
func init() {
	// 如果配置文件中没有相关内容 提供一个默认值
	GlobalObject = &GlobalObj{
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024,
	}

	// 从 conf/zin:x.json 中加载用户自定义的参数
	GlobalObject.LoadConfig()
}
