package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/automation/corredor"
	mailTrigger "github.com/cortezaproject/corteza-server/pkg/automation/mail"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/cortezaproject/corteza-server/system/proto"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	automationRunner struct {
		opt          AutomationRunnerOpt
		logger       *zap.Logger
		runner       corredor.ScriptRunnerClient
		userFinder   automationRunnerUserFinder
		scriptFinder automationScriptsFinder
		jwtEncoder   intAuth.TokenEncoder
	}

	automationRunnerUserFinder interface {
		FindByEmail(string) (*types.User, error)
	}

	automationScriptsFinder interface {
		Watch(ctx context.Context)
		FindRunnableScripts(resource, event string, cc ...automation.TriggerConditionChecker) automation.ScriptSet
	}

	AutomationRunnerOpt struct {
		ApiBaseURLSystem    string
		ApiBaseURLMessaging string
		ApiBaseURLCompose   string
	}
)

func AutomationRunner(opt AutomationRunnerOpt, f automationScriptsFinder, r corredor.ScriptRunnerClient) automationRunner {
	var svc = automationRunner{
		opt: opt,

		userFinder:   DefaultUser,
		scriptFinder: f,
		runner:       r,

		logger:     DefaultLogger.Named("automationRunner"),
		jwtEncoder: intAuth.DefaultJwtHandler,
	}

	return svc
}

func (svc automationRunner) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc automationRunner) Watch(ctx context.Context) {
	svc.scriptFinder.Watch(ctx)
}

func (svc automationRunner) OnReceiveMailMessage(ctx context.Context, mail *types.MailMessage) error {
	return svc.findMailScripts(mail.Header).Walk(
		svc.makeMailScriptRunner(ctx, mail),
	)
}

// Finds all scripts that can process email
func (svc automationRunner) findMailScripts(headers types.MailMessageHeader) automation.ScriptSet {
	uev := func(email string) bool {
		u, err := svc.userFinder.FindByEmail(email)
		return u != nil && err == nil
	}

	ss, _ := svc.scriptFinder.
		FindRunnableScripts("system:mail", "onReceive", mailTrigger.MakeChecker(headers, uev)).
		Filter(func(script *automation.Script) (bool, error) {
			// Filter out user-agent scripts && scripts w/o defined runner.
			return !script.RunInUA && script.RunAsDefined(), nil
		})

	return ss
}

func (svc automationRunner) RecordScriptTester(ctx context.Context, source string, payload interface{}) (err error) {
	// Make record script runner
	// @todo figure out how to convert payload to *types.MailMessage
	// runner := svc.makeMailScriptRunner(ctx, payload)
	//
	// return runner(&automation.Script{
	// 	ID:        0,
	// 	Name:      "test",
	// 	SourceRef: "test",
	// 	Source:    source,
	// 	Async:     false,
	// 	RunAs:     0,
	// 	RunInUA:   false,
	// 	Timeout:   0,
	// 	Critical:  true,
	// 	Enabled:   false,
	// })
	return repository.ErrNotImplemented
}

// Runs record script
//
// We set-up script-running environment: security (definer / invoker), async, critical
// and copying values from the run to the given Record
//
func (svc automationRunner) makeMailScriptRunner(ctx context.Context, mail *types.MailMessage) func(script *automation.Script) error {
	// Static request params (record gets updated
	var req = &corredor.RunMailMessageRequest{
		MailMessage: proto.NewMailMessage(mail),
	}

	svc.log(ctx).Debug("preparing mail script runner", zap.Any("mail", mail))

	return func(script *automation.Script) error {
		if svc.runner == nil {
			return errors.New("can not run corredor script: not connected")
		}

		// This could be executed in a goroutine (by *after triggers,
		// so we need to rewire the sentry panic recovery
		defer sentry.Recover()

		ctx, cancelFn := context.WithTimeout(ctx, time.Second*5)
		defer cancelFn()

		// Add invoker's or defined credentials/jwt
		req.Config = map[string]string{
			"api.jwt": svc.getJWT(ctx, script),

			// Let the script know where the API is
			"api.baseURL.system":    svc.opt.ApiBaseURLSystem,
			"api.baseURL.compose":   svc.opt.ApiBaseURLCompose,
			"api.baseURL.messaging": svc.opt.ApiBaseURLMessaging,
		}

		// Add script info
		req.Script = corredor.FromScript(script)

		_, err := svc.runner.MailMessage(ctx, req, grpc.WaitForReady(script.Critical))

		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				svc.log(ctx).Error("unexpected error type", zap.Error(err))
				return err
			}

			switch s.Code() {
			case codes.FailedPrecondition:
				// Sent on syntax errors:
				err = errors.New(s.Message())
			case codes.Aborted:
				err = errors.New(s.Message())
			case codes.InvalidArgument:
				err = errors.New("invalid argument")
			case codes.Internal:
				err = errors.New("internal corredor error")
			default:
			}

			svc.log(ctx).Info("script executed with errors", zap.Error(err))

			if !script.Critical {
				// This was not a critical call and we do not care about
				// errors from script running service.
				return nil
			}

			return err
		}

		return nil
	}
}

// Creates a new JWT for
func (svc automationRunner) getJWT(ctx context.Context, script *automation.Script) string {
	if script.RunAsDefined() {
		return script.Credentials()
	}

	return svc.jwtEncoder.Encode(intAuth.GetIdentityFromContext(ctx))
}
