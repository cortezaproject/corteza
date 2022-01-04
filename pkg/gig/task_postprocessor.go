package gig

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cast"
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

func PostprocessorNoopParams(_ map[string]interface{}) postprocessorNoop {
	return PostprocessorNoop()
}

func PostprocessorNoop() postprocessorNoop {
	return postprocessorNoop{}
}

func (t postprocessorNoop) Postprocess(ctx context.Context, baseMeta WorkMeta, ss SourceSet) (out SourceSet, meta WorkMeta, err error) {
	out = ss
	meta = baseMeta
	return
}

func (d postprocessorNoop) Ref() string {
	return PostprocessorHandleNoop
}

func (d postprocessorNoop) Params() map[string]interface{} {
	return nil
}

func PostprocessorDiscardParams(_ map[string]interface{}) postprocessorDiscard {
	return PostprocessorDiscard()
}

func PostprocessorDiscard() postprocessorDiscard {
	return postprocessorDiscard{}
}

func (t postprocessorDiscard) Postprocess(ctx context.Context, baseMeta WorkMeta, ss SourceSet) (out SourceSet, meta WorkMeta, err error) {
	meta = baseMeta
	return
}

func (d postprocessorDiscard) Ref() string {
	return PostprocessorHandleDiscard
}

func (d postprocessorDiscard) Params() map[string]interface{} {
	return nil
}

func PostprocessorSaveParams(_ map[string]interface{}) postprocessorSave {
	return PostprocessorSave()
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

func (d postprocessorSave) Ref() string {
	return PostprocessorHandleSave
}

func (d postprocessorSave) Params() map[string]interface{} {
	return nil
}

func PostprocessorArchiveParams(params map[string]interface{}) (out postprocessorArchive) {
	name := cast.ToString(params["name"])
	format := ArchiveZIP

	if f := cast.ToString(params["format"]); f != "" {
		// ignire error; default to zip
		format, _ = archiveFromString(f)
	}

	return PostprocessorArchive(name, format)
}

func PostprocessorArchive(name string, format archive) (out postprocessorArchive) {
	out.name = name
	out.encoding = format

	if out.name == "" {
		out.name = "archive"
	}

	return
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

func (d postprocessorArchive) Ref() string {
	return PostprocessorHandleArchive
}

func (d postprocessorArchive) Params() map[string]interface{} {
	return map[string]interface{}{
		"format": d.encoding.String(),
		"name":   d.name,
	}
}
