package healthcheck

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
)

type (
	checkFn func(ctx context.Context) error

	Meta struct {
		Label       string
		Description string
	}

	check struct {
		fn checkFn
		*Meta
	}

	result struct {
		err error
		*Meta
	}

	results []*result
	checks  struct {
		cc []*check
	}
)

var (
	defaults *checks
)

func init() {
	defaults = New()
}

func Defaults() *checks {
	return defaults
}

func New() *checks {
	return &checks{cc: []*check{}}
}

// Add appends new check
func (c *checks) Add(fn checkFn, label string, description ...string) {
	c.cc = append(c.cc, &check{fn, &Meta{Label: label, Description: strings.Join(description, "")}})
}

func (c checks) Run(ctx context.Context) results {
	var rr = make([]*result, len(c.cc))

	for i, c := range c.cc {
		rr[i] = &result{c.fn(ctx), c.Meta}
	}

	return rr
}

func (rr results) Healthy() bool {
	for _, c := range rr {
		if c.err != nil {
			return false
		}
	}

	return true
}

func (rr results) String() string {
	buf := &bytes.Buffer{}

	rr.WriteTo(buf)

	return buf.String()
}

func (rr results) WriteTo(w io.Writer) (int64, error) {
	var (
		p = func(f string, aa ...interface{}) {
			_, _ = fmt.Fprintf(w, f, aa...)
		}
	)

	for _, r := range rr {
		if r.IsHealthy() {
			p("PASS")
		} else {
			p("FAIL")
		}

		p(" %s", r.Label)

		if !r.IsHealthy() {
			p(": %v", r.Error())
		}

		p("\n")
	}

	return 0, nil
}

func (r *result) IsHealthy() bool {
	return r != nil && r.err == nil
}

func (r *result) Error() string {
	if r == nil || r.err == nil {
		return ""
	}

	return r.err.Error()
}
