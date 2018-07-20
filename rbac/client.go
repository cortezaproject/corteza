package rbac

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
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

func (c *Client) Users() *Users { return &Users{c} }

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
		// @todo: === remove this line ===
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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

func (c *Client) Debug(debug bool) *Client {
	c.isDebug = debug
	return c
}

func (c *Client) Get(url string) (*http.Response, error) {
	return c.Request("GET", url, nil)
}

func (c *Client) Post(url string, body interface{}) (*http.Response, error) {
	return c.Request("POST", url, body)
}

func (c *Client) Delete(url string) (*http.Response, error) {
	return c.Request("DELETE", url, nil)
}

func (c *Client) Request(method string, url string, body interface{}) (*http.Response, error) {
	link := strings.TrimRight(c.config.baseURL, "/") + "/" + strings.TrimLeft(url, "/")

	if c.isDebug {
		fmt.Println("RBAC >>>", method, link)
	}

	request := func() (*http.Request, error) {
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			return http.NewRequest(method, link, bytes.NewBuffer(b))
		}
		return http.NewRequest(method, link, nil)
	}

	req, err := request()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.config.auth)))
	// req.Header.Add("X-TENANT-ID", c.config.tenant)
	req.Header["X-TENANT-ID"] = []string{c.config.tenant}

	if c.isDebug {
		fmt.Println("RBAC >>> (request)")
		b, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Println("RBAC >>> Error:", err)
		} else {
			if b != nil {
				fmt.Println(strings.TrimSpace(string(b)))
			}
		}
		fmt.Println("---")
	}

	resp, err := c.Client.Do(req)
	if c.isDebug {
		fmt.Println("RBAC <<< (response)")
		if err != nil {
			fmt.Println("RBAC <<< Error:", err)
		} else {

			b, err := httputil.DumpResponse(resp, true)
			if err != nil {
				fmt.Println("RBAC <<< Error:", err)
			} else {
				if b != nil {
					fmt.Println(string(b))
				}
			}
		}
		fmt.Println("-----------------")
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}
