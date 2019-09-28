package importer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAutomationScript(t *testing.T) {

	impFixTester(t, "automation_s+t", func(t *testing.T, s *AutomationScript) {
		req := require.New(t)
		req.NotNil(s)
		req.Len(s.set, 1)
		req.NotEmpty(s.set[0].Source)
		req.Len(s.triggers[s.set[0].Name], 2)
		req.False(s.triggers[s.set[0].Name][0].Enabled)
		req.True(s.triggers[s.set[0].Name][1].Enabled)
	})
}
