package app

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
	"github.com/cortezaproject/corteza/server/app/options"
	"github.com/cortezaproject/corteza/server/system"
	"github.com/cortezaproject/corteza/server/compose"
	"github.com/cortezaproject/corteza/server/automation"
	"github.com/cortezaproject/corteza/server/federation"
)

corteza: schema.#platform & {
	handle: "corteza"

	"options": [
		options.DB,
		options.HTTPClient,
		options.HTTPServer,
		options.RBAC,
		options.SCIM,
		options.SMTP,
		options.actionLog,
		options.apigw,
		options.auth,
		options.corredor,
		options.environment,
		options.eventbus,
		options.federation,
		options.limit,
		options.locale,
		options.log,
		options.messagebus,
		options.monitor,
		options.objectStore,
		options.plugins,
		options.provision,
		options.seeder,
		options.sentry,
		options.template,
		options.upgrade,
		options.waitFor,
		options.websocket,
		options.workflow,
		options.discovery,
	]

	components: [
		system.component,
		compose.component,
		automation.component,
		federation.component,
	]
}
