package yaml

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

type (
	yamlEncoder struct {
		cfg *EncoderConfig

		resState  map[resource.Interface]resourceState
		Documents encoderState
	}

	// encodedDocument holds the yaml encoded document sources.
	// Use the reader for further processing (file creation, HTTP, ...)
	encodedDocument struct {
		doc          *Document
		Source       io.Reader
		ResourceType string
		Scope        string
		Identifier   string
	}

	encodedDocumentSet []*encodedDocument

	// encoderState holds all of the generated, encoded documents indexed by resource type
	encoderState map[string]encodedDocumentSet

	// EncoderConfig allows us to configure the resource encoding process
	EncoderConfig struct {
		// Skip if defines a pkg/expr expression when to skip the resource
		SkipIf string
		// Defer is called after the resource is encoded, regardles of the result
		Defer func()
		// DeferOk is called after the resource is encoded, only when successful
		DeferOk func()
		// DeferNok is called after the resource is encoded, only when failed
		// If you return an error, the encoding will terminate.
		// If you return nil (ignore the error), the encoding will continue.
		DeferNok func(error) error

		// Timezone defines what timezone should be used when encoding timestamps
		//
		// If not defined, UTC is used
		Timezone string
		// TimeLayout defines how to format the encoded timestamp
		//
		// If not defined, RFC3339 is used (this one - 2006-01-02T15:04:05Z07:00)
		TimeLayout string

		// // CompactOutput forces the output to be as compact as possible
		// CompactOutput bool
		// MappedOutput forces the sequences to encode as maps (where possible)
		MappedOutput bool

		// @todo different output structuring
	}

	// resourceState holds some intermedia values to help with encoding
	resourceState interface {
		Prepare(ctx context.Context, state *envoy.ResourceState) (err error)
		Encode(ctx context.Context, ds *Document, state *envoy.ResourceState) (err error)
	}
)

var (
	ErrUnknownResource        = errors.New("unknown resource")
	ErrResourceStateUndefined = errors.New("undefined resource state")
	ErrInvalidResourceType    = errors.New("invalid resource state")
)

// NewYamlEncoder initializes a fresh yaml encoder
func NewYamlEncoder(cfg *EncoderConfig) envoy.PrepareEncodeStreammer {
	if cfg == nil {
		cfg = &EncoderConfig{}
	}

	return &yamlEncoder{
		cfg: cfg,

		Documents: make(encoderState),
		resState:  make(map[resource.Interface]resourceState),
	}
}

// Prepare prepares the encoder for the given set of resources
//
// It initializes and prepares the resource state for each provided resource
func (ye *yamlEncoder) Prepare(ctx context.Context, ee ...*envoy.ResourceState) (err error) {
	f := func(rs resourceState, es *envoy.ResourceState) error {
		err = rs.Prepare(ctx, es)
		if err != nil {
			return err
		}

		ye.resState[es.Res] = rs
		return nil
	}

	for _, e := range ee {
		switch res := e.Res.(type) {
		// Compose resources
		case *resource.ComposeNamespace:
			err = f(composeNamespaceFromResource(res, ye.cfg), e)
		case *resource.ComposeModule:
			err = f(composeModuleFromResource(res, ye.cfg), e)
		case *resource.ComposeRecord:
			err = f(composeRecordSetFromResource(res, ye.cfg), e)
		case *resource.ComposePage:
			err = f(composePageFromResource(res, ye.cfg), e)
		case *resource.ComposeChart:
			err = f(composeChartFromResource(res, ye.cfg), e)

		// System resources
		case *resource.Role:
			err = f(roleFromResource(res, ye.cfg), e)
		case *resource.User:
			err = f(userFromResource(res, ye.cfg), e)
		case *resource.Template:
			err = f(templateFromResource(res, ye.cfg), e)
		case *resource.Application:
			err = f(applicationFromResource(res, ye.cfg), e)
		case *resource.Setting:
			err = f(settingFromResource(res, ye.cfg), e)
		case *resource.RbacRule:
			err = f(rbacRuleFromResource(res, ye.cfg), e)

		// Automation resources
		case *resource.AutomationWorkflow:
			err = f(automationWorkflowFromResource(res, ye.cfg), e)

		default:
			err = ErrUnknownResource
		}

		if err != nil {
			return ye.WrapError("prepare", e.Res, err)
		}
	}

	return nil
}

// Encode encodes the resources into a series of yaml documents
//
// The document structure follows the provided configuration as much as possible,
// but some modifications may be permitted (such as splitting compose resources that belong
// to different namespaces).
//
// @todo improve document structuring; the base encodes each resource type into it's own document.
//       This is good enough for now but should be expanded in the near future.
func (ye *yamlEncoder) Encode(ctx context.Context, p envoy.Provider) error {
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

		state := ye.resState[e.Res]
		if state == nil {
			err = ErrResourceStateUndefined
		} else {
			// Determine the document that we should encode into
			//
			// @todo improve flexibility
			c := ye.Documents.byResourceType(e.Res.ResourceType())
			if c.doc.cfg == nil {
				c.doc.cfg = ye.cfg
			}

			err = state.Encode(ctx, c.doc, e)
		}

		if err != nil {
			return ye.WrapError("encode: build doc", e.Res, err)
		}
	}

	// Encode documents into yaml document streams
	for _, d := range ye.Documents.all() {
		if err := d.encode(); err != nil {
			return err
		}

		// The source Document doesn't need to exist after this step
		d.doc = nil
	}

	return nil
}

func (se *yamlEncoder) Stream() []*envoy.Stream {
	ss := make([]*envoy.Stream, 0, 20)

	for _, dd := range se.Documents {
		for _, d := range dd {
			ss = append(ss, &envoy.Stream{
				Resource:   d.ResourceType,
				Identifier: d.Identifier,
				Source:     d.Source,
			})
		}
	}

	return ss
}

// WrapError wraps errors related to yaml encoding
//
// Always wrap your errors.
func (se *yamlEncoder) WrapError(act string, res resource.Interface, err error) error {
	rt := strings.Join(strings.Split(strings.TrimSpace(strings.TrimRight(res.ResourceType(), ":")), ":"), " ")
	return fmt.Errorf("yaml encoder %s %s %v: %s", act, rt, res.Identifiers().StringSlice(), err)
}

func encoderErrInvalidResource(exp, got string) error {
	return fmt.Errorf("invalid resource type: expecting %s, got %s", exp, got)
}

// byResourceType returns the first document used for a specific resource type
func (ec encoderState) byResourceType(rt string) *encodedDocument {
	if _, has := ec[rt]; !has {
		ec[rt] = make(encodedDocumentSet, 0, 2)
	}

	if ec[rt] == nil || len(ec[rt]) == 0 {
		e := &encodedDocument{
			doc:          &Document{},
			ResourceType: rt,
		}
		ec[rt] = append(ec[rt], e)

		return e
	}

	return (ec[rt])[0]
}

// all returns all of the encodedDocument structs stored in the state.
//
// Use this when encoding everything.
// You can use specific state fields when encoding specific resources only.
func (ec encoderState) all() encodedDocumentSet {
	dd := make(encodedDocumentSet, 0, 100)

	for _, ss := range ec {
		dd = append(dd, ss...)
	}

	return dd
}

// Little helper to encode the raw Document into a io.Reader
func (d *encodedDocument) encode() error {
	b := bytes.Buffer{}

	err := yaml.NewEncoder(&b).Encode(d.doc)
	if err != nil {
		return err
	}

	d.Source = &b
	return nil
}
