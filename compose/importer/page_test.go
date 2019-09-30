package importer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
)

func TestPageImport_CastSet(t *testing.T) {
	impFixTester(t, "page_full", func(t *testing.T, page *Page) {
		req := require.New(t)

		req.Len(page.set, 4)

		root := page.set.FindByHandle("root")
		req.NotNil(root)
		req.Equal(root.Title, "Root page")

		sub1 := page.set.FindByHandle("sub1")
		req.NotNil(sub1)
		req.Equal(sub1.Title, "Sub page 1")
		req.Equal(sub1.Blocks, types.PageBlocks{
			{
				Title:   "B1",
				Options: nil,
				Style: types.PageBlockStyle{
					Variants: map[string]string{"v1": "V1"},
				},
				Kind: "TheTestingKind",
				XYWH: [4]int{1, 2, 3, 4},
			},
			{
				Title:   "B2",
				Options: nil,
				Style:   types.PageBlockStyle{},
				Kind:    "TheTestingKind",
				XYWH:    [4]int{11, 12, 13, 14},
			},
		})

		sub2 := page.set.FindByHandle("sub2")
		req.NotNil(sub2)
		req.Equal(sub2.Title, "Sub page 2")

		sub21 := page.set.FindByHandle("sub21")
		req.NotNil(sub21)
		req.Equal(sub21.Title, "Sub-sub page 2.1")

		req.Equal(page.tree, map[string][]string{
			"":     {"root"},
			"root": {"sub1", "sub2"},
			"sub2": {"sub21"},
		})

	})
}
