package gig

import (
	"context"

	automationTypes "github.com/cortezaproject/corteza-server/automation/types"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
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

	envoyDirection int

	workerEnvoy struct {
		direction envoyDirection
		store     store.Storer

		resources resource.InterfaceSet
		filter    *es.DecodeFilter
	}

	storeDecoder interface {
		Decode(context.Context, store.Storer, *es.DecodeFilter) ([]resource.Interface, error)
	}
)

var (
	WorkerHandleEnvoy = "envoy"

	envoyDirectionIn  envoyDirection = 0
	envoyDirectionOut envoyDirection = 1

	ApigwRouteResourceType  = systemTypes.ApigwRouteResourceType
	ApplicationResourceType = systemTypes.ApplicationResourceType
	ReportResourceType      = systemTypes.ReportResourceType
	RoleResourceType        = systemTypes.RoleResourceType
	TemplateResourceType    = systemTypes.TemplateResourceType
	UserResourceType        = systemTypes.UserResourceType

	ComposeChartResourceType       = composeTypes.ChartResourceType
	ComposeModuleResourceType      = composeTypes.ModuleResourceType
	ComposeModuleFieldResourceType = composeTypes.ModuleFieldResourceType
	ComposeNamespaceResourceType   = composeTypes.NamespaceResourceType
	ComposePageResourceType        = composeTypes.PageResourceType

	AutomationWorkflowResourceType = automationTypes.WorkflowResourceType

	SettingsResourceType            = resource.SettingsResourceType
	RbacResourceType                = resource.RbacResourceType
	ResourceTranslationResourceType = resource.ResourceTranslationType
)

func WorkerEnvoy(s store.Storer) *workerEnvoy {
	return &workerEnvoy{
		store: s,
	}
}

func WorkerImport(s store.Storer) *workerEnvoy {
	return &workerEnvoy{
		direction: envoyDirectionIn,
		store:     s,
	}
}

func WorkerExport(s store.Storer) *workerEnvoy {
	return &workerEnvoy{
		store:     s,
		direction: envoyDirectionOut,
	}
}

func (w *workerEnvoy) MarkImport() {
	w.direction = envoyDirectionIn
}

func (w *workerEnvoy) MarkExport() {
	w.direction = envoyDirectionOut
}

func (w *workerEnvoy) Ref() string {
	return WorkerHandleEnvoy
}

func (w *workerEnvoy) Prepare(ctx context.Context, sources ...Source) error {
	if w.direction == envoyDirectionIn {
		return w.prepareImport(ctx, sources...)
	}
	return w.prepareExport(ctx, sources...)
}

func (w *workerEnvoy) Exec(ctx context.Context) (output SourceSet, meta WorkMeta, err error) {
	if len(w.resources) == 0 {
		return
	}

	if w.direction == envoyDirectionIn {
		return w.execImport(ctx)
	}
	return w.execExport(ctx)
}

// @todo ...
func (w *workerEnvoy) collectMeta() (meta WorkMeta) {
	meta = make(WorkMeta)

	return
}

func (w *workerEnvoy) State(context.Context) (WorkerState, error) {
	out := WorkerStateEnvoy{
		Resources: make([]envoyResourceWrap, len(w.resources)),
	}
	for i, r := range w.resources {
		out.Resources[i].ResourceType = r.ResourceType()
		out.Resources[i].Identifier = r.Identifiers().First()
		out.Resources[i].Identifiers = r.Identifiers().StringSlice()
		out.Resources[i].Raw = r.Resource()
	}

	return out, nil
}

func (w *workerEnvoy) Cleanup(context.Context) error {
	w.resources = nil
	return nil
}

func (w *workerEnvoy) getSourceDecoders() []sourceDecoder {
	return []sourceDecoder{
		yaml.Decoder(),
		csv.Decoder(),
	}
}

func (w *workerEnvoy) getStoreDecoders() storeDecoder {
	return es.Decoder()
}

func (w *workerEnvoy) getStoreEncoder() envoy.PrepareEncoder {
	return es.NewStoreEncoder(w.store, &es.EncoderConfig{})
}

func (w *workerEnvoy) getYamlEncoder() envoy.PrepareEncodeStreamer {
	return yaml.NewYamlEncoder(&yaml.EncoderConfig{
		MappedOutput: false,
	})
}

func (d envoyDirection) String() string {
	if d == envoyDirectionOut {
		return "export"
	}
	return "import"
}
