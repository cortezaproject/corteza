package registry

import (
	"fmt"
	"net/http"

	"github.com/cortezaproject/corteza/server/automation/service"
	"github.com/cortezaproject/corteza/server/pkg/apigw/filter"
	"github.com/cortezaproject/corteza/server/pkg/apigw/filter/proxy"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
)

type (
	Registry struct {
		cfg types.Config
		h   map[string]types.Handler
	}

	secureStorageTodo struct{}
)

func NewRegistry(cfg types.Config) *Registry {
	return &Registry{
		h:   map[string]types.Handler{},
		cfg: cfg,
	}
}

func (r *Registry) Add(n string, h types.Handler) {
	r.h[n] = h
}

func (r *Registry) Merge(h types.Handler, b []byte, cfg types.Config) (hh types.Handler, err error) {
	hh, err = h.Merge(b, cfg)
	return
}

func (r *Registry) Get(identifier string) (types.Handler, error) {
	var (
		ok bool
		f  types.Handler
	)

	if f, ok = r.h[identifier]; !ok {
		return nil, fmt.Errorf("could not get element from registry: %s", identifier)
	}

	return f.New(r.cfg), nil
}

func (r *Registry) All() (list types.FilterMetaList) {
	for _, handler := range r.h {
		if !handler.Enabled() {
			continue
		}

		meta := handler.Meta()
		list = append(list, &meta)
	}

	return
}

func (r *Registry) Preload() {
	// prefilters
	r.Add("queryParam", filter.NewQueryParam(r.cfg))
	r.Add("header", filter.NewHeader(r.cfg))
	r.Add("profiler", filter.NewProfiler(r.cfg))

	// processers
	r.Add("workflow", filter.NewWorkflow(r.cfg, NewWorkflow()))
	r.Add("proxy", proxy.New(r.cfg, service.DefaultLogger, http.DefaultClient, secureStorageTodo{}))
	r.Add("payload", filter.NewPayload(r.cfg, service.DefaultLogger))

	// postfilters
	r.Add("redirection", filter.NewRedirection(r.cfg))
	r.Add("response", filter.NewResponse(r.cfg, service.Registry()))
	r.Add("defaultJsonResponse", filter.NewDefaultJsonResponse(r.cfg))
}

func NewWorkflow() (wf filter.WfExecer) {
	return service.DefaultWorkflow
}
