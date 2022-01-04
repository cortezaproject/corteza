package reporter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/compose/rest"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/pkg/api/server"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza-server/pkg/envoy/directory"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/objstore/plain"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/report"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/go-chi/chi/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/afero"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
)

type (
	helper struct {
		t *testing.T
		a *require.Assertions

		cUser  *sysTypes.User
		roleID uint64
		token  string
	}

	auxReport struct {
		*types.Report

		Frames   report.FrameDefinitionSet `json:"frames"`
		Describe []string                  `json:"describe"`
	}

	valueDef map[string][]string
)

var (
	testApp *app.CortezaApp
	r       chi.Router

	eventBus = eventbus.New()
)

func init() {
	helpers.RecursiveDotEnvLoad()
}

// random string, 10 chars long by default
func rs(a ...int) string {
	var l = 10
	if len(a) > 0 {
		l = a[0]
	}

	return string(rand.Bytes(l))
}

func InitTestApp() {
	if testApp == nil {
		ctx := cli.Context()

		testApp = helpers.NewIntegrationTestApp(ctx, func(app *app.CortezaApp) (err error) {
			service.DefaultStore = app.Store
			service.DefaultObjectStore, err = plain.NewWithAfero(afero.NewMemMapFs(), "test")
			if err != nil {
				return err
			}

			eventbus.Set(eventBus)
			return nil
		})

	}

	if r == nil {
		r = chi.NewRouter()
		r.Use(server.BaseMiddleware(false, logger.Default())...)
		helpers.BindAuthMiddleware(r)
		rest.MountRoutes(r)
	}
}

func TestMain(m *testing.M) {
	InitTestApp()
	os.Exit(m.Run())
}

func newHelper(t *testing.T) helper {
	h := helper{
		t:      t,
		a:      require.New(t),
		roleID: id.Next(),
		cUser: &sysTypes.User{
			ID: id.Next(),
		},
	}

	h.cUser.SetRoles(h.roleID)
	helpers.UpdateRBAC(h.roleID)

	var err error
	h.token, err = auth.DefaultJwtHandler.Generate(context.Background(), h.cUser)
	if err != nil {
		panic(err)
	}

	return h
}

func (h helper) MyRole() uint64 {
	return h.roleID
}

// Returns context w/ security details
func (h helper) secCtx() context.Context {
	return auth.SetIdentityToContext(context.Background(), h.cUser)
}

// apitest basics, initialize, set handler, add auth
func (h helper) apiInit() *apitest.APITest {
	InitTestApp()

	return apitest.
		New().
		Handler(r).
		Intercept(helpers.ReqHeaderRawAuthBearer(h.token))
}
func (h helper) mockPermissions(rules ...*rbac.Rule) {
	h.noError(rbac.Global().Grant(
		// TestService we use does not have any backend storage,
		context.Background(),
		rules...,
	))
}

// Unwraps error before it passes it to the tester
func (h helper) noError(err error) {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	h.a.NoError(err)
}

func setup(t *testing.T) (context.Context, helper, store.Storer) {
	h := newHelper(t)
	s := service.DefaultStore

	u := &sysTypes.User{
		ID: id.Next(),
	}
	u.SetRoles(auth.BypassRoles().IDs()...)

	ctx := auth.SetIdentityToContext(context.Background(), u)

	return ctx, h, s
}

func loadNoErr(ctx context.Context, h helper, m report.M, dd ...*report.FrameDefinition) (ff []*report.Frame) {
	ff, err := m.Load(ctx, dd...)
	h.a.NoError(err)
	return
}

func loadErr(ctx context.Context, h helper, m report.M, d *report.FrameDefinition, msg string) {
	_, err := m.Load(ctx, d)
	h.a.Error(err)
	h.a.Contains(err.Error(), msg)
}

// loadNoErrMulti is a little wrapper that does some preprocessing on the frame definitions.
// It is a copy from the system/service/report.
//
// @todo make this cleaner!
func loadNoErrMulti(ctx context.Context, h helper, m report.M, dd ...*report.FrameDefinition) (ff []*report.Frame) {
	ff = make([]*report.Frame, 0, len(dd))

	auxdd := make([]*report.FrameDefinition, 0, len(dd))
	for i, d := range dd {
		// first one; nothing special needed
		if i == 0 {
			auxdd = append(auxdd, d)
			continue
		}

		stp := m.GetStep(d.Source)
		if stp == nil {
			h.a.FailNow(fmt.Sprintf("unknown source: %s", d.Source))
		}

		// if the current source matches the prev. source, and they both define references,
		// they fall into the same chunk.
		if stp.Def().Join != nil && (d.Source == dd[i-1].Source) && (d.Ref != "" && dd[i-1].Ref != "") {
			auxdd = append(auxdd, d)
			continue
		}

		// if the current one doesn't fall into the current chunk, process
		// the chunk and reset it
		aux, err := m.Load(ctx, auxdd...)
		h.a.NoError(err)
		ff = append(ff, aux...)

		auxdd = make([]*report.FrameDefinition, 0, len(dd))
		auxdd = append(auxdd, d)
	}

	if len(auxdd) > 0 {
		aux, err := m.Load(ctx, auxdd...)
		h.a.NoError(err)
		ff = append(ff, aux...)
	}

	return
}

