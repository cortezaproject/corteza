package agentic

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/crusttech/human/server/compose/service/values"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/require"
)

// A module whose expressions cannot work is still stored, so the write result
// is the only place the caller can learn about it — otherwise the failure lands
// on the next compose_record_create, worded as if the record were at fault.
func TestModuleWriteResult_issues(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	broken := &cmpTypes.Module{
		ID: 1, NamespaceID: 2, Name: "m", Handle: "m",
		Fields: cmpTypes.ModuleFieldSet{
			&cmpTypes.ModuleField{Name: "stage", Kind: "Select"},
			&cmpTypes.ModuleField{
				Name:        "is_active",
				Kind:        "Bool",
				Expressions: cmpTypes.ModuleFieldExpr{ValueExpr: `record.values.stage != "hired"`},
			},
		},
	}

	extra := moduleExpressionExtras(ctx, broken)
	req.Contains(extra, "issues")
	req.Contains(extra, "note")

	ii, ok := extra["issues"].([]values.ExpressionIssue)
	req.True(ok)
	req.Len(ii, 1)
	req.Equal("is_active", ii[0].Field)
	req.Equal(values.IssueSlotValue, ii[0].Slot)
	req.Equal(values.IssueSeverityError, ii[0].Severity)

	// the extras ride beside the module, so they must serialise as siblings
	enc, err := json.Marshal(extra)
	req.NoError(err)
	req.Contains(string(enc), `"field":"is_active"`)
}

func TestModuleWriteResult_cleanModuleHasNoIssuesKey(t *testing.T) {
	req := require.New(t)

	clean := &cmpTypes.Module{
		ID: 1, NamespaceID: 2, Name: "m", Handle: "m",
		Fields: cmpTypes.ModuleFieldSet{
			&cmpTypes.ModuleField{Name: "stage", Kind: "Select"},
			&cmpTypes.ModuleField{
				Name:        "is_active",
				Kind:        "Bool",
				Expressions: cmpTypes.ModuleFieldExpr{ValueExpr: `stage != "hired"`},
			},
		},
	}

	req.Empty(moduleExpressionExtras(context.Background(), clean))
}
