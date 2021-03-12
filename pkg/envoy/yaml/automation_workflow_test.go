package yaml

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestAutomationWorkflow_UnmarshalYAML(t *testing.T) {
	var (
		parseString = func(src string) (*automationWorkflow, error) {
			w := &automationWorkflow{}
			return w, yaml.Unmarshal([]byte(src), w)
		}
	)

	t.Run("empty", func(t *testing.T) {
		req := require.New(t)

		w, err := parseString(``)
		req.NoError(err)
		req.NotNil(w)
		req.Nil(w.res)
	})

	t.Run("workflow 1", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("workflow_1")
		req.NoError(err)
		req.NotNil(doc)
		req.Len(doc.automation.Workflows, 1)
		req.NotNil(doc.automation.Workflows[0])
		req.Len(doc.automation.Workflows[0].triggers, 1)
		req.Len(doc.automation.Workflows[0].steps, 1)
		req.Len(doc.automation.Workflows[0].paths, 1)
	})
}
