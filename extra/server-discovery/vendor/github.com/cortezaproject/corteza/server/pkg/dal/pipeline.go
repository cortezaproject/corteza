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
		Analysis() map[string]OpAnalysis
	}

	// clobberableStep can have other steps offload their work to it.
	// This is primarily used to offload the work to the database.
	clobberableStep interface {
		clobber(PipelineStep) bool
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
		Label               string
		IsPrimary           bool
		IsSystem            bool
		IsFilterable        bool
		IsSortable          bool
		IsMultivalue        bool
		MultivalueDelimiter string
		Nullable            bool
		Type                Type
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
	// Find root
	var r PipelineStep
	for _, p := range pp {
		if p.Identifier() == ident {
			r = p
			break
		}
	}

	return pp.slice(r)
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
