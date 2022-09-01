package dal

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	// Aggregate produces a series of aggregated rows from the provided sources based
	// on the specified group.
	Aggregate struct {
		Ident     string
		RelSource string
		Filter    filter.Filter
		filter    internalFilter

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

func (def *Aggregate) Attributes() [][]AttributeMapping {
	return [][]AttributeMapping{append(def.Group, def.OutAttributes...)}
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

func (def *Aggregate) init(ctx context.Context) (err error) {
	if def.Filter != nil {
		def.filter, err = toInternalFilter(def.Filter)
		if err != nil {
			return
		}
	}

	if len(def.SourceAttributes) == 0 {
		def.SourceAttributes = collectAttributes(def.rel)
	}

	err = def.validate()
	if err != nil {
		return
	}

	return nil
}

func (def *Aggregate) exec(ctx context.Context, src Iterator) (out Iterator, err error) {
	exec := &aggregate{
		def:    *def,
		filter: def.filter,
		source: src,
	}

	return exec, exec.init(ctx)
}

func (def *Aggregate) validate() (err error) {
	err = func() (err error) {
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
