package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/internal/repository"
	"github.com/cortezaproject/corteza-server/compose/proto"
	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/internal/store"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
	proto2 "github.com/cortezaproject/corteza-server/system/proto"
)

type (
	permissionServicer interface {
		accessControlPermissionServicer
		Watch(ctx context.Context)
	}

	Config struct {
		Storage          options.StorageOpt
		Corredor         options.CorredorOpt
		GRPCClientSystem options.GRPCServerOpt
	}
)

var (
	DefaultLogger *zap.Logger

	// DefaultPermissions Retrives & stores permissions
	DefaultPermissions permissionServicer

	// DefaultAccessControl Access control checking
	DefaultAccessControl *accessControl

	// DefaultAutomationScriptManager manages scripts
	DefaultAutomationScriptManager automationScript

	// DefaultAutomationTriggerManager manages triggerManager
	DefaultAutomationTriggerManager automationTrigger

	// DefaultAutomationRunner runs automation scripts by listening to triggerManager and invoking Corredor service
	DefaultAutomationRunner automationRunner

	DefaultNamespace NamespaceService
	DefaultRecord    RecordService
	DefaultModule    ModuleService
	DefaultChart     ChartService
	DefaultPage      PageService

	DefaultAttachment   AttachmentService
	DefaultNotification NotificationService

	DefaultSystemUser *systemUser
)

func Init(ctx context.Context, log *zap.Logger, c Config) (err error) {
	var db = repository.DB(ctx)

	DefaultLogger = log.Named("service")

	fs, err := store.New(c.Storage.Path)
	log.Info("initializing store", zap.String("path", c.Storage.Path), zap.Error(err))
	if err != nil {
		return err
	}

	// Permissions, access control
	DefaultPermissions = permissions.Service(
		ctx,
		DefaultLogger,
		permissions.Repository(db, "compose_permission_rules"))

	DefaultAccessControl = AccessControl(DefaultPermissions)

	{
		systemClientConn, err := NewSystemGRPCClient(ctx, c.GRPCClientSystem, DefaultLogger)
		if err != nil {
			return err
		}

		DefaultSystemUser = SystemUser(proto2.NewUsersClient(systemClientConn))
	}

	// ias: Internal Automatinon Service
	// handles script & trigger management & keeping runnables cripts in internal cache
	ias := automation.Service(automation.AutomationServiceConfig{
		Logger:        DefaultLogger,
		DbTablePrefix: "compose",
		DB:            db,
		TokenMaker: func(ctx context.Context, userID uint64) (s string, e error) {
			ctx = auth.SetSuperUserContext(ctx)
			return DefaultSystemUser.MakeJWT(ctx, userID)
		},
	})

	// Pass automation manager to
	DefaultAutomationScriptManager = AutomationScript(ias)
	DefaultAutomationTriggerManager = AutomationTrigger(ias)

	{
		var scriptRunnerClient proto.ScriptRunnerClient

		if c.Corredor.Enabled {
			corredor, err := automation.Corredor(ctx, c.Corredor, DefaultLogger)

			log.Info("initializing corredor connection", zap.String("addr", c.Corredor.Addr), zap.Error(err))
			if err != nil {
				return err
			}

			scriptRunnerClient = proto.NewScriptRunnerClient(corredor)
		}

		DefaultAutomationRunner = AutomationRunner(ias, scriptRunnerClient)
	}

	// Compose internals:
	DefaultNamespace = Namespace()
	DefaultRecord = Record()
	DefaultModule = Module()
	DefaultPage = Page()
	DefaultChart = Chart()
	DefaultNotification = Notification()
	DefaultAttachment = Attachment(fs)

	return nil
}

func Watchers(ctx context.Context) {
	// Reloading automation scripts on change
	DefaultAutomationRunner.Watch(ctx)

	// Reloading permissions on change
	DefaultPermissions.Watch(ctx)
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
