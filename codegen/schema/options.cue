package schema

import (
	"strings"
)

#_ENV:     =~ "^[A-Z][A-Z0-9_]*[A-Z0-9]?$"
//#_optName: =~ "^[a-zA-Z][a-zA-Z0-9\\s]*[a-zA-Z0-9]+$"

#optionsGroup: #_base & {
	imports: [...string] | *([])

  handle: #handle

  title?: string
  description?: string

  env: #_ENV | *(strings.ToUpper(strings.Replace(handle, "-", "_", -1)))
	_envPrefix: env

	options: {
		[_opt=_]: #option & {
		  handle: _opt
			env: #_ENV | *(_envPrefix + "_" + strings.ToUpper(strings.Replace(handle, " ", "_", -1)))
		}
	}
}

#option: {
  handle: #handle
	_words: strings.Replace(strings.Replace(strings.Replace(handle, "-", " ", -1), "_", " ", -1), ".", " ", -1)

	// lowercased (unexported, golang) identifier
	ident: #ident | *strings.ToCamel(strings.Replace(strings.ToTitle(_words), " ", "", -1))

	// upercased (exported, golang) identifier
	expIdent: #expIdent | *strings.Replace(strings.ToTitle(_words), " ", "", -1)

	type: string | *"string"
	description?: string
	default?: string
	env?: #_ENV
}
