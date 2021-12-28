package schema

import (
	"strings"
	"list"
)


#locale: {
	resourceExpIdent: #expIdent

	// @todo we need a better name here!
	skipSvc: bool | *false

	extended: bool | *false

	resource: {
		// @todo merge with RBAC res-ref and move 2 levels lower.
		references: [ ...string] | *["ID"]
		type: string
		const: string | *("\(resourceExpIdent)ResourceTranslationType")
	}

	keys: {
		[key=_]: #localeKey & {
			name: key
			_resourceExpIdent: resourceExpIdent
	  }
	}
}

#localeKey: {
	name: #handle
	_resourceExpIdent: #expIdent

	path: [...(#ident | { part: #ident, var: bool | *false })] | *([name])

  expandedPath: [for p in path {
	  if (p & { "p": #ident }) != _|_ { p, var: p.var }
	  if (p & string) != _|_ { "part": p, var: false }
  }]


  _suffix: strings.Join([for p in expandedPath { strings.ToTitle(p.part) }], "")

	struct:     string | *("LocaleKey" + _resourceExpIdent + _suffix)

	// As soon as we use vars in the path,
	// custom handler must be present
  _hasVars: list.Contains([for p in path { p.var | false }], true)
	customHandler: bool | *_hasVars

	if customHandler {
	  decodeFunc: string | *("decodeTranslations" + _suffix)
	  encodeFunc: string | *("encodeTranslations" + _suffix)
	  serviceFunc: string | *("handle" + _resourceExpIdent + _suffix)
	}
}
