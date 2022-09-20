package dal

import (
	"context"
	"fmt"
)

type (
	// Pipeline defines a series of steps performed over the data
	//
	// General outline of how a pipeline is used:
	// 1. Populate the pipeline with steps to define the desired outcome
	// 2. Analyze the pipeline to collect statistics and assist optimizers
	// 3. Run structure optimization which may reorder/drop steps
	// 4. Run step optimization which may re-configure individual steps such as
	//    push parts of a filter to lower level, request specific sort order, ...
	// 5. Use the pipeline as an iterator to pull the data from.
	Pipeline []PipelineStep

	// PipelineStep defines an operation performed over the data in the pipeline
	PipelineStep interface {
		Identifier() string
		Sources() []string
		Attributes() [][]AttributeMapping

		Analyze(ctx context.Context) error
		Analysis() stepAnalysis
	}

	// Attribute mapping outlines specific attributes within a pipeline
	// @todo reconsider this interface; potentially remove it or split it up
	AttributeMapping interface {
		Identifier() string
		Expression() (expression string)
		Source() (ident string)
		Properties() MapProperties
	}

	// MapProperties describe the attribute such as it's type and constraints
	MapProperties struct {
		Label     string
		IsPrimary bool
		Nullable  bool
		Type      Type
	}
)

// LinkSteps links related steps into a tree structure
//
// @todo make it return a new slice and not mutate the original
func (pp Pipeline) LinkSteps() (err error) {
	// map steps by identifiers
	steps := make(map[string]PipelineStep)
	for _, s := range pp {
		steps[s.Identifier()] = s
	}

	// Link up...
	err = func() (err error) {
		for _, s := range pp {
			switch rs := s.(type) {
			case *Aggregate:
				rs.rel = steps[rs.RelSource]
				if rs.rel == nil {
					return fmt.Errorf("aggregate: missing source relation %s", rs.RelSource)
				}

			case *Join:
				rs.relLeft = steps[rs.RelLeft]
				rs.relRight = steps[rs.RelRight]
				if rs.relLeft == nil {
					return fmt.Errorf("join: missing left relation %s", rs.relLeft)
				}
				if rs.relRight == nil {
					return fmt.Errorf("join: missing right relation %s", rs.relRight)
				}

			case *Link:
				rs.relLeft = steps[rs.RelLeft]
				rs.relRight = steps[rs.RelRight]
				if rs.relLeft == nil {
					return fmt.Errorf("link: missing left relation %s", rs.relLeft)
				}
				if rs.relRight == nil {
					return fmt.Errorf("link: missing right relation %s", rs.relRight)
				}
			}
		}
		return
	}()
	if err != nil {
		return fmt.Errorf("unable to link steps: %w", err)
	}

	return nil
}

// Analyze runs analysis over each step in the pipeline
//
// Step analysis hints to the optimizers as to how expensive specific operations
// are and the general dataset size involved.
func (pp Pipeline) Analyze(ctx context.Context) (err error) {
	for _, p := range pp {
		err = p.Analyze(ctx)
		if err != nil {
			return
		}
	}
	return
}

// Optimize runs all optimization and returns an optimized pipeline
func (base Pipeline) Optimize(ctx context.Context) (optimized Pipeline, err error) {
	base, err = base.OptimizeStructure(ctx)
	if err != nil {
		return
	}
	return base.OptimizeSteps(ctx)
}

// OptimizeStructure performs general pipeline structure optimizations such as
// restructuring and clobbering steps onto the datasource layer
func (base Pipeline) OptimizeStructure(ctx context.Context) (optimized Pipeline, err error) {
	optimized = base.Clone()

	return optimized, optimized.walkSubtrees(optimized.root(), func(step PipelineStep, isRoot bool) (out PipelineStep, err error) {
		out = step
		for _, opt := range pipelineOptimizers {
			out, err = opt(out, isRoot)
			if err != nil {
				return
			}
		}
		return
	})
}

