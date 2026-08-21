package agentic

import (
	"fmt"
	"strings"
	"testing"

	autoTypes "github.com/crusttech/human/server/automation/types"
	"github.com/stretchr/testify/require"
)

func layoutStep(id uint64, kind autoTypes.WorkflowStepKind) *autoTypes.WorkflowStep {
	return &autoTypes.WorkflowStep{ID: id, Kind: kind}
}

func layoutPath(parent, child uint64) *autoTypes.WorkflowPath {
	return &autoTypes.WorkflowPath{ParentID: parent, ChildID: child}
}

// pos is where a step ended up, and fails the test rather than the assertion
// when the step was never placed at all.
func pos(t *testing.T, steps autoTypes.WorkflowStepSet, id uint64) (int, int) {
	t.Helper()
	for _, s := range steps {
		if s.ID != id {
			continue
		}
		x, y, ok := visualPosition(s.Meta.Visual)
		require.True(t, ok, "step %d carries no position", id)
		return x, y
	}
	t.Fatalf("step %d not in the set", id)
	return 0, 0
}

func TestLayoutWorkflowSteps(t *testing.T) {
	t.Run("a chain runs straight down one column", func(t *testing.T) {
		steps := autoTypes.WorkflowStepSet{
			layoutStep(1, autoTypes.WorkflowStepKindExpressions),
			layoutStep(2, autoTypes.WorkflowStepKindExpressions),
			layoutStep(3, autoTypes.WorkflowStepKindTermination),
		}
		layoutWorkflowSteps(steps, autoTypes.WorkflowPathSet{layoutPath(1, 2), layoutPath(2, 3)})

		x1, y1 := pos(t, steps, 1)
		x2, y2 := pos(t, steps, 2)
		x3, y3 := pos(t, steps, 3)

		require.Equal(t, layoutOriginX, x1)
		require.Equal(t, layoutOriginY, y1)
		require.Equal(t, x1, x2, "a single successor keeps the column")
		require.Equal(t, x1, x3)
		require.Equal(t, layoutPitchY, y2-y1)
		require.Equal(t, layoutPitchY, y3-y2)
	})

	t.Run("every node is given the standard size", func(t *testing.T) {
		steps := autoTypes.WorkflowStepSet{layoutStep(1, autoTypes.WorkflowStepKindExpressions)}
		layoutWorkflowSteps(steps, nil)

		xywh := steps[0].Meta.Visual["xywh"].([]interface{})
		require.Equal(t, layoutNodeWidth, xywh[2])
		require.Equal(t, layoutNodeHeight, xywh[3])
		require.Equal(t, "1", steps[0].Meta.Visual["id"])
		require.Equal(t, "1", steps[0].Meta.Visual["parent"])
	})

	t.Run("a gateway's arms take their own columns", func(t *testing.T) {
		steps := autoTypes.WorkflowStepSet{
			layoutStep(1, autoTypes.WorkflowStepKindGateway),
			layoutStep(2, autoTypes.WorkflowStepKindExpressions),
			layoutStep(3, autoTypes.WorkflowStepKindExpressions),
			layoutStep(4, autoTypes.WorkflowStepKindExpressions),
		}
		layoutWorkflowSteps(steps, autoTypes.WorkflowPathSet{
			layoutPath(1, 2), layoutPath(1, 3), layoutPath(1, 4),
		})

		x1, y1 := pos(t, steps, 1)
		x2, y2 := pos(t, steps, 2)
		x3, _ := pos(t, steps, 3)
		x4, _ := pos(t, steps, 4)

		require.Equal(t, x1, x2, "the first arm carries the column")
		require.Equal(t, x1+layoutPitchX, x3)
		require.Equal(t, x1+2*layoutPitchX, x4)
		require.Equal(t, layoutPitchY, y2-y1, "every arm drops one row")
	})

	t.Run("an iterator sends its body aside and carries its exit down", func(t *testing.T) {
		steps := autoTypes.WorkflowStepSet{
			layoutStep(1, autoTypes.WorkflowStepKindIterator),
			layoutStep(2, autoTypes.WorkflowStepKindExpressions), // body
			layoutStep(3, autoTypes.WorkflowStepKindTermination), // exit
		}
		layoutWorkflowSteps(steps, autoTypes.WorkflowPathSet{layoutPath(1, 2), layoutPath(1, 3)})

		x1, _ := pos(t, steps, 1)
		body, _ := pos(t, steps, 2)
		exit, _ := pos(t, steps, 3)

		require.Equal(t, x1+layoutPitchX, body, "the body branches away")
		require.Equal(t, x1, exit, "the exit stays on the column")
	})

	t.Run("a loop back into the iterator terminates", func(t *testing.T) {
		steps := autoTypes.WorkflowStepSet{
			layoutStep(1, autoTypes.WorkflowStepKindIterator),
			layoutStep(2, autoTypes.WorkflowStepKindExpressions),
			layoutStep(3, autoTypes.WorkflowStepKindTermination),
		}
		// The body loops back to the iterator, which is what an iterator body does.
		layoutWorkflowSteps(steps, autoTypes.WorkflowPathSet{
			layoutPath(1, 2), layoutPath(1, 3), layoutPath(2, 1),
		})

		for _, id := range []uint64{1, 2, 3} {
			pos(t, steps, id)
		}
	})

	t.Run("a step the paths never reach is still placed", func(t *testing.T) {
		steps := autoTypes.WorkflowStepSet{
			layoutStep(1, autoTypes.WorkflowStepKindExpressions),
			layoutStep(9, autoTypes.WorkflowStepKindExpressions),
		}
		layoutWorkflowSteps(steps, nil)

		x1, _ := pos(t, steps, 1)
		x9, _ := pos(t, steps, 9)
		require.NotEqual(t, x1, x9, "two entry points do not share a column")
	})

	t.Run("a step that carries a position is left alone", func(t *testing.T) {
		fixed := layoutStep(1, autoTypes.WorkflowStepKindExpressions)
		fixed.Meta.Visual = map[string]interface{}{
			"id":     "1",
			"parent": "1",
			"xywh":   []interface{}{999, 777, 200, 80},
		}
		steps := autoTypes.WorkflowStepSet{fixed, layoutStep(2, autoTypes.WorkflowStepKindExpressions)}
		layoutWorkflowSteps(steps, autoTypes.WorkflowPathSet{layoutPath(1, 2)})

		x1, y1 := pos(t, steps, 1)
		require.Equal(t, 999, x1)
		require.Equal(t, 777, y1)

		x2, _ := pos(t, steps, 2)
		require.Greater(t, x2, x1, "the computed graph is laid out clear of it")
	})

	t.Run("a fully placed graph is not touched", func(t *testing.T) {
		one := layoutStep(1, autoTypes.WorkflowStepKindExpressions)
		one.Meta.Visual = map[string]interface{}{"xywh": []interface{}{10, 20, 200, 80}}
		two := layoutStep(2, autoTypes.WorkflowStepKindExpressions)
		two.Meta.Visual = map[string]interface{}{"xywh": []interface{}{30, 40, 200, 80}}

		steps := autoTypes.WorkflowStepSet{one, two}
		layoutWorkflowSteps(steps, autoTypes.WorkflowPathSet{layoutPath(1, 2)})

		x1, y1 := pos(t, steps, 1)
		x2, y2 := pos(t, steps, 2)
		require.Equal(t, []int{10, 20, 30, 40}, []int{x1, y1, x2, y2})
	})

	t.Run("a zero-sized box counts as unplaced", func(t *testing.T) {
		// What the editor writes for a node nothing ever positioned.
		s := layoutStep(1, autoTypes.WorkflowStepKindExpressions)
		s.Meta.Visual = map[string]interface{}{"xywh": []interface{}{0, 0, 0, 0}}
		steps := autoTypes.WorkflowStepSet{s}
		layoutWorkflowSteps(steps, nil)

		x, y := pos(t, steps, 1)
		require.Equal(t, layoutOriginX, x)
		require.Equal(t, layoutOriginY, y)
	})

	t.Run("no steps is not a panic", func(t *testing.T) {
		require.NotPanics(t, func() { layoutWorkflowSteps(nil, nil) })
		require.NotPanics(t, func() {
			layoutWorkflowSteps(autoTypes.WorkflowStepSet{nil}, autoTypes.WorkflowPathSet{nil})
		})
	})
}

