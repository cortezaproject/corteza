package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type (
	Config struct {
		BaseURL string
		Timeout int
	}

	Client struct {
		Transport *http.Transport
		Client    *http.Client

		debugLevel DebugLevel
		config     *Config
	}

	Request http.Request

	DebugLevel string
)

const (
	INFO DebugLevel = "info"
	FULL DebugLevel = "full"
)

func New(flags *Config) (*Client, error) {
	timeout := time.Duration(flags.Timeout) * time.Second

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: timeout,
		}).Dial,
		TLSHandshakeTimeout: timeout,
	}

	// @todo migrate to http.DefaultClient & http.DefaultTransport, see internal/http.SetupDefaults
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
		config:    flags,
	}, nil
}

func (c *Client) Debug(level DebugLevel) *Client {
	c.debugLevel = level
	return c
}

func (c *Client) Get(url string) (*http.Request, error) {
	return c.Request("GET", url, nil)
}

func (c *Client) Post(url string, body interface{}) (*http.Request, error) {
	return c.Request("POST", url, body)
}

func (c *Client) Patch(url string, body interface{}) (*http.Request, error) {
	return c.Request("PATCH", url, body)
}

func (c *Client) Delete(url string) (*http.Request, error) {
	return c.Request("DELETE", url, nil)
}

func (c *Client) Request(method, url string, body interface{}) (*http.Request, error) {
	if c.config.BaseURL != "" {
		url = strings.TrimRight(c.config.BaseURL, "/") + "/" + strings.TrimLeft(url, "/")
	}

	if c.debugLevel == "info" {
		fmt.Println("HTTP >>>", method, url)
	}

	request := func() (*http.Request, error) {
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			return http.NewRequest(method, url, bytes.NewBuffer(b))
		}
		return http.NewRequest(method, url, nil)
	}

	req, err := request()
	if err != nil {
		return nil, errors.Wrap(err, "creating request failed")
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.debugLevel == FULL {
		fmt.Println("HTTP >>> (request)")
		b, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Println("HTTP >>> Error:", err)
		} else {
			if b != nil {
				fmt.Println(strings.TrimSpace(string(b)))
			}
		}
		fmt.Println("---")
	}

	resp, err := c.Client.Do(req)
	if c.debugLevel == FULL {
		fmt.Println("HTTP <<< (response)")
		if err != nil {
			fmt.Println("HTTP <<< Error:", err)
		} else {

			b, err := httputil.DumpResponse(resp, true)
			if err != nil {
				fmt.Println("HTTP <<< Error:", err)
			} else {
				if b != nil {
					fmt.Println(string(b))
				}
			}
		}
		fmt.Println("-----------------")
	}
	if err != nil {
		if c.debugLevel == INFO {
			fmt.Println("HTTP <<< Response error", err)
		}
		return nil, errors.Wrap(err, "request failed")
	}
	if c.debugLevel == INFO {
		fmt.Println("HTTP <<< Response", resp.StatusCode)
	}
	return resp, nil
}

func ToError(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if body == nil || err != nil {
		return errors.Errorf("unexpected response (%d, %s)", resp.StatusCode, err)
	}
	return errors.New(string(body))
}
