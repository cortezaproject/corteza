package dal

import (
	"context"
	"fmt"
)

type (
	// Join produces a series of joined rows from the provided sources based
	// on the JoinPredicate.
	//
	// The join step produces an SQL left join-like output where all right rows
	// have a corresponding left row.
	Join struct {
		Ident    string
		RelLeft  string
		RelRight string
		// @todo allow multiple join predicates; for now (for easier indexing)
		// only allow one (this is the same as we had before)
		On     JoinPredicate
		Filter internalFilter

		OutAttributes   []AttributeMapping
		LeftAttributes  []AttributeMapping
		RightAttributes []AttributeMapping

		relLeft  PipelineStep
		relRight PipelineStep
		plan     joinPlan
		analysis stepAnalysis
	}

	// JoinPredicate determines the attributes the two datasets should get joined on
	JoinPredicate struct {
		Left  string
		Right string
	}

	// joinPlan outlines how the optimizer determined the two datasets should be
	// joined on.
	joinPlan struct {
		// @todo add strategy when we have different strategies implemented
		// strategy string

		// partialScan indicates we can partially pull data from the two sources
		// as the data is provided in the correct order.
		partialScan bool
	}
)

func (def *Join) Identifier() string {
	return def.Ident
}

func (def *Join) Sources() []string {
	return []string{def.RelLeft, def.RelRight}
}

func (def *Join) Analyze(ctx context.Context) (err error) {
	def.analysis = stepAnalysis{
		scanCost:   costUnknown,
		searchCost: costUnknown,
		filterCost: costUnknown,
		sortCost:   costUnknown,
		outputSize: sizeUnknown,
	}
	return
}

func (def *Join) Analysis() stepAnalysis {
	return def.analysis
}

func (def *Join) Optimize(req internalFilter) (res internalFilter, err error) {
	err = fmt.Errorf("not implemented")
	return
}

func (def *Join) Initialize(ctx context.Context, ii ...Iterator) (out Iterator, err error) {
	err = def.validate(ii)
	if err != nil {
		return
	}

	// @todo adjust the used exec based on other strategies when added
	exec := &joinLeft{
		def:         *def,
		filter:      def.Filter,
		leftSource:  ii[0],
		rightSource: ii[1],
	}

	return exec, exec.init(ctx)
}

func (def *Join) validate(ii []Iterator) (err error) {
	err = func() (err error) {
		if len(ii) != 2 {
			return fmt.Errorf("expected 2 iterators, got %d", len(ii))
		}

		if len(def.OutAttributes) == 0 {
			return fmt.Errorf("no attributes specified")
		}
		if len(def.LeftAttributes) == 0 {
			return fmt.Errorf("no left attributes specified")
		}
		if len(def.RightAttributes) == 0 {
			return fmt.Errorf("no right attributes specified")
		}

		if def.On.Left == "" {
			return fmt.Errorf("no left attribute in the join predicate specified")
		}
		if def.On.Right == "" {
			return fmt.Errorf("no right attribute in the join predicate specified")
		}

		return
	}()
	if err != nil {
		return fmt.Errorf("invalid definition: %v", err)
	}

	return
}
