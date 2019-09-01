package service

import (
	"context"
	"encoding/json"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	intAuth "github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/automation/corredor"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/cortezaproject/corteza-server/system/proto"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	automationRunner struct {
		opt          AutomationRunnerOpt
		logger       *zap.Logger
		runner       corredor.ScriptRunnerClient
		scriptFinder automationScriptsFinder
		jwtEncoder   intAuth.TokenEncoder
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

	TriggerCondition struct {
		MatchAll bool                            `json:"matchAll"`
		Headers  []TriggerConditionHeaderMatcher `json:"headers"`
	}

	TriggerConditionHeaderMatcher struct {
		Name  string `json:"name"`
		Match string `json:"match"`
		Op    string `json:"op"`
	}
)

const (
	AutomationResourceRecord = "compose:record"
)

func AutomationRunner(opt AutomationRunnerOpt, f automationScriptsFinder, r corredor.ScriptRunnerClient) automationRunner {
	var svc = automationRunner{
		opt: opt,

		scriptFinder: f,
		runner:       r,

		logger:     DefaultLogger.Named("automationRunner"),
		jwtEncoder: intAuth.DefaultJwtHandler,
	}

	return svc
}

func (svc automationRunner) Watch(ctx context.Context) {
	svc.scriptFinder.Watch(ctx)
}

func (svc automationRunner) OnReceivedMailMessage(ctx context.Context, mail *types.MailMessage) error {
	return svc.findMailScripts(mail.Header).Walk(
		svc.makeMailScriptRunner(ctx, mail),
	)
}

// Finds all scripts that can process email
func (svc automationRunner) findMailScripts(headers types.MailMessageHeader) automation.ScriptSet {
	ss, _ := svc.scriptFinder.FindRunnableScripts("system:mail", "onReceived", svc.makeMailHeaderChecker(headers)).
		Filter(func(script *automation.Script) (bool, error) {
			// Filter out user-agent scripts
			return !script.RunInUA, nil
		})

	return ss
}

func (svc automationRunner) RecordScriptTester(ctx context.Context, source string, mail *types.MailMessage) (err error) {
	// Make record script runner and
	runner := svc.makeMailScriptRunner(ctx, mail)

	return runner(&automation.Script{
		ID:        0,
		Name:      "test",
		SourceRef: "test",
		Source:    source,
		Async:     false,
		RunAs:     0,
		RunInUA:   false,
		Timeout:   0,
		Critical:  true,
		Enabled:   false,
	})
}

// Runs record script
//
// We set-up script-running environment: security (definer / invoker), async, critical
// and copying values from the run to the given Record
//
func (svc automationRunner) makeMailScriptRunner(ctx context.Context, mail *types.MailMessage) func(script *automation.Script) error {
	// Static request params (record gets updated
	var req = &corredor.RunMailMessageRequest{
		MailMessage: proto.FromMailMessage(mail),
	}

	svc.logger.Debug("executing script", zap.Any("mail", mail))

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
				svc.logger.Error("unexpected error type", zap.Error(err))
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

			svc.logger.Info("script executed with errors", zap.Error(err))

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

func (svc automationRunner) makeMailHeaderChecker(headers types.MailMessageHeader) automation.TriggerConditionChecker {
	return func(c string) bool {
		var (
			err   error
			tc    = TriggerCondition{}
			re    *regexp.Regexp
			match bool
		)

		if err := json.Unmarshal([]byte(c), &tc); err != nil {
			panic(err) // @todo replace with log
			return false
		}

		for _, m := range tc.Headers {
			if m.Op == "regex" {
				if re, err = regexp.Compile(m.Match); err != nil {
					// Invalid re
					continue
				}
			}

			for name, vv := range headers.Raw {
				name = strings.ToLower(name)
				if strings.ToLower(m.Name) != name {
					continue
				}

				for _, v := range vv {
					switch name {
					case "from",
						"to",
						"cc",
						"bcc",
						"reply-to":
						a, _ := mail.ParseAddress(v)
						v = a.Address
					}

					switch m.Op {
					case "regex":
						match = re.MatchString(v)
					case "ci":
						match = strings.ToLower(v) == strings.ToLower(m.Match)
					default:
						match = v == m.Match
					}

					if tc.MatchAll && !match {
						return false
					} else if !tc.MatchAll && match {
						return true
					}
				}
			}
		}

		return match
	}
}