func TestLayoutTriggerVisual(t *testing.T) {
	steps := autoTypes.WorkflowStepSet{layoutStep(4, autoTypes.WorkflowStepKindPrompt)}
	layoutWorkflowSteps(steps, nil)
	stepX, stepY := pos(t, steps, 4)

	t.Run("sits one row above the step it starts", func(t *testing.T) {
		v := layoutTriggerVisual(nil, 0, "onManual", 4, steps)

		x, y, ok := visualPosition(v)
		require.True(t, ok)
		require.Equal(t, stepX, x)
		require.Equal(t, stepY-layoutPitchY, y)
	})

	t.Run("carries an id out of the step id space", func(t *testing.T) {
		require.Equal(t, "1000", layoutTriggerVisual(nil, 0, "onManual", 4, steps)["id"])
		require.Equal(t, "1001", layoutTriggerVisual(nil, 1, "onManual", 4, steps)["id"])
	})

	t.Run("draws the edge to its step", func(t *testing.T) {
		v := layoutTriggerVisual(nil, 0, "onManual", 4, steps)

		edges, ok := v["edges"].([]interface{})
		require.True(t, ok)
		require.Len(t, edges, 1)

		edge := edges[0].(map[string]interface{})
		require.Equal(t, "1000", edge["parentID"])
		require.Equal(t, "4", edge["childID"])

		vis := edge["meta"].(map[string]interface{})["visual"].(map[string]interface{})
		require.Equal(t, "e-1000-4", vis["id"])
		require.Equal(t, layoutEdgeStyle, vis["style"])
	})

	t.Run("a trigger with no step gets no edge", func(t *testing.T) {
		v := layoutTriggerVisual(nil, 0, "onManual", 0, steps)
		require.Nil(t, v["edges"])
		_, _, ok := visualPosition(v)
		require.True(t, ok, "it is still placed")
	})

	t.Run("a visual the caller supplied is left alone", func(t *testing.T) {
		given := map[string]interface{}{"id": "42", "xywh": []interface{}{5, 6, 200, 80}}
		v := layoutTriggerVisual(given, 0, "onManual", 4, steps)

		require.Equal(t, "42", v["id"])
		x, y, _ := visualPosition(v)
		require.Equal(t, []int{5, 6}, []int{x, y})
	})
}

