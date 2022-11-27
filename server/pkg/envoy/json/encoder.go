package json

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
)

type (
	bulkRecordEncoder struct {
		cfg *EncoderConfig

		resState map[resource.Interface]*encoderState
	}

	encoderState struct {
		res          resourceState
		source       io.ReadWriter
		resourceType string
		Scope        string
		identifier   string
	}

	// EncoderConfig allows us to configure the resource encoding process
	EncoderConfig struct {
		// Timezone defines what timezone should be used when encoding timestamps
		//
		// If not defined, UTC is used
		Timezone string
		// TimeLayout defines how to format the encoded timestamp
		//
		// If not defined, RFC3339 is used (this one - 2006-01-02T15:04:05Z07:00)
		TimeLayout string
		// Fields specifies what fields we wish to include in the export
		Fields map[string]bool
	}

	// resourceState holds some intermedia values to help with encoding
	resourceState interface {
		Prepare(ctx context.Context, state *envoy.ResourceState) (err error)
		Encode(ctx context.Context, w io.Writer, state *envoy.ResourceState) (err error)
	}
)

var (
	ErrUnknownResource        = fmt.Errorf("unknown resource")
	ErrResourceStateUndefined = fmt.Errorf("undefined resource state")
)

func NewBulkRecordEncoder(cfg *EncoderConfig) envoy.PrepareEncodeStreamer {
	if cfg == nil {
		cfg = &EncoderConfig{}
	}

	return &bulkRecordEncoder{
		cfg: cfg,

		resState: make(map[resource.Interface]*encoderState),
	}
}

// Prepare prepares the encoder for the given set of resources
//
// It initializes and prepares the resource state for each provided resource
func (se *bulkRecordEncoder) Prepare(ctx context.Context, ee ...*envoy.ResourceState) (err error) {
	f := func(rs resourceState, es *envoy.ResourceState) error {
		err = rs.Prepare(ctx, es)
		if err != nil {
			return err
		}

		se.resState[es.Res] = &encoderState{
			res:          rs,
			source:       &bytes.Buffer{},
			resourceType: es.Res.ResourceType(),
			identifier:   es.Res.Identifiers().First(),
		}
		return nil
	}

	for _, e := range ee {
		switch res := e.Res.(type) {
		// @todo other resources; we'll only do records for now
		case *resource.ComposeRecord:
			err = f(bulkComposeRecordEncoderFromResource(res, se.cfg), e)

		default:
			err = ErrUnknownResource
		}

		if err != nil {
			return se.WrapError("prepare", e.Res, err)
		}
	}

	return nil
}

// Encode encodes the resources into a series of jsonl documents
func (se *bulkRecordEncoder) Encode(ctx context.Context, p envoy.Provider) error {
	var e *envoy.ResourceState
	var err error

	// Encode the resources into document structs
	for {
		e, err = p.NextInverted(ctx)

		if err != nil {
			return err
		}
		if e == nil {
			break
		}

		state := se.resState[e.Res]
		if state == nil {
			err = ErrResourceStateUndefined
		} else {
			err = state.res.Encode(ctx, state.source, e)
		}

		if err != nil {
			return se.WrapError("encode: build doc", e.Res, err)
		}
	}

	for _, s := range se.resState {
		s.res = nil
	}

	return nil
}

func (se *bulkRecordEncoder) Stream() []*envoy.Stream {
	ss := make([]*envoy.Stream, 0, 20)

	for _, s := range se.resState {
		ss = append(ss, &envoy.Stream{
			Resource:   s.resourceType,
			Identifier: s.identifier,
			Source:     s.source,
		})
	}

	return ss
}

// WrapError wraps errors related to json encoding
//
// Always wrap your errors.
func (se *bulkRecordEncoder) WrapError(act string, res resource.Interface, err error) error {
	rt := strings.Join(strings.Split(strings.TrimSpace(strings.TrimRight(res.ResourceType(), ":")), ":"), " ")
	return fmt.Errorf("json encoder %s %s %v: %s", act, rt, res.Identifiers().StringSlice(), err)
}

func encoderErrInvalidResource(exp, got string) error {
	return fmt.Errorf("invalid resource type: expecting %s, got %s", exp, got)
}
