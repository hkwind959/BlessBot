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
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)
	for _, user := range accounts.Users {
		for _, node := range user.Nodes {
			wg.Add(1)
			go func(u model.User, n model.UserNode) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				nodeModel := api.NewNodeModel(n.NodeID, n.Proxy, u.UserToken, n.HardwareID, u.Remark)
				err2 := nodeModel.CheckProxy()
				if err2 != nil {
					logs.Log().Error("CheckProxy", zap.String("NodeId", n.NodeID), zap.Error(err2))
					return
				}
				err := nodeModel.CheckNode()
				if err != nil {
					logs.Log().Error("CheckNode", zap.String("NodeId", n.NodeID), zap.Error(err))
					return
				}
				if !nodeModel.IsConnected {
					// 注册
					nodeModel.RegisterNode()
					// ping
					err := nodeModel.PingNode()
					if err != nil {
						logs.Log().Error("PingNode", zap.String("NodeId", n.NodeID), zap.Error(err))
						return
					}
					// 开启会话
					err = nodeModel.StartSession()
					if err != nil {
						logs.Log().Error("StartSession", zap.String("NodeId", n.NodeID), zap.Error(err))
						return
					}
				}
				go processCheckNode(nodeModel)
				go processNode(nodeModel)
			}(user, node)
		}
	}
	wg.Wait()
	select {}
}

// processNode 处理节点
func processNode(nodeModel *api.NodeModel) {
	for {
		time.Sleep(10 * time.Minute)
		err := nodeModel.HeathCheck()
		if err != nil {
			logs.Log().Error("CheckNode", zap.String("NodeId", nodeModel.NodeID), zap.Error(err))
			continue
		}
		err = nodeModel.CheckNode()
		if err != nil {
			logs.Log().Error("CheckNode", zap.String("NodeId", nodeModel.NodeID), zap.Error(err))
			continue
		}
		if !nodeModel.IsConnected {
			// 注册
			nodeModel.RegisterNode()
			// ping
			err := nodeModel.PingNode()
			if err != nil {
				logs.Log().Error("PingNode", zap.String("NodeId", nodeModel.NodeID), zap.Error(err))
				continue
			}
			// 开启会话
			err = nodeModel.StartSession()
			if err != nil {
				logs.Log().Error("StartSession", zap.String("NodeId", nodeModel.NodeID), zap.Error(err))
				continue
			}
		}
	}
}

// processCheckNode 处理检测节点
func processCheckNode(nodeModel *api.NodeModel) {
	for {
		time.Sleep(1 * time.Minute)
		err := nodeModel.CheckNode()
		if err != nil {
			logs.Log().Error("CheckNode", zap.String("NodeId", nodeModel.NodeID), zap.Error(err))
			continue
		}
	}
}
