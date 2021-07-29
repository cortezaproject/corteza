package apigw

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/jsenv"
	"go.uber.org/zap"
)

var (
	hopHeaders = []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}
)

type (
	WfExecer interface {
		Exec(ctx context.Context, workflowID uint64, p atypes.WorkflowExecParams) (*expr.Vars, atypes.Stacktrace, error)
	}

	processerWorkflow struct {
		functionMeta
		d WfExecer

		params struct {
			Workflow uint64 `json:"workflow,string"`
		}
	}

	processerProxy struct {
		functionMeta
		a   AuthServicer
		c   *http.Client
		log *zap.Logger

		params struct {
			Location string     `json:"location"`
			Auth     authParams `json:"auth"`
		}
	}

	processerPayload struct {
		functionMeta
		vm  jsenv.Vm
		log *zap.Logger

		params struct {
			Func   string `json:"func"`
			Encode bool   `json:"encode"`
		}
	}
)

func NewProcesserWorkflow(wf WfExecer) (p *processerWorkflow) {
	p = &processerWorkflow{}

	p.d = wf

	p.Step = 2
	p.Name = "processerWorkflow"
	p.Label = "Workflow processer"
	p.Kind = FunctionKindProcesser

	p.Args = []*functionMetaArg{
		{
			Type:    "workflow",
			Label:   "workflow",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h processerWorkflow) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h processerWorkflow) Meta() functionMeta {
	return h.functionMeta
}

func (f *processerWorkflow) Merge(params []byte) (Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&f.params)

	return f, err
}

func (h processerWorkflow) Exec(ctx context.Context, scope *scp) error {
	var (
		err error
	)

	// // pp := map[string]interface{}{
	// // 	"payload": "test",
	// // }

	// // ppe, err := expr.NewVars(pp)

	// payload, err := scope.Get("payload")

	// if err != nil {
	// 	return err
	// }

	// // setup scope for workflow
	// vv := map[string]interface{}{
	// 	"payload": payload,
	// }

	// // temp
	// // for i, v := range map[string]interface{}(*scope) {
	// // 	vv[i] = v
	// // }

	// // get the request data and put it into vars

	// in, err := expr.NewVars(vv)
	// // rq, err := automation.NewRequest(*req)

	// if err != nil {
	// 	return err
	// }

	// // if err != nil {
	// // 	return err
	// // }

	// wp := atypes.WorkflowExecParams{
	// 	Trace: true,
	// 	// todo depending on settings per-route
	// 	Async: false,
	// 	// todo depending on settings per-route
	// 	Wait:  true,
	// 	Input: in,
	// }

	// out, _, err := h.d.Exec(ctx, h.params.Workflow, wp)

	// spew.Dump("OUTTT", out)

	// if err != nil {
	// 	return err
	// }

	// // merge out with scope
	// merged, err := in.Merge(out)

	// if err != nil {
	// 	return err
	// }

	// mm, err := expr.CastToVars(merged)

	// for k, v := range mm {
	// 	scope.Set(k, v)
	// }

	// // spew.Dump("MMMM", mm)

	// ss := scope.Filter(func(k string, v interface{}) bool {
	// 	if k == "eventType" || k == "resourceType" {
	// 		return false
	// 	}

	// 	return true
	// })

	// // spew.Dump(scope.Get("payload"))
	// trgt, _ := ss.Get("trgt")

	// scope.Writer().Write([]byte(trgt.(*expr.String).GetValue()))

	return err
}

func NewProcesserProxy(l *zap.Logger, c *http.Client) (p *processerProxy) {
	p = &processerProxy{}

	p.c = c
	p.log = l

	p.Step = 2
	p.Name = "processerProxy"
	p.Label = "Proxy processer"
	p.Kind = FunctionKindProcesser

	p.Args = []*functionMetaArg{
		{
			Type:    "text",
			Label:   "location",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h processerProxy) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h processerProxy) Meta() functionMeta {
	return h.functionMeta
}

func (f *processerProxy) Merge(params []byte) (Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&f.params)

	if err != nil {
		return nil, err
	}

	// get the auth mechanism
	f.a, err = NewAuthServicer(f.c, f.params.Auth)

	if err != nil {
		return nil, fmt.Errorf("could not load auth servicer for proxying: %s", err)
	}

	return f, err
}

