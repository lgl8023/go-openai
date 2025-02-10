package openai

import (
	"net/http"
	"net/url"
)

// Http Client Transport Support
func (c *Client) SetProxyURL(proxyURL string) error {
	parseURL, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}
	// 创建一个自定义的 HTTPDoer 实例，不传 Transport 参数
	customHTTPClient := NewCustomHTTPClient(nil)
	// 创建 ClientConfig 实例，并将自定义的 HTTPDoer 赋值给 HTTPClient 字段
	c.config.HTTPClient = customHTTPClient
	newTransport := &http.Transport{
		Proxy: http.ProxyURL(parseURL),
	}
	customHTTPClient.SetTransport(newTransport)
	return nil
}
