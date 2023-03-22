package dal

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/ql"
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

		analysis   map[string]OpAnalysis
		connection *ConnectionWrap
		model      *Model

		// clobbered lists all of the steps that are offloaded into the datasource.
		// The list is provided in order; the first step is the first step which should execute
		// in the report pipeline.
		//
		// @todo change to generic step; doing this because I can for now
		clobbered []*Aggregate

		// provided in the init step so we can omit some code in the exec step
		// @todo consider removing this
		auxIter Iterator
	}
)

func (def *Datasource) Identifier() string {
	return def.Ident
}

func (def *Datasource) Sources() []string {
	return []string{}
}

func (def *Datasource) Attributes() [][]AttributeMapping {
	if len(def.clobbered) > 0 {
		return def.offloadedAttributes()
	}

	return def.ownAttributes()
}

func (def *Datasource) Analyze(ctx context.Context) (err error) {
	a, err := def.connection.connection.Analyze(ctx, def.model)
	if err != nil {
		return
	}

	def.analysis = map[string]OpAnalysis{
		OpAnalysisIterate: {
			ScanCost:   CostUnknown,
			SearchCost: CostUnknown,
			FilterCost: CostUnknown,
			SortCost:   CostUnknown,
			OutputSize: SizeUnknown,
		},
	}

	if _, ok := a[OpAnalysisAggregate]; ok {
		def.analysis[OpAnalysisAggregate] = a[OpAnalysisAggregate]
	}
	if _, ok := a[OpAnalysisJoin]; ok {
		def.analysis[OpAnalysisJoin] = a[OpAnalysisJoin]
	}

	return
}

func (def *Datasource) Analysis() map[string]OpAnalysis {
	return def.analysis
}

func (def *Datasource) Optimize(req internalFilter) (res internalFilter, err error) {
	return internalFilter{}, fmt.Errorf("optimization not implemented")
}

func (def *Datasource) init(ctx context.Context) (err error) {
	if def.model == nil {
		return fmt.Errorf("cannot initialize datasource: model not set")
	}
	if def.connection == nil {
		return fmt.Errorf("cannot initialize datasource: connection not set")
	}

	if def.Filter != nil {
		def.filter, err = toInternalFilter(def.Filter)
		if err != nil {
			return
		}
	}

	if len(def.OutAttributes) == 0 {
		def.OutAttributes = def.outAttrsFromModel(def.model)
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

	// Firstly validate the base
	err = def.validate()
	if err != nil {
		return
	}

	if len(def.clobbered) > 0 {
		// @todo currently, we can only do one; change this when we tweak the DB
		//       offloading code.
		ag := def.clobbered[0]

		// Invoke the aggregation's init to perform the validation and preparation logic
		var wa *aggregate
		wa, err = ag.init(ctx, nil)
		if err != nil {
			return
		}

		// Preprocess the filters to conform to connection API
		f, having, err := def.getAggregationFilters(def.filter, wa.filter)
		if err != nil {
			return err
		}

		def.auxIter, err = def.connection.connection.Aggregate(ctx, def.model, f, wa.groupDefs, wa.aggregateDefs, having)
		return err
	}

	def.auxIter, err = def.connection.connection.Search(ctx, def.model, def.filter)
	return err
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
				IsPrimary:    attr.PrimaryKey,
				IsSystem:     attr.System,
				IsFilterable: attr.Filterable,
				IsSortable:   attr.Sortable,
				Nullable:     attr.Type.IsNullable(),
				Type:         attr.Type,
			},
		})
	}

	return
}

func (def *Datasource) shouldClobberAggregation() bool {
	costs, ok := def.analysis[OpAnalysisAggregate]
	if !ok {
		return false
	}

	// @todo more to it; check and compare costs; for now we know that all rdbms
	//       offloads will be faster.

	_ = costs

	return true
}

func (def *Datasource) ownAttributes() [][]AttributeMapping {
	return [][]AttributeMapping{def.OutAttributes}
}

