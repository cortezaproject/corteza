package schema

import (
	"strings"
)

#gigTask: {
	kind:        string
	expKind:     string | *strings.ToTitle(kind)
	goInterface: string | *expKind

	ident:        string
	expIdent:     string | *strings.ToTitle(ident)
	description?: string
	goType:       string | *"\(kind)\(expIdent)"
	goConst:      string | *"\(expKind)Handle\(expIdent)"

	_constructorBase:  "\(expKind)\(expIdent)"
	constructor:       string | *"\(_constructorBase)"
	constructorParams: string | *"\(_constructorBase)Params"
	transformer?:      string

	struct: {
		[key=_]: {
			ident:        string | *key
			expIdent:     string | *strings.ToTitle(ident)
			required:     bool | *false
			param:        bool | *true
			description?: string

			goType:   string | *"string"
			exprType: string | *strings.ToTitle(goType)
			castFunc: string | *"cast.To\(strings.ToTitle(goType))"

			constructor: string | *"\(_constructorBase)\(expIdent)"
		}
	}
}
