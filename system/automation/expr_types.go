package automation

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cast"
)

type (
	renderedDocument struct {
		Document io.Reader
		Name     string
		Type     string
	}
)

func CastToUser(val interface{}) (out *types.User, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = &types.User{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToUser(out, k, v)
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case *types.User:
		return val, nil
	case map[string]interface{}:
		out = &types.User{}
		m, _ := json.Marshal(val)
		_ = json.Unmarshal(m, out)
		return out, nil
	case nil:
		return &types.User{}, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToRole(val interface{}) (out *types.Role, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = &types.Role{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToRole(out, k, v)
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case *types.Role:
		return val, nil
	case map[string]interface{}:
		out = &types.Role{}
		m, _ := json.Marshal(val)
		_ = json.Unmarshal(m, out)
		return out, nil
	case nil:
		return &types.Role{}, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToTemplate(val interface{}) (out *types.Template, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = &types.Template{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToTemplate(out, k, v)
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case *types.Template:
		return val, nil
	case map[string]interface{}:
		out = &types.Template{}
		m, _ := json.Marshal(val)
		_ = json.Unmarshal(m, out)
		return out, nil
	case nil:
		return &types.Template{}, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToTemplateMeta(val interface{}) (out types.TemplateMeta, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = types.TemplateMeta{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToTemplateMeta(out, k, v)
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case types.TemplateMeta:
		return val, nil
	case nil:
		return types.TemplateMeta{}, nil
	default:
		return types.TemplateMeta{}, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToRenderedDocument(val interface{}) (out *renderedDocument, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = &renderedDocument{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToRenderedDocument(out, k, v)
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case *renderedDocument:
		return val, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToDocumentType(val interface{}) (out types.DocumentType, err error) {
	switch val := val.(type) {
	case string:
		return types.DocumentType(val), nil
	case *expr.String:
		return types.DocumentType(val.GetValue()), nil
	default:
		return "", fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToRenderOptions(val interface{}) (out map[string]string, err error) {
	switch val := expr.UntypedValue(val).(type) {
	case map[string]string:
		return val, nil
	case nil:
		return make(map[string]string), nil
	default:
		out, err = cast.ToStringMapStringE(val)
		if err != nil {
			return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
		}
		return out, nil
	}
}

func (t *RenderOptions) AssignFieldValue(key string, val expr.TypedValue) error {
	if t.value == nil {
		t.value = make(map[string]string)
	}

	str, err := cast.ToStringE(expr.UntypedValue(val))
	t.value[key] = str
	return err
}

func CastToQueueMessage(val interface{}) (out *types.QueueMessage, err error) {
	switch val := val.(type) {
	case expr.Iterator:
		out = &types.QueueMessage{}
		return out, val.Each(func(k string, v expr.TypedValue) error {
			return assignToQueueMessage(out, k, v)
		})
	}

	switch val := expr.UntypedValue(val).(type) {
	case *types.QueueMessage:
		return val, nil
	case nil:
		return &types.QueueMessage{}, nil
	default:
		return &types.QueueMessage{}, fmt.Errorf("unable to cast type %T to %T", val, out)
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

func (doc renderedDocument) String() string {
	aux, _ := ioutil.ReadAll(doc.Document)
	return string(aux)
}

func CastToRbacResource(val interface{}) (out rbac.Resource, err error) {
	switch val := expr.UntypedValue(val).(type) {
	case rbac.Resource:
		return val, nil
	case RbacResource:
		return val.value, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}

func CastToAction(val interface{}) (out *actionlog.Action, err error) {
	switch val := expr.UntypedValue(val).(type) {
	case *actionlog.Action:
		return val, nil
	default:
		return nil, fmt.Errorf("unable to cast type %T to %T", val, out)
	}
}
