package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/pkg/version"
	"time"
)

type (
	DBOpt struct {
		DSN string `env:"DB_DSN"`
	}

	HTTPClientOpt struct {
		TlsInsecure bool          `env:"HTTP_CLIENT_TLS_INSECURE"`
		Timeout     time.Duration `env:"HTTP_CLIENT_TIMEOUT"`
	}

	HttpServerOpt struct {
		Domain                 string `env:"DOMAIN"`
		DomainWebapp           string `env:"DOMAIN_WEBAPP"`
		Addr                   string `env:"HTTP_ADDR"`
		LogRequest             bool   `env:"HTTP_LOG_REQUEST"`
		LogResponse            bool   `env:"HTTP_LOG_RESPONSE"`
		Tracing                bool   `env:"HTTP_ERROR_TRACING"`
		EnableHealthcheckRoute bool   `env:"HTTP_ENABLE_HEALTHCHECK_ROUTE"`
		EnableVersionRoute     bool   `env:"HTTP_ENABLE_VERSION_ROUTE"`
		EnableDebugRoute       bool   `env:"HTTP_ENABLE_DEBUG_ROUTE"`
		EnableMetrics          bool   `env:"HTTP_METRICS"`
		MetricsServiceLabel    string `env:"HTTP_METRICS_NAME"`
		MetricsUsername        string `env:"HTTP_METRICS_USERNAME"`
		MetricsPassword        string `env:"HTTP_METRICS_PASSWORD"`
		EnablePanicReporting   bool   `env:"HTTP_REPORT_PANIC"`
		BaseUrl                string `env:"HTTP_BASE_URL"`
		ApiEnabled             bool   `env:"HTTP_API_ENABLED"`
		ApiBaseUrl             string `env:"HTTP_API_BASE_URL"`
		WebappEnabled          bool   `env:"HTTP_WEBAPP_ENABLED"`
		WebappBaseUrl          string `env:"HTTP_WEBAPP_BASE_URL"`
		WebappBaseDir          string `env:"HTTP_WEBAPP_BASE_DIR"`
		WebappList             string `env:"HTTP_WEBAPP_LIST"`
		SslTerminated          bool   `env:"HTTP_SSL_TERMINATED"`
		AssetsPath             string `env:"HTTP_SERVER_ASSETS_PATH"`
		WebConsoleEnabled      bool   `env:"HTTP_SERVER_WEB_CONSOLE_ENABLED"`
		WebConsoleUsername     string `env:"HTTP_SERVER_WEB_CONSOLE_USERNAME"`
		WebConsolePassword     string `env:"HTTP_SERVER_WEB_CONSOLE_PASSWORD"`
	}

	RbacOpt struct {
		Log                bool   `env:"RBAC_LOG"`
		ServiceUser        string `env:"RBAC_SERVICE_USER"`
		BypassRoles        string `env:"RBAC_BYPASS_ROLES"`
		AuthenticatedRoles string `env:"RBAC_AUTHENTICATED_ROLES"`
		AnonymousRoles     string `env:"RBAC_ANONYMOUS_ROLES"`
	}

	SCIMOpt struct {
		Enabled              bool   `env:"SCIM_ENABLED"`
		BaseURL              string `env:"SCIM_BASE_URL"`
		Secret               string `env:"SCIM_SECRET"`
		ExternalIdAsPrimary  bool   `env:"SCIM_EXTERNAL_ID_AS_PRIMARY"`
		ExternalIdValidation string `env:"SCIM_EXTERNAL_ID_VALIDATION"`
	}

	SMTPOpt struct {
		Host          string `env:"SMTP_HOST"`
		Port          int    `env:"SMTP_PORT"`
		User          string `env:"SMTP_USER"`
		Pass          string `env:"SMTP_PASS"`
		From          string `env:"SMTP_FROM"`
		TlsInsecure   bool   `env:"SMTP_TLS_INSECURE"`
		TlsServerName string `env:"SMTP_TLS_SERVER_NAME"`
	}

	ActionLogOpt struct {
		Enabled                  bool `env:"ACTIONLOG_ENABLED"`
		Debug                    bool `env:"ACTIONLOG_DEBUG"`
		WorkflowFunctionsEnabled bool `env:"ACTIONLOG_WORKFLOW_FUNCTIONS_ENABLED"`
	}

	ApigwOpt struct {
		Enabled              bool          `env:"APIGW_ENABLED"`
		Debug                bool          `env:"APIGW_DEBUG"`
		LogEnabled           bool          `env:"APIGW_LOG_ENABLED"`
		ProfilerEnabled      bool          `env:"APIGW_PROFILER_ENABLED"`
		ProfilerGlobal       bool          `env:"APIGW_PROFILER_GLOBAL"`
		LogRequestBody       bool          `env:"APIGW_LOG_REQUEST_BODY"`
		ProxyEnableDebugLog  bool          `env:"APIGW_PROXY_ENABLE_DEBUG_LOG"`
		ProxyFollowRedirects bool          `env:"APIGW_PROXY_FOLLOW_REDIRECTS"`
		ProxyOutboundTimeout time.Duration `env:"APIGW_PROXY_OUTBOUND_TIMEOUT"`
	}

	AuthOpt struct {
		LogEnabled               bool          `env:"AUTH_LOG_ENABLED"`
		PasswordSecurity         bool          `env:"AUTH_PASSWORD_SECURITY"`
		JwtAlgorithm             string        `env:"AUTH_JWT_ALGORITHM"`
		Secret                   string        `env:"AUTH_JWT_SECRET"`
		JwtKey                   string        `env:"AUTH_JWT_KEY"`
		AccessTokenLifetime      time.Duration `env:"AUTH_OAUTH2_ACCESS_TOKEN_LIFETIME"`
		RefreshTokenLifetime     time.Duration `env:"AUTH_OAUTH2_REFRESH_TOKEN_LIFETIME"`
		ExternalRedirectURL      string        `env:"AUTH_EXTERNAL_REDIRECT_URL"`
		ExternalCookieSecret     string        `env:"AUTH_EXTERNAL_COOKIE_SECRET"`
		BaseURL                  string        `env:"AUTH_BASE_URL"`
		SessionCookieName        string        `env:"AUTH_SESSION_COOKIE_NAME"`
		SessionCookiePath        string        `env:"AUTH_SESSION_COOKIE_PATH"`
		SessionCookieDomain      string        `env:"AUTH_SESSION_COOKIE_DOMAIN"`
		SessionCookieSecure      bool          `env:"AUTH_SESSION_COOKIE_SECURE"`
		SessionLifetime          time.Duration `env:"AUTH_SESSION_LIFETIME"`
		SessionPermLifetime      time.Duration `env:"AUTH_SESSION_PERM_LIFETIME"`
		GarbageCollectorInterval time.Duration `env:"AUTH_GARBAGE_COLLECTOR_INTERVAL"`
		RequestRateLimit         int           `env:"AUTH_REQUEST_RATE_LIMIT"`
		RequestRateWindowLength  time.Duration `env:"AUTH_REQUEST_RATE_WINDOW_LENGTH"`
		CsrfSecret               string        `env:"AUTH_CSRF_SECRET"`
		CsrfEnabled              bool          `env:"AUTH_CSRF_ENABLED"`
		CsrfFieldName            string        `env:"AUTH_CSRF_FIELD_NAME"`
		CsrfCookieName           string        `env:"AUTH_CSRF_COOKIE_NAME"`
		DefaultClient            string        `env:"AUTH_DEFAULT_CLIENT"`
		AssetsPath               string        `env:"AUTH_ASSETS_PATH"`
		DevelopmentMode          bool          `env:"AUTH_DEVELOPMENT_MODE"`
		ProvisionSuperUser       string        `env:"AUTH_PROVISION_SUPER_USER"`
	}

	CorredorOpt struct {
		Enabled               bool          `env:"CORREDOR_ENABLED"`
		Addr                  string        `env:"CORREDOR_ADDR"`
		MaxBackoffDelay       time.Duration `env:"CORREDOR_MAX_BACKOFF_DELAY"`
		MaxReceiveMessageSize int           `env:"CORREDOR_MAX_RECEIVE_MESSAGE_SIZE"`
		DefaultExecTimeout    time.Duration `env:"CORREDOR_DEFAULT_EXEC_TIMEOUT"`
		ListTimeout           time.Duration `env:"CORREDOR_LIST_TIMEOUT"`
		ListRefresh           time.Duration `env:"CORREDOR_LIST_REFRESH"`
		RunAsEnabled          bool          `env:"CORREDOR_RUN_AS_ENABLED"`
		TlsCertEnabled        bool          `env:"CORREDOR_CLIENT_CERTIFICATES_ENABLED"`
		TlsCertPath           string        `env:"CORREDOR_CLIENT_CERTIFICATES_PATH"`
		TlsCertCA             string        `env:"CORREDOR_CLIENT_CERTIFICATES_CA"`
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
		Enabled                  bool          `env:"FEDERATION_ENABLED"`
		Label                    string        `env:"FEDERATION_LABEL"`
		Host                     string        `env:"FEDERATION_HOST"`
		StructureMonitorInterval time.Duration `env:"FEDERATION_SYNC_STRUCTURE_MONITOR_INTERVAL"`
		StructurePageSize        int           `env:"FEDERATION_SYNC_STRUCTURE_PAGE_SIZE"`
		DataMonitorInterval      time.Duration `env:"FEDERATION_SYNC_DATA_MONITOR_INTERVAL"`
		DataPageSize             int           `env:"FEDERATION_SYNC_DATA_PAGE_SIZE"`
	}

	LimitOpt struct {
		SystemUsers             int `env:"LIMIT_SYSTEM_USERS"`
		RecordCountPerNamespace int `env:"LIMIT_RECORD_COUNT_PER_NAMESPACE"`
	}

	LocaleOpt struct {
		Languages                   string `env:"LOCALE_LANGUAGES"`
		Path                        string `env:"LOCALE_PATH"`
		QueryStringParam            string `env:"LOCALE_QUERY_STRING_PARAM"`
		ResourceTranslationsEnabled bool   `env:"LOCALE_RESOURCE_TRANSLATIONS_ENABLED"`
		Log                         bool   `env:"LOCALE_LOG"`
		DevelopmentMode             bool   `env:"LOCALE_DEVELOPMENT_MODE"`
	}

	LogOpt struct {
		Debug           bool   `env:"LOG_DEBUG"`
		Level           string `env:"LOG_LEVEL"`
		Filter          string `env:"LOG_FILTER"`
		IncludeCaller   bool   `env:"LOG_INCLUDE_CALLER"`
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
		Path            string `env:"STORAGE_PATH"`
		MinioEndpoint   string `env:"MINIO_ENDPOINT"`
		MinioSecure     bool   `env:"MINIO_SECURE"`
		MinioAccessKey  string `env:"MINIO_ACCESS_KEY"`
		MinioSecretKey  string `env:"MINIO_SECRET_KEY"`
		MinioSSECKey    string `env:"MINIO_SSEC_KEY"`
		MinioBucket     string `env:"MINIO_BUCKET"`
		MinioPathPrefix string `env:"MINIO_PATH_PREFIX"`
		MinioStrict     bool   `env:"MINIO_STRICT"`
	}

	ProvisionOpt struct {
		Always bool   `env:"PROVISION_ALWAYS"`
		Path   string `env:"PROVISION_PATH"`
	}

	SentryOpt struct {
		DSN              string  `env:"SENTRY_DSN"`
		Debug            bool    `env:"SENTRY_DEBUG"`
		AttachStacktrace bool    `env:"SENTRY_ATTACH_STACKTRACE"`
		SampleRate       float64 `env:"SENTRY_SAMPLE_RATE"`
		MaxBreadcrumbs   int     `env:"SENTRY_MAX_BREADCRUMBS"`
		ServerName       string  `env:"SENTRY_SERVERNAME"`
		Release          string  `env:"SENTRY_RELEASE"`
		Dist             string  `env:"SENTRY_DIST"`
		Environment      string  `env:"SENTRY_ENVIRONMENT"`
		WebappDSN        string  `env:"SENTRY_WEBAPP_DSN"`
	}

	TemplateOpt struct {
		RendererGotenbergAddress string `env:"TEMPLATE_RENDERER_GOTENBERG_ADDRESS"`
		RendererGotenbergEnabled bool   `env:"TEMPLATE_RENDERER_GOTENBERG_ENABLED"`
	}

	UpgradeOpt struct {
		Debug  bool `env:"UPGRADE_DEBUG"`
		Always bool `env:"UPGRADE_ALWAYS"`
	}

	WaitForOpt struct {
		Delay                 time.Duration `env:"WAIT_FOR"`
		StatusPage            bool          `env:"WAIT_FOR_STATUS_PAGE"`
		Services              string        `env:"WAIT_FOR_SERVICES"`
		ServicesTimeout       time.Duration `env:"WAIT_FOR_SERVICES_TIMEOUT"`
		ServicesProbeTimeout  time.Duration `env:"WAIT_FOR_SERVICES_PROBE_TIMEOUT"`
		ServicesProbeInterval time.Duration `env:"WAIT_FOR_SERVICES_PROBE_INTERVAL"`
	}

	WebsocketOpt struct {
		LogEnabled  bool          `env:"WEBSOCKET_LOG_ENABLED"`
		Timeout     time.Duration `env:"WEBSOCKET_TIMEOUT"`
		PingTimeout time.Duration `env:"WEBSOCKET_PING_TIMEOUT"`
		PingPeriod  time.Duration `env:"WEBSOCKET_PING_PERIOD"`
	}

	WorkflowOpt struct {
		Register          bool `env:"WORKFLOW_REGISTER"`
		ExecDebug         bool `env:"WORKFLOW_EXEC_DEBUG"`
		CallStackSize     int  `env:"WORKFLOW_CALL_STACK_SIZE"`
		StackTraceEnabled bool `env:"WORKFLOW_STACK_TRACE_ENABLED"`
		StackTraceFull    bool `env:"WORKFLOW_STACK_TRACE_FULL"`
	}

	DiscoveryOpt struct {
		Enabled       bool   `env:"DISCOVERY_ENABLED"`
		Debug         bool   `env:"DISCOVERY_DEBUG"`
		CortezaDomain string `env:"DISCOVERY_CORTEZA_DOMAIN"`
		BaseUrl       string `env:"DISCOVERY_BASE_URL"`
	}

	AttachmentOpt struct {
		AvatarMaxFileSize             int64  `env:"ATTACHMENT_AVATAR_MAX_FILE_SIZE"`
		AvatarInitialsFontPath        string `env:"AVATAR_INITIALS_FONT_PATH"`
		AvatarInitialsBackgroundColor string `env:"AVATAR_INITIALS_BACKGROUND_COLOR"`
		AvatarInitialsColor           string `env:"AVATAR_INITIALS_COLOR"`
	}

	WebappOpt struct {
		ScssDirPath string `env:"WEBAPP_SCSS_DIR_PATH"`
	}
)

