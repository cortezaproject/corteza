package registry

import (
	"fmt"
	"net/http"

	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/pkg/apigw/filter"
	"github.com/cortezaproject/corteza-server/pkg/apigw/filter/proxy"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"go.uber.org/zap"
)

type (
	Registry struct {
		h map[string]types.Handler
	}

	secureStorageTodo struct{}
)

func NewRegistry() *Registry {
	return &Registry{
		h: map[string]types.Handler{},
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

	return f, nil
}

func (r *Registry) All() (list types.FilterMetaList) {
	for _, handler := range r.h {
		meta := handler.Meta()
		list = append(list, &meta)
	}

	return
}

func (r *Registry) Preload() {
	// prefilters
	r.Add("queryParam", filter.NewQueryParam())
	r.Add("header", filter.NewHeader())

	// processers
	r.Add("workflow", filter.NewWorkflow(NewWorkflow()))
	r.Add("proxy", proxy.New(service.DefaultLogger, http.DefaultClient, secureStorageTodo{}))
	r.Add("payload", filter.NewPayload(service.DefaultLogger))

	// postfilters
	r.Add("redirection", filter.NewRedirection())
	r.Add("defaultJsonResponse", filter.NewDefaultJsonResponse())
}

func NewWorkflow() (wf filter.WfExecer) {
	return service.Workflow(&zap.Logger{}, options.CorredorOpt{})
}
