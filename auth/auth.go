package auth

import (
	"context"
	"embed"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/cortezaproject/corteza-server/auth/external"
	"github.com/cortezaproject/corteza-server/auth/handlers"
	"github.com/cortezaproject/corteza-server/auth/oauth2"
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"github.com/cortezaproject/corteza-server/store"
	systemService "github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi"
	oauth2def "github.com/go-oauth2/oauth2/v4"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	service struct {
		handlers *handlers.AuthHandlers
		log      *zap.Logger
		opt      options.AuthOpt
		settings *settings.Settings
		store    store.Storer
	}
)

//go:embed assets/public
var publicAssets embed.FS

// New initializes Auth service that orchestrates session manager, oauth2 manager and http request handlers
func New(ctx context.Context, log *zap.Logger, s store.Storer, opt options.AuthOpt) (svc *service, err error) {
	var (
		tpls      templateExecutor
		defClient *types.AuthClient
	)

	log = log.Named("auth")
	ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_Auth)

	svc = &service{
		opt:      opt,
		log:      log,
		store:    s,
		settings: &settings.Settings{ /* all disabled by default. */ },
	}

	// use modified log ger for the resrt
	if opt.LogEnabled {
		log = log.WithOptions(zap.AddStacktrace(zap.PanicLevel))
	} else {
		log = zap.NewNop()
	}

	sesManager := request.NewSessionManager(s, opt, log)

	oauth2Manager := oauth2.NewManager(
		opt,
		&oauth2.ContextClientStore{},
		&oauth2.CortezaTokenStore{Store: s},
	)

	oauth2Server := oauth2.NewServer(oauth2Manager)

	// Called after oauth2 authorization request is validated
	// We'll try to get valid user out of the session or redirect user to login page
	oauth2Server.SetUserAuthorizationHandler(oauth2.NewUserAuthorizer(
		sesManager,
		handlers.GetLinks().Login,
		handlers.GetLinks().OAuth2AuthorizeClient,
	))

	oauth2Server.SetClientAuthorizedHandler(func(id string, grant oauth2def.GrantType) (bool, error) {
		// this is a bit silly and a bad design of the oauth2 server lib
		// why do we need to keep on load the client??
		var (
			clientID uint64
			client   *types.AuthClient
			err      error
		)

		clientID, err = strconv.ParseUint(id, 10, 64)
		if err != nil {
			return false, fmt.Errorf("could not authorize client: %w", err)
		}

		client, err = store.LookupAuthClientByID(ctx, s, clientID)
		if err != nil {
			return false, fmt.Errorf("could not authorize client: %w", err)
		}

		// each client only has 1 valid grant type (+ refresh_token)!
		if client.ValidGrant != grant.String() && oauth2def.Refreshing != grant {
			return false, fmt.Errorf("client does not support %s flow", grant)
		}

		return true, nil
	})

	oauth2Server.SetClientScopeHandler(func(id, ss string) (allowed bool, err error) {
		// this is a bit silly and a bad design of the oauth2 server lib
		// why do we need to keep on load the client??
		var (
			clientID uint64
			client   *types.AuthClient
		)

		clientID, err = strconv.ParseUint(id, 10, 64)
		if err != nil {
			return false, fmt.Errorf("could not authorize client: %w", err)
		}

		client, err = store.LookupAuthClientByID(ctx, s, clientID)
		if err != nil {
			return false, fmt.Errorf("could not authorize client: %w", err)
		}

		// ensure all requested scopes are allowed on a client
		for _, scope := range strings.Split(ss, " ") {
			if !auth.CheckScope(client.Scope, scope) {
				return false, fmt.Errorf("client does not allow use of '%s' scope", scope)
			}
		}

		return true, nil
	})

	oauth2Server.SetExtensionFieldsHandler(func(ti oauth2def.TokenInfo) (fieldsValue map[string]interface{}) {
		fieldsValue = make(map[string]interface{})
		handlers.SubSplit(ti, fieldsValue)
		fieldsValue["refresh_token_expires_in"] = int(ti.GetRefreshExpiresIn() / time.Second)
		if err = handlers.Profile(ctx, ti, fieldsValue); err != nil {
			log.Error("failed to add profile data", zap.Error(err))
		}

		return
	})

	if opt.DefaultClient != "" {
		// default client will help streamline authorization with default clients
		defClient, err = store.LookupAuthClientByHandle(ctx, s, opt.DefaultClient)
		if err != nil {
			return nil, fmt.Errorf("cannot load default client: %w", err)
		}
	}

	var (
		tplBase = template.New("").
			Funcs(sprig.FuncMap()).
			Funcs(template.FuncMap{
				"version":   func() string { return version.Version },
				"buildtime": func() string { return version.BuildTime },
				"links":     handlers.GetLinks,
			})
		tplLoader templateLoader
	)

	if len(opt.AssetsPath) > 0 {
		tplLoader = func(t *template.Template) (tpl *template.Template, err error) {
			if tpl, err = t.Clone(); err != nil {
				return nil, fmt.Errorf("can not clone templates: %w", err)
			} else {
				return tpl.ParseGlob(opt.AssetsPath + "/templates/*.tpl")
			}
		}
		log.Info("loading assets from filesystem", zap.String("path", opt.AssetsPath))
	} else {
		tplLoader = EmbeddedTemplates
		log.Info("using embedded assets")
	}

	if !opt.DevelopmentMode || len(opt.AssetsPath) == 0 {
		log.Info("initializing templates without reloading (production mode)")
		tpls, err = NewStaticTemplates(tplBase, tplLoader)
		if err != nil {
			return nil, fmt.Errorf("can not load templates: %w", err)
		}
	} else {
		log.Info("initializing reloadable templates (development mode)")
		tpls = NewReloadableTemplates(tplBase, tplLoader)
	}

	svc.handlers = &handlers.AuthHandlers{
		Log:            log,
		Templates:      tpls,
		SessionManager: sesManager,
		OAuth2:         oauth2Server,
		AuthService:    systemService.DefaultAuth,
		UserService:    systemService.DefaultUser,
		ClientService:  &clientService{s},
		TokenService:   &tokenService{s},
		DefaultClient:  defClient,
		Opt:            svc.opt,
		Settings:       svc.settings,
	}

	external.Init(log, sesManager.Store())

	return
}

