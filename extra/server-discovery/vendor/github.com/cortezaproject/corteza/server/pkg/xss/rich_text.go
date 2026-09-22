package xss

import (
	"html"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

var (
	p *bluemonday.Policy
)

func init() {
	// use standard html escaping policy
	p = bluemonday.UGCPolicy()

	// match only colors for html editor elements on style attr
	p.AllowAttrs("style").OnElements("span", "p")
	// Support both hex colors (#ff0000, #f00) and RGB colors (rgb(255, 0, 0))
	hexPattern := regexp.MustCompile(`(?i)^#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})$`)
	rgbPattern := regexp.MustCompile(`(?i)^rgb\(\s*\d+\s*,\s*\d+\s*,\s*\d+\s*\)$`)
	p.AllowStyles("color").Matching(hexPattern).Globally()
	p.AllowStyles("color").Matching(rgbPattern).Globally()
	p.AllowStyles("background-color").Matching(hexPattern).Globally()
	p.AllowStyles("background-color").Matching(rgbPattern).Globally()

	// allow text alignment
	p.AllowStyles("text-align").Matching(regexp.MustCompile("^(left|center|right|justify)$")).OnElements("span", "p", "h1", "h2", "h3", "h4", "h5", "h6", "h7")

	// allow task lists (TipTap format)
	p.AllowAttrs("data-type").Matching(regexp.MustCompile(`^(taskList|taskItem|mention)$`)).OnElements("ul", "li", "span")
	p.AllowAttrs("data-checked").Matching(regexp.MustCompile(`^(true|false)$`)).OnElements("li")
	p.AllowAttrs("checked").OnElements("input")
	p.AllowAttrs("type").Matching(regexp.MustCompile(`^checkbox$`)).OnElements("input")
	p.AllowAttrs("aria-label").OnElements("input")
	p.AllowAttrs("contenteditable").Matching(regexp.MustCompile(`^(true|false)$`)).OnElements("span", "div", "label")
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^(todo-checkbox|todo-content|ProseMirror-trailingBreak|mention|ProseMirror-selectednode|ProseMirror-gapcursor)$`)).OnElements("span", "div", "label")

	// allow mention attributes
	p.AllowAttrs("data-id").Matching(regexp.MustCompile(`^\d+$`)).OnElements("span")
	p.AllowAttrs("data-label").OnElements("span")
	p.AllowAttrs("data-mention-suggestion-char").Matching(regexp.MustCompile(`^@$`)).OnElements("span")

	// allow table attributes (TableKit extension)
	p.AllowAttrs("colspan").Matching(regexp.MustCompile(`^\d+$`)).OnElements("td", "th")
	p.AllowAttrs("rowspan").Matching(regexp.MustCompile(`^\d+$`)).OnElements("td", "th")
	p.AllowAttrs("scope").Matching(regexp.MustCompile(`^(col|row|colgroup|rowgroup)$`)).OnElements("th")

	// some link specifics we need; allow target but assure safety
	p.AllowAttrs("target").OnElements("a")
	p.AddTargetBlankToFullyQualifiedLinks(false)
}

// RichText assures safe HTML content
func RichText(in string) string {
	sanitized := p.Sanitize(in)

	// handle escaped strings and unescape them
	// all the dangerous chars should have been stripped
	// by now
	return html.UnescapeString(sanitized)
}
