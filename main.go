package main

import (
	"BlessBot/api"
	"BlessBot/config"
	"BlessBot/logs"
	"BlessBot/model"
	"go.uber.org/zap"
	"sync"
	"time"
)

func init() {
	config.Init()
}

func main() {

	// 获取token
	accounts := config.GetConfig()
	var models []model.RegisterNoe
	logs.Log().Info("封装参数中....")
	for _, user := range accounts.Users {
		for _, node := range user.Nodes {
			var registerNoe model.RegisterNoe
			proxy := api.CheckProxy(node.Proxy)
			registerNoe.NodeID = node.NodeID
			registerNoe.HardwareID = node.HardwareID
			registerNoe.Proxy = node.Proxy
			registerNoe.Token = user.UserToken
			registerNoe.Remark = user.Remark
			if proxy["ip"] == "" {
				registerNoe.IpAddress = "0.0.0.0"
			} else {
				registerNoe.IpAddress = proxy["ip"]
			}
			models = append(models, registerNoe)
		}
	}
	logs.Log().Info("封装参数完成....")
	for _, registerNode := range models {
		go processNode(registerNode)
	}
	select {}
}

// processNode 处理节点
func processNode(modelNode model.RegisterNoe) {
	var lock sync.Mutex
	var isConnected bool
	var err error
	isConnected, err = api.CheckNode(&modelNode)
	if err != nil {
		logs.Log().Error("CheckNode", zap.String("NodeId", modelNode.NodeID), zap.Error(err))
		return
	}
	if !isConnected {
		// 注册
		api.RegisterNode(&modelNode)
		// ping
		api.PingNode(&modelNode, isConnected)
		// 开启会话
		api.StartSession(&modelNode)
	}

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			lock.Lock()
			isConnected, err = api.CheckNode(&modelNode)
			if err != nil {
				logs.Log().Error("CheckNode", zap.String("NodeId", modelNode.NodeID), zap.Error(err))
				continue
			}
			lock.Unlock()
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Minute)
			api.HeathCheck(&modelNode)
			lock.Lock()
			isConnected, err = api.CheckNode(&modelNode)
			if err != nil {
				logs.Log().Error("CheckNode", zap.String("NodeId", modelNode.NodeID), zap.Error(err))
				continue
			}
			lock.Unlock()
			if !isConnected {
				// 注册
				api.RegisterNode(&modelNode)
				// ping
				api.PingNode(&modelNode, isConnected)
				// 开启会话
				api.StartSession(&modelNode)
			}
		}
	}()
}
