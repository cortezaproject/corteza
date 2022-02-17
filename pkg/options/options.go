package options

type (
	Options struct {
		Environment EnvironmentOpt
		ActionLog   ActionLogOpt
		SMTP        SMTPOpt
		Auth        AuthOpt
		HTTPClient  HTTPClientOpt
		DB          DBOpt
		Template    TemplateOpt
		Upgrade     UpgradeOpt
		Provision   ProvisionOpt
		Sentry      SentryOpt
		ObjStore    ObjectStoreOpt
		Corredor    CorredorOpt
		Monitor     MonitorOpt
		WaitFor     WaitForOpt
		HTTPServer  HttpServerOpt
		Websocket   WebsocketOpt
		Eventbus    EventbusOpt
		Messagebus  MessagebusOpt
		Federation  FederationOpt
		SCIM        SCIMOpt
		Workflow    WorkflowOpt
		RBAC        RbacOpt
		Locale      LocaleOpt
		Limit       LimitOpt
		Plugins     PluginsOpt
		Discovery   DiscoveryOpt
		Apigw       ApigwOpt
	}
)

func Init() *Options {
	return &Options{
		Environment: *Environment(),
		ActionLog:   *ActionLog(),
		Auth:        *Auth(),
		SMTP:        *SMTP(),
		HTTPClient:  *HTTPClient(),
		DB:          *DB(),
		Template:    *Template(),
		Upgrade:     *Upgrade(),
		Provision:   *Provision(),
		Sentry:      *Sentry(),
		ObjStore:    *ObjectStore(),
		Corredor:    *Corredor(),
		Monitor:     *Monitor(),
		WaitFor:     *WaitFor(),
		HTTPServer:  *HttpServer(),
		Websocket:   *Websocket(),
		Eventbus:    *Eventbus(),
		Messagebus:  *Messagebus(),
		Federation:  *Federation(),
		SCIM:        *SCIM(),
		Workflow:    *Workflow(),
		RBAC:        *Rbac(),
		Locale:      *Locale(),
		Limit:       *Limit(),
		Plugins:     *Plugins(),
		Discovery:   *Discovery(),
		Apigw:       *Apigw(),
	}
}
