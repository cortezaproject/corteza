package dal

import (
	"context"
	"fmt"
)

type (
	// Aggregate produces a series of aggregated rows from the provided sources based
	// on the specified group.
	Aggregate struct {
		Ident     string
		RelSource string
		Filter    internalFilter

		Group         []AttributeMapping
		OutAttributes []AttributeMapping

		SourceAttributes []AttributeMapping

		rel      PipelineStep
		plan     aggregatePlan
		analysis stepAnalysis
	}

	// aggregatePlan outlines how the optimizer determined the dataset should be
	// aggregated
	aggregatePlan struct {
		// partialScan indicates we can partially pull data from the two sources
		// as the data is provided in the correct order.
		partialScan bool
	}

	groupKey       []any
	aggregateGroup struct {
		key groupKey
		agg *aggregator
	}
)

func (def *Aggregate) Identifier() string {
	return def.Ident
}

func (def *Aggregate) Sources() []string {
	return []string{def.RelSource}
}

func (def *Aggregate) Analyze(ctx context.Context) (err error) {
	// @todo proper analysis; for now we'll leave this as defaults
	def.analysis = stepAnalysis{
		scanCost:   costUnknown,
		searchCost: costUnknown,
		filterCost: costUnknown,
		sortCost:   costUnknown,
		outputSize: sizeUnknown,
	}
	return
}

func (def *Aggregate) Analysis() stepAnalysis {
	return def.analysis
}

func (def *Aggregate) Optimize(reqFilter internalFilter) (rspFilter internalFilter, err error) {
	err = fmt.Errorf("not implemented")
	return
}

func (def *Aggregate) Initialize(ctx context.Context, ii ...Iterator) (out Iterator, err error) {
	err = def.validate(ii)
	if err != nil {
		return
	}

	exec := &aggregate{
		def:    *def,
		filter: def.Filter,
		source: ii[0],
	}

	return exec, exec.init(ctx)
}

func (def *Aggregate) validate(ii []Iterator) (err error) {
	err = func() (err error) {
		if len(ii) != 1 {
			return fmt.Errorf("expected 1 iterator, got %d", len(ii))
		}

		if len(def.Group) == 0 {
			return fmt.Errorf("no group attributes specified")
		}

		if len(def.OutAttributes) == 0 {
			return fmt.Errorf("no output attributes specified")
		}

		if len(def.SourceAttributes) == 0 {
			return fmt.Errorf("no source attributes specified")
		}

		return
	}()
	if err != nil {
		return fmt.Errorf("invalid definition: %v", err)
	}

	return
}
