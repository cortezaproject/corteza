package tools

import (
	"os"
	"path/filepath"
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

// TestWebappURLPrefersTheCheckoutsOwnEnv pins the case a worktree creates: each
// one runs vite on its own port and gets an .env.e2e naming it, while
// playwright.config.ts holds the same literal 5173 fallback in every checkout.
// Reading the config alone answers 5173 everywhere, so a browser check run
// inside a worktree exercises the primary's webapp and reports a pass for code
// it never loaded.
func TestWebappURLPrefersTheCheckoutsOwnEnv(t *testing.T) {
	const config = "baseURL: process.env.E2E_BASE_URL || 'http://localhost:5173',\n"

	cases := []struct {
		name   string
		env    string
		config string
		want   string
	}{
		{
			name:   "worktree env wins over the config fallback",
			env:    "E2E_BASE_URL=http://localhost:5176\nE2E_USER=agent@local.dev\n",
			config: config,
			want:   "http://localhost:5176",
		},
		{
			name:   "config fallback when there is no env file",
			config: config,
			want:   "http://localhost:5173",
		},
		{
			name: "built-in default when there is neither",
			want: "http://localhost:5173",
		},
		{
			name:   "an empty assignment does not shadow the config",
			env:    "E2E_BASE_URL=\n",
			config: config,
			want:   "http://localhost:5173",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()
			unify := filepath.Join(root, "client", "web", "unify")
			if err := os.MkdirAll(unify, 0o755); err != nil {
				t.Fatal(err)
			}
			if c.env != "" {
				if err := os.WriteFile(filepath.Join(unify, ".env.e2e"), []byte(c.env), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if c.config != "" {
				if err := os.WriteFile(filepath.Join(unify, "playwright.config.ts"), []byte(c.config), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			if got := webappURL(root); got != c.want {
				t.Fatalf("webappURL = %q, want %q", got, c.want)
			}
		})
	}
}