func (def *Datasource) offloadedAttributes() [][]AttributeMapping {
	// The last offloaded step's attributes is what this step's iterator returns
	return def.clobbered[len(def.clobbered)-1].Attributes()
}

func (def *Datasource) clobber(s PipelineStep) (ok bool) {
	switch cs := s.(type) {
	case *Aggregate:
		if !def.shouldClobberAggregation() {
			return false
		}

		def.clobbered = append(def.clobbered, cs)
		return true
	}

	return
}

// getAggregationFilters transforms the base and the aggregation filters to conform
// to the store/dal's model.Aggregate API
//
// - The full filter.Filter return parameter is applied to the base dataset
// - The QL node return parameter is applied only to the aggregated dataset
//
// @todo should we change the store's API to accept 2x filter.Filter? I think
//       that would allow the underlaying driver to decide how to handle them
//       instead of relying on what SQL does.
func (def *Datasource) getAggregationFilters(base, agg internalFilter) (filter internalFilter, having *ql.ASTNode, err error) {
	var typedV expr.TypedValue
	filter = base

	// Move the aggregation's order and limit to the base filter because those
	// are applied to the final output (per SQL)
	//
	// Leave base filtering parameters (filter, constraints, ...) as those are apploed
	// to the dataset before aggregation (which is correct).
	filter.orderBy = agg.orderBy
	filter.limit = agg.limit

	// Convert the rest of the filtering parameters defined on the aggregation's filter
	// to the QL node.
	// ALl of that is applied after the aggregation which is correct.
	var nConstraints *ql.ASTNode
	if len(agg.constraints) > 0 {
		nConstraints = &ql.ASTNode{
			Ref: "and",
		}
		for k, c := range agg.constraints {
			arg := &ql.ASTNode{
				Ref: "or",
			}

			for _, v := range c {
				typedV, err = expr.Typify(v)
				if err != nil {
					return
				}
				arg.Args = append(arg.Args, &ql.ASTNode{
					Symbol: k,
					Value:  ql.WrapValue(typedV),
				})
			}

			nConstraints.Args = append(nConstraints.Args, &ql.ASTNode{
				Ref:  "group",
				Args: ql.ASTNodeSet{arg},
			})
		}
	}

	var nStateConstraints *ql.ASTNode
	if len(agg.stateConstraints) > 0 {
		nStateConstraints = &ql.ASTNode{
			Ref: "and",
		}
		for k, c := range agg.stateConstraints {
			typedV, err = expr.Typify(c)
			if err != nil {
				return
			}

			nStateConstraints.Args = append(nStateConstraints.Args, &ql.ASTNode{
				Symbol: k,
				Value:  ql.WrapValue(typedV),
			})
		}
	}

	var nMetaConstraints *ql.ASTNode
	if len(agg.stateConstraints) > 0 {
		nMetaConstraints = &ql.ASTNode{
			Ref: "and",
		}
		for k, c := range agg.stateConstraints {
			typedV, err = expr.Typify(c)
			if err != nil {
				return
			}

			nMetaConstraints.Args = append(nMetaConstraints.Args, &ql.ASTNode{
				Symbol: k,
				Value:  ql.WrapValue(typedV),
			})
		}
	}

	nExpression := agg.expParsed

	var nCursor *ql.ASTNode
	if agg.cursor != nil {
		nCursor, err = agg.cursor.ToAST(nil, nil)
		if err != nil {
			return
		}
	}

	having = &ql.ASTNode{
		Ref: "and",
	}

	if nConstraints != nil {
		having.Args = append(having.Args, nConstraints)
	}
	if nStateConstraints != nil {
		having.Args = append(having.Args, nStateConstraints)
	}
	if nMetaConstraints != nil {
		having.Args = append(having.Args, nMetaConstraints)
	}
	if nExpression != nil {
		having.Args = append(having.Args, nExpression)
	}
	if nCursor != nil {
		having.Args = append(having.Args, nCursor)
	}

	// In case everything is empty, no need for the having part
	if len(having.Ref) == 0 {
		having = nil
	}

	return
}
