package tools

import (
	"strings"
	"testing"
)

// TestUINoteReportsLeavingTheRequestedPath pins the case this check exists for:
// the webapp's router redirects a path it does not recognise to the home
// section, which renders cleanly, so a mistyped path reports a pass with a
// screenshot of a page the caller never asked about unless the note catches it.
func TestUINoteReportsLeavingTheRequestedPath(t *testing.T) {
	cases := []struct {
		name string
		out  uiResult
		want string // substring the note must carry; empty means it must not mention the path
	}{
		{
			name: "unmatched path redirected home",
			out: uiResult{
				RequestedPath: "/compose/ns/test-crm/pages/123",
				LandedPath:    "/",
			},
			want: "did not stay on /compose/ns/test-crm/pages/123",
		},
		{
			// Ranked above console errors: they describe the wrong page.
			name: "console errors on a page that was never reached",
			out: uiResult{
				RequestedPath: "/compose/namespace/x/pages/1",
				LandedPath:    "/",
				ConsoleErrors: []string{"TypeError: x is not a function"},
			},
			want: "did not stay on",
		},
		{
			name: "arrived, with a query string the app appended",
			out: uiResult{
				RequestedPath: "/compose/namespace/x/pages/1",
				LandedPath:    "/compose/namespace/x/pages/1?limit=20",
			},
		},
		{
			name: "arrived, trailing slash only",
			out: uiResult{
				RequestedPath: "/compose/namespaces",
				LandedPath:    "/compose/namespaces/",
			},
		},
		{
			// A failed step is the more specific finding and keeps its place.
			name: "failed step outranks the path",
			out: uiResult{
				RequestedPath: "/a",
				LandedPath:    "/",
				Steps:         []uiStep{{Selector: "button", OK: false}},
			},
			want: "failed",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			note := uiNote(c.out)

			if c.want == "" {
				if strings.Contains(note, "did not stay on") {
					t.Errorf("note claims the app left the path when it arrived: %q", note)
				}
				return
			}

			if !strings.Contains(note, c.want) {
				t.Errorf("note %q does not mention %q", note, c.want)
			}
		})
	}
}
