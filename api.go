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
	c.config.HTTPClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(parseURL),
	}
	return nil
}
