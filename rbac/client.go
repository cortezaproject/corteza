package rbac

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type (
	Client struct {
		Transport *http.Transport
		Client    *http.Client

		isDebug bool
		config  configuration
	}
)

func New() (*Client, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	timeout := time.Duration(config.timeout) * time.Second

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: timeout,
		}).Dial,
		TLSHandshakeTimeout: timeout,
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Client{
		Transport: transport,
		Client:    client,
		isDebug:   false,
		config:    config,
	}, nil
}

func (c *Client) Get(url string) (*http.Response, error) {
	return c.Request("GET", url)
}

func (c *Client) Request(method string, url string) (*http.Response, error) {
	link := strings.TrimRight(c.config.baseURL, "/") + "/" + strings.TrimLeft(url, "/")

	if c.isDebug {
		fmt.Println("RBAC >>> ", link)
	}

	req, err := http.NewRequest(method, link, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.config.auth)))
	req.Header.Add("X-TENANT-ID", c.config.tenant)

	resp, err := c.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}
