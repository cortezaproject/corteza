package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server-discovery/searcher"
	"github.com/cortezaproject/corteza-server-discovery/searcher/rest/request"
	"github.com/jmoiron/sqlx/types"
	"net/http"
)

type (
	search struct{}

	cResponse struct {
		Response struct {
			Set []struct {
				NamespaceID uint64 `json:",string"`
				Slug        string `json:"slug"`

				Name     string         `json:"name"`
				ModuleID uint64         `json:",string"`
				Handle   string         `json:"handle"`
				Meta     types.JSONText `json:"meta"`
			} `json:"set,omitempty"`
		} `json:"response,omitempty"`
	}

	moduleMeta struct {
		Discovery ModuleMeta `json:"discovery"`
	}

	ModuleMeta struct {
		Public struct {
			Result []Result `json:"result"`
		} `json:"public"`
		Private struct {
			Result []Result `json:"result"`
		} `json:"private"`
		Protected struct {
			Result []Result `json:"result"`
		} `json:"protected"`
	}

	Result struct {
		Lang   string   `json:"lang"`
		Fields []string `json:"fields"`
		// @todo? TBD? excludeModuleFields, includeModuleFields <- if passed filter module field accordingly.
	}

	nsMeta struct {
		Name   string
		Handle string
	}
	mMeta struct {
		Name   string
		Handle string
	}
)

func Search() *search {
	return &search{}
}

