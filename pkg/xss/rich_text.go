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

	// allow text alignment
	p.AllowStyles("text-align").Matching(regexp.MustCompile("^(left|center|right|justify)$")).OnElements("span", "p")

	// allow checklists
	p.AllowAttrs("data-type").Matching(regexp.MustCompile("^(todo_list|todo_item)$")).OnElements("ul", "li")
	p.AllowAttrs("data-done").Matching(regexp.MustCompile("^(true|false)$")).OnElements("ul", "li")
	p.AllowAttrs("data-drag-handle").OnElements("li")
	p.AllowAttrs("contenteditable").OnElements("span", "div")
	p.AllowAttrs("class").Matching(regexp.MustCompile("^(todo-checkbox|todo-content)$"))

	// some link specifics we need; allow target but assure safety
	p.AllowAttrs("target").OnElements("a")
	p.AddTargetBlankToFullyQualifiedLinks(false)

	sanitized := p.Sanitize(in)

	// handle escaped strings and unescape them
	// all the dangerous chars should have been stripped
	// by now
	return html.UnescapeString(sanitized)
}
