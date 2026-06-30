package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIssueDetail_roundTrip(t *testing.T) {
	req := require.New(t)

	in := &NgAutomationIssue{
		Code:     IssueCodeScopeUnknown,
		Severity: NgAutomationSeverityError,
		Message:  `Step "step3" arguments reference unknown scope "step2".`,
		Details: []*NgAutomationIssueDetail{
			NewDetailMissingReference("scope", "step2", 3, "arguments", 0),
		},
	}

	raw, err := json.Marshal(in)
	req.NoError(err)
	req.Contains(string(raw), `"type":"missingReference"`)
	req.Contains(string(raw), `"stepID":"3"`)

	var out NgAutomationIssue
	req.NoError(json.Unmarshal(raw, &out))
	req.Equal(in, &out)
}

func TestIssueDetail_cycleIDsAsStrings(t *testing.T) {
	req := require.New(t)

	d := NewDetailCycle(3, 5)
	raw, err := json.Marshal(d)
	req.NoError(err)
	req.Contains(string(raw), `"stepIDs":["3","5"]`)

	var out NgAutomationIssueDetail
	req.NoError(json.Unmarshal(raw, &out))
	req.NotNil(out.Cycle)
	req.Equal(IssueIDs{3, 5}, out.Cycle.StepIDs)
}

func TestIssueDetail_marshalRejectsZeroOrMultiVariant(t *testing.T) {
	req := require.New(t)

	_, err := json.Marshal(&NgAutomationIssueDetail{})
	req.Error(err)

	_, err = json.Marshal(&NgAutomationIssueDetail{
		Cycle:      &DetailCycle{},
		EmptyField: &DetailEmptyField{},
	})
	req.Error(err)
}

func TestIssueDetail_unmarshalUnknownType(t *testing.T) {
	req := require.New(t)

	var out NgAutomationIssueDetail
	err := json.Unmarshal([]byte(`{"type":"bogus"}`), &out)
	req.Error(err)
}
