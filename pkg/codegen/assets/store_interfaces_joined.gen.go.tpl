package store

// This file is auto-generated.
//
// Template:	pkg/codegen/assets/store_interfaces_joined.gen.go.tpl
// Definitions:
{{- range .Definitions }}
//  - {{ .Source }}
{{- end }}
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
)

type (
	Transactionable interface {
		Tx(context.Context, func(context.Context, Storable) error) error
	}

	// Sortable interface combines interfaces of all supported store interfaces
	Storable interface {
		Transactionable

	{{ range .Definitions -}}
		{{ export .Types.Plural }}
	{{ end }}
	}
)
