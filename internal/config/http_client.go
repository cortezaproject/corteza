package config

import (
	"github.com/namsral/flag"
)

type (
	HTTPClient struct {
		BaseURL string
		Timeout int
	}
)

var httpClient *HTTPClient

func (c *HTTPClient) Validate() error {
	return nil
}

func (*HTTPClient) Init(prefix ...string) *HTTPClient {
	if httpClient != nil {
		return httpClient
	}
	httpClient = new(HTTPClient)
	flag.StringVar(&httpClient.BaseURL, "http-client-base-url", "", "HTTP Client Base URL")
	flag.IntVar(&httpClient.Timeout, "http-client-timeout", 5, "HTTP Client request timeout (seconds)")
	return httpClient
}
