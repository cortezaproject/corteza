package locale

import (
	"html"

	"github.com/microcosm-cc/bluemonday"
)

var (
	stripHtml = bluemonday.StripTagsPolicy().Sanitize
)

func SanitizeMessage(in string) string {
	return html.UnescapeString(stripHtml(in))
}
