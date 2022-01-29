package codegen

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

// List of all codegen jobs for the entire platform
//
// How to run it?
// @todo when this gets into
// cue eval codegen/*.cue --out json -e platform | go run codegen/tool/*.go -v
platform: [...schema.#codegen] &
	rbacAccessControl+
	rbacTypes+
	localeTypes+
	envoyRBAC+
	options+
	[] // placeholder
