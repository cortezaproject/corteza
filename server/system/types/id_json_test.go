package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIDsMarshalAsStrings(t *testing.T) {
	const big = 516687706984677377

	cases := map[string]struct {
		v    any
		want string
	}{
		"setting value":   {SettingValue{Name: "a", Value: []byte(`1`), UpdatedBy: big}, `"updatedBy":"516687706984677377"`},
		"settings filter": {SettingsFilter{OwnedBy: big}, `"ownedBy":"516687706984677377"`},
		"reminder filter": {ReminderFilter{AssignedTo: big, ReminderID: []uint64{big}}, `"assignedTo":"516687706984677377","scheduledFrom"`},
		"reminder IDs":    {ReminderFilter{ReminderID: []uint64{big}}, `"reminderID":["516687706984677377"]`},
		"task filter":     {ProjectTaskFilter{TaskID: []uint64{big}}, `"taskID":["516687706984677377"]`},
		"auth client":     {AuthClient{OwnedBy: big, CreatedBy: big}, `"ownedBy":"516687706984677377"`},
		"report":          {Report{CreatedBy: big}, `"createdBy":"516687706984677377"`},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			out, err := json.Marshal(c.v)
			require.NoError(t, err)
			require.Contains(t, string(out), c.want)
		})
	}
}

func TestSettingValueReadsBothIDForms(t *testing.T) {
	for _, in := range []string{
		`{"name":"a","value":1,"updatedBy":516687706984677377}`,
		`{"name":"a","value":1,"updatedBy":"516687706984677377"}`,
	} {
		var v SettingValue
		require.NoError(t, json.Unmarshal([]byte(in), &v), in)
		require.Equal(t, uint64(516687706984677377), v.UpdatedBy, in)
		require.Equal(t, "a", v.Name, in)
		require.JSONEq(t, `1`, string(v.Value), in)
	}
}
