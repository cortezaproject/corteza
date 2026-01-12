package rest

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/api/drivers"
	"github.com/spf13/cast"
)

type (
	defaultOps struct {
		Create struct {
			Endpoint string `json:"endpoint"`
			Method   string `json:"method"`

			// ...
		} `json:"create"`

		Delete struct {
			Endpoint string `json:"endpoint"`
			Method   string `json:"method"`

			// ...
		} `json:"delete"`

		Read struct {
			Endpoint string `json:"endpoint"`
			Method   string `json:"method"`

			// ...
		} `json:"read"`

		Search struct {
			Endpoint string `json:"endpoint"`
			Method   string `json:"method"`

			// ...
			Payload struct {
				DataPath string `json:"dataPath"`

				Meta struct {
					TotalCountPath string `json:"totalCountPath"`
				} `json:"meta"`
			} `json:"payload"`

			Sorting struct {
				Location string `json:"location"`
				KeyTpl   string `json:"keyTpl"`
				// ValueTpl string `json:"valueTpl"`
				FieldTpl  string `json:"fieldTpl"`
				AscLabel  string `json:"ascLabel"`
				DescLabel string `json:"descLabel"`
				MergeTpl  string `json:"mergeTpl"`
			} `json:"sorting"`

			Limit struct {
				Location string `json:"location"`
				Value    string `json:"value"`
			} `json:"limit"`
		} `json:"search"`

		Update struct {
			Endpoint string `json:"endpoint"`
			Method   string `json:"method"`

			// ...
		} `json:"update"`
	}

	apiDialect struct {
		defaultOps defaultOps
	}
)

var (
	_ drivers.Dialect = &apiDialect{}

	dialect = &apiDialect{}

	nuances = drivers.Nuances{}
)

func init() {
}

func Dialect() *apiDialect {
	return dialect
}

func (d apiDialect) TypeWrap(dt dal.Type) drivers.Type {
	// Any exception to general type-wrap implementation in the drivers package
	// should be placed here
	switch c := dt.(type) {
	case *dal.TypeTimestamp:
		return &drivers.TypeTimestamp{&dal.TypeTimestamp{
			Nullable: c.Nullable,

			// sqlserver does not support timezone
			Timezone: false,

			// sqlserver does not support precision
			Precision: 0,
		}}
	}

	return drivers.TypeWrap(dt)
}

func (d apiDialect) MapOpParams(op string) (method string, endpoint string, err error) {
	switch op {
	case "create":
		method = d.defaultOps.Create.Method
		endpoint = d.defaultOps.Create.Endpoint

	case "delete":
		method = d.defaultOps.Delete.Method
		endpoint = d.defaultOps.Delete.Endpoint

	case "read":
		method = d.defaultOps.Read.Method
		endpoint = d.defaultOps.Read.Endpoint

	case "search":
		method = d.defaultOps.Search.Method
		endpoint = d.defaultOps.Search.Endpoint

	case "update":
		method = d.defaultOps.Update.Method
		endpoint = d.defaultOps.Update.Endpoint

	default:
		err = fmt.Errorf("unknown operation: %s", op)
		return
	}

	return
}

func (d apiDialect) AddSort(req drivers.XRequest, field string, desc bool) (out drivers.XRequest, err error) {
	out = req

	var (
		mergeTpl  = d.defaultOps.Search.Sorting.MergeTpl
		location  = d.defaultOps.Search.Sorting.Location
		keyTpl    = d.defaultOps.Search.Sorting.KeyTpl
		fieldTpl  = d.defaultOps.Search.Sorting.FieldTpl
		ascLabel  = d.defaultOps.Search.Sorting.AscLabel
		descLabel = d.defaultOps.Search.Sorting.DescLabel
	)

	// Defaults
	{
		if location == "" {
			location = "query"
		}

		if keyTpl == "" {
			keyTpl = "sort"
		}

		if fieldTpl == "" {
			fieldTpl = ""
		}

		if ascLabel == "" {
			ascLabel = "asc"
		}

		if descLabel == "" {
			descLabel = "desc"
		}
	}

	// Sort label
	sort := ascLabel
	if desc {
		sort = descLabel
	}

	{
		rpl := strings.NewReplacer(
			"{{field}}", field,
			"{{sort}}", sort,
		)

		field = rpl.Replace(fieldTpl)
	}

	// @todo process the key template when it becomes needed
	key := keyTpl

	switch location {
	case "query":
		out, err = d.addSortQuery(out, key, field, mergeTpl)
		if err != nil {
			return
		}

	default:
		err = fmt.Errorf("can not add sort: location not supported: %s", location)
		return
	}

	return
}

func (d apiDialect) addSortQuery(req drivers.XRequest, key string, field string, mergeTpl string) (out drivers.XRequest, err error) {
	out = req

	if out.Query == nil {
		out.Query = map[string][]string{}
	}

	if mergeTpl == "" {
		out.Query[key] = append(out.Query[key], field)
		return
	}

	prev := ""
	if len(out.Query[key]) > 0 {
		prev = out.Query[key][len(out.Query[key])-1]

		prev = strings.ReplaceAll(
			mergeTpl,
			"{{old}}",
			prev,
		)

		prev = strings.ReplaceAll(
			prev,
			"{{new}}",
			field,
		)

		out.Query[key][0] = prev
	} else {
		out.Query[key] = append(out.Query[key], field)
	}

	return
}

func (d apiDialect) AddLimit(req drivers.XRequest, count uint) (out drivers.XRequest, err error) {
	// @todo different locations?
	out = req

	if out.Query == nil {
		out.Query = map[string][]string{}
	}

	out.Query["limit"] = []string{cast.ToString(count)}

	return
}

func (d apiDialect) SearchDataPath() string {
	return d.defaultOps.Search.Payload.DataPath
}

func (d apiDialect) SearchMetaPath() drivers.BodyMetaPath {
	return drivers.BodyMetaPath{
		TotalCount: d.defaultOps.Search.Payload.Meta.TotalCountPath,
	}
}

func (d apiDialect) EncrichEndpoint(endpoint string, xr drivers.XRequest) (out string) {
	values := url.Values(xr.Query)

	enc := values.Encode()
	if len(enc) == 0 {
		return endpoint
	}

	return fmt.Sprintf("%s?%s", endpoint, enc)
}
