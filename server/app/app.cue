package app

import (
	"github.com/crusttech/human/server/codegen/schema"
	"github.com/crusttech/human/server/app/options"
	"github.com/crusttech/human/server/system"
	"github.com/crusttech/human/server/compose"
	"github.com/crusttech/human/server/automation"
	"github.com/crusttech/human/server/federation"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/corredor"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	discoveryTypes "github.com/crusttech/human/server/discovery/types"
)

human: schema.#platform & {
	"ident": "corteza"

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
		options.provision,
		options.sentry,
		options.template,
		options.upgrade,
		options.waitFor,
		options.websocket,
		options.workflow,
		options.discovery,
		options.attachment,
		options.webapp,
		options.observability,
		options.agentic,
	]

	// platform resources
	"resources": resources

	"components": [
		system.component,
		compose.component,
		automation.component,
		federation.component,
	]

	"bundles": [
		actionlog.bundle,
		corredor.bundle,
		labelTypes.bundle,
		discoveryTypes.bundle,
	]
}
