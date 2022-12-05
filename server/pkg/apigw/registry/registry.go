package registry

import (
	"fmt"
	"net/http"

	"github.com/cortezaproject/corteza/server/automation/service"
	"github.com/cortezaproject/corteza/server/pkg/apigw/filter"
	"github.com/cortezaproject/corteza/server/pkg/apigw/filter/proxy"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	"github.com/cortezaproject/corteza/server/pkg/options"
)

type (
	Registry struct {
		opts options.ApigwOpt
		h    map[string]types.Handler
	}

	secureStorageTodo struct{}
)

func NewRegistry(opts options.ApigwOpt) *Registry {
	return &Registry{
		h:    map[string]types.Handler{},
		opts: opts,
	}
}

func (r *Registry) Add(n string, h types.Handler) {
	r.h[n] = h
}

func (r *Registry) Merge(h types.Handler, b []byte) (hh types.Handler, err error) {
	hh, err = h.Merge(b)
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

	return f.New(r.opts), nil
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
	r.Add("queryParam", filter.NewQueryParam(r.opts))
	r.Add("header", filter.NewHeader(r.opts))
	r.Add("profiler", filter.NewProfiler(r.opts))

	// processers
	r.Add("workflow", filter.NewWorkflow(r.opts, NewWorkflow()))
	r.Add("proxy", proxy.New(r.opts, service.DefaultLogger, http.DefaultClient, secureStorageTodo{}))
	r.Add("payload", filter.NewPayload(r.opts, service.DefaultLogger))

	// postfilters
	r.Add("redirection", filter.NewRedirection(r.opts))
	r.Add("response", filter.NewResponse(r.opts, service.Registry()))
	r.Add("defaultJsonResponse", filter.NewDefaultJsonResponse(r.opts))
}

func NewWorkflow() (wf filter.WfExecer) {
	return service.DefaultWorkflow
}