// OptimizeSteps performs step specific optimizations such as pushing filters
// on the lower levels, determining step-specific plans, ...
func (base Pipeline) OptimizeSteps(ctx context.Context) (optimized Pipeline, err error) {
	optimized = base.Clone()
	return optimized, optimized.optimizeSteps(base.root(), internalFilter{})
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities

func (pp Pipeline) root() PipelineStep {
	ix := make(map[string]PipelineStep)

	for _, s := range pp {
		ix[s.Identifier()] = s
	}

	for _, s := range pp {
		for _, src := range s.Sources() {
			delete(ix, src)
		}
	}

	for _, s := range ix {
		return ix[s.Identifier()]
	}

	panic("impossible state: no root-level pipeline step")
}

// Slice returns a new pipeline with all steps of the subtree with ident as root
func (pp Pipeline) Slice(ident string) (out Pipeline) {
	// Make a copy so we can assure the caller can go ham over the pipeline
	ppc := pp.Clone()

	// Find root
	var r PipelineStep
	for _, p := range ppc {
		if p.Identifier() == ident {
			r = p
			break
		}
	}

	return ppc.slice(r)
}

// slice is the recursive counterpart for the Slice method
func (pp Pipeline) slice(s PipelineStep) (out Pipeline) {
	out = append(out, s)

	switch n := s.(type) {
	case *Datasource:
		return

	case *Aggregate:
		return append(out, pp.slice(n.rel)...)

	case *Join:
		return append(out, append(pp.slice(n.relLeft), pp.slice(n.relRight)...)...)

	case *Link:
		return append(out, append(pp.slice(n.relLeft), pp.slice(n.relRight)...)...)
	}

	return
}

// optimizeSteps is the recursive counterpart to the .OptimizeSteps method
func (p Pipeline) optimizeSteps(node PipelineStep, inF internalFilter) (err error) {
	switch n := node.(type) {
	case *Datasource:
		inF, err = n.Optimize(inF)
		if err != nil {
			return
		}
		if !inF.empty() {
			return fmt.Errorf("a datasource can not offload optimizations")
		}
		return

	case *Aggregate:
		inF, err = n.Optimize(inF)
		if err != nil {
			return
		}
		return p.optimizeSteps(n.rel, inF)

	case *Join:
		inF, err = n.Optimize(inF)
		if err != nil {
			return
		}
		err = p.optimizeSteps(n.relLeft, inF)
		if err != nil {
			return
		}
		err = p.optimizeSteps(n.relRight, inF)
		if err != nil {
			return
		}

	case *Link:
		inF, err = n.Optimize(inF)
		if err != nil {
			return
		}
		err = p.optimizeSteps(n.relLeft, inF)
		if err != nil {
			return
		}
		err = p.optimizeSteps(n.relRight, inF)
		if err != nil {
			return
		}
	}

	return
}

// walkSubtrees performs a DFS and invokes fn for every sub-tree node in the returning order
func (p Pipeline) walkSubtrees(root PipelineStep, fn func(step PipelineStep, isRoot bool) (PipelineStep, error)) (err error) {
	err = p.walkSubtreesRec(root, true, fn)
	if err != nil {
		return
	}

	// dfsRec reports only subtrees; in case no sub tree was there, do it
	switch root.(type) {
	case *Join, *Link:
		n, err := fn(root, true)
		if err != nil {
			return err
		}
		p.replace(root, n)
	}

	return
}

// walkSubtreesRec is the recursive counterpart to the .walkSubtrees method
func (p Pipeline) walkSubtreesRec(root PipelineStep, isRoot bool, fn func(step PipelineStep, isRoot bool) (PipelineStep, error)) (err error) {
	var n PipelineStep

	switch s := root.(type) {
	case *Datasource:
		// this one doesn't have anything under it
		return

	case *Aggregate:
		return p.walkSubtreesRec(s.rel, false, fn)

	case *Join:
		err = p.walkSubtreesRec(s.relLeft, false, fn)
		if err != nil {
			return
		}
		err = p.walkSubtreesRec(s.relRight, false, fn)
		if err != nil {
			return
		}
		n, err = fn(root, isRoot)
		if err != nil {
			return
		}
		p.replace(root, n)

	case *Link:
		err = p.walkSubtreesRec(s.relLeft, false, fn)
		if err != nil {
			return
		}
		err = p.walkSubtreesRec(s.relRight, false, fn)
		if err != nil {
			return
		}
		n, err = fn(root, isRoot)
		if err != nil {
			return
		}
		p.replace(root, n)
	}

	return
}

func (pp Pipeline) replace(o, n PipelineStep) {
	for i, p := range pp {
		if p == o {
			pp[i] = n
			return
		}
	}
}

func (p Pipeline) Clone() (out Pipeline) {
	out = make(Pipeline, 0, len(p))
	for _, s := range p {
		switch s := s.(type) {
		case *Aggregate:
			aux := *s
			out = append(out, &aux)

		case *Join:
			aux := *s
			out = append(out, &aux)

		case *Link:
			aux := *s
			out = append(out, &aux)

		case *Datasource:
			aux := *s
			out = append(out, &aux)

		default:
			panic("unsupported step")
		}
	}
	return
}
