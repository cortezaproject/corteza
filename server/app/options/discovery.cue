package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

discovery: schema.#optionsGroup & {
	handle: "discovery"
	options: {
		enabled: {
			type:          "bool"
			defaultGoExpr: "false"
			description:   "Enable discovery endpoints"
		}
		debug: {
			type:          "bool"
			defaultGoExpr: "false"
			description:   "Enable discovery related activity info"
		}
		corteza_domain: {
			type:        "string"
			description: "Indicates host of corteza compose webapp"
		}
		base_url: {
			type:        "string"
			description: "Indicates host of corteza discovery server"
		}
		embeddings_enabled: {
			type:          "bool"
			defaultGoExpr: "false"
			description:   "Enable discovery embeddings generation"
		}
		embeddings_dimension: {
      type:        "int"
      description: "Embeddings dimension"
      defaultGoExpr: "384"
    }
    hnsw_ef_construction: {
      type:        "int"
      description: "HNSW ef construction parameter"
      defaultGoExpr: "128"
    }
    hnsw_m: {
      type:        "int"
      description: "HNSW m parameter"
      defaultGoExpr: "16"
    }

	}
	title: "Discovery"
}
