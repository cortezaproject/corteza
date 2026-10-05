// Package federation_e2e runs federation between two real server processes,
// each with its own SQLite database: pairing over the handshake endpoints, then
// structure and data sync driven by the servers' own sync workers.
//
// It builds the server binary, so it is opt-in: HUMAN_FEDERATION_E2E=1.
package federation_e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const adminEmail = "admin@federation-e2e.test"

type node struct {
	t     *testing.T
	name  string
	bin   string
	dir   string
	port  int
	env   []string
	proc  *exec.Cmd
	token string
}

func TestFederationSyncBetweenTwoServers(t *testing.T) {
	if os.Getenv("HUMAN_FEDERATION_E2E") == "" {
		t.Skip("set HUMAN_FEDERATION_E2E=1 to build the server and run two of them")
	}

	req := require.New(t)
	bin := buildServer(t)

	origin := startNode(t, bin, "origin")
	dest := startNode(t, bin, "dest")

	// Origin data: a module and two records, one carrying an ID as text
	originNs := origin.createNamespace("fed_origin")
	originMod := origin.createModule(originNs, "contact", "Name", "ExternalRef")
	originRec1 := origin.createRecord(originNs, originMod, map[string]string{"Name": "Ada", "ExternalRef": "516687706987429889"})
	originRec2 := origin.createRecord(originNs, originMod, map[string]string{"Name": "Grace", "ExternalRef": ""})

	// Destination module the shared module is mapped onto
	destNs := dest.createNamespace("fed_dest")
	destMod := dest.createModule(destNs, "person", "FullName", "Ref")

	// Pairing: origin creates a node for dest and a pairing URI, dest registers
	// from it and asks to pair, origin confirms
	originNodeID := origin.data("POST", "/federation/nodes/", url.Values{
		"name":    {"dest"},
		"baseURL": {dest.federationURL()},
	}).str("nodeID")
	pairingURI := origin.data("POST", fmt.Sprintf("/federation/nodes/%s/uri", originNodeID), nil).raw

	destNodeID := dest.data("POST", "/federation/nodes/", url.Values{"pairingURI": {pairingURI}}).str("nodeID")
	// The pairing URI is always https; these servers speak plain http
	dest.data("POST", "/federation/nodes/"+destNodeID, url.Values{
		"name":    {"origin"},
		"baseURL": {origin.federationURL()},
	})
	dest.data("POST", fmt.Sprintf("/federation/nodes/%s/pair", destNodeID), nil)
	origin.data("POST", fmt.Sprintf("/federation/nodes/%s/handshake-confirm", originNodeID), nil)

	req.Equal("paired", origin.data("GET", "/federation/nodes/"+originNodeID, nil).str("status"))
	req.Equal("paired", dest.data("GET", "/federation/nodes/"+destNodeID, nil).str("status"))

	// Origin exposes the module to dest
	field := func(name string) map[string]any {
		return map[string]any{"kind": "String", "name": name, "label": name, "isMulti": false}
	}
	origin.json("PUT", fmt.Sprintf("/federation/nodes/%s/modules/", originNodeID), map[string]any{
		"composeModuleID":    originMod,
		"composeNamespaceID": originNs,
		"name":               "Contact",
		"handle":             "contact",
		"fields":             []any{field("Name"), field("ExternalRef")},
	})

	// The paired node reads the exposed records as the origin's federation
	// role, which the admin grants access to them
	origin.grantFederationRole(
		"corteza::compose:namespace/"+originNs, "read",
		"corteza::compose:module/"+originNs+"/"+originMod, "read",
		"corteza::compose:module/"+originNs+"/"+originMod, "records.search",
		"corteza::compose:record/"+originNs+"/"+originMod+"/*", "read",
		"corteza::compose:module-field/"+originNs+"/"+originMod+"/*", "record.value.read",
	)

	// Structure sync: the shared module appears on dest
	var sharedModuleID string
	req.Eventually(func() bool {
		for _, m := range dest.data("GET", fmt.Sprintf("/federation/nodes/%s/modules/?shared=true", destNodeID), nil).list() {
			sharedModuleID, _ = m["moduleID"].(string)
		}
		return sharedModuleID != ""
	}, time.Minute, time.Second, "shared module never arrived on dest")

	// Dest maps the shared module onto its own module
	dest.json("PUT", fmt.Sprintf("/federation/nodes/%s/modules/%s/mapped", destNodeID, sharedModuleID), map[string]any{
		"composeModuleID":    destMod,
		"composeNamespaceID": destNs,
		"fields": []any{
			map[string]any{"origin": field("Name"), "destination": field("FullName")},
			map[string]any{"origin": field("ExternalRef"), "destination": field("Ref")},
		},
	})

	// Data sync: both records land on dest, each pointing back at its exact
	// origin record ID
	byOrigin := map[string]map[string]string{}
	sources := map[string]bool{}
	req.Eventually(func() bool {
		byOrigin = map[string]map[string]string{}
		for _, r := range dest.records(destNs, destMod) {
			meta, _ := r["meta"].(map[string]any)
			ext, _ := meta["federation_extrecord"].(string)
			src, _ := meta["federation"].(string)
			byOrigin[ext] = recordValues(r)
			sources[src] = true
		}
		return len(byOrigin) == 2
	}, time.Minute, time.Second, "records never arrived on dest")

	req.Contains(byOrigin, originRec1, "federation_extrecord must be the exact origin record ID")
	req.Contains(byOrigin, originRec2)
	req.Equal("Ada", byOrigin[originRec1]["FullName"])
	req.Equal("516687706987429889", byOrigin[originRec1]["Ref"])
	req.Equal("Grace", byOrigin[originRec2]["FullName"])
	req.Equal(map[string]bool{origin.federationURL(): true}, sources, "each record names the node it came from")
}

