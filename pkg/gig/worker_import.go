package gig

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	envoyStore "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	workerImport struct {
		store     store.Storer
		resources resource.InterfaceSet
	}
)

func WorkerImport(s store.Storer) Worker {
	return &workerImport{
		store: s,
	}
}

func (w *workerImport) Ref() string {
	return WorkerHandleImport
}

func (w *workerImport) Prepare(ctx context.Context, sources ...Source) error {
	return w.prepareImport(ctx, sources...)
}

func (w *workerImport) Exec(ctx context.Context) (output SourceSet, meta WorkMeta, err error) {
	if len(w.resources) == 0 {
		return
	}

	return w.execImport(ctx)
}

// @todo ...
func (w *workerImport) collectMeta() (meta WorkMeta) {
	meta = make(WorkMeta)

	return
}

func (w *workerImport) State(context.Context) (WorkerState, error) {
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

func (w *workerImport) Cleanup(context.Context) error {
	w.resources = nil
	return nil
}

func getSourceDecoders() []sourceDecoder {
	return []sourceDecoder{
		yaml.Decoder(),
		csv.Decoder(),
	}
}

func getStoreDecoders() storeDecoder {
	return envoyStore.Decoder()
}

func (w *workerImport) getStoreEncoder() envoy.PrepareEncoder {
	return envoyStore.NewStoreEncoder(w.store, &envoyStore.EncoderConfig{})
}

func (w *workerImport) prepareImport(ctx context.Context, sources ...Source) error {
	res, err := parseSources(ctx, sources...)
	if err != nil {
		return err
	}

	w.resources = append(w.resources, res...)
	return nil
}

func (w *workerImport) execImport(ctx context.Context) (output SourceSet, meta WorkMeta, err error) {
	if len(w.resources) == 0 {
		return
	}

	enc := w.getStoreEncoder()

	bld := envoy.NewBuilder(enc)
	g, err := bld.Build(ctx, w.resources...)
	if err != nil {
		return
	}

	err = envoy.Encode(ctx, g, enc)
	if err != nil {
		return
	}

	meta = w.collectMeta()

	return
}
