package apigw

import (
	"fmt"
	"net/http"

	"github.com/cortezaproject/corteza-server/automation/service"
)

type (
	registry struct {
		h map[string]Handler
	}
)

func NewRegistry() *registry {
	return &registry{
		h: map[string]Handler{},
	}
}

func (r *registry) Add(n string, h Handler) {
	r.h[n] = h
}

func (r *registry) Merge(h Handler, b []byte) (hh Handler, err error) {
	hh, err = h.Merge(b)
	return
}

func (r *registry) Get(identifier string) (Handler, error) {
	var (
		ok bool
		f  Handler
	)

	if f, ok = r.h[identifier]; !ok {
		return nil, fmt.Errorf("could not get element from registry: %s", identifier)
	}

	return f, nil
}

func (r *registry) All() (list functionMetaList) {
	for _, handler := range r.h {
		meta := handler.Meta()
		list = append(list, &meta)
	}

	return
}

func (r *registry) Preload() {
	r.Add("verifierQueryParam", NewVerifierQueryParam())
	r.Add("verifierOrigin", NewVerifierOrigin())
	r.Add("validatorHeader", NewValidatorHeader())
	r.Add("expediterRedirection", NewExpediterRedirection())
	r.Add("processerWorkflow", NewProcesserWorkflow(NewWorkflow()))
	r.Add("processerProxy", NewProcesserProxy(service.DefaultLogger, http.DefaultClient))
}
