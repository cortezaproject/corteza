package report

import (
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/qlng"
	"github.com/spf13/cast"
)

type (
	// frameLoadCtx encapsulates some loading metadata to make it easier to work with
	frameLoadCtx struct {
		// initLoader initializes a fresh loader in case where it can't be evaluated beforehand
		initLoader func(int, *Filter) (Loader, Closer, error)

		loader Loader
		closer Closer

		// General state for easier work
		metaInitialized bool
		sorting         filter.SortExprSet
		sortColumns     []int

		keyCol      string
		keyColIndex int
	}
)

// keys returns the unique key values based on the key column and the provided frames
func (bl *frameLoadCtx) keys(ff []*Frame) (keys []string, err error) {
	keys = make([]string, 0, defaultPageSize)
	keySet := make(map[string]bool)
	var k string

	for _, f := range ff {
		err = f.WalkRows(func(i int, r FrameRow) error {
			k, err = cast.ToStringE(r[bl.keyColIndex].Get())
			if ok := keySet[k]; !ok {
				keys = append(keys, k)
				keySet[k] = true
			}
			return err
		})
		if err != nil {
			return
		}
	}

	return
}

// keyFilter prepares the filter that should be used when fetching related rows.
//
// @todo do some compression, ie "id > x && id < y"
//       this will return more stuff but it could be faster then the current thing
func (bl *frameLoadCtx) keyFilter(keys []string) *Filter {
	aa := make(qlng.ASTNodeSet, len(keys))

	for i, k := range keys {
		aa[i] = &qlng.ASTNode{
			Ref: "eq",
			Args: qlng.ASTNodeSet{
				&qlng.ASTNode{Symbol: bl.keyCol},
				&qlng.ASTNode{Value: qlng.MakeValueOf("String", k)},
			},
		}
	}

	return &Filter{
		ASTNode: &qlng.ASTNode{
			Ref: "group",
			Args: qlng.ASTNodeSet{
				&qlng.ASTNode{
					Ref:  "or",
					Args: aa,
				},
			},
		},
	}
}
