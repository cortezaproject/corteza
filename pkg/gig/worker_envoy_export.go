package gig

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
)

func (w *workerEnvoy) prepareExport(ctx context.Context, sources ...Source) (err error) {
	return w.parseSources(ctx, sources...)
}

func (w *workerEnvoy) execExport(ctx context.Context) (output SourceSet, meta WorkMeta, err error) {
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
