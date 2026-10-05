package wfexec

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFrameJSON(t *testing.T) {
	f := Frame{SessionID: 516687706984677377, StateID: 516687706984677378, ParentID: 2, StepID: 3, NextSteps: []uint64{516687706984677379}, Action: "x"}

	out, err := json.Marshal(f)
	require.NoError(t, err)
	require.Contains(t, string(out), `"sessionID":"516687706984677377"`)
	require.Contains(t, string(out), `"stateID":"516687706984677378"`)
	require.Contains(t, string(out), `"parentID":"2"`)
	require.Contains(t, string(out), `"stepID":"3"`)
	require.Contains(t, string(out), `"nextSteps":["516687706984677379"]`)
	require.Contains(t, string(out), `"action":"x"`)

	var back Frame
	require.NoError(t, json.Unmarshal(out, &back))
	require.Equal(t, f.SessionID, back.SessionID)
	require.Equal(t, f.StateID, back.StateID)
	require.Equal(t, f.StepID, back.StepID)
	require.Equal(t, f.Action, back.Action)

	// a stacktrace stored with numeric IDs still reads
	var stored Frame
	require.NoError(t, json.Unmarshal([]byte(`{"sessionID":516687706984677377,"stepID":3,"nextSteps":[4]}`), &stored))
	require.Equal(t, uint64(516687706984677377), stored.SessionID)
	require.Equal(t, uint64(3), stored.StepID)
	require.Equal(t, []uint64{4}, []uint64(stored.NextSteps))
}
