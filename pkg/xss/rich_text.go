package xss

import (
	"html"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

// RichText assures safe HTML content
func RichText(in string) string {
	// use standard html escaping policy
	p := bluemonday.UGCPolicy()

	// match only colors for html editor elements on style attr
	p.AllowAttrs("style").OnElements("span", "p")
	p.AllowStyles("color").Matching(regexp.MustCompile("(?i)^#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})$")).Globally()
	p.AllowStyles("background-color").Matching(regexp.MustCompile("(?i)^#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})$")).Globally()

	// some link specifics we need; allow target but assure safety
	p.AllowAttrs("target").OnElements("a")
	p.AddTargetBlankToFullyQualifiedLinks(false)

	sanitized := p.Sanitize(in)

	// handle escaped strings and unescape them
	// all the dangerous chars should have been stripped
	// by now
	return html.UnescapeString(sanitized)
}
