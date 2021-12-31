package app

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
	"github.com/cortezaproject/corteza-server/system"
	"github.com/cortezaproject/corteza-server/compose"
)

corteza: schema.#platform & {
	ident: "corteza"

	components: [
		system.component,
		compose.component,
	]
}
