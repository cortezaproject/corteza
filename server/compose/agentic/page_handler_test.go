package agentic

import (
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
)

// blocksOf builds a block set from kind/xywh pairs, since kind and layout are
// all these tests care about.
func blocksOf(pairs ...any) cmpTypes.PageBlocks {
	out := make(cmpTypes.PageBlocks, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, cmpTypes.PageBlock{
			Kind: pairs[i].(string),
			XYWH: pairs[i+1].([4]int),
		})
	}
	return out
}

func TestAutoLayoutBlocksKeepsExplicitLayout(t *testing.T) {
	// A row of four tiles above a full-width list: the layout the tool
	// description tells an agent to send, which must arrive unchanged.
	in := blocksOf(
		"Metric", [4]int{0, 0, 12, 20},
		"Metric", [4]int{12, 0, 12, 20},
		"Metric", [4]int{24, 0, 12, 20},
		"Metric", [4]int{36, 0, 12, 20},
		"RecordList", [4]int{0, 20, 48, 30},
	)
	want := []([4]int){{0, 0, 12, 20}, {12, 0, 12, 20}, {24, 0, 12, 20}, {36, 0, 12, 20}, {0, 20, 48, 30}}

	for i, b := range autoLayoutBlocks(in) {
		if b.XYWH != want[i] {
			t.Errorf("block %d (%s): got xywh %v, want %v", i, b.Kind, b.XYWH, want[i])
		}
	}
}

func TestAutoLayoutBlocksFlowsBlocksWithoutWidth(t *testing.T) {
	// Nothing positioned: tiles flow four to a row, charts two to a row, and a
	// RecordList takes the full width on a row of its own.
	in := blocksOf(
		"Metric", [4]int{},
		"Metric", [4]int{},
		"Metric", [4]int{},
		"Metric", [4]int{},
		"Chart", [4]int{},
		"Chart", [4]int{},
		"RecordList", [4]int{},
	)
	want := []([4]int){
		{0, 0, 12, 20}, {12, 0, 12, 20}, {24, 0, 12, 20}, {36, 0, 12, 20},
		{0, 20, 24, 30}, {24, 20, 24, 30},
		{0, 50, 48, 30},
	}

	for i, b := range autoLayoutBlocks(in) {
		if b.XYWH != want[i] {
			t.Errorf("block %d (%s): got xywh %v, want %v", i, b.Kind, b.XYWH, want[i])
		}
	}
}

func TestAutoLayoutBlocksWrapsWhenRowFills(t *testing.T) {
	// Five quarter-width tiles: the fifth cannot fit beside the other four, so
	// it starts a new row below the tallest block of the one it left.
	got := autoLayoutBlocks(blocksOf(
		"Metric", [4]int{},
		"Metric", [4]int{},
		"Metric", [4]int{},
		"Metric", [4]int{},
		"Metric", [4]int{},
	))

	if got[4].XYWH != [4]int{0, 20, 12, 20} {
		t.Errorf("fifth tile: got xywh %v, want [0 20 12 20]", got[4].XYWH)
	}
}

func TestAutoLayoutBlocksPlacesAutoBlocksBelowExplicitOnes(t *testing.T) {
	// Mixed page: the auto-placed chart must clear the positioned list rather
	// than land on top of it, wherever the two sit in the array.
	got := autoLayoutBlocks(blocksOf(
		"Chart", [4]int{},
		"RecordList", [4]int{0, 0, 48, 30},
	))

	if got[0].XYWH != [4]int{0, 30, 24, 30} {
		t.Errorf("auto chart: got xywh %v, want [0 30 24 30]", got[0].XYWH)
	}
	if got[1].XYWH != [4]int{0, 0, 48, 30} {
		t.Errorf("explicit list moved: got xywh %v, want [0 0 48 30]", got[1].XYWH)
	}
}

