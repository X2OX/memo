package dandelion

import (
	"net/http"
	"net/url"
)

// HTTPClient is the type needed for the bot to perform HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	PostForm(url string, data url.Values) (*http.Response, error)
}

type ProxyClient struct {
	client *http.Client
}

func (p *ProxyClient) Do(req *http.Request) (*http.Response, error) {
	return p.client.Do(req)
}

func (p *ProxyClient) PostForm(url string, data url.Values) (*http.Response, error) {
	return p.client.PostForm(url, data)
}

func NewProxyClient(proxyUrl string) (HTTPClient, error) {
	u, err := url.Parse(proxyUrl)
	if err != nil {
		return nil, err
	}

	return &ProxyClient{client: &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(u),
		},
	}}, nil
}