func TestLayoutWorkflowPaths(t *testing.T) {
	// The graph the layout produces for a gateway that splits and re-joins.
	build := func() (autoTypes.WorkflowStepSet, autoTypes.WorkflowPathSet) {
		steps := autoTypes.WorkflowStepSet{
			layoutStep(1, autoTypes.WorkflowStepKindGateway),
			layoutStep(2, autoTypes.WorkflowStepKindExpressions),
			layoutStep(3, autoTypes.WorkflowStepKindExpressions),
			layoutStep(4, autoTypes.WorkflowStepKindTermination),
		}
		paths := autoTypes.WorkflowPathSet{
			layoutPath(1, 2), layoutPath(1, 3), layoutPath(2, 4), layoutPath(3, 4),
		}
		layoutWorkflowSteps(steps, paths)
		layoutWorkflowPaths(paths, steps)
		return steps, paths
	}

	t.Run("no two edges claim the same handle on a node", func(t *testing.T) {
		_, paths := build()

		seen := map[string]bool{}
		for _, p := range paths {
			style, _ := p.Meta.Visual["style"].(string)
			require.NotEmpty(t, style)

			exit := style[:strings.Index(style, ";exitDx")]
			entry := style[strings.Index(style, "entryX"):strings.Index(style, ";entryDx")]

			out := fmt.Sprintf("%d/%s", p.ParentID, exit)
			in := fmt.Sprintf("%d/%s", p.ChildID, entry)
			require.False(t, seen[out], "two edges leave step %d by the same handle", p.ParentID)
			require.False(t, seen[in], "two edges enter step %d by the same handle", p.ChildID)
			seen[out], seen[in] = true, true
		}
	})

	t.Run("an edge into the next column leaves by the side", func(t *testing.T) {
		_, paths := build()

		// 1->2 carries the column, 1->3 crosses into the one on its right.
		require.Contains(t, paths[0].Meta.Visual["style"], "exitX=0.5;exitY=1")
		require.Contains(t, paths[1].Meta.Visual["style"], "exitX=1;exitY=0.5")
	})

	t.Run("every path gets an id and a connector", func(t *testing.T) {
		paths := autoTypes.WorkflowPathSet{layoutPath(1, 2), layoutPath(2, 5)}
		layoutWorkflowPaths(paths, nil)

		require.Equal(t, "p-1-2", paths[0].Meta.Visual["id"])
		require.Equal(t, "p-2-5", paths[1].Meta.Visual["id"])
		for _, p := range paths {
			require.Equal(t, "1", p.Meta.Visual["parent"])
			require.NotEmpty(t, p.Meta.Visual["style"])
		}
	})

	t.Run("a style the caller supplied is left alone", func(t *testing.T) {
		p := layoutPath(1, 2)
		p.Meta.Visual = map[string]interface{}{"style": "exitX=1;exitY=0.5;entryX=0;entryY=0.5;"}
		layoutWorkflowPaths(autoTypes.WorkflowPathSet{p}, nil)

		require.Equal(t, "exitX=1;exitY=0.5;entryX=0;entryY=0.5;", p.Meta.Visual["style"])
	})

	t.Run("no paths is not a panic", func(t *testing.T) {
		require.NotPanics(t, func() { layoutWorkflowPaths(nil, nil) })
		require.NotPanics(t, func() { layoutWorkflowPaths(autoTypes.WorkflowPathSet{nil}, nil) })
	})
}