func TestAutoLayoutBlocksFillsMissingHeightAndClampsOffGrid(t *testing.T) {
	got := autoLayoutBlocks(blocksOf(
		"Metric", [4]int{0, 0, 12, 0}, // width given, height missing
		"RecordList", [4]int{-2, 0, 96, 30}, // negative x, wider than the grid
	))

	if got[0].XYWH[3] != 20 {
		t.Errorf("missing height: got h %d, want the Metric default 20", got[0].XYWH[3])
	}
	if got[1].XYWH[0] != 0 || got[1].XYWH[2] != gridColumns {
		t.Errorf("off-grid block: got xywh %v, want x 0 and w %d", got[1].XYWH, gridColumns)
	}
}

func TestInheritBlockLayoutKeepsPositionOfEditedBlocks(t *testing.T) {
	existing := cmpTypes.PageBlocks{
		{BlockID: 1, Kind: "Metric", XYWH: [4]int{12, 40, 12, 20}},
		{BlockID: 2, Kind: "Chart", XYWH: [4]int{0, 0, 24, 30}},
	}
	// One block edited without layout, one moved deliberately, one new.
	incoming := cmpTypes.PageBlocks{
		{BlockID: 1, Kind: "Metric"},
		{BlockID: 2, Kind: "Chart", XYWH: [4]int{24, 0, 24, 30}},
		{Kind: "RecordList"},
	}

	inheritBlockLayout(existing, incoming)

	if incoming[0].XYWH != [4]int{12, 40, 12, 20} {
		t.Errorf("edited block moved: got xywh %v, want the stored [12 40 12 20]", incoming[0].XYWH)
	}
	if incoming[1].XYWH != [4]int{24, 0, 24, 30} {
		t.Errorf("explicit move lost: got xywh %v, want [24 0 24 30]", incoming[1].XYWH)
	}
	if incoming[2].XYWH != [4]int{} {
		t.Errorf("new block inherited a layout it has no claim to: %v", incoming[2].XYWH)
	}
}

// TestFoldBlockOptionAliases pins the rename half of ref resolution. The value
// being a real ID is not enough: the webapp reads options.moduleID and
// options.chartID, so a resolved ID left under "module"/"chart" renders as
// "No module selected" / "invalid ID" with the right ID sitting next to it.
func TestFoldBlockOptionAliases(t *testing.T) {
	tests := []struct {
		name    string
		options map[string]any
		want    map[string]any
	}{
		{
			name:    "module alias",
			options: map[string]any{"module": "expense"},
			want:    map[string]any{"moduleID": "expense"},
		},
		{
			name:    "chart alias",
			options: map[string]any{"chart": "spend_by_category"},
			want:    map[string]any{"chartID": "spend_by_category"},
		},
		{
			// The canonical key is the one the caller meant; the alias goes.
			name:    "both keys, canonical wins",
			options: map[string]any{"module": "expense", "moduleID": "12345"},
			want:    map[string]any{"moduleID": "12345"},
		},
		{
			name:    "canonical key alone is untouched",
			options: map[string]any{"moduleID": "12345"},
			want:    map[string]any{"moduleID": "12345"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			foldBlockOptionAliases(tc.options)

			if len(tc.options) != len(tc.want) {
				t.Fatalf("got options %v, want %v", tc.options, tc.want)
			}
			for k, v := range tc.want {
				if tc.options[k] != v {
					t.Errorf("option %q: got %v, want %v", k, tc.options[k], v)
				}
			}
		})
	}
}

func TestParsePageConfigRejectsUnrenderableIcons(t *testing.T) {
	// The shape the tool description used to hand out, which renders as a
	// broken image in the navigation.
	_, err := parsePageConfig(`{"navItem":{"icon":{"type":"library","src":"font-awesome://home"}}}`)
	if err == nil {
		t.Fatal("expected a library icon to be rejected, got no error")
	}

	for _, cfg := range []string{
		`{"navItem":{"expanded":true}}`,
		`{"navItem":{"icon":{"type":"link","src":"https://example.test/i.png"}}}`,
		`{"navItem":{"icon":{"type":"attachment","src":"/compose/namespace/1/page/2/attachment/3/original/i.png"}}}`,
	} {
		if _, err := parsePageConfig(cfg); err != nil {
			t.Errorf("config %s: unexpected error %v", cfg, err)
		}
	}
}
