package dal

import (
	"context"
	"testing"
	"time"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
)

type (
	helper struct {
		t *testing.T
		b *testing.B
		a *require.Assertions

		cUser  *types.User
		roleID uint64
		token  []byte
	}
)

var (
	testUser *types.User
)

func newHelperT(t *testing.T) helper {
	h := newHelper(t, require.New(t))
	h.t = t
	return h
}

func newHelperB(b *testing.B) helper {
	h := newHelper(b, require.New(b))
	h.b = b
	return h
}

func newHelper(t require.TestingT, a *require.Assertions) helper {
	var (
		h = helper{
			roleID: id.Next(),
			a:      a,
		}
		ctx = context.Background()

		err error
	)

	if testUser == nil {
		testUser = &types.User{
			Handle: "test_user",
			Name:   "test_user",
			ID:     id.Next(),
		}

		err = store.CreateUser(ctx, service.DefaultStore, testUser)
		if err != nil {
			panic(err)
		}

	}
	h.cUser = testUser

	h.cUser.SetRoles(h.roleID)
	helpers.UpdateRBAC(h.roleID)
	h.mockPermissionsWithAccess()

	h.token, err = auth.TokenIssuer.Issue(ctx, auth.WithIdentity(h.cUser))
	if err != nil {
		panic(err)
	}

	return h
}

// apitest basics, initialize, set handler, add auth
func (h helper) apiInit() *apitest.APITest {
	InitTestApp()

	return apitest.
		New().
		Handler(r).
		Intercept(helpers.ReqHeaderRawAuthBearer(h.token))
}

func (h helper) MyRole() uint64 {
	return h.roleID
}

// Returns context w/ security details
func (h helper) secCtx() context.Context {
	return auth.SetIdentityToContext(context.Background(), h.cUser)
}

func (h helper) mockPermissions(rules ...*rbac.Rule) {
	h.a.NoError(rbac.Global().Grant(
		// TestService we use does not have any backend storage,
		context.Background(),
		rules...,
	))
}

// Prepends allow access rule for system service for everyone
func (h helper) mockPermissionsWithAccess(rules ...*rbac.Rule) {
	h.mockPermissions(rules...)
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Resource utilities

func (h helper) createNamespace(name string) *composeTypes.Namespace {
	ns := &composeTypes.Namespace{Name: name, Slug: name}
	ns.ID = id.Next()
	ns.CreatedAt = time.Now()
	h.a.NoError(store.CreateComposeNamespace(context.Background(), service.DefaultStore, ns))
	return ns
}

func (h helper) createSensitivityLevel(res *types.DalSensitivityLevel) *types.DalSensitivityLevel {
	if res.ID == 0 {
		res.ID = id.Next()
	}

	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	h.a.NoError(service.DefaultStore.CreateDalSensitivityLevel(context.Background(), res))
	h.a.NoError(service.DefaultDalSensitivityLevel.ReloadSensitivityLevels(context.Background(), service.DefaultStore))
	return res
}

func (h helper) createDalConnection(res *types.DalConnection) *types.DalConnection {
	if res.ID == 0 {
		res.ID = id.Next()
	}

	if res.Name == "" {
		res.Name = "Test Connection"
	}
	if res.Handle == "" {
		res.Handle = "test_connection"
	}
	if res.Type == "" {
		res.Type = types.DalConnectionResourceType
	}
	if res.Ownership == "" {
		res.Ownership = "tester"
	}

	if res.Config.DefaultModelIdent == "" {
		res.Config.DefaultModelIdent = "compose_records"
	}
	if res.Config.DefaultAttributeIdent == "" {
		res.Config.DefaultAttributeIdent = "values"
	}
	if res.Config.DefaultPartitionFormat == "" {
		res.Config.DefaultPartitionFormat = "compose_records_{{namespace}}_{{module}}"
	}
	if res.Config.PartitionFormatValidator == "" {
		res.Config.PartitionFormatValidator = ""
	}
	if res.Config.Connection.Params == nil {
		res.Config.Connection = dal.NewDSNConnection("sqlite3://file::memory:?cache=shared&mode=memory")
	}

	if len(res.Capabilities.Enforced) == 0 {
		res.Capabilities.Enforced = capabilities.FullCapabilities()
	}

	if len(res.Capabilities.Supported) == 0 {
		res.Capabilities.Supported = capabilities.Set{}
	}

	if len(res.Capabilities.Unsupported) == 0 {
		res.Capabilities.Unsupported = capabilities.Set{}
	}

	if len(res.Capabilities.Enabled) == 0 {
		res.Capabilities.Enabled = capabilities.Set{}
	}

	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	if res.CreatedBy == 0 {
		res.CreatedBy = h.cUser.ID
	}

	h.a.NoError(service.DefaultStore.CreateDalConnection(context.Background(), res))
	h.a.NoError(service.DefaultDalConnection.ReloadConnections(context.Background()))
	return res
}

func (h helper) getPrimaryConnection() *types.DalConnection {
	cc, _, err := store.SearchDalConnections(context.Background(), service.DefaultStore, types.DalConnectionFilter{Type: types.DalPrimaryConnectionResourceType})
	h.a.NoError(err)

	if len(cc) != 1 {
		h.a.FailNow("invalid state: no or too many primary connections")
	}

	return cc[0]
}

func makeConnectionDefinition(dsn string) *types.DalConnection {
	return &types.DalConnection{
		ID:   id.Next(),
		Type: types.DalConnectionResourceType,
		Config: types.ConnectionConfig{
			DefaultModelIdent:     "compose_record",
			DefaultAttributeIdent: "values",

			DefaultPartitionFormat: "compose_record_{{namespace}}_{{module}}",

			PartitionFormatValidator: "",

			Connection: dal.NewDSNConnection(dsn),
		},
		Capabilities: types.ConnectionCapabilities{
			Supported: capabilities.FullCapabilities(),
		},
	}
}

// // // // // // // // // // // // // // // // // // // // // // // // //
