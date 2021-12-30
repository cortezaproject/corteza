package gig

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type (
	Postprocessor interface {
		Postprocess(context.Context, WorkMeta, SourceSet) (SourceSet, WorkMeta, error)
	}

	PostprocessorWrap struct {
		Ref    postprocessor   `json:"ref"`
		Params json.RawMessage `json:"params"`
	}
	PostprocessorWrapSet []PostprocessorWrap

	savedSource struct {
		Name string
		Size int64
		URI  string
	}

	postprocessor string
)

func PostprocessorNoop() (Postprocessor, error) {
	return postprocessorNoop{}, nil
}

func (t postprocessorNoop) Postprocess(ctx context.Context, baseMeta WorkMeta, ss SourceSet) (out SourceSet, meta WorkMeta, err error) {
	out = ss
	meta = baseMeta
	return
}

func PostprocessorDiscard() (Postprocessor, error) {
	return postprocessorDiscard{}, nil
}

func (t postprocessorDiscard) Postprocess(ctx context.Context, baseMeta WorkMeta, ss SourceSet) (out SourceSet, meta WorkMeta, err error) {
	meta = baseMeta
	return
}

func PostprocessorSave() (Postprocessor, error) {
	return nil, fmt.Errorf("postprocessor not yet defined: %s", PostprocessorHandleSave)
	// return postprocessorSave{}, nil
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

func PostprocessorArchive(format archive, name string) (out postprocessorArchive, err error) {
	out.Encoding = format
	out.Name = name

	if out.Name == "" {
		out.Name = "archive"
	}

	return
}

func (t postprocessorArchive) Postprocess(ctx context.Context, baseMeta WorkMeta, ss SourceSet) (out SourceSet, meta WorkMeta, err error) {
	p, name, err := compressTarGz(ctx, ss, t.Name)
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

func ParsePostprocessorWrap(ss []string) (out PostprocessorWrapSet, err error) {
	for _, s := range ss {
		aux := make(PostprocessorWrapSet, 0, 2)
		err = json.Unmarshal([]byte(s), &aux)
		if err != nil {
			return
		}

		out = append(out, aux...)
	}
	return
}
