package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"time"
)

type (
	DBOpt struct {
		DSN string `env:"DB_DSN"`
	}

	HTTPClientOpt struct {
		Timeout     time.Duration `env:"HTTP_CLIENT_TIMEOUT"`
		TlsInsecure bool          `env:"HTTP_CLIENT_TLS_INSECURE"`
	}

	HTTPServerOpt struct {
		Addr                   string `env:"HTTP_ADDR"`
		ApiBaseUrl             string `env:"HTTP_API_BASE_URL"`
		ApiEnabled             bool   `env:"HTTP_API_ENABLED"`
		BaseUrl                string `env:"HTTP_BASE_URL"`
		EnableDebugRoute       bool   `env:"HTTP_ENABLE_DEBUG_ROUTE"`
		EnableHealthcheckRoute bool   `env:"HTTP_ENABLE_HEALTHCHECK_ROUTE"`
		EnableMetrics          bool   `env:"HTTP_METRICS"`
		EnablePanicReporting   bool   `env:"HTTP_REPORT_PANIC"`
		EnableVersionRoute     bool   `env:"HTTP_ENABLE_VERSION_ROUTE"`
		LogRequest             bool   `env:"HTTP_LOG_REQUEST"`
		LogResponse            bool   `env:"HTTP_LOG_RESPONSE"`
		MetricsPassword        string `env:"HTTP_METRICS_PASSWORD"`
		MetricsServiceLabel    string `env:"HTTP_METRICS_NAME"`
		MetricsUsername        string `env:"HTTP_METRICS_USERNAME"`
		SslTerminated          bool   `env:"HTTP_SSL_TERMINATED"`
		Tracing                bool   `env:"HTTP_ERROR_TRACING"`
		WebappBaseDir          string `env:"HTTP_WEBAPP_BASE_DIR"`
		WebappBaseUrl          string `env:"HTTP_WEBAPP_BASE_URL"`
		WebappEnabled          bool   `env:"HTTP_WEBAPP_ENABLED"`
		WebappList             string `env:"HTTP_WEBAPP_LIST"`
	}

	RBACOpt struct {
		AnonymousRoles     string `env:"RBAC_ANONYMOUS_ROLES"`
		AuthenticatedRoles string `env:"RBAC_AUTHENTICATED_ROLES"`
		BypassRoles        string `env:"RBAC_BYPASS_ROLES"`
		Log                bool   `env:"RBAC_LOG"`
		ServiceUser        string `env:"RBAC_SERVICE_USER"`
	}

	SCIMOpt struct {
		BaseURL              string `env:"SCIM_BASE_URL"`
		Enabled              bool   `env:"SCIM_ENABLED"`
		ExternalIdAsPrimary  bool   `env:"SCIM_EXTERNAL_ID_AS_PRIMARY"`
		ExternalIdValidation string `env:"SCIM_EXTERNAL_ID_VALIDATION"`
		Secret               string `env:"SCIM_SECRET"`
	}

	SMTPOpt struct {
		From          string `env:"SMTP_FROM"`
		Host          string `env:"SMTP_HOST"`
		Pass          string `env:"SMTP_PASS"`
		Port          int    `env:"SMTP_PORT"`
		TlsInsecure   bool   `env:"SMTP_TLS_INSECURE"`
		TlsServerName string `env:"SMTP_TLS_SERVER_NAME"`
		User          string `env:"SMTP_USER"`
	}

	ActionLogOpt struct {
		Debug                    bool `env:"ACTIONLOG_DEBUG"`
		Enabled                  bool `env:"ACTIONLOG_ENABLED"`
		WorkflowFunctionsEnabled bool `env:"ACTIONLOG_WORKFLOW_FUNCTIONS_ENABLED"`
	}

	ApigwOpt struct {
		Debug                bool          `env:"APIGW_DEBUG"`
		Enabled              bool          `env:"APIGW_ENABLED"`
		LogEnabled           bool          `env:"APIGW_LOG_ENABLED"`
		LogRequestBody       bool          `env:"APIGW_LOG_REQUEST_BODY"`
		ProxyEnableDebugLog  bool          `env:"APIGW_PROXY_ENABLE_DEBUG_LOG"`
		ProxyFollowRedirects bool          `env:"APIGW_PROXY_FOLLOW_REDIRECTS"`
		ProxyOutboundTimeout time.Duration `env:"APIGW_PROXY_OUTBOUND_TIMEOUT"`
	}

	AuthOpt struct {
		AccessTokenLifetime      time.Duration `env:"AUTH_OAUTH2_ACCESS_TOKEN_LIFETIME"`
		AssetsPath               string        `env:"AUTH_ASSETS_PATH"`
		BaseURL                  string        `env:"AUTH_BASE_URL"`
		CsrfCookieName           string        `env:"AUTH_CSRF_COOKIE_NAME"`
		CsrfEnabled              bool          `env:"AUTH_CSRF_ENABLED"`
		CsrfFieldName            string        `env:"AUTH_CSRF_FIELD_NAME"`
		CsrfSecret               string        `env:"AUTH_CSRF_SECRET"`
		DefaultClient            string        `env:"AUTH_DEFAULT_CLIENT"`
		DevelopmentMode          bool          `env:"AUTH_DEVELOPMENT_MODE"`
		Expiry                   time.Duration `env:"AUTH_JWT_EXPIRY"`
		ExternalCookieSecret     string        `env:"AUTH_EXTERNAL_COOKIE_SECRET"`
		ExternalRedirectURL      string        `env:"AUTH_EXTERNAL_REDIRECT_URL"`
		GarbageCollectorInterval time.Duration `env:"AUTH_GARBAGE_COLLECTOR_INTERVAL"`
		LogEnabled               bool          `env:"AUTH_LOG_ENABLED"`
		PasswordSecurity         bool          `env:"AUTH_PASSWORD_SECURITY"`
		RefreshTokenLifetime     time.Duration `env:"AUTH_OAUTH2_REFRESH_TOKEN_LIFETIME"`
		RequestRateLimit         int           `env:"AUTH_REQUEST_RATE_LIMIT"`
		RequestRateWindowLength  time.Duration `env:"AUTH_REQUEST_RATE_WINDOW_LENGTH"`
		Secret                   string        `env:"AUTH_JWT_SECRET"`
		SessionCookieDomain      string        `env:"AUTH_SESSION_COOKIE_DOMAIN"`
		SessionCookieName        string        `env:"AUTH_SESSION_COOKIE_NAME"`
		SessionCookiePath        string        `env:"AUTH_SESSION_COOKIE_PATH"`
		SessionCookieSecure      bool          `env:"AUTH_SESSION_COOKIE_SECURE"`
		SessionLifetime          time.Duration `env:"AUTH_SESSION_LIFETIME"`
		SessionPermLifetime      time.Duration `env:"AUTH_SESSION_PERM_LIFETIME"`
	}

	CorredorOpt struct {
		Addr                  string        `env:"CORREDOR_ADDR"`
		DefaultExecTimeout    time.Duration `env:"CORREDOR_DEFAULT_EXEC_TIMEOUT"`
		Enabled               bool          `env:"CORREDOR_ENABLED"`
		ListRefresh           time.Duration `env:"CORREDOR_LIST_REFRESH"`
		ListTimeout           time.Duration `env:"CORREDOR_LIST_TIMEOUT"`
		MaxBackoffDelay       time.Duration `env:"CORREDOR_MAX_BACKOFF_DELAY"`
		MaxReceiveMessageSize int           `env:"CORREDOR_MAX_RECEIVE_MESSAGE_SIZE"`
		RunAsEnabled          bool          `env:"CORREDOR_RUN_AS_ENABLED"`
		TlsCertCA             string        `env:"CORREDOR_CLIENT_CERTIFICATES_CA"`
		TlsCertEnabled        bool          `env:"CORREDOR_CLIENT_CERTIFICATES_ENABLED"`
		TlsCertPath           string        `env:"CORREDOR_CLIENT_CERTIFICATES_PATH"`
		TlsCertPrivate        string        `env:"CORREDOR_CLIENT_CERTIFICATES_PRIVATE"`
		TlsCertPublic         string        `env:"CORREDOR_CLIENT_CERTIFICATES_PUBLIC"`
		TlsServerName         string        `env:"CORREDOR_CLIENT_CERTIFICATES_SERVER_NAME"`
	}

	EnvironmentOpt struct {
		Environment string `env:"ENVIRONMENT"`
	}

	EventbusOpt struct {
		SchedulerEnabled  bool          `env:"EVENTBUS_SCHEDULER_ENABLED"`
		SchedulerInterval time.Duration `env:"EVENTBUS_SCHEDULER_INTERVAL"`
	}

	FederationOpt struct {
		DataMonitorInterval      time.Duration `env:"FEDERATION_SYNC_DATA_MONITOR_INTERVAL"`
		DataPageSize             int           `env:"FEDERATION_SYNC_DATA_PAGE_SIZE"`
		Enabled                  bool          `env:"FEDERATION_ENABLED"`
		Host                     string        `env:"FEDERATION_HOST"`
		Label                    string        `env:"FEDERATION_LABEL"`
		StructureMonitorInterval time.Duration `env:"FEDERATION_SYNC_STRUCTURE_MONITOR_INTERVAL"`
		StructurePageSize        int           `env:"FEDERATION_SYNC_STRUCTURE_PAGE_SIZE"`
	}

	LimitOpt struct {
		SystemUsers int `env:"LIMIT_SYSTEM_USERS"`
	}

	LocaleOpt struct {
		DevelopmentMode             bool   `env:"LOCALE_DEVELOPMENT_MODE"`
		Languages                   string `env:"LOCALE_LANGUAGES"`
		Log                         bool   `env:"LOCALE_LOG"`
		Path                        string `env:"LOCALE_PATH"`
		QueryStringParam            string `env:"LOCALE_QUERY_STRING_PARAM"`
		ResourceTranslationsEnabled bool   `env:"LOCALE_RESOURCE_TRANSLATIONS_ENABLED"`
	}

	LogOpt struct {
		Debug           bool   `env:"LOG_DEBUG"`
		Filter          string `env:"LOG_FILTER"`
		IncludeCaller   bool   `env:"LOG_INCLUDE_CALLER"`
		Level           string `env:"LOG_LEVEL"`
		StacktraceLevel string `env:"LOG_STACKTRACE_LEVEL"`
	}

	MessagebusOpt struct {
		Enabled    bool `env:"MESSAGEBUS_ENABLED"`
		LogEnabled bool `env:"MESSAGEBUS_LOG_ENABLED"`
	}

	MonitorOpt struct {
		Interval time.Duration `env:"MONITOR_INTERVAL"`
	}

	ObjectStoreOpt struct {
		MinioAccessKey  string `env:"MINIO_ACCESS_KEY"`
		MinioBucket     string `env:"MINIO_BUCKET"`
		MinioEndpoint   string `env:"MINIO_ENDPOINT"`
		MinioPathPrefix string `env:"MINIO_PATH_PREFIX"`
		MinioSSECKey    string `env:"MINIO_SSEC_KEY"`
		MinioSecretKey  string `env:"MINIO_SECRET_KEY"`
		MinioSecure     bool   `env:"MINIO_SECURE"`
		MinioStrict     bool   `env:"MINIO_STRICT"`
		Path            string `env:"STORAGE_PATH"`
	}

	PluginsOpt struct {
		Enabled bool   `env:"PLUGINS_ENABLED"`
		Paths   string `env:"PLUGINS_PATHS"`
	}

	ProvisionOpt struct {
		Always bool   `env:"PROVISION_ALWAYS"`
		Path   string `env:"PROVISION_PATH"`
	}

	SeederOpt struct {
		LogEnabled bool `env:"SEEDER_LOG_ENABLED"`
	}

	SentryOpt struct {
		DSN              string  `env:"SENTRY_DSN"`
		AttachStacktrace bool    `env:"SENTRY_ATTACH_STACKTRACE"`
		Debug            bool    `env:"SENTRY_DEBUG"`
		Dist             string  `env:"SENTRY_DIST"`
		Environment      string  `env:"SENTRY_ENVIRONMENT"`
		MaxBreadcrumbs   int     `env:"SENTRY_MAX_BREADCRUMBS"`
		Release          string  `env:"SENTRY_RELEASE"`
		SampleRate       float64 `env:"SENTRY_SAMPLE_RATE"`
		ServerName       string  `env:"SENTRY_SERVERNAME"`
	}

	TemplateOpt struct {
		RendererGotenbergAddress string `env:"TEMPLATE_RENDERER_GOTENBERG_ADDRESS"`
		RendererGotenbergEnabled bool   `env:"TEMPLATE_RENDERER_GOTENBERG_ENABLED"`
	}

	UpgradeOpt struct {
		Always bool `env:"UPGRADE_ALWAYS"`
		Debug  bool `env:"UPGRADE_DEBUG"`
	}

	WaitForOpt struct {
		Delay                 time.Duration `env:"WAIT_FOR"`
		Services              string        `env:"WAIT_FOR_SERVICES"`
		ServicesProbeInterval time.Duration `env:"WAIT_FOR_SERVICES_PROBE_INTERVAL"`
		ServicesProbeTimeout  time.Duration `env:"WAIT_FOR_SERVICES_PROBE_TIMEOUT"`
		ServicesTimeout       time.Duration `env:"WAIT_FOR_SERVICES_TIMEOUT"`
		StatusPage            bool          `env:"WAIT_FOR_STATUS_PAGE"`
	}

	WebsocketOpt struct {
		LogEnabled  bool          `env:"WEBSOCKET_LOG_ENABLED"`
		PingPeriod  time.Duration `env:"WEBSOCKET_PING_PERIOD"`
		PingTimeout time.Duration `env:"WEBSOCKET_PING_TIMEOUT"`
		Timeout     time.Duration `env:"WEBSOCKET_TIMEOUT"`
	}

	WorkflowOpt struct {
		CallStackSize int  `env:"WORKFLOW_CALL_STACK_SIZE"`
		ExecDebug     bool `env:"WORKFLOW_EXEC_DEBUG"`
		Register      bool `env:"WORKFLOW_REGISTER"`
	}
)

