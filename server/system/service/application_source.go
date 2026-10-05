package service

// The sandbox a custom application runs in refuses a good deal of ordinary
// HTML, silently: a page that reaches for the network, for storage, for a
// dialog or for another document simply does nothing there. Refusing such a
// page while it is being stored is the only moment its author is still around
// to hear about it, so this runs for every way a page arrives — the MCP tool,
// the REST endpoint, the editor.

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"
)

// The origin every page may load scripts from, listed or not.
const sourceDefaultOrigin = "https://cdnjs.cloudflare.com"

// SourceOriginsMax caps how many origins one page lists.
const SourceOriginsMax = 20

// NormalizeSourceOrigins reduces what a page lists to exact https origins —
// scheme, host and port, nothing else — without repeats. Anything else is
// refused, saying what to write: an origin is a place the sandbox will let
// the page reach, so it is never a pattern.
func NormalizeSourceOrigins(in []string) ([]string, error) {
	out := make([]string, 0, len(in))
	for _, raw := range in {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		u, err := url.Parse(raw)
		switch {
		case err != nil || u.Host == "":
			return nil, fmt.Errorf("the origin %q is not an address; write one as https://host, for example https://cdn.jsdelivr.net", raw)
		case u.Scheme != "https":
			return nil, fmt.Errorf("the origin %q is not https; the sandbox loads nothing over plain http", raw)
		case strings.Contains(u.Host, "*"):
			return nil, fmt.Errorf("the origin %q has a wildcard; list each host on its own", raw)
		case u.User != nil || (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "":
			return nil, fmt.Errorf("the origin %q carries more than scheme and host; write it as https://%s", raw, u.Host)
		}

		origin := "https://" + strings.ToLower(u.Host)
		if origin == sourceDefaultOrigin || slices.Contains(out, origin) {
			continue
		}
		out = append(out, origin)
	}

	if len(out) > SourceOriginsMax {
		return nil, fmt.Errorf("the page lists %d origins and the limit is %d", len(out), SourceOriginsMax)
	}

	return out, nil
}

// originAllowed reports whether a URL the page loads is on an origin it may
// reach: cdnjs, or one it lists.
func originAllowed(raw string, origins []string) bool {
	if strings.HasPrefix(raw, "//") {
		raw = "https:" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return false
	}
	origin := "https://" + strings.ToLower(u.Host)
	return origin == sourceDefaultOrigin || slices.Contains(origins, origin)
}

// checkApplicationSource refuses a document the app sandbox cannot run.
//
// The sandbox serves the HTML with a CSP of default-src 'none', connect-src
// 'none' and script, style, font and image sources limited to inline content,
// cdnjs (scripts) and the origins the page lists: an ES module never
// loads, and every network call fails with nothing to catch it on. Both faults
// surface as a blank frame in front of a user, minutes after this call
// returned success, so they are refused here where the wording can say what to
// write instead.
func CheckApplicationSource(source string, origins []string) error {
	if strings.TrimSpace(source) == "" {
		return fmt.Errorf("the source is empty; a custom application is one whole HTML document")
	}

	if len(source) > ApplicationSourceMaxSize {
		return fmt.Errorf(
			"the source is %d bytes and the limit is %d; keep a custom application to a single document with inline script, well under that",
			len(source), ApplicationSourceMaxSize,
		)
	}

	const plainHTML = "; write one plain HTML document with inline <script>, loading libraries from https://cdnjs.cloudflare.com or an origin the page lists"

	if strings.Contains(source, `<script type="module"`) || strings.Contains(source, "<script type='module'") {
		return fmt.Errorf(`the source has a <script type="module">, which the sandbox's script-src never loads` + plainHTML)
	}

	if strings.Contains(source, "from 'react'") || strings.Contains(source, `from "react"`) {
		return fmt.Errorf("the source imports React, which the sandbox cannot load and which needs a build step" + plainHTML)
	}

	for i, line := range strings.Split(source, "\n") {
		trimmed := strings.TrimSpace(line)
		for _, keyword := range []string{"import ", "export "} {
			if strings.HasPrefix(trimmed, keyword) {
				return fmt.Errorf("line %d is an ES module statement (%q), which the sandbox never evaluates%s", i+1, strings.TrimSpace(keyword), plainHTML)
			}
		}
	}

	for _, api := range []string{"fetch(", "XMLHttpRequest", "WebSocket"} {
		if usesIdentifier(source, api) {
			return fmt.Errorf(
				"the source uses %s, and the sandbox is served with connect-src 'none' — the call is blocked at runtime and returns nothing to render; read data through the bridge (human.records.list, human.records.read, human.records.report) instead",
				strings.TrimSuffix(api, "("),
			)
		}
	}

	for _, store := range []string{"localStorage", "sessionStorage", "indexedDB", "document.cookie"} {
		if usesIdentifier(source, store) {
			return fmt.Errorf(
				"the source uses %s; the sandbox's origin is opaque, so every storage API throws SecurityError there — keep state in page variables, it lasts as long as the page is open",
				store,
			)
		}
	}

	for _, dialog := range []string{"alert", "confirm", "prompt"} {
		if callsGlobal(source, dialog) {
			return fmt.Errorf(
				"the source calls %s(), which the sandbox suppresses without a word (confirm() always answers false, so what it guards never runs); show an in-page dialog instead",
				dialog,
			)
		}
	}

	if strings.Contains(source, "window.open(") {
		return fmt.Errorf("the source calls window.open(), and the sandbox allows no new windows; show the content in the page, for example in a detail pane")
	}

	if newTab.MatchString(source) {
		return fmt.Errorf(`the source has target="_blank", and the sandbox allows no new windows, so the link does nothing; show the content in the page instead`)
	}

	if downloadLink.MatchString(source) || downloadProperty.MatchString(source) {
		return fmt.Errorf("the source offers a file through a download link, which the sandbox blocks; hand the file to human.download(name, text) instead and Human saves it")
	}

	for _, m := range linkHref.FindAllStringSubmatch(source, -1) {
		href := strings.TrimSpace(m[1] + m[2] + m[3])
		if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(strings.ToLower(href), "javascript:") {
			continue
		}
		return fmt.Errorf(
			"the source links to %s; in Human a link that leaves the page does nothing, though it works in a preview — the app is one document, so link only to #anchors within it, and show an email address or phone number as text",
			href,
		)
	}

	const listIt = "add its origin to the page's allowed origins, or "

	for _, m := range externalScript.FindAllStringSubmatch(source, -1) {
		if !originAllowed(m[1], origins) {
			return fmt.Errorf("the source loads a script from %s; the sandbox admits scripts from https://cdnjs.cloudflare.com and the page's allowed origins only — %sload the library from cdnjs, pinned to an exact version", m[1], listIt)
		}
	}

	for _, m := range externalLink.FindAllStringSubmatch(source, -1) {
		if !originAllowed(m[1], origins) {
			return fmt.Errorf("the source links %s; the sandbox loads stylesheets and fonts only from the page's allowed origins — %sinline the CSS in <style> and use a system font stack", m[1], listIt)
		}
	}

	for _, m := range externalCSSURL.FindAllStringSubmatch(source, -1) {
		if !originAllowed(m[1], origins) {
			return fmt.Errorf("the source's CSS reaches %s; the sandbox loads nothing from outside the page's allowed origins — %sinline it, or use a data: URL", m[1], listIt)
		}
	}

	for _, m := range externalImage.FindAllStringSubmatch(source, -1) {
		if !originAllowed(m[1], origins) {
			return fmt.Errorf("the source shows the image %s; the sandbox shows images from data: URLs and the page's allowed origins only — %sembed it as data: or drop it", m[1], listIt)
		}
	}

	return nil
}

var (
	newTab           = regexp.MustCompile(`(?i)target\s*=\s*["']?_blank`)
	linkHref         = regexp.MustCompile(`(?i)<a\b[^>]*\shref\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s>]*))`)
	downloadLink     = regexp.MustCompile(`(?i)<a\b[^>]*\sdownload\b`)
	downloadProperty = regexp.MustCompile(`\.download\s*=[^=]`)
	externalScript   = regexp.MustCompile(`(?i)<script\b[^>]*\ssrc\s*=\s*["']?((?:https?:)?//[^"'\s>]+)`)
	externalLink     = regexp.MustCompile(`(?i)<link\b[^>]*\shref\s*=\s*["']?((?:https?:)?//[^"'\s>]+)`)
	externalCSSURL   = regexp.MustCompile(`(?i)(?:@import\s+(?:url\()?|url\()\s*["']?((?:https?:)?//[^"')\s]+)`)
	externalImage    = regexp.MustCompile(`(?i)<img\b[^>]*\ssrc\s*=\s*["']?((?:https?:)?//[^"'\s>]+)`)
)

// callsGlobal reports whether the source calls the global name — `name(` or
// `window.name(` — rather than a method of that name on some object of its own.
func callsGlobal(source, name string) bool {
	call := name + "("
	for i := 0; i < len(source); {
		at := strings.Index(source[i:], call)
		if at < 0 {
			return false
		}
		at += i
		switch {
		case at == 0:
			return true
		case strings.HasSuffix(source[:at], "window."):
			return true
		case source[at-1] != '.' && !isIdentifierByte(source[at-1]):
			return true
		}
		i = at + len(call)
	}
	return false
}

// usesIdentifier reports whether the source uses name as an identifier of its
// own rather than as the tail of a longer one — "prefetch(" is not "fetch(".
func usesIdentifier(source, name string) bool {
	for i := 0; i < len(source); {
		at := strings.Index(source[i:], name)
		if at < 0 {
			return false
		}
		at += i
		if at == 0 || !isIdentifierByte(source[at-1]) {
			return true
		}
		i = at + len(name)
	}
	return false
}

func isIdentifierByte(b byte) bool {
	return b == '_' || b == '$' ||
		(b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}
