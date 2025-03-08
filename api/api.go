package api

import (
	"BlessBot/common"
	"BlessBot/constant"
	"BlessBot/logs"
	"BlessBot/model"
	"BlessBot/utils"
	"encoding/json"
	"go.uber.org/zap"
)

// RegisterNode 注册节点
func RegisterNode(node *model.RegisterNoe) {
	registerUrl := constant.BaseURL + "/nodes/" + node.NodeID
	// 创建连接
	client := utils.NewHttpClient(node.Proxy)
	header := map[string]string{
		"Content-Type":          "application/json; charset=utf-8",
		"Authorization":         "Bearer " + node.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	body := map[string]interface{}{
		"extensionVersion": constant.ExtensionVersion,
		"hardwareId":       node.HardwareID,
		"ipAddress":        node.IpAddress,
		"hardwareInfo":     common.GenerateRandomHardwareInfo(),
	}
	resp, err := client.Post(registerUrl, header, body, nil)
	if err != nil {
		logs.Log().Error("注册节点失败", zap.Error(err))
		return
	}
	// 处理响应
	if resp.StatusCode() == 200 {
		logs.Log().Info("注册节点成功 :", zap.String("UserMail", node.Remark), zap.String("NodeID", node.NodeID), zap.String("IP", node.IpAddress))
	} else {
		logs.Log().Error("注册节点失败 :", zap.String("NodeID", node.NodeID))
	}
}

// StartSession 启动会话
func StartSession(node *model.RegisterNoe) {
	startSessionUrl := constant.BaseURL + "/nodes/" + node.NodeID + "/start-session"
	// 创建连接
	client := utils.NewHttpClient(node.Proxy)

	header := map[string]string{
		"Authorization":         "Bearer " + node.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	resp, err := client.Post(startSessionUrl, header, nil, nil)
	if err != nil {
		logs.Log().Error("启动会话失败", zap.Error(err))
		return
	}
	logs.Log().Info("启动会话成功", zap.String("UserMail", node.Remark), zap.String("NodeID", node.NodeID), zap.String("接口返回 :", resp.String()))
}

// StopSession 停止会话
func StopSession(node *model.RegisterNoe) {
	stopSessionUrl := constant.BaseURL + "/nodes/" + node.NodeID + "/stop-session"
	// 创建连接
	client := utils.NewHttpClient(node.Proxy)

	header := map[string]string{
		"Authorization":         "Bearer " + node.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	resp, err := client.Post(stopSessionUrl, header, nil, nil)
	if err != nil {
		logs.Log().Error("停止会话失败", zap.Error(err))
		return
	}
	logs.Log().Info("停止会话成功", zap.String("UserMail", node.Remark), zap.String("NodeID", node.NodeID), zap.String("接口返回 : ", resp.String()))
}

// PingNode 维持心跳
func PingNode(node *model.RegisterNoe, isB7SConnected bool) error {
	pingUrl := constant.BaseURL + "/nodes/" + node.NodeID + "/ping"
	// 创建连接
	client := utils.NewHttpClient(node.Proxy)
	header := map[string]string{
		"Content-Type":          "application/json; charset=utf-8",
		"Authorization":         "Bearer " + node.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	body := map[string]interface{}{
		"isB7SConnected": isB7SConnected,
	}
	resp, err := client.Post(pingUrl, header, body, nil)
	if err != nil {
		logs.Log().Error("心跳失败", zap.Error(err))
		return err
	}
	logs.Log().Info("心跳成功", zap.String("UserMail", node.Remark), zap.String("NodeID", node.NodeID), zap.String("Resp", resp.String()))
	return nil
}

// CheckNode 检查节点
func CheckNode(node *model.RegisterNoe) (bool, error) {
	var isConnected bool
	checkNodeUrl := constant.BaseURL + "/nodes/" + node.NodeID
	// 创建连接
	client := utils.NewHttpClient(node.Proxy)
	header := map[string]string{
		"Content-Type":          "application/json; charset=utf-8",
		"Authorization":         "Bearer " + node.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	resp, err := client.Get(checkNodeUrl, nil, header, nil)
	if err != nil {
		logs.Log().Error("检查节点", zap.Error(err))
		return false, err
	}
	respMap := map[string]interface{}{}
	_ = json.Unmarshal(resp.Body(), &respMap)
	if respMap["isConnected"] == nil {
		isConnected = false
	} else {
		isConnected = respMap["isConnected"].(bool)
	}
	return isConnected, nil
}

// HeathCheck 健康检查
func HeathCheck(node *model.RegisterNoe) {
	checkUrl := "https://gateway-run.bls.dev/health"
	client := utils.NewHttpClient(node.Proxy)
	header := map[string]string{
		"Authorization":         "Bearer " + node.Token,
		"User-Agent":            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Origin":                "chrome-extension://pljbjcehnhcnofmkdbjolghdcjnmekia",
		"X-Extension-Version":   constant.ExtensionVersion,
		"X-Extension-Signature": constant.ExtensionSignature,
	}
	_, err := client.Get(checkUrl, nil, header, nil)
	if err != nil {
		logs.Log().Error("健康检查", zap.Error(err))
	}
}

// CheckProxy 检查代理
func CheckProxy(proxy string) map[string]string {
	if proxy == "" {
		logs.Log().Error("代理为空")
		return nil
	}
	client := utils.NewHttpClient(proxy)
	resp, err := client.Get(constant.CheckProxyURL, nil, nil, nil)
	if err != nil {
		logs.Log().Error("代理异常：", zap.Error(err))
		return nil
	}
	var result map[string]string
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		logs.Log().Error("解析代理返回失败：", zap.Error(err))
		return nil
	}
	return result
}
