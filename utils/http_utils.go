package utils

import (
	"BlessBot/logs"
	"context"
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	defaultTimeout      = 100 * time.Second
	tlsHandshakeTimeout = 50 * time.Second
)

type HttpUtil struct {
	client *resty.Client
	mu     sync.Mutex
}

func NewHttpClient(proxyIP string) *HttpUtil {
	h := &HttpUtil{}
	if proxyIP == "" {
		h.newClient()
	} else {
		err := h.newProxyClient(proxyIP)
		if err != nil {
			logs.Log().Fatal("创建代理失败：", zap.Any(proxyIP, err))
		}
	}
	return h
}

// newProxyClient 创建代理客户端
func (h *HttpUtil) newProxyClient(proxyIp string) error {
	// SOCKS5 代理地址
	proxyURL, err := url.Parse(proxyIp)
	if err != nil {
		logs.Log().Error("解析代理地址失败:", zap.Error(err))
	}
	var transport *http.Transport
	if strings.HasPrefix(proxyIp, "socks") {
		password, _ := proxyURL.User.Password()
		auth := &proxy.Auth{
			User:     proxyURL.User.Username(),
			Password: password,
		}
		dialer, err := proxy.SOCKS5("tcp", strings.TrimPrefix(proxyURL.Host, "socks5://"), auth, proxy.Direct)
		if err != nil {
			logs.Log().Warn("创建代理连接失败:", zap.Any(proxyURL.Host, err))
			return err
		}
		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
	} else {
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			// 设置超时时间
			TLSHandshakeTimeout: tlsHandshakeTimeout, // TLS 握手超时
		}
	}
	client := resty.New().SetTimeout(defaultTimeout)
	client.SetTransport(transport)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	h.client = client
	return nil
}

// newClient 创建客户端
func (h *HttpUtil) newClient() {
	// 创建自定义的 Transport
	transport := &http.Transport{
		// 设置超时时间
		TLSHandshakeTimeout: tlsHandshakeTimeout, // TLS 握手超时
	}
	client := resty.New().SetTimeout(defaultTimeout)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTransport(transport)
	h.client = client
	logs.Log().Info("默认客户端创建成功")
}

func (h *HttpUtil) Get(url string, body interface{}, headers map[string]string, result interface{}) (*resty.Response, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	resp, err := h.client.R().SetHeaders(headers).SetBody(body).SetResult(result).Get(url)
	if err != nil {
		logs.Log().Error("Get 获取信息出错：", zap.Error(err))
		return nil, err
	}
	return resp, nil
}

func (h *HttpUtil) Post(url string, headers map[string]string, body interface{}, result interface{}) (*resty.Response, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	resp, err := h.client.R().SetBody(body).SetResult(result).SetHeaders(headers).Post(url)
	if err != nil {
		logs.Log().Error("Post 获取信息出错：", zap.Error(err))
		return nil, err
	}
	return resp, nil
}