func (h processerProxy) Exec(ctx context.Context, scope *scp) (err error) {
	ctx, cancel := context.WithTimeout(ctx, scope.Opts().ProxyOutboundTimeout)
	defer cancel()

	req := scope.Request()
	log := h.log.With(zap.String("ref", h.Name))

	outreq := req.Clone(ctx)

	l, err := url.ParseRequestURI(h.params.Location)

	if err != nil {
		return fmt.Errorf("could not parse destination location for proxying: %s", err)
	}

	// should we preserve query params? headers? post data?
	outreq.URL = l
	outreq.RequestURI = ""
	outreq.Method = req.Method
	outreq.Host = l.Hostname()

	// use authservicer, set any additional headers
	err = h.a.Do(outreq)

	if err != nil {
		return fmt.Errorf("errors setting auth for proxying: %s", err)
	}

	// merge the old query params to the new request
	// do not overwrite old ones
	// do it after the authServicer, since we also may add them there
	mergeQueryParams(req, outreq)

	if scope.Opts().ProxyEnableDebugLog {
		o, _ := httputil.DumpRequestOut(outreq, false)
		log.Debug("proxy outbound request", zap.Any("request", string(o)))
	}

	// temporary metrics before the proper functionality
	startTime := time.Now()

	// todo - disable / enable follow redirects, already
	// added to options
	resp, err := h.c.Do(outreq)

	if err != nil {
		return fmt.Errorf("could not proxy request: %s", err)
	}

	if scope.Opts().ProxyEnableDebugLog {
		o, _ := httputil.DumpResponse(resp, false)
		log.Debug("proxy outbound response", zap.Any("request", string(o)), zap.Duration("duration", time.Since(startTime)))
	}

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("could not read get body on proxy request: %s", err)
	}

	mergeHeaders(resp.Header, scope.Writer().Header())

	// add to writer
	scope.Writer().Write(b)

	return nil
}

func NewProcesserPayload(l *zap.Logger) (p *processerPayload) {
	p = &processerPayload{}

	// todo - check the consequences of doing this here
	p.vm = jsenv.New(jsenv.NewTransformer(jsenv.LoaderJS, jsenv.TargetES2016))
	p.log = l

	p.Step = 2
	p.Name = "processerPayload"
	p.Label = "Payload processer"
	p.Kind = FunctionKindProcesser

	p.Args = []*functionMetaArg{
		{
			Type:    "jsfunc",
			Label:   "jsfunc",
			Options: map[string]interface{}{},
		},
	}

	// register a request body reader
	// since it's a readcloser, it can be read only once
	p.vm.Register("readRequestBody", func(rc io.ReadCloser) string {
		b, _ := io.ReadAll(rc)
		return string(b)
	})

	return
}

func (h processerPayload) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h processerPayload) Meta() functionMeta {
	return h.functionMeta
}

func (f *processerPayload) Merge(params []byte) (Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&f.params)

	if err != nil {
		return nil, err
	}

	fn, err := base64.StdEncoding.DecodeString(f.params.Func)

	if err != nil {
		return nil, fmt.Errorf("could not decode js func: %s", err)
	}

	f.params.Func = string(fn)

	return f, err
}

func (h processerPayload) Exec(ctx context.Context, scope *scp) (err error) {
	log := h.log.With(zap.String("function", h.String()))

	if h.params.Func == "" {
		err = errors.New("function body empty")
		log.Debug("could not register function", zap.Error(err))
		return
	}

	fn, err := h.vm.RegisterFunction(h.params.Func)

	if err != nil {
		log.Debug("could not register function", zap.Error(err))
		return
	}

	out, err := fn.Exec(h.vm.New(scope))

	if err != nil {
		log.Debug("could not exec payload function", zap.Error(err))
		return
	}

	// add to scope, so next steps can get the structure
	scope.Set("payload", out)

	// check if string
	switch out.(type) {
	case string:
		// handling the newline, to keep the consistency with the json encoder
		// which automatically appends the newline
		_, err = scope.Writer().Write([]byte(fmt.Sprintf("%s\n", out)))
	default:
		err = json.NewEncoder(scope.Writer()).Encode(out)
	}

	if err != nil {
		log.Debug("could not write to body", zap.Error(err))
		return
	}

	return
}

func (h processerPayload) VM() jsenv.Vm {
	return h.vm
}

func mergeHeaders(orig, dest http.Header) {
OUTER:
	for name, values := range orig {
		// skip headers that need to be omitted
		// when proxying
		for _, v := range hopHeaders {
			if v == name {
				continue OUTER
			}
		}
		dest[name] = values
	}
}

func mergeQueryParams(orig, dest *http.Request) {
	origValues := dest.URL.Query()

	for k, qp := range orig.URL.Query() {
		// skip existing
		if dest.URL.Query().Get(k) != "" {
			continue
		}

		for _, v := range qp {
			origValues.Add(k, v)
		}
	}

	dest.URL.RawQuery = origValues.Encode()
}
