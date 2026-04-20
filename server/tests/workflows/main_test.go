package workflows

import (
	"context"
	"fmt"

	// "fmt"
	"path"
	"testing"

	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/automation/service"
	autTypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/envoy"
	"github.com/crusttech/human/server/pkg/envoy/csv"
	"github.com/crusttech/human/server/pkg/envoy/directory"
	"github.com/crusttech/human/server/pkg/envoy/json"
	envoyStore "github.com/crusttech/human/server/pkg/envoy/store"
	"github.com/crusttech/human/server/pkg/envoy/yaml"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/store"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/crusttech/human/server/tests/helpers"
	"github.com/stretchr/testify/require"
)

var (
	defApp   *app.HumanApp
	defStore store.Storer
	defDal   dal.FullService
	eventBus = eventbus.New()
)

func init() {
	helpers.RecursiveDotEnvLoad()
}

func TestMain(m *testing.M) {
	logger.SetDefault(logger.MakeDebugLogger())
	ctx := context.Background()

	defApp = helpers.NewIntegrationTestApp(ctx, func(app *app.HumanApp) (err error) {
		// some test suites require action-log enabled
		app.Opt.ActionLog.WorkflowFunctionsEnabled = true
		defStore = app.Store
		eventbus.Set(eventBus)

		return nil
	})

	defDal = dal.Service()

	if err := defApp.Activate(ctx); err != nil {
		panic(fmt.Errorf("could not activate human: %v", err))
	}

	if err := defApp.InitExpr(ctx); err != nil {
		panic(fmt.Errorf("could not init expressions: %v", err))
	}

	m.Run()
}

func cleanup(t *testing.T) {
	var (
		ctx = context.Background()
	)

	if err := defStore.TruncateAutomationWorkflows(ctx); err != nil {
		t.Fatalf("failed to decode scenario data: %v", err)
	}
}

func truncateRecords(ctx context.Context) error {
	models, err := defDal.SearchModels(ctx)
	if err != nil {
		return err
	}
	for _, model := range models {
		err = defDal.Truncate(ctx, model.ToFilter(), nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadScenario(ctx context.Context, t *testing.T) {
	loadScenarioWithName(ctx, t, "S"+t.Name()[4:])
}

// 1st step in migration to workflow testdata w/o number prefix
//
// When all old scenarios are renamed, replace it with loadScenario.
func loadNewScenario(ctx context.Context, t *testing.T) {
	loadScenarioWithName(ctx, t, t.Name()[5:])
}

func loadScenarioWithName(ctx context.Context, t *testing.T, scenario string) {
	var (
		err error
	)

	cleanup(t)

	decoded, err := directory.Decode(
		ctx,
		path.Join("testdata", scenario),
		yaml.Decoder(),
		csv.Decoder(),
		json.Decoder(),
	)
	if err != nil {
		t.Fatalf("failed to decode scenario data: %v", err)
	}

	storeEnc := envoyStore.NewStoreEncoder(defStore, dal.Service(), &envoyStore.EncoderConfig{})

	b := envoy.NewBuilder(storeEnc)
	g, err := b.Build(ctx, decoded...)
	if err != nil {
		t.Fatalf("failed to build structure graph: %v", err)
	}

	if err = envoy.Encode(ctx, g, storeEnc); err != nil {
		t.Fatalf("failed to build structure graph: %v", err)
	}

	// Reload and register workflows
	if err = service.DefaultWorkflow.Load(ctx); err != nil {
		t.Fatalf("failed to reload workflows: %v", err)
	}
}

func bypassRBAC(ctx context.Context) context.Context {
	u := &sysTypes.User{
		ID: id.Next(),
	}

	if err := defStore.CreateUser(ctx, u); err != nil {
		panic(err)
	}

	u.SetRoles(auth.BypassRoles().IDs()...)
	return auth.SetIdentityToContext(ctx, u)
}

func execWorkflow(ctx context.Context, name string, p autTypes.WorkflowExecParams) (*expr.Vars, uint64, autTypes.Stacktrace, error) {
	wf, err := defStore.LookupAutomationWorkflowByHandle(ctx, name)
	if err != nil {
		return nil, 0, nil, err
	}

	return service.DefaultWorkflow.Exec(ctx, wf.ID, p)
}

func mustExecWorkflow(ctx context.Context, t *testing.T, name string, p autTypes.WorkflowExecParams) (vars *expr.Vars, strace autTypes.Stacktrace) {
	var err error
	vars, _, strace, err = execWorkflow(ctx, name, p)
	if err != nil {
		if issues, is := err.(autTypes.WorkflowIssueSet); is {
			for _, i := range issues {
				t.Logf("issue: %s", i.Description)
				t.Logf("       %v", i.Culprit)
			}
		}

		t.Fatalf("could not exec %q: %v", name, err)

	}

	return
}

func addRoleMember(ctx context.Context, req *require.Assertions, r string, uu ...string) {
	role, err := store.LookupRoleByHandle(ctx, defStore, r)
	req.NoError(err)

	rr := make([]*sysTypes.RoleMember, len(uu))
	for i, u := range uu {
		usr, err := store.LookupUserByHandle(ctx, defStore, u)
		req.NoError(err)

		rr[i] = &sysTypes.RoleMember{
			RoleID:   role.ID,
			Resource: fmt.Sprintf("corteza::system:user/%d", usr.ID),
		}
	}

	req.NoError(store.CreateRoleMember(ctx, defStore, rr...))
}