// buildServer compiles the server once into the test's temp dir.
func buildServer(t *testing.T) string {
	bin := filepath.Join(t.TempDir(), "human")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/human")
	cmd.Dir = filepath.Join("..", "..")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "go build: %s", out)
	return bin
}

func startNode(t *testing.T, bin, name string) *node {
	n := &node{t: t, name: name, bin: bin, dir: t.TempDir(), port: freePort(t)}
	n.env = append(os.Environ(),
		"DB_DSN=sqlite3://file:"+filepath.Join(n.dir, "db.sqlite")+"?_busy_timeout=10000&_journal_mode=WAL",
		fmt.Sprintf("HTTP_ADDR=127.0.0.1:%d", n.port),
		fmt.Sprintf("DOMAIN=127.0.0.1:%d", n.port),
		"HTTP_API_BASE_URL=/api",
		"HTTP_WEBAPP_ENABLED=false",
		"CORREDOR_ENABLED=false",
		"ENVIRONMENT=test",
		"LOCALE_PATH="+repoPath(t, "locale"),
		"PROVISION_PATH="+repoPath(t, "server", "provision")+"/*",
		"LOCALE_LANGUAGES=en",
		"LOG_LEVEL=info",
		"AUTH_JWT_SECRET=federation-e2e-"+name,
		"FEDERATION_ENABLED=true",
		"FEDERATION_LABEL="+name,
		fmt.Sprintf("FEDERATION_HOST=127.0.0.1:%d", n.port),
		"FEDERATION_SYNC_STRUCTURE_MONITOR_INTERVAL=2s",
		"FEDERATION_SYNC_DATA_MONITOR_INTERVAL=2s",
		// one record per page, so a second record needs the next-page request
		"FEDERATION_SYNC_DATA_PAGE_SIZE=1",
	)

	logFile, err := os.Create(filepath.Join(n.dir, "server.log"))
	require.NoError(t, err)

	n.proc = exec.Command(bin, "serve-api")
	n.proc.Dir = n.dir
	n.proc.Env = n.env
	n.proc.Stdout = logFile
	n.proc.Stderr = logFile
	require.NoError(t, n.proc.Start())

	exited := make(chan struct{})
	go func() {
		_ = n.proc.Wait()
		close(exited)
	}()

	t.Cleanup(func() {
		_ = n.proc.Process.Kill()
		<-exited
		if t.Failed() {
			log, _ := os.ReadFile(logFile.Name())
			t.Logf("%s server log, without successful actions:\n%s", name, serviceLog(string(log), 60))
		}
	})

	require.Eventually(t, func() bool {
		select {
		case <-exited:
			require.FailNow(t, name+" server exited during startup")
		default:
		}

		rsp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/healthcheck", n.port))
		if err != nil {
			return false
		}
		rsp.Body.Close()
		return rsp.StatusCode == http.StatusOK
	}, 2*time.Minute, 500*time.Millisecond, "%s server never became healthy", name)

	n.cli("users", "add", adminEmail, "--password", "Federation-e2e-1", "--role", "admin")
	out := n.cli("auth", "jwt", adminEmail, "--scope", "profile", "--scope", "api")
	for _, f := range strings.Fields(out) {
		if strings.Count(f, ".") == 2 && strings.HasPrefix(f, "ey") {
			n.token = f
		}
	}
	require.NotEmpty(t, n.token, "%s: no JWT in %q", name, out)
	return n
}

func (n *node) cli(args ...string) string {
	cmd := exec.Command(n.bin, args...)
	cmd.Dir = n.dir
	cmd.Env = n.env
	out, err := cmd.CombinedOutput()
	require.NoError(n.t, err, "%s %v: %s", n.name, args, out)
	return string(out)
}

func (n *node) apiURL(path string) string {
	return fmt.Sprintf("http://127.0.0.1:%d/api%s", n.port, path)
}

func (n *node) federationURL() string {
	return n.apiURL("/federation")
}

type response struct {
	t   *testing.T
	raw string
	obj map[string]any
	arr []any
}

// data calls the API with form values and returns its "response", failing the
// test on an error payload.
func (n *node) data(method, path string, form url.Values) response {
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}

	r, err := http.NewRequest(method, n.apiURL(path), body)
	require.NoError(n.t, err)
	r.Header.Set("Authorization", "Bearer "+n.token)
	r.Header.Set("Accept", "application/json")
	if form != nil {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return n.do(r)
}

