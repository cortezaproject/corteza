package app

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
	"github.com/cortezaproject/corteza-server/system"
	"github.com/cortezaproject/corteza-server/compose"
	"github.com/cortezaproject/corteza-server/automation"
	"github.com/cortezaproject/corteza-server/federation"
)

corteza: schema.#platform & {
	ident: "corteza"

	components: [
		system.component,
		compose.component,
		automation.component,
		federation.component,
	]
}
