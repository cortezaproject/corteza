package dal

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	// Datasource is a simple passthrough step for underlaying datasources.
	// It exists primarily to make operations consistent.
	Datasource struct {
		Ident    string
		Filter   filter.Filter
		ModelRef ModelRef
		filter   internalFilter

		OutAttributes []AttributeMapping

		analysis stepAnalysis

		// provided in the init step so we can omit some code in the exec step
		auxIter Iterator
	}

	iterProvider func(ctx context.Context, mf ModelRef, f filter.Filter) (iter Iterator, model *Model, err error)
)

func (def *Datasource) Identifier() string {
	return def.Ident
}

func (def *Datasource) Sources() []string {
	return []string{}
}

func (def *Datasource) Attributes() [][]AttributeMapping {
	return [][]AttributeMapping{def.OutAttributes}
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

func (def *Datasource) init(ctx context.Context, s iterProvider) (err error) {
	if def.Filter != nil {
		def.filter, err = toInternalFilter(def.Filter)
		if err != nil {
			return
		}
	}

	// Do the first init to get the model.
	// For now, this is a ~free operation, but it should change when we allow things
	// like nesting reports, etc.
	// @todo refactor datasource descriptors to avoid this
	_, model, err := s(ctx, def.ModelRef, def.filter)
	if err != nil {
		return
	}

	if len(def.OutAttributes) == 0 {
		def.OutAttributes = def.outAttrsFromModel(model)
	}

	pp := make([]string, 0, len(def.OutAttributes)/2+1)
	for _, a := range def.OutAttributes {
		if a.Properties().IsPrimary {
			pp = append(pp, a.Identifier())
		}
	}

	// Assure and attempt to correct the provided sort to conform with the data set and the
	// paging cursor (if any)
	def.filter, err = assureSort(def.filter, pp)
	if err != nil {
		return
	}

	// Get the iterator for actual use
	def.auxIter, _, err = s(ctx, def.ModelRef, def.filter)
	if err != nil {
		return
	}

	err = def.validate()
	if err != nil {
		return
	}

	return nil
}

func (def *Datasource) exec(ctx context.Context) (out Iterator, err error) {
	if def.auxIter == nil {
		return nil, fmt.Errorf("datasource not initialized")
	}

	return def.auxIter, nil
}

func (def *Datasource) validate() (err error) {
	err = func() (err error) {
		if len(def.OutAttributes) == 0 {
			return fmt.Errorf("no attributes specified")
		}

		return
	}()
	if err != nil {
		return fmt.Errorf("invalid definition: %v", err)
	}

	return
}

func (def *Datasource) outAttrsFromModel(model *Model) (attrs []AttributeMapping) {
	for _, attr := range model.Attributes {
		if attr.Type == nil {
			panic(fmt.Sprintf("impossible state: attribute %s has no type", attr.Ident))
		}

		attrs = append(attrs, SimpleAttr{
			Ident: attr.Ident,
			Src:   attr.Ident,
			Props: MapProperties{
				IsPrimary: attr.PrimaryKey,
				IsSystem:  attr.System,
				Nullable:  attr.Type.IsNullable(),
				Type:      attr.Type,
			},
		})
	}

	return
}
