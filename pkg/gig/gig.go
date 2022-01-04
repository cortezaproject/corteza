package gig

import "time"

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
	}
)

const (
	onDemand completion = iota
	onExec
	onOutput
)

func newGig(w Worker) Gig {
	return Gig{
		ID:        nextID(),
		Worker:    w,
		CreatedAt: *now(),
	}
}

func (set GigSet) FindByID(id uint64) *Gig {
	for _, s := range set {
		if s.ID == id {
			return s
		}
	}

	return nil
}

func (c completion) String() string {
	switch c {
	case onExec:
		return "onExec"
	case onOutput:
		return "onOutput"
	default:
		return "onDemand"
	}
}
