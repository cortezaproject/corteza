package rest

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/gig"
	"github.com/cortezaproject/corteza-server/system/rest/conv"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
)

type (
	Gig struct {
		svc  service.GigService
		conv conv.Gig
	}

	gigSourcePayload struct {
		ID       uint64            `json:"sourceID,string"`
		Name     string            `json:"name"`
		Size     int64             `json:"size"`
		MimeType string            `json:"mimeType"`
		Checksum string            `json:"checksum"`
		Decoders conv.ParamWrapSet `json:"decoders"`
	}
	gigWorkerPayload struct {
		Ref string `json:"ref"`
	}

	gigPayload struct {
		ID        uint64             `json:"gigID,string"`
		Signature string             `json:"signature"`
		Sources   []gigSourcePayload `json:"sources"`
		Worker    gigWorkerPayload   `json:"worker"`

		Preprocess  conv.ParamWrapSet `json:"preprocess"`
		Postprocess conv.ParamWrapSet `json:"postprocess"`

		State gig.WorkerState `json:"state,omitempty"`
	}

	gigTaskPayload struct {
		Set gig.TaskDefSet `json:"set"`
	}
)

func (Gig) New() *Gig {
	return &Gig{
		svc:  service.DefaultGig,
		conv: conv.Gig{},
	}
}

func (ctrl Gig) Create(ctx context.Context, r *request.GigCreate) (interface{}, error) {
	pre, err := ctrl.conv.UnwrapPreprocessorSet(r.Preprocessors)
	if err != nil {
		return nil, err
	}
	post, err := ctrl.conv.UnwrapPostprocessorSet(r.Postprocessors)
	if err != nil {
		return nil, err
	}

	g, err := ctrl.svc.Create(ctx, r.Worker, gig.UpdatePayload{
		Preprocess:  pre,
		Postprocess: post,
	})
	return ctrl.makeGigPayload(ctx, g, err)
}

func (ctrl Gig) Read(ctx context.Context, r *request.GigRead) (interface{}, error) {
	g, err := ctrl.svc.Read(ctx, r.GigID)
	return ctrl.makeGigPayload(ctx, g, err)
}

func (ctrl Gig) Update(ctx context.Context, r *request.GigUpdate) (interface{}, error) {
	decode, err := ctrl.conv.UnwrapDecoderSet(r.Decoders)
	if err != nil {
		return nil, err
	}
	pre, err := ctrl.conv.UnwrapPreprocessorSet(r.Preprocessors)
	if err != nil {
		return nil, err
	}
	post, err := ctrl.conv.UnwrapPostprocessorSet(r.Postprocessors)
	if err != nil {
		return nil, err
	}

	g, err := ctrl.svc.Update(ctx, r.GigID, gig.UpdatePayload{
		Decode:      decode,
		Preprocess:  pre,
		Postprocess: post,
	})
	return ctrl.makeGigPayload(ctx, g, err)
}

func (ctrl Gig) Delete(ctx context.Context, r *request.GigDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.Delete(ctx, r.GigID)
}

func (ctrl Gig) Undelete(ctx context.Context, r *request.GigUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.Undelete(ctx, r.GigID)
}

func (ctrl Gig) AddSource(ctx context.Context, r *request.GigAddSource) (interface{}, error) {
	decode, err := ctrl.conv.UnwrapDecoderSet(r.Decoders)
	if err != nil {
		return nil, err
	}
	src, err := ctrl.prepareSources(r.Upload, r.Uri)
	if err != nil {
		return nil, err
	}

	g, err := ctrl.svc.AddSources(ctx, r.GigID, gig.UpdatePayload{
		Decode:  decode,
		Sources: src,
	})
	return ctrl.makeGigPayload(ctx, g, err)
}

func (ctrl Gig) RemoveSource(ctx context.Context, r *request.GigRemoveSource) (interface{}, error) {
	g, err := ctrl.svc.RemoveSources(ctx, r.GigID, r.SourceID)
	return ctrl.makeGigPayload(ctx, g, err)
}

