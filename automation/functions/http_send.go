package functions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	. "github.com/cortezaproject/corteza-server/automation/types/fn"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	stdHttpSendParameters = []*Param{
		NewParam("url", String, Required),
		NewParam("header", KV),
		NewParam("headerAuthBearer", String),
		NewParam("headerAuthUsername", String),
		NewParam("headerAuthPassword", String),
		NewParam("headerContentType", String),
		NewParam("timeout", Duration),
	}

	stdHttpPayloadParameters = []*Param{
		NewParam("form", KV),
		NewParam("body", String, Reader, Any),
	}

	stdHttpSendResults = []*Param{
		NewParam("status", String),
		NewParam("statusCode", Int),
		NewParam("header", KV),
		NewParam("contentLength", Int),
		NewParam("body", String),
		NewParam("parsedJson", Any),
	}
)

func makeHttpRequest(ctx context.Context, in wfexec.Variables) (req *http.Request, err error) {
	var (
		body   io.Reader
		header = make(http.Header)
		method = strings.ToUpper(in.String("method"))
	)

	if method == "" && in.Any("form", "body") {
		// when no method is set and form or body are passed
		method = "POST"
	}

	err = func() error {
		switch method {
		case "POST", "PUT", "PATCH":
		default:
			return nil
		}

		// @todo handle (multiple) file upload as well

		if in.Has("form") {
			if in.Has("body") {
				return fmt.Errorf("can not not use form and body parameters at the same time")
			}

			var form url.Values
			form, err = cast.ToStringMapStringSliceE(in["form"])
			if err != nil {
				return fmt.Errorf("failed to resolve form values: %w", err)
			}

			header.Add("Content-Type", "application/x-www-form-urlencoded")
			payload := &bytes.Buffer{}
			if _, err = payload.WriteString(form.Encode()); err != nil {
				return err
			}

			body = payload
			return nil
		}

		if !in.Has("body") {
			return nil
		}

		switch payload := in["body"].(type) {
		case string:
			body = strings.NewReader(payload)
		case []byte:
			body = bytes.NewReader(payload)
		case io.Reader:
			body = payload
		default:
			aux := &bytes.Buffer{}
			payload = aux
			return json.NewEncoder(aux).Encode(in["body"])
		}

		return nil
	}()
	if err != nil {
		return nil, err
	}

	if in.Has("timeout") {
		var tfn context.CancelFunc
		ctx, tfn = context.WithTimeout(ctx, in.Duration("timeout"))
		defer tfn()
	}

	req, err = http.NewRequestWithContext(ctx, method, in.String("url"), body)
	if err != nil {
		return nil, err
	}

	header.Set("User-Agent", in.String("headerUserAgent", "Corteza-Automation-Client/"+version.Version))

	if in.Has("header") {
		for k, v := range in["header"].(map[string]string) {
			header.Add(k, v)
		}
	}

	switch {
	case in.Has("headerAuthBearer"):
		header.Add("Authorization", "Bearer "+in.String("headerAuthBearer"))
	case in.Any("headerAuthUsername", "headerAuthPassword"):
		req.SetBasicAuth(
			in.String("headerAuthUsername"),
			in.String("headerAuthPassword"),
		)
	}

	if in.Has("headerContentType") {
		header.Add("Content-Type", in.String("headerContentType"))
	}

	req.Header = header

	return
}

func httpSend(ctx context.Context, in wfexec.Variables) (out wfexec.Variables, err error) {
	var (
		req *http.Request
		rsp *http.Response
	)

	req, err = makeHttpRequest(ctx, in)
	if err != nil {
		return nil, err
	}

	rsp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	out = wfexec.Variables{
		"status":        rsp.Status,
		"statusCode":    rsp.StatusCode,
		"header":        rsp.Header,
		"contentLength": rsp.ContentLength,
		"body":          rsp.Body,
	}

	return
}

func httpSendRequest() *Function {
	return &Function{
		Ref:  "core.http.send", // @todo figure out naming convention
		Meta: FunctionMeta{},
		Parameters: append(append(
			[]*Param{NewParam("method", String, Required)},
			stdHttpSendParameters...),
			stdHttpPayloadParameters...,
		),
		Results: stdHttpSendResults,
		Handler: httpSend,
	}
}

func httpSendGetRequest() *Function {
	return &Function{
		Ref:        "core.http.send.get", // @todo figure out naming convention
		Meta:       FunctionMeta{},
		Parameters: stdHttpSendParameters,
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in wfexec.Variables) (out wfexec.Variables, err error) {
			return httpSend(ctx, in.Merge(wfexec.Variables{"method": "get"}))
		},
	}
}

func httpSendPostRequest() *Function {
	return &Function{
		Ref:        "core.http.send.post", // @todo figure out naming convention
		Meta:       FunctionMeta{},
		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in wfexec.Variables) (out wfexec.Variables, err error) {
			return httpSend(ctx, in.Merge(wfexec.Variables{"method": "post"}))
		},
	}
}

func httpSendPutRequest() *Function {
	return &Function{
		Ref:        "core.http.send.put", // @todo figure out naming convention
		Meta:       FunctionMeta{},
		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in wfexec.Variables) (out wfexec.Variables, err error) {
			return httpSend(ctx, in.Merge(wfexec.Variables{"method": "put"}))
		},
	}
}

func httpSendPatchRequest() *Function {
	return &Function{
		Ref:        "core.http.send.patch", // @todo figure out naming convention
		Meta:       FunctionMeta{},
		Parameters: append(stdHttpSendParameters, stdHttpPayloadParameters...),
		Results:    stdHttpSendResults,
		Handler: func(ctx context.Context, in wfexec.Variables) (out wfexec.Variables, err error) {
			return httpSend(ctx, in.Merge(wfexec.Variables{"method": "patch"}))
		},
	}
}
