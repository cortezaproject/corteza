package automation

import (
	"context"
	"io"

	atypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/expr"
)

type (
	constructSvc interface {
		AddFunctions(ff ...atypes.ConstructFunction)
		AddTriggers(tt ...atypes.ConstructTrigger)
	}

	typeRegistry interface {
		Type(ref string) expr.Type
	}

	ngHttpRequestHandler struct {
		reg  constructSvc
		tReg typeRegistry
		h    *httpRequestHandler
	}
)

func NgHttpRequestHandler(reg constructSvc, tReg typeRegistry) *ngHttpRequestHandler {
	ngh := &ngHttpRequestHandler{
		reg:  reg,
		tReg: tReg,
		h:    &httpRequestHandler{},
	}

	ngh.register()
	return ngh
}

func (h ngHttpRequestHandler) register() {}

func (h ngHttpRequestHandler) Send() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "httpRequestSend",
		Kind:   "function",
		Groups: []string{"HTTP"},
		Labels: map[string]string{"http request": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Send HTTP Request",
			Description: "Sends HTTP requests",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "link"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "url", Types: []string{"String"}, Required: true},
			{ArgumentName: "method", Types: []string{"String"}, Required: true},
			{ArgumentName: "params", Types: []string{"KVV"}},
			{ArgumentName: "headers", Types: []string{"KVV"}},
			{ArgumentName: "headerAuthBearer", Types: []string{"String"}},
			{ArgumentName: "headerAuthUsername", Types: []string{"String"}},
			{ArgumentName: "headerAuthPassword", Types: []string{"String"}},
			{ArgumentName: "headerUserAgent", Types: []string{"String"}},
			{ArgumentName: "headerContentType", Types: []string{"String"}},
			{ArgumentName: "timeout", Types: []string{"Duration"}},
			{ArgumentName: "form", Types: []string{"KVV"}},
			{ArgumentName: "body", Types: []string{"String", "Reader", "Any"}},
		},

		Results: []*atypes.Param{
			{ArgumentName: "status", Types: []string{"String"}},
			{ArgumentName: "statusCode", Types: []string{"Integer"}},
			{ArgumentName: "headers", Types: []string{"KVV"}},
			{ArgumentName: "contentLength", Types: []string{"Integer"}},
			{ArgumentName: "contentType", Types: []string{"String"}},
			{ArgumentName: "body", Types: []string{"Reader"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "URL", Argument: "url"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Method", Argument: "method"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Query Params (KVV)", Argument: "params"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Headers (KVV)", Argument: "headers"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Auth Bearer", Argument: "headerAuthBearer"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Auth Username", Argument: "headerAuthUsername"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Auth Password", Argument: "headerAuthPassword"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "User Agent", Argument: "headerUserAgent"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Content Type", Argument: "headerContentType"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Timeout", Argument: "timeout"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Form (KVV)", Argument: "form"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Body", Argument: "body"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &httpRequestSendArgs{
					hasUrl:                in.Has("url"),
					hasMethod:             in.Has("method"),
					hasParams:             in.Has("params"),
					hasHeaders:            in.Has("headers"),
					hasHeaderAuthBearer:   in.Has("headerAuthBearer"),
					hasHeaderAuthUsername: in.Has("headerAuthUsername"),
					hasHeaderAuthPassword: in.Has("headerAuthPassword"),
					hasHeaderUserAgent:    in.Has("headerUserAgent"),
					hasHeaderContentType:  in.Has("headerContentType"),
					hasTimeout:            in.Has("timeout"),
					hasForm:               in.Has("form"),
					hasBody:               in.Has("body"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Body argument
			if args.hasBody {
				aux := expr.Must(expr.Select(in, "body"))
				switch aux.Type() {
				case h.tReg.Type("String").Type():
					args.bodyString = aux.Get().(string)
				case h.tReg.Type("Reader").Type():
					args.bodyStream = aux.Get().(io.Reader)
				case h.tReg.Type("Any").Type():
					args.bodyRaw = aux.Get().(interface{})
				}
			}

			var results *httpRequestSendResults
			if results, err = h.h.send(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			if tval, err := h.tReg.Type("String").Cast(results.Status); err == nil {
				_ = expr.Assign(out, "status", tval)
			}
			if tval, err := h.tReg.Type("Integer").Cast(results.StatusCode); err == nil {
				_ = expr.Assign(out, "statusCode", tval)
			}
			if tval, err := h.tReg.Type("KVV").Cast(results.Headers); err == nil {
				_ = expr.Assign(out, "headers", tval)
			}
			if tval, err := h.tReg.Type("Integer").Cast(results.ContentLength); err == nil {
				_ = expr.Assign(out, "contentLength", tval)
			}
			if tval, err := h.tReg.Type("String").Cast(results.ContentType); err == nil {
				_ = expr.Assign(out, "contentType", tval)
			}
			if tval, err := h.tReg.Type("Reader").Cast(results.Body); err == nil {
				_ = expr.Assign(out, "body", tval)
			}

			return
		},
	}
}