func (s search) SearchResources(ctx context.Context, r *request.SearchResources) (out interface{}, err error) {
	var (
		log           = searcher.DefaultLogger
		allowedRoles  = searcher.DefaultConfig.Searcher.AllowedRole
		searchString  = r.GetQuery()
		size          = r.GetSize()
		from          = r.GetFrom()
		namespaceAggs = r.GetNamespaceAggs()
		moduleAggs    = r.GetModuleAggs()
		validDumpRaw  = r.GetDumpRaw() != ""

		page          pagination
		results       *esSearchResponse
		aggregation   *esSearchResponse
		nsAggregation *esSearchResponse
		mAggregation  *esSearchResponse

		nsReq      *http.Request
		nsRes      *http.Response
		mReq       *http.Request
		mRes       *http.Response
		nsResponse cResponse
		mResponse  cResponse
		moduleMap  = make(map[string][]string)

		nsHandleMap = make(map[string]nsMeta)
		mHandleMap  = make(map[string]mMeta)
	)

	esc, err := searcher.DefaultEs.Client()
	if err != nil {
		return nil, err
	}

	results, page, err = esSearch(ctx, log, esc, searchParams{
		title:         "results",
		query:         searchString,
		from:          from,
		size:          size,
		moduleAggs:    moduleAggs,
		namespaceAggs: namespaceAggs,
		dumpRaw:       validDumpRaw,
		allowedRoles:  allowedRoles,
	})
	if err != nil {
		return nil, fmt.Errorf("could not execute search: %w", err)
	}

	if len(searchString) == 0 {
		aggregation, _, err = esSearch(ctx, log, esc, searchParams{
			title: "aggregation",
			//size:          size,
			dumpRaw:       validDumpRaw,
			namespaceAggs: namespaceAggs,
			aggOnly:       true,
			allowedRoles:  allowedRoles,
		})
		if err != nil {
			return nil, fmt.Errorf("could not execute aggregation search: %w", err)
		}
	}

	// append all namespace agg with counts no matter what
	nsAggregation, _, err = esSearch(ctx, log, esc, searchParams{
		title: "nsAggregation",
		//size:    size,
		dumpRaw:      validDumpRaw,
		aggOnly:      true,
		allowedRoles: allowedRoles,
	})
	if err != nil {
		return nil, fmt.Errorf("could not execute namespace aggregation search: %w", err)
	}

	if len(searchString) == 0 {
		if aggregation != nil && nsAggregation != nil {
			aggregation.Aggregations.Namespace = nsAggregation.Aggregations.Namespace
		}
	} else {
		if results != nil && nsAggregation != nil {
			nsMap := make(map[string]struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			})
			for _, bucket := range results.Aggregations.Namespace.Namespace.Buckets {
				nsMap[bucket.Key] = bucket
			}

			var buckets []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			}
			for _, bucket := range nsAggregation.Aggregations.Namespace.Namespace.Buckets {
				val, ok := nsMap[bucket.Key]
				if ok {
					val.DocCount = nsMap[bucket.Key].DocCount
				} else {
					val.Key = bucket.Key
					val.DocCount = 0
				}
				buckets = append(buckets, val)
			}

			results.Aggregations.Namespace.Namespace.Buckets = buckets
		}
	}
	// append namespace agg response which are not in es response
	if results != nil && len(namespaceAggs) > 0 {
		nsMap := make(map[string]struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		})
		var bb []struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		}
		for _, b := range results.Aggregations.Namespace.Namespace.Buckets {
			nsMap = map[string]struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			}{
				b.Key: b,
			}
			bb = append(bb, b)
		}

		for _, agg := range namespaceAggs {
			if _, ok := nsMap[agg]; !ok {
				nsMap = map[string]struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
				}{
					agg: {Key: agg, DocCount: 0},
				}
				bb = append(bb, struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
				}{Key: agg, DocCount: 0})
			}
		}

		if len(bb) > 0 {
			results.Aggregations.Namespace.Namespace.Buckets = bb
		}
	}

	mAggregation, _, err = esSearch(ctx, log, esc, searchParams{
		title: "mAggregation",
		//size:          size,
		dumpRaw:       validDumpRaw,
		query:         searchString,
		namespaceAggs: namespaceAggs,
		aggOnly:       true,
		mAggOnly:      true,
		allowedRoles:  allowedRoles,
	})
	if err != nil {
		return nil, fmt.Errorf("could not execute module aggregation search: %w", err)
	}
	if len(searchString) > 0 {
		if results != nil && mAggregation != nil {
			results.Aggregations.Module = mAggregation.Aggregations.Module
		}
	}

	// append module agg response which are not in es response
	if results != nil && len(moduleAggs) > 0 {
		mMap := make(map[string]struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		})
		var bb []struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		}
		for _, b := range results.Aggregations.Module.Module.Buckets {
			mMap = map[string]struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			}{
				b.Key: b,
			}
			bb = append(bb, b)
		}

		for _, agg := range moduleAggs {
			if _, ok := mMap[agg]; !ok {
				mMap = map[string]struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
				}{
					agg: {Key: agg, DocCount: 0},
				}
				bb = append(bb, struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
				}{Key: agg, DocCount: 0})
			}
		}

		if len(bb) > 0 {
			results.Aggregations.Module.Module.Buckets = bb
		}
	}

	noHits := len(searchString) == 0 && len(moduleAggs) == 0 && len(namespaceAggs) == 0
	//if !noHits {
	// @todo only fetch module from result but that requires another loop to fetch module Id from es response
	// 			TEMP fix, I have solution use elastic for the same but different index
	nsReq, err = searcher.DefaultApiClient.Namespaces()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare namespace request: %w", err)
	} else {
		if nsRes, err = searcher.DefaultApiClient.HttpClient().Do(nsReq.WithContext(ctx)); err != nil {
			return nil, fmt.Errorf("failed to send namespace request: %w", err)
		}
		if nsRes.StatusCode != http.StatusOK {
			fmt.Println("err: ", err)
			return nil, fmt.Errorf("request resulted in an unexpected status: %s: %w", err)
		}
		if err = json.NewDecoder(nsRes.Body).Decode(&nsResponse); err != nil {
			return nil, fmt.Errorf("failed to decode namespace response: %w", err)
		}

		if err = nsRes.Body.Close(); err != nil {
			return nil, fmt.Errorf("failed to close namespace response body: %w", err)
		}

		for _, s := range nsResponse.Response.Set {
			// Get the module handles for aggs response
			nsHandleMap[s.Slug] = nsMeta{
				Name:   s.Name,
				Handle: s.Slug,
			}
			if mReq, err = searcher.DefaultApiClient.Modules(s.NamespaceID); err != nil {
				return nil, fmt.Errorf("failed to prepare module meta request: %w", err)
			}
			if mRes, err = searcher.DefaultApiClient.HttpClient().Do(mReq.WithContext(ctx)); err != nil {
				return nil, fmt.Errorf("failed to send module request: %w", err)
			}
			if mRes.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("request resulted in an unexpected status: %s: %w", err)
			}
			if err = json.NewDecoder(mRes.Body).Decode(&mResponse); err != nil {
				return nil, fmt.Errorf("failed to decode response: %w", err)
			}
			if err = mRes.Body.Close(); err != nil {
				return nil, fmt.Errorf("failed to close response body: %w", err)
			}

			for _, m := range mResponse.Response.Set {
				// Get the module handles for aggs response
				mHandleMap[m.Handle] = mMeta{
					Name:   m.Name,
					Handle: m.Slug,
				}

				var (
					meta moduleMeta
					key  = fmt.Sprintf("%d-%d", s.NamespaceID, m.ModuleID)
				)
				err = json.Unmarshal(m.Meta, &meta)
				if err != nil {
					return nil, fmt.Errorf("failed to unmarshal module meta: %w", err)
				} else if len(meta.Discovery.Private.Result) > 0 && len(meta.Discovery.Private.Result[0].Fields) > 0 {
					moduleMap[key] = meta.Discovery.Private.Result[0].Fields
				}
			}
		}
	}
	//}

	return conv(results, aggregation, noHits, moduleMap, nsHandleMap, mHandleMap, page)
}
