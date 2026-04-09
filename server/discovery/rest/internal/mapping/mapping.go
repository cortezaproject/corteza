package mapping

import "github.com/cortezaproject/corteza/server/pkg/options"

type (
	Mapping struct {
		Index   string               `json:"index"`
		Mapping map[string]*property `json:"mapping"`
	}

	property struct {
		// https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping-types.html
		Type string `json:"type,omitempty"`

		// Boost factor
		// https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping-boost.html
		Boost float32 `json:"boost,omitempty"`

		Analyzer string `json:"analyzer,omitempty"`

		Properties map[string]*property `json:"properties,omitempty"`

		// for vector indexing
		Dimension int            `json:"dimension,omitempty"`
		Method    map[string]any `json:"method,omitempty"`
	}

	Context struct {
		AccessRestriction string
	}
)

func change() *property {
	return &property{
		Type: "nested",
		Properties: map[string]*property{
			"at": {Type: "date"},
			"by": {Type: "nested", Properties: map[string]*property{
				"id":     {Type: "long"},
				"email":  {Type: "keyword"},
				"name":   {Type: "keyword"},
				"handle": {Type: "keyword"},
			}},
		},
	}
}

func security() *property {
	return &property{
		Properties: map[string]*property{
			"allowedRoles": {Type: "long"},
			"deniedRoles":  {Type: "long"},
		},
	}
}

func vector(opts options.DiscoveryOpt) *property {
	return &property{
		Type:      "knn_vector",
		Dimension: opts.EmbeddingsDimension,
		Method: map[string]any{
			"name":   "hnsw",
			"engine": "lucene",
			"parameters": map[string]any{
				"ef_construction": opts.HnswEfConstruction,
				"m":               opts.HnswM,
			},
		},
	}
}
