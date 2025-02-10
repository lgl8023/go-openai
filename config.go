package openai

import (
	"net/http"
	"regexp"
)

const (
	openaiAPIURLv1                 = "https://api.openai.com/v1"
	defaultEmptyMessagesLimit uint = 300

	azureAPIPrefix         = "openai"
	azureDeploymentsPrefix = "deployments"
)

type APIType string

const (
	APITypeOpenAI          APIType = "OPEN_AI"
	APITypeAzure           APIType = "AZURE"
	APITypeAzureAD         APIType = "AZURE_AD"
	APITypeCloudflareAzure APIType = "CLOUDFLARE_AZURE"
)

const AzureAPIKeyHeader = "api-key"

const defaultAssistantVersion = "v2" // upgrade to v2 to support vector store
// CustomHTTPClient 实现了 HTTPDoer 接口
type CustomHTTPClient struct {
	client http.Client
}

// NewCustomHTTPClient 创建一个新的 CustomHTTPClient 实例
func NewCustomHTTPClient(transport http.RoundTripper) *CustomHTTPClient {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &CustomHTTPClient{
		client: http.Client{
			Transport: transport,
		},
	}
}

// Do 实现 HTTPDoer 接口的 Do 方法
func (c *CustomHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// SetTransport 动态设置 Transport
func (c *CustomHTTPClient) SetTransport(transport http.RoundTripper) {
	if transport == nil {
		transport = http.DefaultTransport
	}
	c.client.Transport = transport
}

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken            string
	ProxyURL             string
	BaseURL              string
	OrgID                string
	APIType              APIType
	APIVersion           string // required when APIType is APITypeAzure or APITypeAzureAD
	AssistantVersion     string
	AzureModelMapperFunc func(model string) string // replace model to azure deployment name func
	HTTPClient           HTTPDoer

	EmptyMessagesLimit uint
}

func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		authToken:        authToken,
		BaseURL:          openaiAPIURLv1,
		APIType:          APITypeOpenAI,
		OrgID:            "",
		ProxyURL:         "",
		AssistantVersion: defaultAssistantVersion,

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func DefaultAzureConfig(apiKey, baseURL string) ClientConfig {
	return ClientConfig{
		authToken:  apiKey,
		BaseURL:    baseURL,
		OrgID:      "",
		APIType:    APITypeAzure,
		APIVersion: "2023-05-15",
		AzureModelMapperFunc: func(model string) string {
			return regexp.MustCompile(`[.:]`).ReplaceAllString(model, "")
		},

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func (ClientConfig) String() string {
	return "<OpenAI API ClientConfig>"
}

func (c ClientConfig) GetAzureDeploymentByModel(model string) string {
	if c.AzureModelMapperFunc != nil {
		return c.AzureModelMapperFunc(model)
	}

	return model
}
