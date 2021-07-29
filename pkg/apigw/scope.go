package apigw

import (
	"fmt"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/options"
)

type (
	scp map[string]interface{}
)

func (s scp) Keys() (kk []string) {
	for i := range s {
		kk = append(kk, i)
	}

	return
}

func (s scp) Request() *http.Request {
	if _, ok := s["request"]; ok {
		return s["request"].(*http.Request)
	}

	return nil
}

func (s scp) Writer() http.ResponseWriter {
	if _, ok := s["writer"]; ok {
		return s["writer"].(http.ResponseWriter)
	}

	return nil
}

func (s scp) Opts() *options.ApigwOpt {
	if _, ok := s["opts"]; ok {
		return s["opts"].(*options.ApigwOpt)
	}

	return nil
}

func (s scp) Set(k string, v interface{}) {
	s[k] = v
}

func (s scp) Get(k string) (v interface{}, err error) {
	var ok bool

	if v, ok = s[k]; !ok {
		err = fmt.Errorf("could not get key on index: %s", k)
		return
	}

	return
}

func (s *scp) Filter(fn func(k string, v interface{}) bool) *scp {
	ss := scp{}

	for k, v := range *s {
		if !fn(k, v) {
			continue
		}

		ss[k] = v
	}

	return &ss
}
