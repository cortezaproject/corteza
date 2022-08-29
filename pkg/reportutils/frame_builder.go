package reportutils

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	// reportFrameBuilder is a helper struct for building report frames from
	// dal iterators
	reportFrameBuilder struct {
		def         *types.ReportFrameDefinition
		frame       *types.ReportFrame
		attrMapping map[string]int
	}
)

func newReportFrameBuilder(def *types.ReportFrameDefinition) *reportFrameBuilder {
	// Index requested columns for easier lookup
	attrMap := make(map[string]int)
	for i, c := range def.Columns {
		attrMap[c.Name] = i
	}

	out := &reportFrameBuilder{
		def:         def,
		attrMapping: attrMap,
	}

	out.freshFrame()

	return out
}

// withRefs includes additional metadata required by the link step
func (b *reportFrameBuilder) withRefs(col string) {
	b.frame.RelColumn = col
}

// addRow adds a new dal.Row to the frame
func (b *reportFrameBuilder) addRow(r *dal.Row) {
	rrow := make(types.ReportFrameRow, len(b.def.Columns))

	for k, cc := range r.CountValues() {
		ix, ok := b.attrMapping[k]
		if !ok {
			continue
		}

		auxRow := make(types.ReportFrameRow, cc)
		for i := uint(0); i < cc; i++ {
			val, _ := r.GetValue(k, i)
			auxRow[i] = b.stringifyVal(k, val)
		}
		rrow[ix] = b.joinMultiVal(k, auxRow)
	}

	b.frame.Rows = append(b.frame.Rows, rrow)

	if b.frame.RelColumn != "" {
		v, _ := r.GetValue(b.frame.RelColumn, 0)
		b.frame.RefValue = b.stringifyVal(b.frame.RelColumn, v)
	}
}

func (b *reportFrameBuilder) stringifyVal(col string, val any) string {
	// @todo nicer formatting and such? V1 didn't do much different
	return fmt.Sprintf("%v", val)
}

func (b *reportFrameBuilder) joinMultiVal(col string, vals []string) string {
	// @todo add delimiter (extend attrs)
	return strings.Join(vals, "\n")
}

// done returns the constructed frame and prepares a new frame with the same
// metadata as the original one
func (b *reportFrameBuilder) done() *types.ReportFrame {
	out := b.frame
	b.freshFrame()

	return out
}

func (b *reportFrameBuilder) freshFrame() {
	if b.frame != nil {
		aux := *b.frame
		b.frame = &aux
		b.frame.Rows = nil
		return
	}

	b.frame = &types.ReportFrame{
		Name:   b.def.Name,
		Source: b.def.Source,
		Ref:    b.def.Ref,

		Columns: b.def.Columns,
		Sort:    b.def.Sort,
		Filter:  b.def.Filter,
	}

	if b.def.Paging != nil {
		b.frame.Paging = &filter.Paging{
			Limit:      b.def.Paging.Limit,
			PageCursor: b.def.Paging.PageCursor,
		}
	}
}
