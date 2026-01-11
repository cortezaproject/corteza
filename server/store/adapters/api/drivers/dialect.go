package drivers

import (
	"fmt"
	"net/url"

	"github.com/cortezaproject/corteza/server/pkg/dal"
)

type (
	Nuances struct {
	}

	XRequest struct {
		Query map[string][]string
	}

	Dialect interface {
		// TypeWrap returns driver's type implementation for a particular attribute type
		TypeWrap(dal.Type) Type

		MapOpParams(op string) (method string, endpoint string, err error)
		AddSort(req XRequest, field string, desc bool) (XRequest, error)
		AddLimit(req XRequest, limit uint) (XRequest, error)

		EncrichEndpoint(endpoint string, xr XRequest) (out string)

		SearchDataPath() string
		SearchMetaPath() string
	}
)

func (r XRequest) EnrichEndpoint(endpoint string) (out string) {
	values := url.Values(r.Query)

	return fmt.Sprintf("%s?%s", endpoint, values.Encode())
}
