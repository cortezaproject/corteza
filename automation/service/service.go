package service

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/automation/automation"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/objstore"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

type (
	RBACServicer interface {
		accessControlRBACServicer
		Watch(ctx context.Context)
	}

	Config struct {
		ActionLog options.ActionLogOpt
		Workflow  options.WorkflowOpt
	}

	userService interface {
		FindByID(ctx context.Context, userID uint64) (*types.User, error)
	}
)

var (
	DefaultObjectStore objstore.Store

	// DefaultStore is an interface to storage backend(s)
	// ng (next-gen) is a temporary prefix
	// so that we can differentiate between it and the file-only store
	DefaultStore store.Storer

	DefaultLogger *zap.Logger

	// DefaultAccessControl Access control checking
	DefaultAccessControl *accessControl

	DefaultActionlog actionlog.Recorder

	DefaultUser     userService
	DefaultWorkflow *workflow
	DefaultTrigger  *trigger
	DefaultSession  *session

	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}

	// wrapper around nextID that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

func Initialize(ctx context.Context, log *zap.Logger, s store.Storer, c Config) (err error) {
	var (
	//hcd = healthcheck.Defaults()
	)

	// we're doing conversion to avoid having
	// store interface exposed or generated inside app package
	DefaultStore = s

	DefaultLogger = log.Named("workflow")

	{
		tee := zap.NewNop()
		policy := actionlog.MakeProductionPolicy()

		if !c.ActionLog.Enabled {
			policy = actionlog.MakeDisabledPolicy()
		} else if c.ActionLog.Debug {
			policy = actionlog.MakeDebugPolicy()
			tee = log
		}

		DefaultActionlog = actionlog.NewService(DefaultStore, log, tee, policy)
	}

	DefaultAccessControl = AccessControl(rbac.Global())

	DefaultSession = Session(DefaultLogger.Named("session"), c.Workflow)
	DefaultWorkflow = Workflow(DefaultLogger.Named("workflow"))
	DefaultTrigger = Trigger(DefaultLogger.Named("trigger"), c.Workflow)

	DefaultWorkflow.triggers = DefaultTrigger

	Registry().AddTypes(
		&expr.Any{},
		&expr.Array{},
		&expr.Boolean{},
		&expr.ID{},
		&expr.Integer{},
		&expr.UnsignedInteger{},
		&expr.Float{},
		&expr.String{},
		&expr.Handle{},
		&expr.DateTime{},
		&expr.Duration{},
		&expr.KV{},
		&expr.KVV{},
		&expr.Reader{},
		&expr.Vars{},

		&automation.EmailMessage{},
	)

	automation.HttpRequestHandler(Registry())
	automation.LogHandler(Registry())
	automation.LoopHandler(Registry(), DefaultWorkflow.parser)
	automation.CorredorHandler(Registry(), corredor.Service())
	automation.EmailHandler(Registry())

	return
}

func Activate(ctx context.Context) (err error) {
	if err = DefaultWorkflow.Load(ctx); err != nil {
		return
	}

	return
}

func Watchers(ctx context.Context) {
	DefaultSession.Watch(ctx)
	return
}

// Data is stale when new date does not match updatedAt or createdAt (before first update)
func isStale(new *time.Time, updatedAt *time.Time, createdAt time.Time) bool {
	if new == nil {
		// Change to true for stale-data-check
		return false
	}

	if updatedAt != nil {
		return !new.Equal(*updatedAt)
	}

	return new.Equal(createdAt)
}

// trim1st removes 1st param and returns only error
func trim1st(_ interface{}, err error) error {
	return err
}
