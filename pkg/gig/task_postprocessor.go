package gig

import (
	"context"
	"fmt"
	"os"
)

type (
	Postprocessor interface {
		Postprocess(context.Context, WorkMeta, SourceSet) (SourceSet, WorkMeta, error)
		Ref() string
		Params() map[string]interface{}
	}
	PostprocessorSet []Postprocessor

	savedSource struct {
		Name string
		Size int64
		URI  string
	}
)

func PostprocessorNoop() postprocessorNoop {
	return postprocessorNoop{}
}

func (t postprocessorNoop) Postprocess(ctx context.Context, baseMeta WorkMeta, ss SourceSet) (out SourceSet, meta WorkMeta, err error) {
	out = ss
	meta = baseMeta
	return
}

func PostprocessorDiscard() postprocessorDiscard {
	return postprocessorDiscard{}
}

func (t postprocessorDiscard) Postprocess(ctx context.Context, baseMeta WorkMeta, ss SourceSet) (out SourceSet, meta WorkMeta, err error) {
	meta = baseMeta
	return
}

func PostprocessorSave() postprocessorSave {
	return postprocessorSave{}
}

func (t postprocessorSave) Postprocess(ctx context.Context, baseMeta WorkMeta, ss SourceSet) (out SourceSet, meta WorkMeta, err error) {
	// @todo
	saved := make([]savedSource, 0, len(ss))
	for _, s := range ss {
		saved = append(saved, savedSource{
			Name: s.Name(),
			Size: s.Size(),
			URI:  fmt.Sprintf("https://domain.tld/%s", s.Name()),
		})
	}

	out = ss
	meta = baseMeta
	meta["saved"] = saved

	return
}

func postprocessorArchiveTransformer(base postprocessorArchive) postprocessorArchive {
	if base.name == "" {
		base.name = "archive"
	}

	return base
}

func (t postprocessorArchive) Postprocess(ctx context.Context, baseMeta WorkMeta, ss SourceSet) (out SourceSet, meta WorkMeta, err error) {
	p, name, err := compressTarGz(ctx, ss, t.name)
	if err != nil {
		return
	}

	f, err := os.Open(p)
	if err != nil {
		return
	}

	src, err := FileSourceFromBlob(ctx, name, f)
	if err != nil {
		return
	}

	out = SourceSet{src}
	return
}
