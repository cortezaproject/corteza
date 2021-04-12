package cockroach

import "github.com/cortezaproject/corteza-server/pkg/ql"

type (
	// QueryEncoder provides query parts encoding rules for CockroachDB
	// see ql.QueryEncoder for mor info
	QueryEncoder struct{}
)

var _ ql.Encoder = &QueryEncoder{}

func (QueryEncoder) CaseInsensitiveLike(neg bool) string {
	if neg {
		return "NOT ILIKE"
	} else {
		return "ILIKE"
	}
}