func describeNoErr(ctx context.Context, h helper, m report.M, dd ...string) (ff report.FrameDescriptionSet) {
	ff = make(report.FrameDescriptionSet, 0, len(dd))

	for _, d := range dd {
		aux, err := m.Describe(ctx, d)
		h.a.NoError(err)

		ff = append(ff, aux...)
	}

	return
}

func loadScenario(ctx context.Context, s store.Storer, t *testing.T, h helper) (report.M, *auxReport, report.FrameDefinitionSet) {
	return loadScenarioWithName(ctx, s, t, h, t.Name()[5:])
}

func loadScenarioWithName(ctx context.Context, s store.Storer, t *testing.T, h helper, scenario string) (report.M, *auxReport, report.FrameDefinitionSet) {
	var (
		providers = map[string]report.DatasourceProvider{
			"composeRecords": service.DefaultRecord,
		}
	)

	cleanup(ctx, h, s)
	parseEnvoy(ctx, s, h, "testdata/data_model")
	rr := parseReport(h, path.Join("testdata", scenario, "report.json"))
	m := modelReport(ctx, h, providers, rr)

	return m, rr, rr.Frames
}

func loadScenarioOwnDM(ctx context.Context, s store.Storer, t *testing.T, h helper) (report.M, *auxReport, report.FrameDefinitionSet) {
	return loadScenarioOwnDMWithName(ctx, s, t, h, t.Name()[5:])
}

func loadScenarioOwnDMWithName(ctx context.Context, s store.Storer, t *testing.T, h helper, scenario string) (report.M, *auxReport, report.FrameDefinitionSet) {
	var (
		providers = map[string]report.DatasourceProvider{
			"composeRecords": service.DefaultRecord,
		}
	)

	cleanup(ctx, h, s)
	parseEnvoy(ctx, s, h, path.Join("testdata", scenario, "data_model"))
	rr := parseReport(h, path.Join("testdata", scenario, "report.json"))
	m := modelReport(ctx, h, providers, rr)

	return m, rr, rr.Frames
}

func cleanup(ctx context.Context, h helper, s store.Storer) {
	h.noError(s.TruncateComposeNamespaces(ctx))
	h.noError(s.TruncateComposeModules(ctx))
	h.noError(s.TruncateComposeModuleFields(ctx))
	h.noError(s.TruncateComposeRecords(ctx, nil))
}

func parseEnvoy(ctx context.Context, s store.Storer, h helper, path string) {
	nn, err := directory.Decode(
		ctx,
		path,
		yaml.Decoder(),
		csv.Decoder(),
	)
	if err != nil {
		h.t.Fatalf("failed to decode scenario data: %v", err)
	}

	crs := resource.ComposeRecordShaper()
	nn, err = resource.Shape(nn, crs)
	h.a.NoError(err)

	// import into the store
	se := es.NewStoreEncoder(s, nil)
	bld := envoy.NewBuilder(se)
	g, err := bld.Build(ctx, nn...)
	h.a.NoError(err)
	err = envoy.Encode(ctx, g, se)
	h.a.NoError(err)
}

func parseReport(h helper, path string) *auxReport {
	f, err := os.Open(path)
	h.a.NoError(err)
	defer f.Close()

	aux := &auxReport{}
	raw, err := ioutil.ReadAll(f)
	h.a.NoError(err)

	err = json.Unmarshal(raw, &aux)
	h.a.NoError(err)

	return aux
}

func modelReport(ctx context.Context, h helper, pp map[string]report.DatasourceProvider, rr *auxReport) report.M {
	ss := rr.Sources.ModelSteps()
	model, err := report.Model(ctx, pp, ss...)
	h.a.NoError(err)
	err = model.Run(ctx)
	h.a.NoError(err)

	return model
}

func checkRows(h helper, f *report.Frame, req ...string) {
	h.a.Equal(len(req), f.Size())
	f.WalkRows(func(i int, r report.FrameRow) error {
		h.a.Contains(r.String(), req[i])
		return nil
	})
}

func indexJoinedResult(ff []*report.Frame) map[string]*report.Frame {
	out := make(map[string]*report.Frame)
	// the first one is the local ds
	for _, f := range ff[1:] {
		out[fmt.Sprintf("%s/%s/%s", f.Ref, f.RelSource, f.RefValue)] = f
	}

	return out
}
