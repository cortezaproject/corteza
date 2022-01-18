package gig

import (
	"time"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	completion int

	Gig struct {
		ID         uint64 `json:"gigID,string"`
		Signature  string
		CompleteOn completion

		Sources SourceSet
		Worker  Worker

		// Tasks
		Preprocess  []Preprocessor
		Postprocess []Postprocessor

		// Worker stuff
		Output SourceSet
		Status WorkerStatus

		// Timestamps
		CreatedAt time.Time
		UpdatedAt *time.Time
		DeletedAt *time.Time

		CompletedAt *time.Time
		PreparedAt  *time.Time
	}
	GigSet []*Gig

	UpdatePayload struct {
		Worker      Worker
		Decode      DecoderSet
		Preprocess  PreprocessorSet
		Postprocess PostprocessorSet
		Sources     []SourceWrap
		CompleteOn  completion
	}
)

const (
	OnDemand completion = iota
	OnExec
	OnOutput
)

func newGig(w Worker) Gig {
	return Gig{
		ID:        nextID(),
		Worker:    w,
		CreatedAt: *now(),
	}
}

func (g Gig) TySystemWrapper() *types.Gig {
	return &types.Gig{
		ID: g.ID,
	}
}

func (c completion) String() string {
	switch c {
	case OnExec:
		return "OnExec"
	case OnOutput:
		return "OnOutput"
	default:
		return "OnDemand"
	}
}
