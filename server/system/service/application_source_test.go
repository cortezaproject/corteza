package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleApp = `<!doctype html>
<html><body><div id="app"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/chart.js/4.4.1/chart.umd.min.js"></script>
<script>
window.SAMPLE = { 'records.list': () => ({ records: [], refs: {} }) }
window.human = window.human || (() => { return {} })()
</script>
</body></html>`

func TestApplicationSourceAcceptsAPlainDocument(t *testing.T) {
	require.NoError(t, CheckApplicationSource(sampleApp, nil))
}

func TestApplicationSourceRefusesWhatTheSandboxCannotRun(t *testing.T) {
	cases := []struct {
		name    string
		source  string
		message string
	}{
		{"module script", `<script type="module">const a = 1</script>`, "script-src"},
		{"single-quoted module script", "<script type='module'>const a = 1</script>", "script-src"},
		{"import statement", "<script>\nimport { h } from 'https://esm.sh/preact'\n</script>", "line 2"},
		{"export statement", "<script>\nexport const render = () => {}\n</script>", "ES module statement"},
		{"react import", `<script>const { useState } = from "react"</script>`, "React"},
		{"fetch", `<script>fetch('/api/compose').then(r => r.json())</script>`, "connect-src"},
		{"window fetch", `<script>window.fetch('/api/compose')</script>`, "connect-src"},
		{"xhr", `<script>const r = new XMLHttpRequest()</script>`, "connect-src"},
		{"websocket", `<script>const s = new WebSocket('wss://example.com')</script>`, "connect-src"},
		{"local storage", `<script>try { localStorage.setItem('tab', 'a') } catch {}</script>`, "SecurityError"},
		{"session storage", `<script>sessionStorage.getItem('x')</script>`, "SecurityError"},
		{"indexed db", `<script>indexedDB.open('db')</script>`, "SecurityError"},
		{"cookie", `<script>document.cookie = 'a=1'</script>`, "SecurityError"},
		{"alert", `<script>alert('saved')</script>`, "in-page dialog"},
		{"window confirm", `<script>if (window.confirm('Sure?')) go()</script>`, "always answers false"},
		{"prompt", `<script>const n = prompt('Name?')</script>`, "in-page dialog"},
		{"window open", `<script>window.open('https://example.com')</script>`, "no new windows"},
		{"new tab link", `<a href="https://example.com" target="_blank">x</a>`, "no new windows"},
		{"download link", `<a href="data:text/csv,a" download="x.csv">Export</a>`, "download"},
		{"download property", `<script>a.download = 'x.csv'; a.click()</script>`, "download"},
		{"script elsewhere", `<script src="https://cdn.jsdelivr.net/npm/chart.js@4"></script>`, "cdn.jsdelivr.net"},
		{"protocol-relative script", `<script src="//unpkg.com/x"></script>`, "cdnjs"},
		{"web font", `<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter">`, "system font"},
		{"css import", `<style>@import url('https://fonts.googleapis.com/css2?family=Inter');</style>`, "inline it"},
		{"css url", `<style>.hero { background: url(https://example.com/a.png) }</style>`, "data: URL"},
		{"remote image", `<img src="https://example.com/logo.png" alt="">`, "allowed origins"},
		{"mailto link", `<a href="mailto:ana@example.com">Ana</a>`, "leaves the page does nothing"},
		{"tel link", `<a href='tel:+38612345'>call</a>`, "tel:+38612345"},
		{"site link", `<a class="x" href=https://example.com>site</a>`, "leaves the page"},
		{"file link", `<a href="details.html">more</a>`, "details.html"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := CheckApplicationSource(c.source, nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(), c.message)
		})
	}
}

// A name that merely ends in one of the refused ones is somebody's own
// function, and refusing it would send the caller looking for a call it never
// wrote.
func TestApplicationSourceAllowsIdentifiersThatOnlyEndInARefusedName(t *testing.T) {
	require.NoError(t, CheckApplicationSource(`<script>const prefetch = () => {}; prefetch()</script>`, nil))
	require.NoError(t, CheckApplicationSource(`<script>function importRows () {}</script>`, nil))
}

// What only looks like a refused call is the page's own: a method of that name
// on some object, a library from the one admitted host, an image carried as
// data.
func TestApplicationSourceAllowsWhatOnlyResemblesARefusal(t *testing.T) {
	for _, source := range []string{
		`<script>dialog.confirm(); modal.alert(); form.prompt()</script>`,
		`<script>function showAlert () {}; showAlert()</script>`,
		`<script>panel.open(); const open = () => {}; open()</script>`,
		`<script>row.downloadCount = 3; if (a.download == b) {}</script>`,
		`<script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/4.4.1/chart.umd.min.js"></script>`,
		`<img src="data:image/png;base64,AAAA" alt="">`,
		`<a href="#detail">Details</a>`,
		`<a href="" onclick="open()">Open</a>`,
		`<a href="javascript:void(0)">Toggle</a>`,
		`<a name="top"></a>`,
	} {
		require.NoError(t, CheckApplicationSource(source, nil), source)
	}
}

func TestApplicationSourceRefusesAnEmptyDocument(t *testing.T) {
	err := CheckApplicationSource("   \n\t ", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

// The service caps the source too; failing here is what keeps a 300 KB argument
// from being carried through the whole write before it is refused.
func TestApplicationSourceRefusesOverTheSizeCap(t *testing.T) {
	oversize := "<html>" + strings.Repeat("x", ApplicationSourceMaxSize) + "</html>"

	err := CheckApplicationSource(oversize, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "the limit is")

	require.NoError(t, CheckApplicationSource("<html>"+strings.Repeat("x", ApplicationSourceMaxSize-20)+"</html>", nil))
}

// The skill hands out a bridge snippet that every custom app pastes in. If the
// guard refused it, the tool would reject the first app written to the skill.

func TestNormalizeSourceOrigins(t *testing.T) {
	got, err := NormalizeSourceOrigins([]string{" https://Fonts.googleapis.com/ ", "https://fonts.gstatic.com", "https://fonts.gstatic.com", "https://cdnjs.cloudflare.com", "https://cdn.test:8443", ""})
	require.NoError(t, err)
	require.Equal(t, []string{"https://fonts.googleapis.com", "https://fonts.gstatic.com", "https://cdn.test:8443"}, got)

	for raw, want := range map[string]string{
		"http://cdn.test":        "not https",
		"https://*.cdn.test":     "wildcard",
		"https://cdn.test/npm/x": "more than scheme and host",
		"https://u:p@cdn.test":   "more than scheme and host",
		"cdn.test":               "not an address",
	} {
		_, err := NormalizeSourceOrigins([]string{raw})
		require.ErrorContains(t, err, want, raw)
	}
}

func TestCheckApplicationSourceAdmitsListedOrigins(t *testing.T) {
	origins := []string{"https://fonts.googleapis.com", "https://img.test"}

	require.NoError(t, CheckApplicationSource(`<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter">`, origins))
	require.NoError(t, CheckApplicationSource(`<img src="https://img.test/logo.png">`, origins))
	require.NoError(t, CheckApplicationSource(`<style>@import url("https://fonts.googleapis.com/css2?family=Inter");</style>`, origins))

	require.ErrorContains(t, CheckApplicationSource(`<img src="https://other.test/x.png">`, origins), "allowed origins")
	require.ErrorContains(t, CheckApplicationSource(`<script src="https://img.test.evil.example/x.js"></script>`, origins), "allowed origins")
	// fetch stays closed whatever is listed.
	require.ErrorContains(t, CheckApplicationSource(`<script>fetch("https://img.test/x")</script>`, origins), "connect-src")
}