func (ctrl Gig) Prepare(ctx context.Context, r *request.GigPrepare) (interface{}, error) {
	return api.OK(), ctrl.svc.Prepare(ctx, r.GigID)
}

func (ctrl Gig) Exec(ctx context.Context, r *request.GigExec) (interface{}, error) {
	return api.OK(), ctrl.svc.Exec(ctx, r.GigID)
}

func (ctrl Gig) Output(ctx context.Context, r *request.GigOutput) (interface{}, error) {
	out, err := ctrl.svc.Output(ctx, r.GigID)
	if err != nil {
		return nil, err
	}
	return ctrl.serve(ctx, out, err)
}

func (ctrl Gig) State(ctx context.Context, r *request.GigState) (interface{}, error) {
	return ctrl.svc.State(ctx, r.GigID)
}

func (ctrl Gig) Status(ctx context.Context, r *request.GigStatus) (interface{}, error) {
	return ctrl.svc.Status(ctx, r.GigID)
}

func (ctrl Gig) Complete(ctx context.Context, r *request.GigComplete) (interface{}, error) {
	return api.OK(), ctrl.svc.Complete(ctx, r.GigID)
}

func (ctrl Gig) Tasks(ctx context.Context, r *request.GigTasks) (interface{}, error) {
	return gigTaskPayload{Set: ctrl.svc.Tasks(ctx)}, nil
}

// ...

func (ctrl Gig) prepareSources(blob *multipart.FileHeader, uri string) (out gig.SourceWrapSet, err error) {
	if blob != nil {
		f, err := blob.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()
		out = append(out, gig.SourceWrap{
			Src:  f,
			Name: blob.Filename,
			Mime: blob.Header.Get("content-type"),
			Size: blob.Size,
		})
	}

	if uri != "" {
		out = append(out, gig.SourceWrap{
			Uri:  uri,
			Name: "",
			Mime: "",
			Size: -1,
		})
	}

	return
}

func (ctrl Gig) makeGigPayload(ctx context.Context, g *gig.Gig, err error) (*gigPayload, error) {
	if err != nil {
		return nil, err
	}

	sources := make([]gigSourcePayload, len(g.Sources))
	for i, src := range g.Sources {
		dd := src.Decoders()
		decoders := make(conv.ParamWrapSet, len(dd))
		for j, d := range dd {
			decoders[j] = ctrl.conv.WrapDecoder(d)
			if err != nil {
				return nil, err
			}
		}
		sources[i] = gigSourcePayload{
			ID:       src.ID(),
			Name:     src.FileName(),
			Size:     src.Size(),
			MimeType: src.MimeType(),
			Checksum: src.Checksum(),
			Decoders: decoders,
		}
	}

	pre := make(conv.ParamWrapSet, len(g.Preprocess))
	for i, t := range g.Preprocess {
		pre[i] = ctrl.conv.WrapPreprocessor(t)
		if err != nil {
			return nil, err
		}
	}

	post := make(conv.ParamWrapSet, len(g.Postprocess))
	for i, t := range g.Postprocess {
		post[i] = ctrl.conv.WrapPostprocessor(t)
		if err != nil {
			return nil, err
		}
	}

	return &gigPayload{
		ID:        g.ID,
		Signature: g.Signature,
		Sources:   sources,
		Worker: gigWorkerPayload{
			Ref: g.Worker.Ref(),
		},
		Preprocess:  pre,
		Postprocess: post,
	}, nil
}

func (ctrl Gig) serve(ctx context.Context, sources gig.SourceSet, err error) (interface{}, error) {
	return func(w http.ResponseWriter, req *http.Request) {
		if err != nil {
			// Simplify error handling for now
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if len(sources) > 1 {
			http.Error(w, "unable to download multiple files: compress the output to an archive", http.StatusInternalServerError)
			return
		}

		if len(sources) == 0 {
			return
		}

		src := sources[0]

		var f *os.File
		f, err := os.Open(src.Name())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		name := url.QueryEscape(src.FileName())

		w.Header().Add("Content-Disposition", "attachment; filename="+name)
		http.ServeContent(w, req, name, time.Now(), f)
	}, nil
}
