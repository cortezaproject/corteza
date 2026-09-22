package weburl

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Every path this package can emit, as the route pattern it must match. The
// left side is what a builder produces with the sample arguments below; the
// right side is the vue-router path that has to exist for it to resolve.
//
// This is the whole point of the package having a test: the paths are a copy of
// somebody else's routing table, and a copy nobody checks is a copy that rots.
// The evidence is in the tree — server/discovery emits /compose/ns/{slug}/...,
// a prefix unify has not routed since the apps were merged, so every link the
// indexer produces lands on the catch-all.
var table = []struct {
	name  string
	got   string
	route string
}{
	{"compose namespace", ComposeNamespace("crm"), "/compose/namespace/:slug"},
	{"compose namespace edit", ComposeNamespaceEdit("crm"), "/compose/namespaces/edit/:slug"},
	{"compose page", ComposePage("crm", 12), "/compose/namespace/:slug/pages/:pageID"},
	{"compose page builder", ComposePageBuilder("crm", 12), "/compose/namespace/:slug/admin/pages/:pageID/builder"},
	{"compose module edit", ComposeModuleEdit("crm", 12), "/compose/namespace/:slug/admin/modules/:moduleID/edit"},
	{"compose chart edit", ComposeChartEdit("crm", 12), "/compose/namespace/:slug/admin/charts/:chartID/edit"},
	{"compose record", ComposeRecord("crm", 12, 34), "/compose/namespace/:slug/admin/modules/:moduleID/records/:recordID"},
	{"compose record on page", ComposeRecordOnPage("crm", 12, 34), "/compose/namespace/:slug/pages/:pageID/records/:recordID"},
	{"taq", TAQ(12), "/taq/builder/:id"},
	{"workflow", Workflow(12), "/workflow/:workflowID/edit"},
	{"agent", Agent(12), "/agentic/:agentID/edit"},
	{"chatbot", Chatbot(12), "/chatbot/:chatbotID/edit"},
	{"user", User(12), "/admin/system/users/:userID"},
	{"user group", UserGroup(12), "/admin/system/user-groups/:userGroupID"},
	{"role", Role(12), "/admin/system/roles/:roleID"},
	{"application", Application(12), "/admin/system/applications/:applicationID"},
	{"custom application", CustomApplication(12), "/app/:applicationID"},
	{"auth client", AuthClient(12), "/admin/system/auth-clients/:authClientID"},
}

// sectionPrefix is what a section's index.js prepends to its own absolute
// paths. Compose and admin declare their children unprefixed and add it when
// the section is mounted, so the raw route files do not contain the full path.
var sectionPrefix = map[string]string{
	"compose": "/compose",
	"admin":   "/admin",
}

var routePath = regexp.MustCompile(`path:\s*'([^']*)'`)

// unifyRoutes collects every absolute path declared by the webapp's sections,
// with its section prefix applied.
//
// Caveat worth knowing before trusting a green run: the files read here are
// outside the Go module, so `go test` will not invalidate its cache when a
// route changes. Run this package with -count=1 when the question is whether a
// route moved.
func unifyRoutes(t *testing.T) map[string]bool {
	t.Helper()

	root := filepath.Join("..", "..", "..", "client", "web", "unify", "src", "sections")
	if _, err := os.Stat(root); err != nil {
		t.Skipf("webapp sources are not in this checkout: %v", err)
	}

	out := map[string]bool{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".js") || strings.Contains(path, ".test.") {
			return err
		}

		rel, _ := filepath.Rel(root, path)
		section := strings.Split(filepath.ToSlash(rel), "/")[0]

		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Compose declares its namespace sub-routes as RELATIVE children of
		// /namespace/:slug, so a path here is either absolute or hangs off the
		// last absolute one seen in the file — which is how the file is laid
		// out. This is a drift detector, not a vue-router parser: what matters
		// is that renaming a segment makes the literal disappear and the test
		// fail.
		last := ""
		for _, m := range routePath.FindAllStringSubmatch(string(body), -1) {
			p := m[1]
			if strings.HasPrefix(p, "/") {
				last = sectionPrefix[section] + p
				out[last] = true
				continue
			}
			if last == "" {
				continue
			}
			if p == "" {
				out[last] = true
				continue
			}
			out[last+"/"+p] = true
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the webapp sections: %v", err)
	}

	if len(out) == 0 {
		t.Fatal("found no routes; the walk is looking in the wrong place")
	}

	return out
}

// pattern turns a built URL back into the route shape it should match: the
// base URL goes, a numeric segment becomes its parameter, and so does the
// namespace slug.
func pattern(url, route string) string {
	path := strings.TrimPrefix(url, strings.TrimRight(Base(), "/"))

	gotParts := strings.Split(strings.Trim(path, "/"), "/")
	wantParts := strings.Split(strings.Trim(route, "/"), "/")
	if len(gotParts) != len(wantParts) {
		return path
	}

	for i, w := range wantParts {
		if strings.HasPrefix(w, ":") {
			gotParts[i] = w
		}
	}

	return "/" + strings.Join(gotParts, "/")
}

func TestEveryLinkMatchesAWebappRoute(t *testing.T) {
	routes := unifyRoutes(t)

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got == "" {
				t.Fatal("builder returned no link for valid arguments")
			}

			if !routes[tc.route] {
				t.Fatalf("no webapp route %q — this package links somewhere unify does not serve", tc.route)
			}

			if got := pattern(tc.got, tc.route); got != tc.route {
				t.Errorf("built %q, which is not the route %q", got, tc.route)
			}
		})
	}
}

// A missing identifier must produce no link rather than one pointing at a
// list, which would read as "here is the thing" and show something else.
func TestBuildersReturnNothingWithoutAnIdentifier(t *testing.T) {
	for name, got := range map[string]string{
		"page without a slug":     ComposePage("", 12),
		"page without an id":      ComposePage("crm", 0),
		"namespace without slug":  ComposeNamespace(""),
		"workflow without an id":  Workflow(0),
		"record without a module": ComposeRecord("crm", 0, 34),
	} {
		if got != "" {
			t.Errorf("%s built %q, want no link", name, got)
		}
	}
}

// A link that loses the port is a link to nothing. The webapp is routinely
// served somewhere other than 80 — every dev checkout is — and with nothing
// naming the webapp host separately the base fell through to a bare
// "localhost", so every URL these builders returned pointed at a port nobody
// listens on.
func TestBaseKeepsThePort(t *testing.T) {
	for name, tc := range map[string]struct {
		domain, domainWebapp, webappEnabled, want string
	}{
		"served by the API host":   {domain: "localhost:1043", webappEnabled: "true", want: "http://localhost:1043"},
		"served elsewhere, unsaid": {domain: "localhost:1043", webappEnabled: "false", want: "http://localhost:1043"},
		"served elsewhere, named":  {domain: "localhost:1043", domainWebapp: "localhost:5173", webappEnabled: "false", want: "http://localhost:5173"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("DOMAIN", tc.domain)
			t.Setenv("DOMAIN_WEBAPP", tc.domainWebapp)
			t.Setenv("HTTP_WEBAPP_ENABLED", tc.webappEnabled)
			t.Setenv("HTTP_BASE_URL", "")
			t.Setenv("HTTP_WEBAPP_BASE_URL", "/")

			if got := strings.TrimRight(Base(), "/"); got != tc.want {
				t.Errorf("Base() = %q, want %q", got, tc.want)
			}
			if got := ComposeNamespace("crm"); got != tc.want+"/compose/namespace/crm" {
				t.Errorf("ComposeNamespace = %q, want %q", got, tc.want+"/compose/namespace/crm")
			}
		})
	}
}
