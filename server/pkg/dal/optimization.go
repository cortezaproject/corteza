package dal

import "fmt"

type (
	// opCost provides a general idea of expensive an operation is for a
	// specific pipeline step.
	//
	// dsSize provides a general idea of how large an underlaying dataset is.
	opCost int
	dsSize int

	OpAnalysis struct {
		ScanCost   opCost
		SearchCost opCost
		FilterCost opCost
		SortCost   opCost

		OutputSize dsSize
	}

	ppStepWrap struct {
		step PipelineStep

		parent *ppStepWrap
		child  []*ppStepWrap
	}
)

var (
	pipelineOptimizers = []func(Pipeline) (Pipeline, error){
		pipelineClobberSteps,
	}
)

const (
	// operation computation indicators
	CostUnknown opCost = iota
	CostFree
	CostCheep
	CostAcceptable
	CostExpensive
	CostInfinite
)

const (
	// dataset size indicators
	SizeUnknown dsSize = iota
	SizeTiny
	SizeSmall
	SizeMedium
	SizeLarge
)

const (
	OpAnalysisIterate   string = "iterate"
	OpAnalysisAggregate string = "aggregate"
	OpAnalysisJoin      string = "join"
)

// wrapPpSteps wraps the pipeline steps in a more processing friendly format
// and returns a slice of leave nodes
//
// @todo the pipeline representation might change which will make this obsolete
func wrapPpSteps(pp Pipeline) (leaves []*ppStepWrap) {
	ix := make(map[PipelineStep]*ppStepWrap)

	for _, p := range pp {
		ix[p] = &ppStepWrap{step: p}
	}

	for _, p := range pp {
		switch c := p.(type) {
		case *Aggregate:
			ix[p].child = append(ix[p].child, ix[c.rel])
			ix[c.rel].parent = ix[p]

		case *Join:
			ix[p].child = append(ix[p].child, ix[c.relLeft])
			ix[c.relLeft].parent = ix[p]

			ix[p].child = append(ix[p].child, ix[c.relRight])
			ix[c.relRight].parent = ix[p]

		case *Link:
			ix[p].child = append(ix[p].child, ix[c.relLeft])
			ix[c.relLeft].parent = ix[p]

			ix[p].child = append(ix[p].child, ix[c.relRight])
			ix[c.relRight].parent = ix[p]

		case *Datasource:
			continue

		default:
			panic(fmt.Errorf("impossible state: unknown pipeline step type %v", p))
		}
	}

	for _, a := range ix {
		if len(a.child) == 0 {
			leaves = append(leaves, a)
		}
	}

	return
}

// unwrapPpSteps unwraps the wrapped pipeline step into the classic Pipeline
//
// @todo the pipeline representation might change which will make this obsolete
func unwrapPpSteps(n *ppStepWrap, seen map[*ppStepWrap]bool) (out Pipeline) {
	for i, c := range n.child {
		switch s := n.step.(type) {
		case *Aggregate:
			s.rel = c.step
			s.RelSource = c.step.Identifier()
		case *Join:
			if i == 0 {
				s.relLeft = c.step
				s.RelLeft = c.step.Identifier()
			} else {
				s.relRight = c.step
				s.RelRight = c.step.Identifier()
			}
		case *Link:
			if i == 0 {
				s.relLeft = c.step
				s.RelLeft = c.step.Identifier()
			} else {
				s.relRight = c.step
				s.RelRight = c.step.Identifier()
			}
		}
	}

	// @todo potentially cancancel sooner
	if !seen[n] {
		out = append(out, n.step)
	}
	seen[n] = true

	if n.parent != nil {
		out = append(unwrapPpSteps(n.parent, seen), out...)
	}

	return
}
