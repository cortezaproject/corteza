package dal

import (
	"context"
	"fmt"
)

type (
	// Datasource is a simple passthrough step for underlaying datasources.
	// It exists primarily to make operations consistent.
	Datasource struct {
		Ident  string
		Filter internalFilter

		Attributes []AttributeMapping
		Source     func(context.Context) (Iterator, error)

		analysis stepAnalysis
	}
)

func (def *Datasource) Identifier() string {
	return def.Ident
}

func (def *Datasource) Sources() []string {
	return []string{}
}

func (def *Datasource) Analyze(ctx context.Context) (err error) {
	// @todo probe datasource; for now, RDBMS only so all is cheap
	def.analysis = stepAnalysis{
		scanCost:   costUnknown,
		searchCost: costUnknown,
		filterCost: costUnknown,
		sortCost:   costUnknown,
		outputSize: sizeUnknown,
	}
	return
}

func (def *Datasource) Analysis() stepAnalysis {
	return def.analysis
}

func (def *Datasource) Optimize(req internalFilter) (res internalFilter, err error) {
	return internalFilter{}, fmt.Errorf("optimization not implemented")
}

func (def *Datasource) Initialize(ctx context.Context, ii ...Iterator) (out Iterator, err error) {
	err = def.validate(ii)
	if err != nil {
		return nil, err
	}

	return def.Source(ctx)
}

func (def *Datasource) validate(ii []Iterator) (err error) {
	err = func() (err error) {
		if len(ii) != 0 {
			return fmt.Errorf("expected 0 iterators, got %d", len(ii))
		}

		if len(def.Attributes) == 0 {
			return fmt.Errorf("no attributes specified")
		}

		return
	}()
	if err != nil {
		return fmt.Errorf("invalid definition: %v", err)
	}

	return
}
