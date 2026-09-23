package service

// The sandbox a custom application runs in refuses a good deal of ordinary
// HTML, silently: a page that reaches for the network, for storage, for a
// dialog or for another document simply does nothing there. Refusing such a
// page while it is being stored is the only moment its author is still around
// to hear about it, so this runs for every way a page arrives — the MCP tool,
// the REST endpoint, the editor.

import (
	"fmt"
	"regexp"
	"strings"
)

// checkApplicationSource refuses a document the app sandbox cannot run.
//
// The sandbox serves the HTML with a CSP of default-src 'none', connect-src
// 'none' and script-src limited to inline script and cdnjs: an ES module never
// loads, and every network call fails with nothing to catch it on. Both faults
// surface as a blank frame in front of a user, minutes after this call
// returned success, so they are refused here where the wording can say what to
// write instead.
func CheckApplicationSource(source string) error {
	if strings.TrimSpace(source) == "" {
		return fmt.Errorf("the source is empty; a custom application is one whole HTML document")
	}

	if len(source) > ApplicationSourceMaxSize {
		return fmt.Errorf(
			"the source is %d bytes and the limit is %d; keep a custom application to a single document with inline script, well under that",
			len(source), ApplicationSourceMaxSize,
		)
	}

	const plainHTML = "; write one plain HTML document with inline <script>, loading libraries from https://cdnjs.cloudflare.com"

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

	for _, m := range externalScript.FindAllStringSubmatch(source, -1) {
		if !strings.HasPrefix(strings.ToLower(m[1]), "https://cdnjs.cloudflare.com/") {
			return fmt.Errorf("the source loads a script from %s; the sandbox's script-src admits https://cdnjs.cloudflare.com only — load the library from there, pinned to an exact version", m[1])
		}
	}

	if m := externalLink.FindStringSubmatch(source); m != nil {
		return fmt.Errorf("the source links %s; the sandbox loads no external stylesheet or font — inline the CSS in <style> and use a system font stack", m[1])
	}

	if m := externalCSSURL.FindStringSubmatch(source); m != nil {
		return fmt.Errorf("the source's CSS reaches %s; the sandbox loads nothing from outside — inline it, or use a data: URL", m[1])
	}

	if m := externalImage.FindStringSubmatch(source); m != nil {
		return fmt.Errorf("the source shows the image %s; the sandbox's img-src admits data: URLs only — embed it as data: or drop it", m[1])
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
