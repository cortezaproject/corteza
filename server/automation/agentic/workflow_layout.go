package agentic

import (
	"encoding/json"
	"fmt"

	autoTypes "github.com/crusttech/human/server/automation/types"
)

// The grid the workflow editor's canvas is authored on. A node is 200x80, a
// step sits one pitch below the one it follows, and a branch takes the next
// column. Read off the workflows people arranged by hand; the editor stores
// absolute coordinates and never lays anything out itself, so a workflow built
// through this API arrives with every node on (0,0) unless it is given a place.
const (
	layoutNodeWidth  = 200
	layoutNodeHeight = 80
	layoutPitchY     = 160
	layoutPitchX     = 280
	layoutOriginX    = 120
	layoutOriginY    = 160

	// Triggers and canvas notes are numbered from 1000 so their ids cannot
	// collide with the step ids paths are written against.
	layoutTriggerIDBase = 1000

	// The connector the editor draws between two vertically stacked nodes:
	// out of the bottom edge, into the top.
	layoutEdgeStyle = "exitX=0.5;exitY=1;exitDx=0;exitDy=0;entryX=0.5;entryY=0;entryDx=0;entryDy=0;"
)

// layoutWorkflowSteps gives a place on the canvas to every step that has none.
//
// The walk follows the paths rather than the array: a step's first successor
// continues its column, every later one opens a fresh column to the right, so a
// gateway's arms and an error handler's catch branch separate on their own. An
// iterator is the one inversion — its first path is the body and its second the
// exit — so the body takes the new column and the exit carries the column down.
//
// Steps that already carry coordinates are left exactly as they are, and the
// computed graph is pushed clear of them.
func layoutWorkflowSteps(steps autoTypes.WorkflowStepSet, paths autoTypes.WorkflowPathSet) {
	if len(steps) == 0 {
		return
	}

	byID := make(map[uint64]*autoTypes.WorkflowStep, len(steps))
	placed := make(map[uint64]bool, len(steps))
	rightmost := 0
	anyPlaced := false

	for _, s := range steps {
		if s == nil {
			continue
		}
		byID[s.ID] = s
		if x, _, ok := visualPosition(s.Meta.Visual); ok {
			placed[s.ID] = true
			anyPlaced = true
			if x > rightmost {
				rightmost = x
			}
		}
	}

	missing := false
	for _, s := range steps {
		if s != nil && !placed[s.ID] {
			missing = true
			break
		}
	}
	if !missing {
		return
	}

	children := make(map[uint64][]uint64, len(paths))
	inbound := make(map[uint64]int, len(steps))
	for _, p := range paths {
		if p == nil {
			continue
		}
		children[p.ParentID] = append(children[p.ParentID], p.ChildID)
		inbound[p.ChildID]++
	}

	originX := layoutOriginX
	if anyPlaced {
		originX = rightmost + layoutPitchX
	}

	var (
		seen    = make(map[uint64]bool, len(steps))
		nextCol int
	)

	var walk func(id uint64, row, col int)
	walk = func(id uint64, row, col int) {
		step, ok := byID[id]
		if !ok || seen[id] {
			// A join, or an edge back into the loop that produced it: the first
			// placement wins, and revisiting would not terminate.
			return
		}
		seen[id] = true

		if col > nextCol {
			nextCol = col
		}

		if !placed[id] {
			step.Meta.Visual = mergeVisual(step.Meta.Visual, map[string]interface{}{
				"id":     fmt.Sprintf("%d", id),
				"parent": "1",
				"value":  step.Meta.Name,
			})
			// Assigned rather than merged: an unplaced node may still carry a
			// zero-sized box, which is a value mergeVisual would keep.
			step.Meta.Visual["xywh"] = []interface{}{
				originX + col*layoutPitchX,
				layoutOriginY + row*layoutPitchY,
				layoutNodeWidth,
				layoutNodeHeight,
			}
		}

		kids := children[id]
		carry := 0
		if step.Kind == autoTypes.WorkflowStepKindIterator && len(kids) > 1 {
			// The body branches away and the exit stays on the column, so the
			// loop reads as a detour rather than the main line.
			carry = 1
		}

		for i, kid := range kids {
			if i == carry {
				walk(kid, row+1, col)
				continue
			}
			nextCol++
			walk(kid, row+1, nextCol)
		}
	}

	// Entry points first, in the order they were given, so the main line lands
	// in the leftmost column. Anything the paths never reach is placed after.
	for _, s := range steps {
		if s == nil || inbound[s.ID] > 0 {
			continue
		}
		if seen[s.ID] {
			continue
		}
		col := 0
		if len(seen) > 0 {
			nextCol++
			col = nextCol
		}
		walk(s.ID, 0, col)
	}

	for _, s := range steps {
		if s == nil || seen[s.ID] {
			continue
		}
		nextCol++
		walk(s.ID, 0, nextCol)
	}
}

