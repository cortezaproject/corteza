package mysql

import "github.com/cortezaproject/corteza-server/pkg/ql"

type (
	// QueryEncoder provides query parts encoding rules for MySQL
	// see ql.QueryEncoder for mor info
	QueryEncoder struct{}
)

var _ ql.Encoder = &QueryEncoder{}

func (QueryEncoder) CaseInsensitiveLike(neg bool) string {
	if neg {
		return "COLLATE utf8_unicode_ci NOT LIKE"
	} else {
		return "COLLATE utf8_unicode_ci LIKE"
	}
}
