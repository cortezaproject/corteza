package service

import (
	"testing"

	"github.com/crusttech/human/server/automation/types"
	"github.com/stretchr/testify/require"
)

func TestValidateScopeRefs_deletedStepScope(t *testing.T) {
	req := require.New(t)

	steps := types.NgAutomationStepSet{
		{ID: 3, Handle: "step3", Arguments: []*types.Expr{
			{ArgumentName: "in", Scope: "step2", Source: "out"},
		}},
	}

	issues := validateScopeRefs(steps, nil)

	req.Len(issues, 1)
	req.Equal(types.IssueCodeScopeUnknown, issues[0].Code)
	req.Equal(types.NgAutomationSeverityError, issues[0].Severity)
	req.Contains(issues[0].Message, `unknown scope "step2"`)
	req.Len(issues[0].Details, 1)

	d := issues[0].Details[0].MissingReference
	req.NotNil(d)
	req.Equal("scope", d.RefKind)
	req.Equal("step2", d.Ref)
	req.Equal(uint64(3), d.StepID)
	req.Equal("arguments", d.Field)
	req.Equal(0, d.FieldIndex)
}
