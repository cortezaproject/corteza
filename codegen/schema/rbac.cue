package schema

import (
	"strings"
)

#rbacComponent: {
	resource: {
		type: string
	}

	operations: {
		[key=_]: #rbacOperation & {
			handle: key
		}
	}
}

#rbacResource: {
	resourceExpIdent: #expIdent

	operations: {
		[key=_]: #rbacOperation & {
			handle:            key
			_resourceExpIdent: resourceExpIdent
		}
	}
}

#rbacOperation: {
	handle:             #handle
	description:        string | *handle
	_resourceExpIdent?: string

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

	if _resourceExpIdent == _|_ {
		checkFuncName: #expIdent | *"Can\(_opFinal)"
	}

	if len(_opSplit) > 1 && _resourceExpIdent != _|_ {
		checkFuncName: #expIdent | *"Can\(_opFinal)On\(_resourceExpIdent)"
	}

	if len(_opSplit) <= 1 && _resourceExpIdent != _|_ {
		checkFuncName: #expIdent | *"Can\(_opFinal)\(_resourceExpIdent)"
	}
}
