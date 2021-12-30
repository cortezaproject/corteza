package gig

import (
	"context"
	"io"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	sourceDecoder interface {
		CanDecodeExt(string) bool
		CanDecodeFile(io.Reader) bool
		Decode(context.Context, io.Reader, *envoy.DecoderOpts) ([]resource.Interface, error)
	}
)

func (w *workerEnvoy) prepareImport(ctx context.Context, sources ...Source) error {
	return w.parseSources(ctx, sources...)
}

func (w *workerEnvoy) execImport(ctx context.Context) (output SourceSet, meta WorkMeta, err error) {
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

func (w *workerEnvoy) parseSources(ctx context.Context, sources ...Source) error {
	decoders := w.getSourceDecoders()

	for _, src := range sources {
		for _, d := range decoders {
			r, err := src.Read()
			if err != nil {
				return err
			}
			if d.CanDecodeFile(r) {
				tmp, err := d.Decode(ctx, src.ReadSafe(), &envoy.DecoderOpts{})
				if err != nil {
					return err
				}
				w.resources = append(w.resources, tmp...)
			}
		}
	}

	return nil
}