// Where an edge may attach, in the order the layout tries them. A handle takes
// one edge and no more — the editor disables a used one — so a node with two
// connections on the same side has to spread them, and an edge that cannot find
// a free handle is dropped by the canvas without a word.
var (
	layoutSourceHandles = []string{
		"source-bottom", "source-bottom-right", "source-bottom-left",
		"source-right", "source-left", "source-top-right", "source-top-left", "source-top",
	}
	layoutTargetHandles = []string{
		"target-top", "target-top-left", "target-top-right",
		"target-left", "target-right", "target-bottom", "target-bottom-left", "target-bottom-right",
	}
	layoutHandleMx = map[string][2]string{
		"source-top":          {"0.5", "0"},
		"source-top-left":     {"0.25", "0"},
		"source-top-right":    {"0.75", "0"},
		"source-bottom":       {"0.5", "1"},
		"source-bottom-left":  {"0.25", "1"},
		"source-bottom-right": {"0.75", "1"},
		"source-right":        {"1", "0.5"},
		"source-left":         {"0", "0.5"},
		"target-top":          {"0.5", "0"},
		"target-top-left":     {"0.25", "0"},
		"target-top-right":    {"0.75", "0"},
		"target-left":         {"0", "0.5"},
		"target-right":        {"1", "0.5"},
		"target-bottom":       {"0.5", "1"},
		"target-bottom-left":  {"0.25", "1"},
		"target-bottom-right": {"0.75", "1"},
	}
)

// layoutWorkflowPaths gives every connection a route between two free handles.
//
// The style is not decoration: the editor reads which handle an edge attaches to
// out of it, and a node whose handles are all named — a termination step, for
// one — has nothing for a style-less edge to land on. Two edges may not share a
// handle either, so a gateway's arms leave by different corners and the branches
// meeting again arrive by different ones.
//
// An edge running to a column on the right leaves by that side, which is what
// makes a branch read as a branch rather than as a line crossing the graph.
func layoutWorkflowPaths(paths autoTypes.WorkflowPathSet, steps autoTypes.WorkflowStepSet) {
	at := make(map[uint64]int, len(steps))
	for _, s := range steps {
		if s == nil {
			continue
		}
		if x, _, ok := visualPosition(s.Meta.Visual); ok {
			at[s.ID] = x
		}
	}

	usedSource := make(map[string]bool, len(paths))
	usedTarget := make(map[string]bool, len(paths))

	for _, p := range paths {
		if p == nil {
			continue
		}
		if p.Meta.Visual != nil {
			if style, has := p.Meta.Visual["style"]; has && style != nil && style != "" {
				continue
			}
		}

		from := layoutSourceHandles
		if at[p.ChildID] > at[p.ParentID] {
			from = []string{"source-right", "source-bottom-right", "source-bottom", "source-bottom-left", "source-top-right", "source-left", "source-top-left", "source-top"}
		}

		source := freeHandle(from, p.ParentID, usedSource)
		target := freeHandle(layoutTargetHandles, p.ChildID, usedTarget)

		p.Meta.Visual = mergeVisual(p.Meta.Visual, map[string]interface{}{
			"id":     fmt.Sprintf("p-%d-%d", p.ParentID, p.ChildID),
			"parent": "1",
			"points": []interface{}{},
			"value":  "",
		})
		p.Meta.Visual["style"] = layoutStyle(source, target)
	}
}

