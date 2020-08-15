package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// {{ .Source }}

import (
	"context"
{{- range $import := $.Import }}
    {{ normalizeImport $import }}
{{- end }}
)

type (
	{{ unpubIdent $.Types.Singular }}Create struct {
	    Done chan struct{}
	    res  *{{ $.Types.GoType }}
	    err  error
    }

	{{ unpubIdent $.Types.Singular }}Update struct {
	    Done chan struct{}
	    res  *{{ $.Types.GoType }}
	    err  error
    }

	{{ unpubIdent $.Types.Singular }}Remove struct {
	    Done chan struct{}
	    res  *{{ $.Types.GoType }}
	    err  error
    }
)

// Create{{ pubIdent $.Types.Singular }} creates a new {{ pubIdent $.Types.Singular }}
// create job that can be pushed to store's transaction handler
func Create{{ pubIdent $.Types.Singular }}(res *{{ $.Types.GoType }}) *{{ unpubIdent $.Types.Singular }}Create {
    return &{{ unpubIdent $.Types.Singular }}Create{res: res}
}

// Do Executes {{ unpubIdent $.Types.Singular }}Create job
func (j *{{ unpubIdent $.Types.Singular }}Create) Do(ctx context.Context, s storeInterface) error {
	j.err = s.Create{{ pubIdent $.Types.Singular }}(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// Update{{ pubIdent $.Types.Singular }} creates a new {{ pubIdent $.Types.Singular }}
// update job that can be pushed to store's transaction handler
func Update{{ pubIdent $.Types.Singular }}(res *{{ $.Types.GoType }}) *{{ unpubIdent $.Types.Singular }}Update {
    return &{{ unpubIdent $.Types.Singular }}Update{res: res}
}

// Do Executes {{ unpubIdent $.Types.Singular }}Update job
func (j *{{ unpubIdent $.Types.Singular }}Update) Do(ctx context.Context, s storeInterface) error {
	j.err = s.Update{{ pubIdent $.Types.Singular }}(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// Remove{{ pubIdent $.Types.Singular }} creates a new {{ pubIdent $.Types.Singular }}
// remove job that can be pushed to store's transaction handler
func Remove{{ pubIdent $.Types.Singular }}(res *{{ $.Types.GoType }}) *{{ unpubIdent $.Types.Singular }}Remove {
    return &{{ unpubIdent $.Types.Singular }}Remove{res: res}
}

// Do Executes {{ unpubIdent $.Types.Singular }}Remove job
func (j *{{ unpubIdent $.Types.Singular }}Remove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.Remove{{ pubIdent $.Types.Singular }}(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}
