package app

import (
	"github.com/cortezaproject/corteza-server/pkg/app/options"
)

type (
	Options struct {
		SMTP       options.SMTPOpt
		Auth       options.AuthOpt
		HTTPClient options.HTTPClientOpt
		DB         options.DBOpt
		Upgrade    options.UpgradeOpt
		Provision  options.ProvisionOpt
		Sentry     options.SentryOpt
		Storage    options.StorageOpt
		Corredor   options.CorredorOpt
		Monitor    options.MonitorOpt
		WaitFor    options.WaitForOpt
		HTTPServer options.HTTPServerOpt
		GRPCServer options.GRPCServerOpt
		Websocket  options.WebsocketOpt
	}
)

func NewOptions(prefix ...string) *Options {
	var p = ""
	if len(prefix) > 0 {
		p = prefix[0]
	}

	return &Options{
		Auth:       *options.Auth(),
		SMTP:       *options.SMTP(p),
		HTTPClient: *options.HttpClient(p),
		DB:         *options.DB(p),
		Upgrade:    *options.Upgrade(p),
		Provision:  *options.Provision(p),
		Sentry:     *options.Sentry(p),
		Storage:    *options.Storage(p),
		Corredor:   *options.Corredor(),
		Monitor:    *options.Monitor(p),
		WaitFor:    *options.WaitFor(p),
		HTTPServer: *options.HTTP(p),
		GRPCServer: *options.GRPCServer(p),
		Websocket:  *options.Websocket(p),
	}
}
