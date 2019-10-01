package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/automation/corredor"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/store"
	systemProto "github.com/cortezaproject/corteza-server/system/proto"
)

type (
	permissionServicer interface {
		accessControlPermissionServicer
		Watch(ctx context.Context)
	}

	automationManager interface {
		automationScriptManager
		automationTriggerManager
		automationScriptsFinder
	}

	Config struct {
		Storage          options.StorageOpt
		Corredor         options.CorredorOpt
		GRPCClientSystem options.GRPCServerOpt
	}
)

var (
	DefaultStore store.Store

	DefaultLogger *zap.Logger

	// DefaultPermissions Retrieves & stores permissions
	DefaultPermissions permissionServicer

	// DefaultAccessControl Access control checking
	DefaultAccessControl *accessControl

	// DefaultInternalAutomationManager manages automation scripts, triggers, runnable scripts
	DefaultInternalAutomationManager automationManager

	// DefaultAutomationScriptManager manages compose automation scripts
	DefaultAutomationScriptManager automationScript

	// DefaultAutomationTriggerManager manages compose automation triggers
	DefaultAutomationTriggerManager automationTrigger

	// DefaultAutomationRunner runs automation scripts by listening to triggerManager and invoking Corredor service
	DefaultAutomationRunner automationRunner

	DefaultNamespace     NamespaceService
	DefaultImportSession ImportSessionService
	DefaultRecord        RecordService
	DefaultModule        ModuleService
	DefaultChart         ChartService
	DefaultPage          PageService
	DefaultAttachment    AttachmentService
	DefaultNotification  *notification

	DefaultSystemUser *systemUser
	DefaultSystemRole *systemRole
)

func Init(ctx context.Context, log *zap.Logger, c Config) (err error) {
	var db = repository.DB(ctx)

	DefaultLogger = log.Named("service")

	if DefaultStore == nil {
		DefaultStore, err = store.New(c.Storage.Path)
		log.Info("initializing store", zap.String("path", c.Storage.Path), zap.Error(err))
		if err != nil {
			return err
		}
	}

	// Permissions, access control
	if DefaultPermissions == nil {
		DefaultPermissions = permissions.Service(
			ctx,
			DefaultLogger,
			permissions.Repository(db, "compose_permission_rules"))
	}
	DefaultAccessControl = AccessControl(DefaultPermissions)

	DefaultNamespace = Namespace()
	DefaultModule = Module()

	{
		systemClientConn, err := NewSystemGRPCClient(ctx, c.GRPCClientSystem, DefaultLogger)
		if err != nil {
			return err
		}

		DefaultSystemUser = SystemUser(systemProto.NewUsersClient(systemClientConn))
		DefaultSystemRole = SystemRole(systemProto.NewRolesClient(systemClientConn))
	}

	{
		if DefaultInternalAutomationManager == nil {
			// handles script & trigger management & keeping runnable scripts in internal cache
			DefaultInternalAutomationManager = automation.Service(automation.AutomationServiceConfig{
				Logger:        DefaultLogger,
				DbTablePrefix: "compose",
				DB:            db,
				TokenMaker: func(ctx context.Context, userID uint64) (s string, e error) {
					ctx = auth.SetSuperUserContext(ctx)
					return DefaultSystemUser.MakeJWT(ctx, userID)
				},
			})
		}

		// Pass internal automation manager to compose's script & trigger managers
		DefaultAutomationTriggerManager = AutomationTrigger(DefaultInternalAutomationManager)
		DefaultAutomationScriptManager = AutomationScript(DefaultInternalAutomationManager)

		var scriptRunnerClient corredor.ScriptRunnerClient

		if c.Corredor.Enabled {
			conn, err := corredor.NewConnection(ctx, c.Corredor, DefaultLogger)

			log.Info("initializing corredor connection", zap.String("addr", c.Corredor.Addr), zap.Error(err))
			if err != nil {
				return err
			}

			scriptRunnerClient = corredor.NewScriptRunnerClient(conn)
		}

		DefaultAutomationRunner = AutomationRunner(
			AutomationRunnerOpt{
				ApiBaseURLSystem:    c.Corredor.ApiBaseURLSystem,
				ApiBaseURLMessaging: c.Corredor.ApiBaseURLMessaging,
				ApiBaseURLCompose:   c.Corredor.ApiBaseURLCompose,
			},
			DefaultInternalAutomationManager,
			scriptRunnerClient,
		)
	}

	DefaultImportSession = ImportSession()
	DefaultRecord = Record()
	DefaultPage = Page()
	DefaultChart = Chart()
	DefaultNotification = Notification()
	DefaultAttachment = Attachment(DefaultStore)

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
