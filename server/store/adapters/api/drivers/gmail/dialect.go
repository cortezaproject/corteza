package gmail

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/api/drivers"
)

type (
	gmailDialect struct{}
)

var (
	_ drivers.Dialect = &gmailDialect{}
)

func Dialect() *gmailDialect {
	return &gmailDialect{}
}

func (d gmailDialect) TypeWrap(dt dal.Type) drivers.Type {
	return &drivers.TypeText{TypeText: &dal.TypeText{}}
}

func (d gmailDialect) MapOpParams(op string) (method string, endpoint string, err error) {
	switch op {
	case "search":
		method = "GET"
		endpoint = "/messages"

	case "read":
		method = "GET"
		endpoint = "/messages/{{recordID}}"

	case "create":
		method = "POST"
		endpoint = "/messages/send"

	case "delete":
		method = "DELETE"
		endpoint = "/messages/{{recordID}}"

	default:
		err = fmt.Errorf("operation not supported: %s", op)
	}

	return
}

func (d gmailDialect) AddSort(req drivers.XRequest, field string, desc bool) (drivers.XRequest, error) {
	return req, nil
}

func (d gmailDialect) AddLimit(req drivers.XRequest, limit uint) (drivers.XRequest, error) {
	return req, nil
}

func (d gmailDialect) SearchDataPath() string {
	return ""
}

func (d gmailDialect) SearchMetaPath() drivers.BodyMetaPath {
	return drivers.BodyMetaPath{}
}

func (d gmailDialect) EncrichEndpoint(endpoint string, xr drivers.XRequest) string {
	values := url.Values(xr.Query)
	enc := values.Encode()
	if len(enc) == 0 {
		return endpoint
	}

	if len(endpoint) > 0 {
		for i := range endpoint {
			if endpoint[i] == '?' {
				return fmt.Sprintf("%s&%s", endpoint, enc)
			}
		}
	}

	return fmt.Sprintf("%s?%s", endpoint, enc)
}

func (d gmailDialect) EncodeBodyInsert(ee ...drivers.PayloadEntry) ([]byte, error) {
	row := make(map[string]any)
	for _, e := range ee {
		if len(e.Path) > 0 {
			row[e.Path[0]] = e.Value
		}
	}
	return json.Marshal(row)
}

func (d gmailDialect) ExtractInsertMeta(in []byte) (out map[string]any, err error) {
	var rsp struct {
		ID       string `json:"id"`
		ThreadID string `json:"threadId"`
	}
	if err = json.Unmarshal(in, &rsp); err != nil {
		return
	}
	out = map[string]any{"id": rsp.ID, "threadId": rsp.ThreadID}
	return
}

func (d gmailDialect) DecodeBodySelect(buf []byte) (out map[string]any, err error) {
	out = make(map[string]any)
	if err = json.Unmarshal(buf, &out); err == nil {
		return
	}

	var aux []map[string]any
	if err = json.Unmarshal(buf, &aux); err != nil {
		return nil, err
	}
	if len(aux) == 0 {
		return nil, nil
	}
	out = aux[len(aux)-1]
	return
}
