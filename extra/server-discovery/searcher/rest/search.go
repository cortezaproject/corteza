package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crusttech/human/extra/server-discovery/searcher"
	"github.com/crusttech/human/extra/server-discovery/searcher/rest/request"
	types2 "github.com/crusttech/human/server/compose/types"
)

type (
	search struct {
		embedder   embedderService
		searchMode string
	}

	embedderService interface {
		GenerateEmbeddings(input string) ([]float64, error)
	}

	cResponse struct {
		Response struct {
			Set []struct {
				NamespaceID uint64 `json:",string"`
				Slug        string `json:"slug"`

				Name     string              `json:"name"`
				ModuleID uint64              `json:",string"`
				Handle   string              `json:"handle"`
				Config   types2.ModuleConfig `json:"config"`
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

func Search(embedderSvc embedderService, searchMode string) *search {
	return &search{
		embedder:   embedderSvc,
		searchMode: searchMode,
	}
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
		resourceTypes = r.GetResourceTypes()
		namespaceIDs  = r.GetNamespaceIDs()
		moduleIDs     = r.GetModuleIDs()
		validDumpRaw  = r.GetDumpRaw() != ""

		page          pagination
		results       *esSearchResponse
		aggregation   *esSearchResponse
		nsAggregation *esSearchResponse
		mAggregation  *esSearchResponse

		indexFieldExclusions []indexField

		nsHandleMap = make(map[string]nsMeta)
		mHandleMap  = make(map[string]mMeta)

		searchMode string
	)

	esc := searcher.DefaultEsClient

	if r.GetSearchMode() != "" {
		searchMode = r.GetSearchMode()
	} else {
		searchMode = s.searchMode
	}

	indexFieldExclusions, nsHandleMap, mHandleMap, err = fetchModNamespaceResources(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch namespace/module resources: %w", err)
	}

	results, page, err = esSearch(ctx, log, esc, searchParams{
		title:                "results",
		query:                searchString,
		resourceType:         resourceTypes,
		from:                 from,
		size:                 size,
		moduleAggs:           moduleAggs,
		namespaceAggs:        namespaceAggs,
		namespaceIDs:         namespaceIDs,
		moduleIDs:            moduleIDs,
		indexFieldExclusions: indexFieldExclusions,
		dumpRaw:              validDumpRaw,
		allowedRoles:         allowedRoles,
		embedder:             s.embedder,
		searchMode:           searchMode,
	})
	if err != nil {
		return nil, fmt.Errorf("could not execute search: %w", err)
	}

	if len(searchString) == 0 {
		aggregation, _, err = esSearch(ctx, log, esc, searchParams{
			title: "aggregation",
			// size:          size,
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
		// size:    size,
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
		// size:          size,
		dumpRaw:       validDumpRaw,
		query:         searchString,
		namespaceAggs: namespaceAggs,
		aggOnly:       true,
		mAggOnly:      true,
		allowedRoles:  allowedRoles,
		embedder:      s.embedder,
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

	return conv(results, aggregation, noHits, nsHandleMap, mHandleMap, page)
}

func fetchModNamespaceResources(ctx context.Context) (indexFieldExclusions []indexField, nsHandleMap map[string]nsMeta, mHandleMap map[string]mMeta, err error) {
	nsHandleMap = make(map[string]nsMeta)
	mHandleMap = make(map[string]mMeta)

	// Prepare request to get namespaces
	nsReq, err := searcher.DefaultApiClient.Namespaces()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to prepare namespace request: %w", err)
	}

	nsRes, err := searcher.DefaultApiClient.HttpClient().Do(nsReq.WithContext(ctx))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to send namespace request: %w", err)
	}
	defer nsRes.Body.Close()

	if nsRes.StatusCode != http.StatusOK {
		return nil, nil, nil, fmt.Errorf("request resulted in an unexpected status: %s", nsRes.Status)
	}

	var nsResponse cResponse
	if err = json.NewDecoder(nsRes.Body).Decode(&nsResponse); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode namespace response: %w", err)
	}

	for _, s := range nsResponse.Response.Set {
		// Get the module handles for aggs response
		nsHandleMap[s.Slug] = nsMeta{
			Name:   s.Name,
			Handle: s.Slug,
		}

		mReq, err := searcher.DefaultApiClient.Modules(s.NamespaceID)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to prepare module meta request: %w", err)
		}

		mRes, err := searcher.DefaultApiClient.HttpClient().Do(mReq.WithContext(ctx))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to send module request: %w", err)
		}

		if mRes.StatusCode != http.StatusOK {
			_ = mRes.Body.Close()
			return nil, nil, nil, fmt.Errorf("request resulted in an unexpected status: %s", mRes.Status)
		}

		var mResponse cResponse
		if err = json.NewDecoder(mRes.Body).Decode(&mResponse); err != nil {
			_ = mRes.Body.Close()
			return nil, nil, nil, fmt.Errorf("failed to decode response: %w", err)
		}
		_ = mRes.Body.Close()

		for _, m := range mResponse.Response.Set {
			// Get the module handles for aggs response
			mHandleMap[m.Handle] = mMeta{
				Name:   m.Name,
				Handle: m.Slug,
			}

			key := fmt.Sprintf("human-private-compose-records-%d-%d", s.NamespaceID, m.ModuleID)
			if len(m.Config.Discovery.Private.Result) > 0 && len(m.Config.Discovery.Private.Result[0].Fields) > 0 {
				var fields []string
				for _, field := range m.Config.Discovery.Private.Result[0].Fields {
					fields = append(fields, fmt.Sprintf("values.%s", field))
				}

				indexFieldExclusions = append(indexFieldExclusions, indexField{
					IndexName: key,
					Fields:    fields,
				})
			}
		}
	}

	return indexFieldExclusions, nsHandleMap, mHandleMap, nil
}