func (n *node) json(method, path string, payload any) response {
	b, err := json.Marshal(payload)
	require.NoError(n.t, err)

	r, err := http.NewRequest(method, n.apiURL(path), bytes.NewReader(b))
	require.NoError(n.t, err)
	r.Header.Set("Authorization", "Bearer "+n.token)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	return n.do(r)
}

func (n *node) do(r *http.Request) response {
	rsp, err := http.DefaultClient.Do(r)
	require.NoError(n.t, err)
	defer rsp.Body.Close()

	raw, err := io.ReadAll(rsp.Body)
	require.NoError(n.t, err)

	var env struct {
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
		Response json.RawMessage `json:"response"`
	}
	require.NoError(n.t, json.Unmarshal(raw, &env), "%s %s %s: %s", n.name, r.Method, r.URL.Path, raw)
	require.Nil(n.t, env.Error, "%s %s %s: %s", n.name, r.Method, r.URL.Path, raw)

	out := response{t: n.t}
	_ = json.Unmarshal(env.Response, &out.raw)
	_ = json.Unmarshal(env.Response, &out.obj)
	_ = json.Unmarshal(env.Response, &out.arr)
	return out
}

func (r response) str(key string) string {
	v, ok := r.obj[key].(string)
	require.True(r.t, ok, "%q is not a string in %v", key, r.obj)
	return v
}

func (r response) list() (out []map[string]any) {
	items := r.arr
	if set, ok := r.obj["set"].([]any); ok {
		items = set
	}
	for _, i := range items {
		if m, ok := i.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return
}

// grantFederationRole allows the federation role each resource, operation
// pair given.
func (n *node) grantFederationRole(resOps ...string) {
	var roleID string
	for _, r := range n.data("GET", "/system/roles/?query=federation", nil).list() {
		if r["handle"] == "federation" {
			roleID, _ = r["roleID"].(string)
		}
	}
	require.NotEmpty(n.t, roleID, "%s has no federation role", n.name)

	// each component takes its own rules: corteza::<component>:…
	rules := map[string][]map[string]string{}
	for i := 0; i < len(resOps); i += 2 {
		component := strings.TrimSuffix(strings.SplitN(strings.TrimPrefix(resOps[i], "corteza::"), ":", 2)[0], "/")
		rules[component] = append(rules[component], map[string]string{"resource": resOps[i], "operation": resOps[i+1], "access": "allow"})
	}
	for component, rr := range rules {
		n.json("PATCH", "/"+component+"/permissions/"+roleID+"/rules", map[string]any{"rules": rr})
	}
}

func (n *node) createNamespace(slug string) string {
	return n.json("POST", "/compose/namespace/", map[string]any{"name": slug, "slug": slug, "enabled": true}).str("namespaceID")
}

func (n *node) createModule(nsID, handle string, fields ...string) string {
	ff := make([]map[string]any, len(fields))
	for i, f := range fields {
		ff[i] = map[string]any{"name": f, "label": f, "kind": "String"}
	}
	return n.json("POST", fmt.Sprintf("/compose/namespace/%s/module/", nsID), map[string]any{
		"name": handle, "handle": handle, "fields": ff,
	}).str("moduleID")
}

func (n *node) createRecord(nsID, modID string, values map[string]string) string {
	vv := []map[string]string{}
	for k, v := range values {
		vv = append(vv, map[string]string{"name": k, "value": v})
	}
	return n.json("POST", fmt.Sprintf("/compose/namespace/%s/module/%s/record/", nsID, modID), map[string]any{
		"values": vv,
	}).str("recordID")
}

func (n *node) records(nsID, modID string) []map[string]any {
	return n.data("GET", fmt.Sprintf("/compose/namespace/%s/module/%s/record/", nsID, modID), nil).list()
}

func recordValues(r map[string]any) map[string]string {
	out := map[string]string{}
	vv, _ := r["values"].([]any)
	for _, v := range vv {
		m, _ := v.(map[string]any)
		name, _ := m["name"].(string)
		value, _ := m["value"].(string)
		out[name] = value
	}
	return out
}

// repoPath resolves a path in the repository; the servers run in temp dirs,
// so the translations and provisioning files they load are given absolutely.
func repoPath(t *testing.T, elem ...string) string {
	p, err := filepath.Abs(filepath.Join(append([]string{"..", "..", ".."}, elem...)...))
	require.NoError(t, err)
	return p
}

func freePort(t *testing.T) int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// serviceLog keeps the last n log lines, dropping action log entries that
// record no error.
func serviceLog(log string, n int) string {
	var kept []string
	for _, l := range strings.Split(log, "\n") {
		if l != "" && (!strings.Contains(l, `"logger":"actionlog"`) || !strings.Contains(l, `"error":""`)) {
			kept = append(kept, l)
		}
	}
	if len(kept) > n {
		kept = kept[len(kept)-n:]
	}
	return strings.Join(kept, "\n")
}
