package query

import (
	"fmt"
)

// Wildcard leaf value: { "value": "..." }
type WildcardValue struct {
	Value string `json:"value"`
}

// Wildcard query: { "wildcard": { "<field>": { "value": "..." } } }
type WildcardQuery struct {
	Wildcard map[string]WildcardValue `json:"wildcard"`
}

// TermQuery: { "term": { "<field>": "<value>" } }
type TermQuery struct {
	Term map[string]string `json:"term"`
}

// Terms query: { "terms": { "<field>": [ ... ] } }
type TermsQuery struct {
	Terms map[string][]string `json:"terms"`
}

type TermsIntQuery struct {
	Terms map[string][]int64 `json:"terms"`
}

// Bool query wrapper: { "bool": { "should": [...], "minimum_should_match": 1 } }
type BoolQuery struct {
	Should             []interface{} `json:"should,omitempty"`
	MinimumShouldMatch int           `json:"minimum_should_match,omitempty"`
}

type BoolQueryWrapper struct {
	Bool BoolQuery `json:"bool"`
}

// BoolMust query wrapper: { "bool": { "must": [...] } }
type BoolMustQuery struct {
	Must []interface{} `json:"must"`
}

type BoolMustWrapper struct {
	Bool BoolMustQuery `json:"bool"`
}

// Generic nested query wrapper. Query is interface{} so you can pass
// a BoolQueryWrapper, TermsQuery or any other typed query structure.
type NestedQuery struct {
	Nested nestedQueryBody `json:"nested"`
}

type nestedQueryBody struct {
	Path           string      `json:"path,omitempty"`
	ScoreMode      string      `json:"score_mode,omitempty"`
	Query          interface{} `json:"query,omitempty"`
	IgnoreUnmapped bool        `json:"ignore_unmapped,omitempty"`
}

// KNN / vector query structure
type VectorKNN struct {
	Vector []float64   `json:"vector"`
	K      int         `json:"k"`
	Filter interface{} `json:"filter,omitempty"`
}

type KNNQuery struct {
	KNN map[string]VectorKNN `json:"knn"`
}

func NewWildcard(field, value string) WildcardQuery {
	return WildcardQuery{
		Wildcard: map[string]WildcardValue{
			field: {Value: value},
		},
	}
}

func NewBoolMust(must []interface{}) BoolMustWrapper {
	return BoolMustWrapper{
		Bool: BoolMustQuery{
			Must: must,
		},
	}
}

func NewBoolShould(should []interface{}, minShouldMatch int) BoolQueryWrapper {
	return BoolQueryWrapper{
		Bool: BoolQuery{
			Should:             should,
			MinimumShouldMatch: minShouldMatch,
		},
	}
}

func NewNestedBool(path string, boolWrap BoolQueryWrapper, ignoreUnmapped bool) NestedQuery {
	return NestedQuery{
		Nested: nestedQueryBody{
			Path:           path,
			Query:          boolWrap,
			IgnoreUnmapped: ignoreUnmapped,
		},
	}
}

func NewTerm(field, value string) TermQuery {
	return TermQuery{
		Term: map[string]string{field: value},
	}
}

func NewNestedTerms(path string, field string, ids []int64, scoreMode string) NestedQuery {
	tq := TermsIntQuery{Terms: map[string][]int64{field: ids}}
	return NestedQuery{
		Nested: nestedQueryBody{
			Path:      path,
			ScoreMode: scoreMode,
			Query:     tq,
		},
	}
}

func NewTermsString(field string, vals []string) TermsQuery {
	return TermsQuery{Terms: map[string][]string{field: vals}}
}

func NewKNN(vector []float64, k int) KNNQuery {
	return KNNQuery{
		KNN: map[string]VectorKNN{
			"vectorsValue": {Vector: vector, K: k},
		},
	}
}

func WildcardValueForQ(q string) string {
	return fmt.Sprintf("*%s*", q)
}

func NewIndexScopedExclusion(indexName string, fields []string, value string) BoolMustWrapper {
	shouldClauses := make([]interface{}, 0, len(fields))
	for _, field := range fields {
		wc := NewWildcard(field, WildcardValueForQ(value))
		shouldClauses = append(shouldClauses, wc)
	}

	nestedBool := NewBoolShould(shouldClauses, 1)
	nestedQuery := NewNestedBool("values", nestedBool, true)
	indexTerm := NewTerm("_index", indexName)

	return NewBoolMust([]interface{}{
		indexTerm,
		nestedQuery,
	})
}
