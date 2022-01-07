package schema

import (
	"strings"
)

#gigWorker: {
	ident:        string
	expIdent:     string | *strings.ToTitle(ident)
	description?: string
	goType:       string | *"worker\(expIdent)"
	goConst:      string | *"Worker\(expIdent)Handle"

	tasks: [...({kind: "preprocessor"} & #gigTask)] | *[]
}
