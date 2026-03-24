package gcalendar

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/api/drivers"
)

type (
	gcalendarDialect struct {
		calendarID string
	}
)

var (
	_ drivers.Dialect = &gcalendarDialect{}
)

func Dialect(calendarID string) *gcalendarDialect {
	return &gcalendarDialect{calendarID: calendarID}
}

func (d gcalendarDialect) TypeWrap(dt dal.Type) drivers.Type {
	return &drivers.TypeText{TypeText: &dal.TypeText{}}
}

func (d gcalendarDialect) MapOpParams(op string) (method string, endpoint string, err error) {
	switch op {
	case "search":
		method = "GET"
		endpoint = "/events"

	case "read":
		method = "GET"
		endpoint = "/events/{{recordID}}"

	case "create":
		method = "POST"
		endpoint = "/events"

	case "update":
		method = "PUT"
		endpoint = "/events/{{recordID}}"

	case "delete":
		method = "DELETE"
		endpoint = "/events/{{recordID}}"

	default:
		err = fmt.Errorf("operation not supported: %s", op)
	}

	return
}

func (d gcalendarDialect) AddSort(req drivers.XRequest, field string, desc bool) (drivers.XRequest, error) {
	return req, nil
}

func (d gcalendarDialect) AddLimit(req drivers.XRequest, limit uint) (drivers.XRequest, error) {
	return req, nil
}

func (d gcalendarDialect) SearchDataPath() string {
	return ""
}

func (d gcalendarDialect) SearchMetaPath() drivers.BodyMetaPath {
	return drivers.BodyMetaPath{}
}

func (d gcalendarDialect) EncrichEndpoint(endpoint string, xr drivers.XRequest) string {
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

func (d gcalendarDialect) EncodeBodyInsert(ee ...drivers.PayloadEntry) ([]byte, error) {
	row := make(map[string]any)
	for _, e := range ee {
		if len(e.Path) > 0 {
			// Deeply nested paths are not fully supported yet in this basic implementation
			row[e.Path[0]] = e.Value
		}
	}
	return json.Marshal(row)
}

func (d gcalendarDialect) ExtractInsertMeta(in []byte) (out map[string]any, err error) {
	var rsp struct {
		ID string `json:"id"`
	}
	if err = json.Unmarshal(in, &rsp); err != nil {
		return
	}
	out = map[string]any{"id": rsp.ID}
	return
}

func (d gcalendarDialect) DecodeBodySelect(buf []byte) (out map[string]any, err error) {
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
