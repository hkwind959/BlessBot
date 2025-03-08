package main

import (
	"BlessBot/api"
	"BlessBot/config"
	"BlessBot/logs"
	"BlessBot/model"
	"go.uber.org/zap"
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
	processNode(models)
	select {}
}

// processNode 处理节点
func processNode(models []model.RegisterNoe) {
	for _, registerNoe := range models {
		isConnected, err := api.CheckNode(&registerNoe)
		if err != nil {
			logs.Log().Error("CheckNode", zap.String("NodeId", registerNoe.NodeID), zap.Error(err))
			continue
		}
		if !isConnected {
			// 注册
			api.RegisterNode(&registerNoe)
			// 开启会话
			api.StartSession(&registerNoe)
		}
		err = api.PingNode(&registerNoe, isConnected)
		if err != nil {
			logs.Log().Error("PingNode", zap.String("NodeId", registerNoe.NodeID), zap.Error(err))
			continue
		}

		go func() {
			for {
				isConnected, err := api.CheckNode(&registerNoe)
				if err != nil {
					logs.Log().Error("CheckNode", zap.String("NodeId", registerNoe.NodeID), zap.Error(err))
					continue
				}
				_ = api.PingNode(&registerNoe, isConnected)
				time.Sleep(10 * time.Minute)
			}
		}()

		go func() {
			for {
				api.HeathCheck(&registerNoe)
				isConnected, err := api.CheckNode(&registerNoe)
				if err != nil {
					logs.Log().Error("CheckNode", zap.String("NodeId", registerNoe.NodeID), zap.Error(err))
					continue
				}
				if !isConnected {
					// 停止
					api.StopSession(&registerNoe)
					// 注册
					api.RegisterNode(&registerNoe)
					// 开启会话
					api.StartSession(&registerNoe)
				}
				time.Sleep(5 * time.Minute)
			}
		}()
		// 发送ping
		go func() {
			for {
				err = api.PingNode(&registerNoe, isConnected)
				if err != nil {
					logs.Log().Error("PingNode", zap.String("NodeId", registerNoe.NodeID), zap.Error(err))
				}
				time.Sleep(1 * time.Minute)
			}
		}()
	}
}
