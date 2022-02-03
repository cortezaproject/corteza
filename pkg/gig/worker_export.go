package gig

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	workerExport struct {
		store     store.Storer
		resources resource.InterfaceSet
	}
)

func WorkerExport(s store.Storer) Worker {
	return &workerExport{
		store: s,
	}
}

func (w *workerExport) Ref() string {
	return WorkerHandleExport
}

func (w *workerExport) Prepare(ctx context.Context, sources ...Source) error {
	return w.prepareExport(ctx, sources...)
}

func (w *workerExport) Exec(ctx context.Context) (output SourceSet, meta WorkMeta, err error) {
	if len(w.resources) == 0 {
		return
	}

	return w.execExport(ctx)
}

// @todo ...
func (w *workerExport) collectMeta() (meta WorkMeta) {
	meta = make(WorkMeta)

	return
}

func (w *workerExport) State(context.Context) (WorkerState, error) {
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

func (w *workerExport) Cleanup(context.Context) error {
	w.resources = nil
	return nil
}

func (w *workerExport) getYamlEncoder() envoy.PrepareEncodeStreamer {
	return yaml.NewYamlEncoder(&yaml.EncoderConfig{
		MappedOutput: false,
	})
}

func (w *workerExport) prepareExport(ctx context.Context, sources ...Source) (err error) {
	res, err := parseSources(ctx, sources...)
	if err != nil {
		return err
	}

	w.resources = append(w.resources, res...)
	return nil
}

func (w *workerExport) execExport(ctx context.Context) (output SourceSet, meta WorkMeta, err error) {
	if len(w.resources) == 0 {
		return
	}

	enc := w.getYamlEncoder()

	bld := envoy.NewBuilder(enc)
	g, err := bld.Build(ctx, w.resources...)
	if err != nil {
		return
	}

	err = envoy.Encode(ctx, g, enc)
	if err != nil {
		return
	}

	// create sources
	var src Source
	for _, s := range enc.Stream() {
		src, err = FileSourceFromBlob(ctx, fmt.Sprintf("%s.yaml", s.Resource), s.Source)
		if err != nil {
			return
		}

		output = append(output, src)
	}

	meta = w.collectMeta()
	return
}
