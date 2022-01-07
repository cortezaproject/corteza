package gig

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type (
	Source interface {
		// Base stuff
		ID() uint64
		FileName() string
		Name() string
		MimeType() string
		Ext() string
		Size() int64
		Checksum() string

		Read() (io.Reader, error)
		ReadSafe() io.Reader
		Cleanup() error

		SetDecoders(...Decoder)
		HasDecoders() bool
		Decoders() []Decoder

		// Remaining fileinfo interface fncs
		Mode() os.FileMode
		ModTime() time.Time
		IsDir() bool
		Sys() interface{}
	}
	SourceSet []Source

	SourceWrap struct {
		ID  uint64
		Src io.Reader
		Uri string

		Name     string
		Mime     string
		Size     int64
		Checksum string
		IsDir    bool
	}
	SourceWrapSet []SourceWrap
)

func FromSourceWrap(ctx context.Context, sources ...SourceWrap) (out SourceSet, err error) {
	for _, wrap := range sources {
		var src Source
		var sources SourceSet
		switch {
		case wrap.Src != nil:
			src, err = FileSourceFromBlob(ctx, wrap.Name, wrap.Src)
		case wrap.Uri != "":
			if wrap.IsDir {
				sources, err = PrepareSourceFromDirectory(ctx, wrap.Uri)
			} else {
				src, err = FileSourceFromURI(ctx, wrap.Uri)
			}
		}

		if err != nil {
			return
		}

		if src != nil {
			out = append(out, src)
		}
		out = append(out, sources...)
	}

	return
}

func ToSourceWrap(sources ...Source) (out SourceWrapSet) {
	for _, src := range sources {
		wrap := SourceWrap{
			ID:       src.ID(),
			Src:      src.ReadSafe(),
			Name:     src.FileName(),
			Mime:     src.MimeType(),
			Size:     src.Size(),
			Checksum: src.Checksum(),
		}
		out = append(out, wrap)
	}
	return
}

func (set SourceWrapSet) Has(id uint64) bool {
	for _, s := range set {
		if s.ID == id {
			return true
		}
	}
	return false
}

func (set SourceSet) GetByID(id uint64) (out *Source) {
	for _, s := range set {
		if s.ID() == id {
			return &s
		}
	}

	return
}

func mapSources(sources SourceSet) map[uint64]Source {
	out := make(map[uint64]Source)
	for _, src := range sources {
		out[src.ID()] = src
	}

	return out
}

func runDecoders(ctx context.Context, sources SourceSet) (out SourceSet, err error) {
	if len(sources) == 0 {
		return
	}

	var tmp []Source
	for _, src := range sources {
		if src.HasDecoders() {
			for _, d := range src.Decoders() {
				if !d.CanDecode(src) {
					err = fmt.Errorf("unable to use decoder with source: can't decode @todo error message")
					return
				}
				tmp, err = d.Decode(ctx, src)
				if err != nil {
					return
				}
				out = append(out, tmp...)
			}

		} else {
			out = append(out, src)
		}
	}

	return
}
