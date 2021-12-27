package schema

import (
//	"strings"
)


#locale: {
	// @todo we need a better name here!
	skipSvc: bool | *false

	resource: {
		// @todo merge with RBAC res-ref and move 2 levels lower.
		references: [ ...string] | *["ID"]
	}

	keys: {
		[key=_]: #localeKey & {
			name: key
	  }
	}
}

#localeKey: {
	name: #handle
	path: string | *(name)
	custom?: true
	customHandler?: string
}
