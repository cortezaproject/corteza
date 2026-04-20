package gsheets

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/store/adapters/api/drivers"
)

type (
	gsheetsDialect struct {
		spreadsheetID string
	}
)

var (
	_ drivers.Dialect = &gsheetsDialect{}
)

func Dialect(spreadsheetID string) *gsheetsDialect {
	return &gsheetsDialect{spreadsheetID: spreadsheetID}
}

// TypeWrap returns TypeText for all attribute types since Sheets stores everything as strings.
func (d gsheetsDialect) TypeWrap(dt dal.Type) drivers.Type {
	return &drivers.TypeText{TypeText: &dal.TypeText{}}
}

// MapOpParams maps DAL operations to HTTP method + Sheets v4 endpoint templates.
//
// The {{sheetName}} placeholder is resolved by model.procEndpointM (mapped to moduleID).
// The {{recordID}} placeholder is resolved by model.procEndpointMR for row-level ops
// but for Sheets we handle row targeting differently (see wrapper).
func (d gsheetsDialect) MapOpParams(op string) (method string, endpoint string, err error) {
	switch op {
	case "search":
		method = "GET"
		endpoint = "/values/{{moduleID}}"

	case "read":
		method = "GET"
		endpoint = "/values/{{moduleID}}!1:{{recordID}}"

	case "create":
		method = "POST"
		endpoint = "/values/{{moduleID}}:append?valueInputOption=USER_ENTERED&insertDataOption=INSERT_ROWS"

	case "update":
		method = "PUT"
		endpoint = "/values/{{moduleID}}!{{recordID}}:{{recordID}}?valueInputOption=USER_ENTERED"

	case "delete":
		// Soft delete — updates the row with a deletedAt value
		method = "PUT"
		endpoint = "/values/{{moduleID}}!{{recordID}}:{{recordID}}?valueInputOption=USER_ENTERED"

	default:
		err = fmt.Errorf("operation not supported: %s", op)
	}

	return
}

// AddSort is a no-op; sorting is done in-memory after fetching all rows.
func (d gsheetsDialect) AddSort(req drivers.XRequest, field string, desc bool) (drivers.XRequest, error) {
	return req, nil
}

// AddLimit is a no-op; pagination is done in-memory after fetching all rows.
func (d gsheetsDialect) AddLimit(req drivers.XRequest, limit uint) (drivers.XRequest, error) {
	return req, nil
}

// SearchDataPath returns empty — after the wrapper transforms the response, data is at root level.
func (d gsheetsDialect) SearchDataPath() string {
	return ""
}

// SearchMetaPath returns empty — Sheets has no server-side pagination metadata.
func (d gsheetsDialect) SearchMetaPath() drivers.BodyMetaPath {
	return drivers.BodyMetaPath{}
}

// EncrichEndpoint enriches the endpoint with query parameters.
func (d gsheetsDialect) EncrichEndpoint(endpoint string, xr drivers.XRequest) string {
	values := url.Values(xr.Query)
	enc := values.Encode()
	if len(enc) == 0 {
		return endpoint
	}

	// Sheets endpoints may already contain query params (e.g. valueInputOption)
	if len(endpoint) > 0 {
		for i := range endpoint {
			if endpoint[i] == '?' {
				return fmt.Sprintf("%s&%s", endpoint, enc)
			}
		}
	}

	return fmt.Sprintf("%s?%s", endpoint, enc)
}

func (d gsheetsDialect) EncodeBodyInsert(ee ...drivers.PayloadEntry) ([]byte, error) {
	row := make([]any, 0, len(ee))

	for _, e := range ee {
		row = append(row, e.Value)
	}

	bb, err := json.Marshal(map[string]any{
		"values": [][]any{row},
	})
	if err != nil {
		return nil, err
	}

	return bb, nil
}

func (d gsheetsDialect) ExtractInsertMeta(in []byte) (out map[string]any, err error) {
	var rsp struct {
		Updates struct {
			UpdatedRange string `json:"updatedRange"`
		} `json:"updates"`
	}
	if err = json.Unmarshal(in, &rsp); err != nil {
		return
	}

	// updatedRange format: "sheet!A6:D6" — extract row number
	r := rsp.Updates.UpdatedRange
	if idx := strings.LastIndex(r, ":"); idx > 0 {
		r = r[:idx] // "sheet!A6"
	}
	if idx := strings.LastIndex(r, "!"); idx >= 0 {
		r = r[idx+1:]
	}
	r = strings.TrimLeft(r, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	out = map[string]any{"id": r}
	return
}

func (d gsheetsDialect) DecodeBodySelect(buf []byte) (out map[string]any, err error) {
	// Try single object first
	out = make(map[string]any)
	if err = json.Unmarshal(buf, &out); err == nil {
		return
	}

	// Fall back to array — take the last element
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
