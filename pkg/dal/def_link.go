package dal

import (
	"context"
	"fmt"
)

type (
	// Link produces a series of left and corresponding right rows based on the
	// provided sources and the LinkPredicate.
	//
	// The Link step produces an SQL left join-like output where left and right
	// rows are served separately and the left rows are not duplicated.
	Link struct {
		Ident    string
		RelLeft  string
		RelRight string
		// @todo allow multiple link predicates; for now (for easier indexing)
		// only allow one (this is the same as we had before)
		On     LinkPredicate
		Filter internalFilter

		OutLeftAttributes  []AttributeMapping
		OutRightAttributes []AttributeMapping
		LeftAttributes     []AttributeMapping
		RightAttributes    []AttributeMapping

		relLeft  PipelineStep
		relRight PipelineStep

		plan     linkPlan
		analysis stepAnalysis
	}

	// LinkPredicate determines the attributes the two datasets should get joined on
	LinkPredicate struct {
		Left  string
		Right string
	}

	// linkPlan outlines how the optimizer determined the two datasets should be
	// joined on.
	linkPlan struct {
		// @todo add strategy when we have different strategies implemented
		// strategy string

		// partialScan indicates we can partially pull data from the two sources
		// as the data is provided in the correct order.
		partialScan bool
	}
)

func (def *Link) Identifier() string {
	return def.Ident
}

func (def *Link) Sources() []string {
	return []string{def.RelLeft, def.RelRight}
}

func (def *Link) Analyze(ctx context.Context) (err error) {
	def.analysis = stepAnalysis{
		scanCost:   costUnknown,
		searchCost: costUnknown,
		filterCost: costUnknown,
		sortCost:   costUnknown,
		outputSize: sizeUnknown,
	}
	return
}

func (def *Link) Analysis() stepAnalysis {
	return def.analysis
}

func (def *Link) Optimize(req internalFilter) (res internalFilter, err error) {
	err = fmt.Errorf("not implemented")
	return
}

func (def *Link) Initialize(ctx context.Context, ii ...Iterator) (_ Iterator, err error) {
	err = def.validate(ii)
	if err != nil {
		return
	}

	// @todo adjust the used exec based on other strategies when added
	exec := &linkLeft{
		def:         *def,
		filter:      def.Filter,
		leftSource:  ii[0],
		rightSource: ii[1],
	}

	return exec, exec.init(ctx)
}

func (def *Link) validate(ii []Iterator) (err error) {
	err = func() (err error) {
		if len(ii) != 2 {
			return fmt.Errorf("expected 2 iterators, got %d", len(ii))
		}

		if len(def.OutLeftAttributes) == 0 {
			return fmt.Errorf("no left output attributes specified")
		}
		if len(def.OutRightAttributes) == 0 {
			return fmt.Errorf("no right output attributes specified")
		}
		if len(def.LeftAttributes) == 0 {
			return fmt.Errorf("no left attributes specified")
		}
		if len(def.RightAttributes) == 0 {
			return fmt.Errorf("no right attributes specified")
		}

		if def.On.Left == "" {
			return fmt.Errorf("no left attribute in the link predicate specified")
		}
		if def.On.Right == "" {
			return fmt.Errorf("no right attribute in the link predicate specified")
		}

		return
	}()
	if err != nil {
		return fmt.Errorf("invalid definition: %v", err)
	}
	return
}
