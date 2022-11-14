package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type (
	// EsSearchAggrTerms is aggregations parameter for es search api.
	EsSearchAggrTerms       map[string]esSearchAggr
	EsSearchNestedAggrTerms map[string]esSearchNestedAggr

	esSearchParamsIndex struct {
		Prefix struct {
			Index struct {
				Value string `json:"value"`
			} `json:"_index"`
		} `json:"prefix"`
	}

	esSimpleQueryString struct {
		Wrap struct {
			Query string `json:"query"`
		} `json:"simple_query_string"`
	}

	esDisMax struct {
		TieBreaker float64       `json:"tie_breaker,omitempty"`
		Boost      float64       `json:"boost,omitempty"`
		Queries    []interface{} `json:"queries,omitempty"`
	}

	esDisMaxWrap struct {
		Wrap esDisMax `json:"dis_max,omitempty"`
	}

	esNestedDisMax struct {
		Wrap struct {
			IgnoreUnmapped bool         `json:"ignore_unmapped"`
			Path           string       `json:"path,omitempty"`
			Query          esDisMaxWrap `json:"query,omitempty"`
		} `json:"nested,omitempty"`
	}

	esMultiMatch struct {
		Wrap struct {
			Query string `json:"query"`
			Type  string `json:"type"`
			// Operator string   `json:"operator"`
			Fields []string `json:"fields"`
		} `json:"multi_match"`
	}

	esSearchParams struct {
		Query struct {
			Bool struct {
				// query context
				Must []interface{} `json:"must,omitempty"`

				// filter context
				Filter  []interface{} `json:"filter,omitempty"`
				MustNot []interface{} `json:"must_not,omitempty"`
			} `json:"bool,omitempty"`
		} `json:"query"`

		Aggregations EsSearchNestedAggrTerms `json:"aggs,omitempty"`
	}

	esSearchAggrTerm struct {
		Field string `json:"field,omitempty"`
		Size  int    `json:"size,omitempty"`
	}

	esSearchAggrComposite struct {
		Sources interface{} `json:"sources"` // it can be esSearchAggrTerm,.. (Histogram, Date histogram, GeoTile grid)
		Size    int         `json:"size,omitempty"`
	}

	esSearchNestedAggs struct {
		Path string `json:"path,omitempty"`
	}

	esSearchAggr struct {
		Nested       esSearchNestedAggs `json:"nested,omitempty"`
		Terms        esSearchAggrTerm   `json:"terms,omitempty"`
		Aggregations EsSearchAggrTerms  `json:"aggs,omitempty"`
		// Composite *esSearchAggrComposite `json:"composite"`
	}

	esSearchNestedAggrTerm struct {
		Terms esSearchAggrTerm `json:"terms,omitempty"`
	}
	esSearchNestedAggr struct {
		Nested       esSearchNestedAggs                `json:"nested,omitempty"`
		Aggregations map[string]esSearchNestedAggrTerm `json:"aggs,omitempty"`
		// Composite *esSearchAggrComposite `json:"composite"`
	}

	esSearchResponse struct {
		Took         int                  `json:"took"`
		TimedOut     bool                 `json:"timed_out"`
		Hits         esSearchHits         `json:"hits"`
		Aggregations esSearchAggregations `json:"aggregations"`
	}

	esSearchTotal struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	}

	esSearchHits struct {
		Total esSearchTotal  `json:"total"`
		Hits  []*esSearchHit `json:"hits"`
	}

	esSearchHit struct {
		Index  string          `json:"_index"`
		ID     string          `json:"_id"`
		Source json.RawMessage `json:"_source"`
	}

	esSearchAggregations struct {
		Resource struct {
			DocCountErrorUpperBound int `json:"-"`
			SumOtherDocCount        int `json:"-"`
			Buckets                 []struct {
				Key          string `json:"key"`
				DocCount     int    `json:"doc_count"`
				ResourceName struct {
					DocCountErrorUpperBound int `json:"-"`
					SumOtherDocCount        int `json:"-"`
					Buckets                 []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"resourceName"`
				Namespaces struct {
					DocCountErrorUpperBound int `json:"-"`
					SumOtherDocCount        int `json:"-"`
					Buckets                 []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"namespaces"`
				Modules struct {
					DocCountErrorUpperBound int `json:"-"`
					SumOtherDocCount        int `json:"-"`
					Buckets                 []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"modules"`
			} `json:"buckets"`
		} `json:"resource"`
		Module struct {
			DocCount int `json:"doc_count"`
			Module   struct {
				DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
				SumOtherDocCount        int `json:"sum_other_doc_count"`
				Buckets                 []struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
				} `json:"buckets"`
			} `json:"module"`
		} `json:"module"`

		Namespace struct {
			DocCount  int `json:"doc_count"`
			Namespace struct {
				DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
				SumOtherDocCount        int `json:"sum_other_doc_count"`
				Buckets                 []struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
				} `json:"buckets"`
			} `json:"namespace"`
		} `json:"namespace"`
	}

	searchParams struct {
		title         string
		query         string
		moduleAggs    []string
		namespaceAggs []string
		dumpRaw       bool
		from          int
		size          int

		aggOnly  bool
		mAggOnly bool

		allowedRoles map[interface{}]bool
	}
)

func esSearch(ctx context.Context, log *zap.Logger, esc *elasticsearch.Client, p searchParams) (sr *esSearchResponse, page pagination, err error) {
	var (
		buf          bytes.Buffer
		roles        []string
		userID       uint64
		_, claims, _ = jwtauth.FromContext(ctx)

		allowedRoleExist = false
	)

	if _, has := claims["roles"]; has {
		if rr, is := claims["roles"].([]interface{}); is {

			for _, role := range rr {
				if _, ok := p.allowedRoles[role]; ok {
					allowedRoleExist = true
				}

				roles = append(roles, fmt.Sprintf("%s", role))
			}
		}
	}
	if _, has := claims["sub"]; has {
		if sub, is := claims["sub"].(string); is {
			userID, _ = extractFromSubClaim(sub)
		}
	}

	noQ := len(p.query) == 0
	noNSFilter := len(p.namespaceAggs) == 0
	// noMFilter := len(p.moduleAggs) == 0
	sqs := esSimpleQueryString{}
	sqs.Wrap.Query = p.query

	query := esSearchParams{}
	index := esSearchParamsIndex{}

	// Decide what indexes we can use
	if userID == 0 {
		// Missing, invalid, expired access token (JWT)
		index.Prefix.Index.Value = "corteza-public-"
	} else {
		// Authenticated user
		index.Prefix.Index.Value = "corteza-private-"

		if !allowedRoleExist {
			// Skip all documents that do not have baring roles in to allow list
			query.Query.Bool.Filter = append(query.Query.Bool.Filter, map[string]map[string]interface{}{
				"terms": {"security.allowedRoles": roles},
			})

			// Skip all documents that have baring roles in to deny list
			query.Query.Bool.MustNot = append(query.Query.Bool.MustNot, map[string]map[string]interface{}{
				"terms": {"security.deniedRoles": roles},
			})
		}
	}

	// Query MUST filter
	query.Query.Bool.Must = []interface{}{index}

	// Aggregations V1.0
	// if len(p.aggregations) > 0 {
	//	query.Aggregations = make(map[string]esSearchAggr)
	//
	//	for _, a := range p.aggregations {
	//		query.Aggregations[a] = esSearchAggr{esSearchAggrTerm{Field: a + ".keyword"}}
	//	}
	// }

	// Search string filter
	if !noQ {
		sqs.Wrap.Query = fmt.Sprintf("%s*", sqs.Wrap.Query)
		query.Query.Bool.Must = append(query.Query.Bool.Must, sqs)
		// query.Query.DisMax.Queries = append(query.Query.DisMax.Queries, sqs)
	}

	var (
		mm   = esMultiMatch{}
		mdd  esNestedDisMax
		nsdd esNestedDisMax
	)
	for _, mAggs := range p.moduleAggs {
		mm.Wrap.Query = mAggs
		mm.Wrap.Type = "cross_fields"
		mm.Wrap.Fields = []string{"module.name"}
		// query.Query.Bool.Must = append(query.Query.Bool.Must, mm)
		// query.Query.DisMax.Queries = append(query.Query.DisMax.Queries, mm)

		// dd.Wrap.Queries = append(dd.Wrap.Queries, mm)
		mdd.Wrap.Path = "module"
		nsdd.Wrap.IgnoreUnmapped = true
		mdd.Wrap.Query.Wrap.Queries = append(mdd.Wrap.Query.Wrap.Queries, mm)
	}

	if len(mdd.Wrap.Query.Wrap.Queries) > 0 {
		query.Query.Bool.Must = append(query.Query.Bool.Must, mdd)
	}

	// no need now since we are adding below as filter
	// if p.aggOnly {
	for _, nAggs := range p.namespaceAggs {
		mm.Wrap.Query = nAggs
		mm.Wrap.Type = "cross_fields"
		mm.Wrap.Fields = []string{"namespace.name"}
		// query.Query.Bool.Must = append(query.Query.Bool.Must, mm)
		// query.Query.DisMax.Queries = append(query.Query.DisMax.Queries, mm)

		// dd.Wrap.Queries = append(dd.Wrap.Queries, mm)
		nsdd.Wrap.Path = "namespace"
		nsdd.Wrap.IgnoreUnmapped = true
		nsdd.Wrap.Query.Wrap.Queries = append(nsdd.Wrap.Query.Wrap.Queries, mm)
	}
	// }

	if len(nsdd.Wrap.Query.Wrap.Queries) > 0 {
		query.Query.Bool.Must = append(query.Query.Bool.Must, nsdd)
	}

	// if !p.aggOnly && !noNSFilter {
	// 	nsf := make(map[string]interface{})
	// 	nsf["terms"] = map[string][]string{
	// 		"namespace.name.keyword": p.namespaceAggs,
	// 	}
	// 	query.Query.Bool.Filter = append(query.Query.Bool.Filter, nsf)
	// }

	// Aggregations V1.0 Improved
	// if len(p.aggregations) > 0 {
	//	for _, a := range p.aggregations {
	//		if len(a) > 0 {
	//			sqs = esSimpleQueryString{}
	//			sqs.Wrap.Query = a
	//			query.Query.Bool.Must = append(query.Query.Bool.Must, sqs)
	//		}
	//	}
	// }

	// if noQ == 0 && len(p.moduleAggs) == 0 && len(p.namespaceAggs) == 0 {
	//	query.Query.DisMax.Queries = append(query.Query.DisMax.Queries, index)
	// }
	query.Aggregations = make(map[string]esSearchNestedAggr)
	query.Aggregations["namespace"] = esSearchNestedAggr{
		Nested: esSearchNestedAggs{Path: "namespace"},
		Aggregations: map[string]esSearchNestedAggrTerm{
			"namespace": {
				Terms: esSearchAggrTerm{
					// Field: "namespace.name.keyword",
					Field: "namespace.handle",
					Size:  999,
				},
			},
		},
	}

	if !noQ || !noNSFilter {
		query.Aggregations["module"] = esSearchNestedAggr{
			Nested: esSearchNestedAggs{Path: "module"},
			Aggregations: map[string]esSearchNestedAggrTerm{
				"module": {
					Terms: esSearchAggrTerm{
						// Field: "module.name.keyword",
						Field: "module.handle",
						Size:  999,
					},
				},
			},
		}
	}

	// query.Aggregations["resource"] = esSearchAggr{
	//	Terms: esSearchAggrTerm{
	//		Field: "resourceType.keyword",
	//		Size:  999,
	//	},
	//	Aggregations: EsSearchAggrTerms{
	//		"resourceName": esSearchAggr{
	//			Terms: esSearchAggrTerm{
	//				Field: "name.keyword",
	//				Size:  999,
	//			},
	//		},
	//		"modules": esSearchAggr{
	//			Terms: esSearchAggrTerm{
	//				Field: "module.name.keyword",
	//				Size:  999,
	//			},
	//		},
	//		"namespaces": esSearchAggr{
	//			Terms: esSearchAggrTerm{
	//				Field: "namespace.name.keyword",
	//				Size:  999,
	//			},
	//		},
	//	},
	// }

	// Aggregations V2.0
	// if len(p.aggregations) > 0 {
	//	query.Aggregations = (Aggregations{}).encodeTerms(p.aggregations)
	// }

	if err = json.NewEncoder(&buf).Encode(query); err != nil {
		err = fmt.Errorf("could not encode query: %q", err)
		return
	}

	log.Debug("searching ",
		zap.String("for", p.title),
		// zap.String("open search body", buf.String()),
	)

	// Why set size to 999? default value for size is 10,
	// so we needed to set value till we add (@todo) pagination to search result
	if p.from > 0 && p.from >= p.size {
		p.size = p.from + 100
	}
	if p.size == 0 {
		p.size = 999
	}
	page.From = p.from
	page.Size = p.size

	sReqArgs := []func(*esapi.SearchRequest){
		esc.Search.WithContext(ctx),
		esc.Search.WithBody(&buf),
		esc.Search.WithTrackTotalHits(true),
		// esc.Search.WithScroll(),
		esc.Search.WithSize(p.size),
		esc.Search.WithFrom(p.from),
		// esc.Search.WithExplain(true), // debug
	}

	if p.dumpRaw {
		sReqArgs = append(
			sReqArgs,
			esc.Search.WithSourceExcludes("security"),
			esc.Search.WithPretty(),
		)
	}

	// Perform the search request.
	res, err := esc.Search(sReqArgs...)

	if err != nil {
		return
	}

	if err = validElasticResponse(res, err); err != nil {
		err = fmt.Errorf("invalid search response: %w", err)
		return
	}

	defer res.Body.Close()

	if p.dumpRaw {
		// Copy body buf and then restore it
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		os.Stdout.Write(bodyBytes)
	}

	sr = &esSearchResponse{}
	if err = json.NewDecoder(res.Body).Decode(sr); err != nil {
		return
	}

	// Print the response status, number of results, and request duration.
	log.Debug("search completed",
		zap.String("for", p.title),
		zap.String("query", sqs.Wrap.Query),
		zap.String("indexPrefix", index.Prefix.Index.Value),
		zap.Int("from", p.from),
		zap.Int("size", p.size),
		zap.String("status", res.Status()),
		zap.Int("took", sr.Took),
		zap.Bool("timedOut", sr.TimedOut),
		zap.Int("hits", sr.Hits.Total.Value),
		zap.String("hitsRelation", sr.Hits.Total.Relation),
		zap.Int("namespaceAggs", len(sr.Aggregations.Namespace.Namespace.Buckets)),
		zap.Int("moduleAggs", len(sr.Aggregations.Module.Module.Buckets)),
	)

	return
}

// @todo move this to es service
func validElasticResponse(res *esapi.Response, err error) error {
	if err != nil {
		return fmt.Errorf("failed to get response from search backend: %w", err)
	}

	if res.IsError() {
		defer res.Body.Close()
		var rsp struct {
			Error struct {
				Type   string
				Reason string
			}
		}

		if err := json.NewDecoder(res.Body).Decode(&rsp); err != nil {
			return fmt.Errorf("could not parse response body: %w", err)
		} else {
			return fmt.Errorf("search backend responded with an error: %s (type: %s, status: %s)", rsp.Error.Reason, rsp.Error.Type, res.Status())
		}
	}

	return nil
}

func extractFromSubClaim(sub string) (userID uint64, rr []uint64) {
	parts := strings.Split(sub, " ")
	rr = make([]uint64, len(parts)-1)
	for p := range parts {
		id, _ := strconv.ParseUint(parts[p], 10, 64)
		if p == 0 {
			userID = id
		} else {
			rr[p-1] = id
		}
	}

	return
}
