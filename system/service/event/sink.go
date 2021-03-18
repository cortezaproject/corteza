package event

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/system/types"
	"strings"
)

const (
	sinkMatchRequestGet    = "request.get."
	sinkMatchRequestPost   = "request.post."
	sinkMatchRequestHeader = "request.header."
)

// Match returns false if given conditions do not match event & resource internals
func (res sinkBase) Match(c eventbus.ConstraintMatcher) bool {
	return sinkMatch(res.request, c)
}

// Handles sink's URL matchers
func sinkMatch(r *types.SinkRequest, c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "request.host":
		return c.Match(r.Host)
	case "request.remoteAddress", "request.remote-address":
		return c.Match(r.RemoteAddr)
	case "request.method":
		return c.Match(r.Method)
	case "request.path":
		// match path, right side of "/sink"
		return c.Match(r.Path)
	case "request.username":
		return c.Match(r.Username)
	case "request.password":
		return c.Match(r.Password)
	case "request.content-type":
		return c.Match(r.Header.Get("content-type"))
	}

	// Dynamically check matcher name if it contains request.(get|post|header).*
	// and use value for matcher:
	//
	// to match "&foo=bar" in URL string use .where('request.get.foo', 'bar')
	//
	// It only matches first value (get, post and header can have multiple values)

	if strings.HasPrefix(c.Name(), sinkMatchRequestGet) {
		return c.Match(r.Query.Get(c.Name()[len(sinkMatchRequestGet):]))
	}

	if strings.HasPrefix(c.Name(), sinkMatchRequestPost) {
		return c.Match(r.PostForm.Get(c.Name()[len(sinkMatchRequestPost):]))
	}

	if strings.HasPrefix(c.Name(), sinkMatchRequestHeader) {
		return c.Match(r.Header.Get(c.Name()[len(sinkMatchRequestHeader):]))
	}

	return true
}