// DB initializes and returns a DBOpt with default values
//
// This function is auto-generated
func DB() (o *DBOpt) {
	o = &DBOpt{
		DSN: "sqlite3://file::memory:?cache=shared&mode=memory",
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// HTTPClient initializes and returns a HTTPClientOpt with default values
//
// This function is auto-generated
func HTTPClient() (o *HTTPClientOpt) {
	o = &HTTPClientOpt{
		Timeout: 30 * time.Second,
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// HttpServer initializes and returns a HttpServerOpt with default values
//
// This function is auto-generated
func HttpServer() (o *HttpServerOpt) {
	o = &HttpServerOpt{
		Domain:                 "localhost",
		DomainWebapp:           "localhost",
		Addr:                   ":80",
		EnableHealthcheckRoute: true,
		EnableVersionRoute:     true,
		MetricsServiceLabel:    "corteza",
		MetricsUsername:        "metrics",
		MetricsPassword:        string(rand.Bytes(5)),
		EnablePanicReporting:   true,
		BaseUrl:                "/",
		ApiEnabled:             true,
		ApiBaseUrl:             "/",
		WebappBaseUrl:          "/",
		WebappBaseDir:          "./webapp/public",
		WebappList:             "admin,compose,workflow,reporter,privacy",
		SslTerminated:          isSecure(),
		WebConsoleEnabled:      false,
		WebConsoleUsername:     "admin",
		WebConsolePassword:     string(rand.Bytes(32)),
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Rbac initializes and returns a RbacOpt with default values
//
// This function is auto-generated
func Rbac() (o *RbacOpt) {
	o = &RbacOpt{
		BypassRoles:        "super-admin",
		AuthenticatedRoles: "authenticated",
		AnonymousRoles:     "anonymous",
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
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

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// SMTP initializes and returns a SMTPOpt with default values
//
// This function is auto-generated
func SMTP() (o *SMTPOpt) {
	o = &SMTPOpt{}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// ActionLog initializes and returns a ActionLogOpt with default values
//
// This function is auto-generated
func ActionLog() (o *ActionLogOpt) {
	o = &ActionLogOpt{
		Enabled: true,
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Apigw initializes and returns a ApigwOpt with default values
//
// This function is auto-generated
func Apigw() (o *ApigwOpt) {
	o = &ApigwOpt{
		Enabled:              true,
		ProfilerEnabled:      true,
		ProfilerGlobal:       false,
		ProxyFollowRedirects: true,
		ProxyOutboundTimeout: time.Second * 30,
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Auth initializes and returns a AuthOpt with default values
//
// This function is auto-generated
func Auth() (o *AuthOpt) {
	o = &AuthOpt{
		PasswordSecurity:         true,
		JwtAlgorithm:             "HS512",
		Secret:                   getSecretFromEnv("jwt secret"),
		AccessTokenLifetime:      time.Hour * 2,
		RefreshTokenLifetime:     time.Hour * 24 * 3,
		ExternalRedirectURL:      FullURL("/auth/external/{provider}/callback"),
		ExternalCookieSecret:     getSecretFromEnv("external cookie secret"),
		BaseURL:                  FullURL("/auth"),
		SessionCookieName:        "session",
		SessionCookiePath:        pathPrefix("/auth"),
		SessionCookieDomain:      GuessApiHostname(),
		SessionCookieSecure:      isSecure(),
		SessionLifetime:          24 * time.Hour,
		SessionPermLifetime:      360 * 24 * time.Hour,
		GarbageCollectorInterval: 15 * time.Minute,
		RequestRateLimit:         60,
		RequestRateWindowLength:  time.Minute,
		CsrfSecret:               getSecretFromEnv("csrf secret"),
		CsrfEnabled:              true,
		CsrfFieldName:            "same-site-authenticity-token",
		CsrfCookieName:           "same-site-authenticity-token",
		DefaultClient:            "corteza-webapp",
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
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
		MaxBackoffDelay:       time.Minute,
		MaxReceiveMessageSize: 2 << 23,
		DefaultExecTimeout:    time.Minute,
		ListTimeout:           time.Second * 2,
		ListRefresh:           time.Second * 5,
		RunAsEnabled:          true,
		TlsCertPath:           "/certs/corredor/client",
		TlsCertCA:             "ca.crt",
		TlsCertPrivate:        "private.key",
		TlsCertPublic:         "public.crt",
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
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

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
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

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Federation initializes and returns a FederationOpt with default values
//
// This function is auto-generated
func Federation() (o *FederationOpt) {
	o = &FederationOpt{
		Label:                    "federated",
		Host:                     "local.cortezaproject.org",
		StructureMonitorInterval: time.Minute * 2,
		StructurePageSize:        1,
		DataMonitorInterval:      time.Minute,
		DataPageSize:             100,
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Limit initializes and returns a LimitOpt with default values
//
// This function is auto-generated
func Limit() (o *LimitOpt) {
	o = &LimitOpt{}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
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

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Log initializes and returns a LogOpt with default values
//
// This function is auto-generated
func Log() (o *LogOpt) {
	o = &LogOpt{
		Level:           "warn",
		StacktraceLevel: "dpanic",
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Messagebus initializes and returns a MessagebusOpt with default values
//
// This function is auto-generated
func Messagebus() (o *MessagebusOpt) {
	o = &MessagebusOpt{
		Enabled: true,
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Monitor initializes and returns a MonitorOpt with default values
//
// This function is auto-generated
func Monitor() (o *MonitorOpt) {
	o = &MonitorOpt{
		Interval: 5 * time.Minute,
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// ObjectStore initializes and returns a ObjectStoreOpt with default values
//
// This function is auto-generated
func ObjectStore() (o *ObjectStoreOpt) {
	o = &ObjectStoreOpt{
		Path:        "var/store",
		MinioSecure: true,
		MinioBucket: "{component}",
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
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

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
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

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Template initializes and returns a TemplateOpt with default values
//
// This function is auto-generated
func Template() (o *TemplateOpt) {
	o = &TemplateOpt{}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
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
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// WaitFor initializes and returns a WaitForOpt with default values
//
// This function is auto-generated
func WaitFor() (o *WaitForOpt) {
	o = &WaitForOpt{
		StatusPage:            true,
		ServicesTimeout:       time.Minute,
		ServicesProbeTimeout:  time.Second * 30,
		ServicesProbeInterval: time.Second * 5,
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Websocket initializes and returns a WebsocketOpt with default values
//
// This function is auto-generated
func Websocket() (o *WebsocketOpt) {
	o = &WebsocketOpt{
		Timeout:     15 * time.Second,
		PingTimeout: 120 * time.Second,
		PingPeriod:  ((120 * time.Second) * 9) / 10,
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Workflow initializes and returns a WorkflowOpt with default values
//
// This function is auto-generated
func Workflow() (o *WorkflowOpt) {
	o = &WorkflowOpt{
		Register:          true,
		CallStackSize:     16,
		StackTraceEnabled: true,
		StackTraceFull:    true,
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Discovery initializes and returns a DiscoveryOpt with default values
//
// This function is auto-generated
func Discovery() (o *DiscoveryOpt) {
	o = &DiscoveryOpt{
		Enabled: false,
		Debug:   false,
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Attachment initializes and returns a AttachmentOpt with default values
//
// This function is auto-generated
func Attachment() (o *AttachmentOpt) {
	o = &AttachmentOpt{
		AvatarMaxFileSize:             1000000,
		AvatarInitialsFontPath:        "fonts/Poppins-Regular.ttf",
		AvatarInitialsBackgroundColor: "#F3F3F3",
		AvatarInitialsColor:           "#162425",
	}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}

// Webapp initializes and returns a WebappOpt with default values
//
// This function is auto-generated
func Webapp() (o *WebappOpt) {
	o = &WebappOpt{}

	// Custom defaults
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	fill(o)

	// Custom cleanup
	func(o interface{}) {
		if def, ok := o.(interface{ Cleanup() }); ok {
			def.Cleanup()
		}
	}(o)

	return
}
