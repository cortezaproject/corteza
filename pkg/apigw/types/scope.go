package types

import (
	"fmt"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/options"
)

type (
	Scp map[string]interface{}
)

func (s Scp) Keys() (kk []string) {
	for i := range s {
		kk = append(kk, i)
	}

	return
}

func (s Scp) Request() *http.Request {
	// todo - fix with expr.Request
	if _, ok := s["request"]; ok {
		switch s["request"].(type) {
		case *http.Request:
			return s["request"].(*http.Request)
		case *expr.Any:
			return s["request"].(*expr.Any).Get().(*http.Request)
		}
	}

	return nil
}

func (s Scp) Writer() http.ResponseWriter {
	if _, ok := s["writer"]; ok {
		return s["writer"].(http.ResponseWriter)
	}

	return nil
}

func (s Scp) Opts() *options.ApigwOpt {
	if _, ok := s["opts"]; ok {
		return s["opts"].(*options.ApigwOpt)
	}

	return nil
}

func (s Scp) Set(k string, v interface{}) {
	s[k] = v
}

func (s Scp) Get(k string) (v interface{}, err error) {
	var ok bool

	if v, ok = s[k]; !ok {
		err = fmt.Errorf("could not get key on index: %s", k)
		return
	}

	return
}

func (s *Scp) Dict() map[string]interface{} {
	return *s
}

func (s *Scp) Filter(fn func(k string, v interface{}) bool) *Scp {
	ss := Scp{}

	for k, v := range *s {
		if !fn(k, v) {
			continue
		}

		ss[k] = v
	}

	return &ss
}

func (s Scp) Has(k string) (has bool) {
	_, has = s[k]
	return
}
