package automation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cortezaproject/corteza-server/automation/types"
	atypes "github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"gopkg.in/mail.v2"
)

type (
	emailMessage struct {
		// only basic implementation for now
		// we can manipulate mail.Message internals through
		// specialized wf functions (message, setSubject, setHeaders, ...)
		msg *mail.Message
	}
)

func CastToEmailMessage(val interface{}) (out *emailMessage, err error) {
	switch val := expr.UntypedValue(val).(type) {
	case *emailMessage:
		if val.msg == nil {
			val.msg = mail.NewMessage()
		}

		return val, nil
	case nil:
		return &emailMessage{msg: mail.NewMessage()}, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToHttpRequest(val interface{}) (out *types.HttpRequest, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = &types.HttpRequest{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToHttpRequest(out, k, v)
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case *http.Request:
		rr := &types.HttpRequest{}

		assignToHttpRequest(rr, "Method", val.Method)
		assignToHttpRequest(rr, "URL", val.URL)
		assignToHttpRequest(rr, "Header", val.Header)
		assignToHttpRequest(rr, "Body", val.Body)
		assignToHttpRequest(rr, "Form", val.Form)
		assignToHttpRequest(rr, "PostForm", val.PostForm)

		return rr, nil
	case *types.HttpRequest:
		return val, nil
	case nil:
		return &types.HttpRequest{}, nil
	default:
		return &types.HttpRequest{}, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToHttpRequestBody(val interface{}) (out *types.HttpRequestBody, err error) {
	switch val := val.(type) {
	case io.ReadCloser:
		rr := &types.HttpRequestBody{}
		return rr, assignToHttpRequestBody(rr, "Body", val)
	}

	switch val := expr.UntypedValue(val).(type) {
	case *io.ReadCloser:
		rr := &types.HttpRequestBody{}
		return rr, assignToHttpRequestBody(rr, "Body", val)
	case *types.HttpRequestBody:
		return val, nil
	case nil:
		return &types.HttpRequestBody{}, nil
	default:
		return &types.HttpRequestBody{}, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToUrl(val interface{}) (out *types.Url, err error) {
	switch val := expr.UntypedValue(val).(type) {
	case *url.URL:
		u := &types.Url{}

		m, _ := json.Marshal(val)
		_ = json.Unmarshal(m, u)

		return u, nil
	case *types.Url:
		return val, nil
	case nil:
		return &types.Url{}, nil
	default:
		return &types.Url{}, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func ReadRequestBody(in interface{}) (s string) {
	var (
		b   []byte
		err error
	)

	switch val := in.(type) {
	case *HttpRequest:
		b, err = val.Get().(*types.HttpRequest).ReadBody()
	case *HttpRequestBody:
		b, err = val.Get().(*types.HttpRequestBody).Read()
	case *atypes.HttpRequest:
		b, err = val.ReadBody()
	case *types.HttpRequest:
		b, err = val.ReadBody()
	case *atypes.HttpRequestBody:
		b, err = val.Read()
	case *types.HttpRequestBody:
		b, err = val.Read()
	case io.Reader:
		b, err = io.ReadAll(val)
	default:
		b = []byte{}
	}

	if err != nil {
		return ""
	}

	return string(b)
}