func (svc *service) UpdateSettings(s *settings.Settings) {
	if svc.settings.LocalEnabled != s.LocalEnabled {
		svc.log.Debug("setting changed", zap.Bool("localEnabled", s.LocalEnabled))
	}

	if svc.settings.SignupEnabled != s.SignupEnabled {
		svc.log.Debug("setting changed", zap.Bool("signupEnabled", s.SignupEnabled))
	}

	if svc.settings.EmailConfirmationRequired != s.EmailConfirmationRequired {
		svc.log.Debug("setting changed", zap.Bool("emailConfirmationRequired", s.EmailConfirmationRequired))
	}

	if svc.settings.PasswordResetEnabled != s.PasswordResetEnabled {
		svc.log.Debug("setting changed", zap.Bool("passwordResetEnabled", s.PasswordResetEnabled))
	}

	if svc.settings.ExternalEnabled != s.ExternalEnabled {
		svc.log.Debug("setting changed", zap.Bool("externalEnabled", s.ExternalEnabled))
	}

	if svc.settings.MultiFactor != s.MultiFactor {
		svc.log.Debug("setting changed", zap.Any("mfa", s.MultiFactor))
	}

	if len(svc.settings.Providers) != len(s.Providers) {
		svc.log.Debug("setting changed", zap.Int("providers", len(s.Providers)))
		external.SetupGothProviders(svc.opt.ExternalRedirectURL, s.Providers...)
	}

	svc.settings = s
	svc.handlers.Settings = s
}

func (svc *service) Watch(ctx context.Context) {
	go svc.gc(ctx)
}

func (svc service) gc(ctx context.Context) {
	svc.log.Info("running startup garbage collection")
	go svc.gcSessions(ctx)
	go svc.gcOAuth2Tokens(ctx)

	i := svc.opt.GarbageCollectorInterval
	if i < time.Minute {
		svc.log.Warn("garbage collection interval less than 1 minute, disabling")
	} else {
		svc.log.Info("starting garbage collecting process", zap.Duration("interval", i))
	}

	tck := time.NewTicker(i)
	for {
		select {
		case <-ctx.Done():
			svc.log.Info("stopping gc", zap.Error(ctx.Err()))
			return
		case <-tck.C:
			svc.log.Info("garbage collector")
			go svc.gcSessions(ctx)
			go svc.gcOAuth2Tokens(ctx)
			return
		}

	}
}

func (svc service) gcSessions(ctx context.Context) {
	err := store.DeleteExpiredAuthSessions(ctx, svc.store)
	if err != nil {
		svc.log.Error("failed to collect session garbage", zap.Error(err))
	}
}
func (svc service) gcOAuth2Tokens(ctx context.Context) {
	err := store.DeleteExpiredAuthOA2Tokens(ctx, svc.store)
	if err != nil {
		svc.log.Error("failed to collect oauth2 token garbage", zap.Error(err))
	}
}

func (svc service) MountHttpRoutes(r chi.Router) {
	svc.handlers.MountHttpRoutes(r)

	const uriRoot = "/auth/assets/public"
	if len(svc.opt.AssetsPath) == 0 {
		r.Handle(uriRoot+"/*", http.StripPrefix("/auth/assets", http.FileServer(http.FS(publicAssets))))
	} else {
		var root = strings.TrimRight(svc.opt.AssetsPath, "/") + "/public"
		r.Handle(uriRoot+"/*", http.StripPrefix(uriRoot, http.FileServer(http.Dir(root))))
	}
}

//func (svc service) WellKnownOpenIDConfiguration() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		json.NewEncoder(w).Encode(map[string]interface{}{
//			"issuer":                                svc.opt.BaseURL,
//			"authorization_endpoint":                svc.opt.BaseURL + "/oauth2/authorize",
//			"token_endpoint":                        svc.opt.BaseURL + "/oauth2/token",
//			"jwks_uri":                              svc.opt.BaseURL + "/oauth2/public-keys", // @todo
//			"subject_types_supported":               []string{"public"},
//			"response_types_supported":              []string{"public"},
//			"id_token_signing_alg_values_supported": []string{"RS256", "HS512"},
//		})
//
//		w.Header().Set("Content-Type", "application/json")
//	}
//}
