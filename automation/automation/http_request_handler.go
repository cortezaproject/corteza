package automation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type (
	httpRequestHandler struct {
		reg httpRequestHandlerRegistry
	}
)

func HttpRequestHandler(reg httpRequestHandlerRegistry) *httpRequestHandler {
	h := &httpRequestHandler{
		reg: reg,
	}

	h.register()
	return h
}

func (h httpRequestHandler) send(ctx context.Context, args *httpRequestSendArgs) (r *httpRequestSendResults, err error) {
	var (
		req *http.Request
		rsp *http.Response
	)

	r = &httpRequestSendResults{}

	req, err = h.makeRequest(ctx, args)
	if err != nil {
		return nil, err
	}

	rsp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	r.StatusCode = int64(rsp.StatusCode)
	r.Status = rsp.Status
	r.Headers = rsp.Header
	r.ContentLength = rsp.ContentLength
	r.ContentType = rsp.Header.Get("Content-Type")
	r.Body = rsp.Body

	return
}

func (h httpRequestHandler) makeRequest(ctx context.Context, args *httpRequestSendArgs) (req *http.Request, err error) {
	args.Method = strings.ToUpper(args.Method)

	if args.Method == "" && (len(args.Form) > 0 || args.Body != nil) {
		// when no method is set and form or body are passed
		args.Method = http.MethodPost
	}

	err = func() error {
		switch args.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
		default:
			return nil
		}

		// @todo handle (multiple) file upload as well

		if len(args.Form) > 0 {
			if args.Body != nil {
				return fmt.Errorf("can not not use form and body parameters at the same time")
			}

			if !args.hasHeaderContentType {
				args.HeaderContentType = "application/x-www-form-urlencoded"
			}

			args.bodyStream = &bytes.Buffer{}
			if _, err = args.bodyStream.(*bytes.Buffer).WriteString(args.Form.Encode()); err != nil {
				return err
			}
			return nil
		}

		if args.Body != nil {
			return nil
		}

		if args.hasBody && args.bodyStream == nil {
			if args.bodyString != "" {
				args.bodyStream = strings.NewReader(args.bodyString)
			} else {
				switch raw := args.bodyRaw.(type) {
				case string:
					args.bodyStream = strings.NewReader(raw)
				case []byte:
					args.bodyStream = bytes.NewReader(raw)
				case io.Reader:
					args.bodyStream = raw
				default:
					args.bodyStream = &bytes.Buffer{}
					return json.NewEncoder(args.bodyStream.(*bytes.Buffer)).Encode(args.bodyRaw)
				}
			}
		}

		return nil
	}()
	if err != nil {
		return nil, err
	}

	if args.Timeout > 0 {
		var tfn context.CancelFunc
		ctx, tfn = context.WithTimeout(ctx, args.Timeout)
		defer tfn()
	}

	if args.hasParams {
		purl, err := url.Parse(args.Url)
		if err != nil {
			return nil, err
		}

		purl.RawQuery = args.Params.Encode()
		args.Url = purl.String()
	}

	req, err = http.NewRequestWithContext(ctx, args.Method, args.Url, args.bodyStream)
	if err != nil {
		return nil, err
	}

	if args.Headers == nil {
		args.Headers = make(http.Header)
	}

	if args.HeaderUserAgent == "" {
		args.HeaderUserAgent = "Corteza-Automation-Client/" + version.Version
	}

	args.Headers.Set("User-Agent", args.HeaderUserAgent)

	switch {
	case len(args.HeaderAuthBearer) > 0:
		args.Headers.Add("Authorization", "Bearer "+args.HeaderAuthBearer)
	case len(args.HeaderAuthPassword+args.HeaderAuthPassword) > 0:
		req.SetBasicAuth(
			args.HeaderAuthPassword,
			args.HeaderAuthPassword,
		)
	}

	if len(args.HeaderContentType) > 0 {
		args.Headers.Add("Content-Type", args.HeaderContentType)
	}

	req.Header = args.Headers

	return
}