// freeHandle takes the first handle on this node nobody has claimed, and falls
// back to the preferred one when every candidate is spoken for — an overlapping
// edge still draws, where an unattached one does not.
func freeHandle(candidates []string, node uint64, used map[string]bool) string {
	for _, h := range candidates {
		key := fmt.Sprintf("%d/%s", node, h)
		if used[key] {
			continue
		}
		used[key] = true
		return h
	}
	return candidates[0]
}

func layoutStyle(source, target string) string {
	src, tgt := layoutHandleMx[source], layoutHandleMx[target]
	return fmt.Sprintf(
		"exitX=%s;exitY=%s;exitDx=0;exitDy=0;entryX=%s;entryY=%s;entryDx=0;entryDy=0;",
		src[0], src[1], tgt[0], tgt[1],
	)
}

// layoutTriggerVisual places a trigger one pitch above the step it starts, and
// draws the connector between them. Without a visual the editor decodes the
// trigger's node id from a missing field — it becomes the string "undefined" —
// and finds no edge to its step, so the Start node arrives detached.
func layoutTriggerVisual(meta map[string]interface{}, index int, label string, stepID uint64, steps autoTypes.WorkflowStepSet) map[string]interface{} {
	if _, _, ok := visualPosition(meta); ok {
		return meta
	}

	id := fmt.Sprintf("%d", layoutTriggerIDBase+index)

	x, y := layoutOriginX, layoutOriginY-layoutPitchY
	for _, s := range steps {
		if s == nil || s.ID != stepID {
			continue
		}
		if sx, sy, ok := visualPosition(s.Meta.Visual); ok {
			x, y = sx, sy-layoutPitchY
		}
		break
	}

	visual := mergeVisual(meta, map[string]interface{}{
		"id":     id,
		"parent": "1",
		"value":  label,
	})
	visual["xywh"] = []interface{}{x, y, layoutNodeWidth, layoutNodeHeight}

	if stepID == 0 {
		return visual
	}

	child := fmt.Sprintf("%d", stepID)
	visual["edges"] = []interface{}{
		map[string]interface{}{
			"parentID": id,
			"childID":  child,
			"meta": map[string]interface{}{
				"description": "",
				"label":       "",
				"visual": map[string]interface{}{
					"id":     fmt.Sprintf("e-%s-%s", id, child),
					"parent": "1",
					"points": []interface{}{},
					"style":  layoutEdgeStyle,
					"value":  "",
				},
			},
		},
	}

	return visual
}

// visualPosition reads x and y off a visual, and reports whether it carries a
// usable position at all. A zero-sized box is what the editor writes for a node
// nothing ever placed, so it counts as absent.
func visualPosition(visual map[string]interface{}) (x, y int, ok bool) {
	if visual == nil {
		return 0, 0, false
	}

	raw, is := visual["xywh"].([]interface{})
	if !is || len(raw) < 4 {
		return 0, 0, false
	}

	nums := make([]int, 4)
	for i, v := range raw[:4] {
		n, is := toInt(v)
		if !is {
			return 0, 0, false
		}
		nums[i] = n
	}

	if nums[2] == 0 && nums[3] == 0 {
		return 0, 0, false
	}

	return nums[0], nums[1], true
}

func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	case json.Number:
		i, err := n.Int64()
		return int(i), err == nil
	}
	return 0, false
}

// mergeVisual keeps whatever the caller already put on the visual — a label, a
// parent, an edge list — and fills in only the keys it has not set.
func mergeVisual(visual map[string]interface{}, defaults map[string]interface{}) map[string]interface{} {
	if visual == nil {
		visual = make(map[string]interface{}, len(defaults))
	}

	for k, v := range defaults {
		if existing, has := visual[k]; has && existing != nil && existing != "" {
			continue
		}
		visual[k] = v
	}

	return visual
}
