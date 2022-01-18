package gig

import (
	"context"
	"io"

	automationTypes "github.com/cortezaproject/corteza-server/automation/types"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	WorkerStateEnvoy struct {
		Resources []envoyResourceWrap
	}
	envoyResourceWrap struct {
		ResourceType string
		Identifier   string
		Identifiers  []string
		Raw          interface{}
	}

	storeDecoder interface {
		Decode(context.Context, store.Storer, *es.DecodeFilter) ([]resource.Interface, error)
	}

	sourceDecoder interface {
		CanDecodeExt(string) bool
		CanDecodeFile(io.Reader) bool
		Decode(context.Context, io.Reader, *envoy.DecoderOpts) ([]resource.Interface, error)
	}
)

var (
	ComposeNamespaceResourceType   = composeTypes.NamespaceResourceType
	AutomationWorkflowResourceType = automationTypes.WorkflowResourceType
)

func parseSources(ctx context.Context, sources ...Source) (resource.InterfaceSet, error) {
	decoders := getSourceDecoders()
	out := make(resource.InterfaceSet, 0, 16)

	for _, src := range sources {
		for _, d := range decoders {
			r, err := src.Read()
			if err != nil {
				return nil, err
			}
			if d.CanDecodeFile(r) {
				tmp, err := d.Decode(ctx, src.ReadSafe(), &envoy.DecoderOpts{})
				if err != nil {
					return nil, err
				}
				out = append(out, tmp...)
			}
		}
	}

	return out, nil
}