// DB initializes and returns a DBOpt with default values
//
// This function is auto-generated
func DB() (o *DBOpt) {
	o = &DBOpt{
		DSN: "sqlite3://file::memory:?cache=shared&mode=memory",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *DBOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// HTTPClient initializes and returns a HTTPClientOpt with default values
//
// This function is auto-generated
func HTTPClient() (o *HTTPClientOpt) {
	o = &HTTPClientOpt{
		Timeout:     30 * time.Second,
		TlsInsecure: false,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *HTTPClientOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// HTTPServer initializes and returns a HTTPServerOpt with default values
//
// This function is auto-generated
func HTTPServer() (o *HTTPServerOpt) {
	o = &HTTPServerOpt{
		Addr:                   ":80",
		ApiBaseUrl:             "/",
		ApiEnabled:             true,
		BaseUrl:                "/",
		EnableDebugRoute:       false,
		EnableHealthcheckRoute: true,
		EnableMetrics:          false,
		EnablePanicReporting:   true,
		EnableVersionRoute:     true,
		LogRequest:             false,
		LogResponse:            false,
		MetricsPassword:        string(rand.Bytes(5)),
		MetricsServiceLabel:    "corteza",
		MetricsUsername:        "metrics",
		SslTerminated:          isSecure(),
		Tracing:                false,
		WebappBaseDir:          "./webapp/public",
		WebappBaseUrl:          "/",
		WebappEnabled:          false,
		WebappList:             "admin,compose,workflow,reporter",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *HTTPServerOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// RBAC initializes and returns a RBACOpt with default values
//
// This function is auto-generated
func RBAC() (o *RBACOpt) {
	o = &RBACOpt{
		AnonymousRoles:     "anonymous",
		AuthenticatedRoles: "authenticated",
		BypassRoles:        "super-admin",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *RBACOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// SCIM initializes and returns a SCIMOpt with default values
//
// This function is auto-generated
func SCIM() (o *SCIMOpt) {
	o = &SCIMOpt{
		BaseURL:              "/scim",
		ExternalIdValidation: "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *SCIMOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// SMTP initializes and returns a SMTPOpt with default values
//
// This function is auto-generated
func SMTP() (o *SMTPOpt) {
	o = &SMTPOpt{
		Host:        "localhost",
		Port:        25,
		TlsInsecure: false,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *SMTPOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// ActionLog initializes and returns a ActionLogOpt with default values
//
// This function is auto-generated
func ActionLog() (o *ActionLogOpt) {
	o = &ActionLogOpt{
		Debug:                    false,
		Enabled:                  true,
		WorkflowFunctionsEnabled: false,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *ActionLogOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Apigw initializes and returns a ApigwOpt with default values
//
// This function is auto-generated
func Apigw() (o *ApigwOpt) {
	o = &ApigwOpt{
		Debug:                false,
		Enabled:              true,
		LogEnabled:           false,
		LogRequestBody:       false,
		ProxyEnableDebugLog:  false,
		ProxyFollowRedirects: true,
		ProxyOutboundTimeout: time.Second * 30,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *ApigwOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Auth initializes and returns a AuthOpt with default values
//
// This function is auto-generated
func Auth() (o *AuthOpt) {
	o = &AuthOpt{
		AccessTokenLifetime:      time.Hour * 2,
		BaseURL:                  fullURL("/auth"),
		CsrfCookieName:           "same-site-authenticity-token",
		CsrfEnabled:              true,
		CsrfFieldName:            "same-site-authenticity-token",
		CsrfSecret:               getSecretFromEnv("csrf secret"),
		DefaultClient:            "corteza-webapp",
		Expiry:                   time.Hour * 24 * 30,
		ExternalCookieSecret:     getSecretFromEnv("external cookie secret"),
		ExternalRedirectURL:      fullURL("/auth/external/{provider}/callback"),
		GarbageCollectorInterval: 15 * time.Minute,
		PasswordSecurity:         true,
		RefreshTokenLifetime:     time.Hour * 24 * 3,
		RequestRateLimit:         60,
		RequestRateWindowLength:  time.Minute,
		Secret:                   getSecretFromEnv("jwt secret"),
		SessionCookieDomain:      guessHostname(),
		SessionCookieName:        "session",
		SessionCookiePath:        pathPrefix("/auth"),
		SessionCookieSecure:      isSecure(),
		SessionLifetime:          24 * time.Hour,
		SessionPermLifetime:      360 * 24 * time.Hour,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *AuthOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Corredor initializes and returns a CorredorOpt with default values
//
// This function is auto-generated
func Corredor() (o *CorredorOpt) {
	o = &CorredorOpt{
		Addr:                  "localhost:50051",
		DefaultExecTimeout:    time.Minute,
		Enabled:               false,
		ListRefresh:           time.Second * 5,
		ListTimeout:           time.Second * 2,
		MaxBackoffDelay:       time.Minute,
		MaxReceiveMessageSize: 2 << 23,
		RunAsEnabled:          true,
		TlsCertCA:             "ca.crt",
		TlsCertEnabled:        false,
		TlsCertPath:           "/certs/corredor/client",
		TlsCertPrivate:        "private.key",
		TlsCertPublic:         "public.crt",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *CorredorOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Environment initializes and returns a EnvironmentOpt with default values
//
// This function is auto-generated
func Environment() (o *EnvironmentOpt) {
	o = &EnvironmentOpt{
		Environment: "production",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *EnvironmentOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Eventbus initializes and returns a EventbusOpt with default values
//
// This function is auto-generated
func Eventbus() (o *EventbusOpt) {
	o = &EventbusOpt{
		SchedulerEnabled:  true,
		SchedulerInterval: time.Minute,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *EventbusOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Federation initializes and returns a FederationOpt with default values
//
// This function is auto-generated
func Federation() (o *FederationOpt) {
	o = &FederationOpt{
		DataMonitorInterval:      time.Second * 60,
		DataPageSize:             100,
		Enabled:                  false,
		Host:                     "local.cortezaproject.org",
		Label:                    "federated",
		StructureMonitorInterval: time.Minute * 2,
		StructurePageSize:        1,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *FederationOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Limit initializes and returns a LimitOpt with default values
//
// This function is auto-generated
func Limit() (o *LimitOpt) {
	o = &LimitOpt{}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *LimitOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Locale initializes and returns a LocaleOpt with default values
//
// This function is auto-generated
func Locale() (o *LocaleOpt) {
	o = &LocaleOpt{
		Languages:        "en",
		QueryStringParam: "lng",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *LocaleOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Log initializes and returns a LogOpt with default values
//
// This function is auto-generated
func Log() (o *LogOpt) {
	o = &LogOpt{
		IncludeCaller:   false,
		Level:           "warn",
		StacktraceLevel: "dpanic",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *LogOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Messagebus initializes and returns a MessagebusOpt with default values
//
// This function is auto-generated
func Messagebus() (o *MessagebusOpt) {
	o = &MessagebusOpt{
		Enabled:    true,
		LogEnabled: false,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *MessagebusOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Monitor initializes and returns a MonitorOpt with default values
//
// This function is auto-generated
func Monitor() (o *MonitorOpt) {
	o = &MonitorOpt{
		Interval: 300 * time.Second,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *MonitorOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// ObjectStore initializes and returns a ObjectStoreOpt with default values
//
// This function is auto-generated
func ObjectStore() (o *ObjectStoreOpt) {
	o = &ObjectStoreOpt{
		MinioBucket: "{component}",
		MinioSecure: true,
		MinioStrict: false,
		Path:        "var/store",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *ObjectStoreOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Plugins initializes and returns a PluginsOpt with default values
//
// This function is auto-generated
func Plugins() (o *PluginsOpt) {
	o = &PluginsOpt{
		Enabled: true,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *PluginsOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Provision initializes and returns a ProvisionOpt with default values
//
// This function is auto-generated
func Provision() (o *ProvisionOpt) {
	o = &ProvisionOpt{
		Always: true,
		Path:   "provision/*",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *ProvisionOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Seeder initializes and returns a SeederOpt with default values
//
// This function is auto-generated
func Seeder() (o *SeederOpt) {
	o = &SeederOpt{}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *SeederOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Sentry initializes and returns a SentryOpt with default values
//
// This function is auto-generated
func Sentry() (o *SentryOpt) {
	o = &SentryOpt{
		AttachStacktrace: true,
		MaxBreadcrumbs:   0,
		Release:          version.Version,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *SentryOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Template initializes and returns a TemplateOpt with default values
//
// This function is auto-generated
func Template() (o *TemplateOpt) {
	o = &TemplateOpt{
		RendererGotenbergEnabled: false,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *TemplateOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Upgrade initializes and returns a UpgradeOpt with default values
//
// This function is auto-generated
func Upgrade() (o *UpgradeOpt) {
	o = &UpgradeOpt{
		Always: true,
		Debug:  false,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *UpgradeOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// WaitFor initializes and returns a WaitForOpt with default values
//
// This function is auto-generated
func WaitFor() (o *WaitForOpt) {
	o = &WaitForOpt{
		Delay:                 0,
		ServicesProbeInterval: time.Second * 5,
		ServicesProbeTimeout:  time.Second * 30,
		ServicesTimeout:       time.Minute,
		StatusPage:            true,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *WaitForOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Websocket initializes and returns a WebsocketOpt with default values
//
// This function is auto-generated
func Websocket() (o *WebsocketOpt) {
	o = &WebsocketOpt{
		PingPeriod:  ((120 * time.Second) * 9) / 10,
		PingTimeout: 120 * time.Second,
		Timeout:     15 * time.Second,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *WebsocketOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

// Workflow initializes and returns a WorkflowOpt with default values
//
// This function is auto-generated
func Workflow() (o *WorkflowOpt) {
	o = &WorkflowOpt{
		CallStackSize: 16,
		ExecDebug:     false,
		Register:      true,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *WorkflowOpt) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
