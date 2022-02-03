package schema

import (
	"strings"
)

#gigWorker: {
	ident:        string
	expIdent:     string | *strings.ToTitle(ident)
	description?: string
	goType:       string | *"worker\(expIdent)"
	goConst:      string | *"WorkerHandle\(expIdent)"

	tasks: [...({kind: "preprocessor"} & #gigTask)] | *[]
}
