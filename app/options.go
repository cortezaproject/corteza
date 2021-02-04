package app

import (
	"github.com/cortezaproject/corteza-server/pkg/options"
)

type (
	Options struct {
		Environment options.EnvironmentOpt
		ActionLog   options.ActionLogOpt
		SMTP        options.SMTPOpt
		Auth        options.AuthOpt
		HTTPClient  options.HTTPClientOpt
		DB          options.DBOpt
		Template    options.TemplateOpt
		Upgrade     options.UpgradeOpt
		Provision   options.ProvisionOpt
		Sentry      options.SentryOpt
		ObjStore    options.ObjectStoreOpt
		Corredor    options.CorredorOpt
		Monitor     options.MonitorOpt
		WaitFor     options.WaitForOpt
		HTTPServer  options.HTTPServerOpt
		Websocket   options.WebsocketOpt
		Eventbus    options.EventbusOpt
		Federation  options.FederationOpt
		SCIM        options.SCIMOpt
	}
)

func NewOptions() *Options {
	return &Options{
		Environment: *options.Environment(),
		ActionLog:   *options.ActionLog(),
		Auth:        *options.Auth(),
		SMTP:        *options.SMTP(),
		HTTPClient:  *options.HTTPClient(),
		DB:          *options.DB(),
		Template:    *options.Template(),
		Upgrade:     *options.Upgrade(),
		Provision:   *options.Provision(),
		Sentry:      *options.Sentry(),
		ObjStore:    *options.ObjectStore(),
		Corredor:    *options.Corredor(),
		Monitor:     *options.Monitor(),
		WaitFor:     *options.WaitFor(),
		HTTPServer:  *options.HTTPServer(),
		Websocket:   *options.Websocket(),
		Eventbus:    *options.Eventbus(),
		Federation:  *options.Federation(),
		SCIM:        *options.SCIM(),
	}
}
