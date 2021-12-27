package schema

import (
	"strings"
)

#rbacComponent: {
	resource: {
		type: string
	}

	operations: {
		[key=_]: #rbacOperation & {handle: key}
	}
}

#rbacResource: {
	resource: {
		type:       string
		expIdent:   #expIdent
		references: [ ...string] | *["ID"]
	}

	operations: {
		[key=_]: #rbacOperation & {
			handle:           key
			resourceExpIdent: resource.expIdent
			description:      string | *(strings.ToTitle(key) + " " + resource.type)
		}
	}
}

#rbacOperation: {
	handle:            #handle
	description:       string
	resourceExpIdent?: string

	_isComponent: resourceExpIdent == _|_

	// Some string manipulation that will result in
	// more pronouncable access-control check function name

	// When check function name is not explicitly defined we try
	// to use resource and operation name and generate easy-to-read name
	//
	// <res> + <op>              => Can<Op><Res>
	// <res> + <op:foo.bar.verb> => Can<Verb><Foo><Bar>On<Res>

	_operation: strings.Replace(strings.Replace(handle, "-", " ", -1), "_", " ", -1)
	_opSplit:   strings.Split(_operation, ".")
	_opFlip:    [_opSplit[len(_opSplit)-1]] + _opSplit[0:len(_opSplit)-1]
	_opFinal:   strings.Replace(strings.ToTitle(strings.Join(_opFlip, " ")), " ", "", -1)

	if _isComponent {
		checkFuncName: #expIdent | *("Can" + _opFinal)
	}

	if !_isComponent {
		checkFuncName: #expIdent | *("Can" + _opFinal + resourceExpIdent)
	}
}
