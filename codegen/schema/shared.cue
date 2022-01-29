package schema

import (
	"strings"
)

// Identifier
#ident: =~"^[a-z][a-zA-Z0-9_]*$"

// Exported identifier
#expIdent: =~"^[A-Z][a-zA-Z0-9_]*$"

// More liberal then identifier, allows underscores and dots
#handle: =~"^[A-Za-z][a-zA-Z0-9_\\-\\.]*[a-zA-Z0-9]+$"

// More liberal then identifier, allows underscores and dots
#baseHandle: =~"^[a-z][a-z0-9-]*[a-z0-9]+$"


#_base: {
	// lowercase dash-separated words
	// used to build ident and exported identifiers
	handle: #baseHandle | *"base"
	_words: strings.Replace(strings.Replace(strings.Replace(handle, "-", " ", -1), "_", " ", -1), ".", " ", -1)

	// lowercased (unexported, golang) identifier
	ident: #ident | *strings.ToCamel(strings.Replace(strings.ToTitle(_words), " ", "", -1))

	// upercased (exported, golang) identifier
	expIdent: #expIdent | *strings.Replace(strings.ToTitle(_words), " ", "", -1)

	...
}
