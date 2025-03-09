package api

import (
	"BlessBot/common"
	"BlessBot/constant"
	"BlessBot/logs"
	"BlessBot/utils"
	"encoding/json"
	"go.uber.org/zap"
)

type NodeModel struct {
	NodeID      string
	Proxy       string
	Client      *utils.HttpUtil
	Token       string
	HardwareID  string
	IpAddress   string
	Remark      string
	IsConnected bool
}

func NewNodeModel(nodeId string, proxy string, token string, hardwareID string, remark string) *NodeModel {
	return &NodeModel{
		NodeID:      nodeId,
		Client:      utils.NewHttpClient(proxy),
		Token:       token,
		HardwareID:  hardwareID,
		Remark:      remark,
		IsConnected: false,
		IpAddress:   "",
	}
}

// RegisterNode 注册节点
func (n *NodeModel) RegisterNode() {
	registerUrl := constant.BaseURL + "/nodes/" + n.NodeID
	// 创建连接
	header := map[string]string{
		"Content-Type":          "application/json; charset=utf-8",
		"Authorization":         "Bearer " + n.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	body := map[string]interface{}{
		"extensionVersion": constant.ExtensionVersion,
		"hardwareId":       n.HardwareID,
		"ipAddress":        n.IpAddress,
		"hardwareInfo":     common.GenerateRandomHardwareInfo(),
	}
	resp, err := n.Client.Post(registerUrl, header, body, nil)
	if err != nil {
		logs.Log().Error("注册节点失败", zap.Error(err))
		return
	}
	// 处理响应
	if resp.StatusCode() == 200 {
		logs.Log().Info("注册节点成功 :", zap.String("UserMail", n.Remark), zap.String("NodeID", n.NodeID), zap.String("IP", n.IpAddress))
	} else {
		logs.Log().Error("注册节点失败 :", zap.String("NodeID", n.NodeID))
	}
}

// StartSession 启动会话
func (n *NodeModel) StartSession() error {
	startSessionUrl := constant.BaseURL + "/nodes/" + n.NodeID + "/start-session"
	header := map[string]string{
		"Authorization":         "Bearer " + n.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	resp, err := n.Client.Post(startSessionUrl, header, nil, nil)
	if err != nil {
		logs.Log().Error("启动会话失败", zap.Error(err))
		return err
	}
	logs.Log().Info("启动会话成功", zap.String("UserMail", n.Remark), zap.String("NodeID", n.NodeID), zap.String("接口返回 :", resp.String()))
	return nil
}

// StopSession 停止会话
func (n *NodeModel) StopSession() {
	stopSessionUrl := constant.BaseURL + "/nodes/" + n.NodeID + "/stop-session"

	header := map[string]string{
		"Authorization":         "Bearer " + n.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	resp, err := n.Client.Post(stopSessionUrl, header, nil, nil)
	if err != nil {
		logs.Log().Error("停止会话失败", zap.Error(err))
		return
	}
	logs.Log().Info("停止会话成功", zap.String("UserMail", n.Remark), zap.String("NodeID", n.NodeID), zap.String("接口返回 : ", resp.String()))
}

// PingNode 维持心跳
func (n *NodeModel) PingNode() error {
	pingUrl := constant.BaseURL + "/nodes/" + n.NodeID + "/ping"
	header := map[string]string{
		"Content-Type":          "application/json; charset=utf-8",
		"Authorization":         "Bearer " + n.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	body := map[string]interface{}{
		"isB7SConnected": n.IsConnected,
	}
	resp, err := n.Client.Post(pingUrl, header, body, nil)
	if err != nil {
		logs.Log().Error("心跳失败", zap.Error(err))
		return err
	}
	logs.Log().Info("心跳成功", zap.String("UserMail", n.Remark), zap.String("NodeID", n.NodeID), zap.String("Resp", resp.String()))
	return nil
}

// CheckNode 检查节点
func (n *NodeModel) CheckNode() error {
	var isConnected bool
	checkNodeUrl := constant.BaseURL + "/nodes/" + n.NodeID
	header := map[string]string{
		"Content-Type":          "application/json; charset=utf-8",
		"Authorization":         "Bearer " + n.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	resp, err := n.Client.Get(checkNodeUrl, nil, header, nil)
	if err != nil {
		logs.Log().Error("检查节点", zap.Error(err))
		return err
	}
	respMap := map[string]interface{}{}
	_ = json.Unmarshal(resp.Body(), &respMap)
	if respMap["isConnected"] == nil {
		isConnected = false
	} else {
		isConnected = respMap["isConnected"].(bool)
	}
	n.IsConnected = isConnected
	return nil
}

// HeathCheck 健康检查
func (n *NodeModel) HeathCheck() error {
	checkUrl := "https://gateway-run.bls.dev/health"
	header := map[string]string{
		"Authorization":         "Bearer " + n.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	_, err := n.Client.Get(checkUrl, nil, header, nil)
	if err != nil {
		logs.Log().Error("健康检查", zap.Error(err))
		return err
	}
	return nil
}

// CheckProxy 检查代理
func (n *NodeModel) CheckProxy() error {
	resp, err := n.Client.Get(constant.CheckProxyURL, nil, nil, nil)
	if err != nil {
		logs.Log().Error("代理异常：", zap.Error(err))
		return err
	}
	var result map[string]string
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		logs.Log().Error("解析代理返回失败：", zap.Error(err))
		return err
	}
	n.IpAddress = result["ip"]
	return nil
}
