package apigw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type (
	expediterRedirection struct {
		functionMeta
		params struct {
			HTTPStatus int    `json:"status,string"`
			Location   string `json:"location"`
		}
	}

	errorHandler struct {
		name   string
		args   []string
		weight int
		step   int
	}
)

func NewExpediterRedirection() (e *expediterRedirection) {
	e = &expediterRedirection{}

	e.Step = 3
	e.Name = "expediterRedirection"
	e.Label = "Redirection expediter"
	e.Kind = FunctionKindExpediter

	e.Args = []*functionMetaArg{
		{
			Type:    "status",
			Label:   "status",
			Options: map[string]interface{}{},
		},
		{
			Type:    "text",
			Label:   "location",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h expediterRedirection) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h expediterRedirection) Meta() functionMeta {
	return h.functionMeta
}

func (f *expediterRedirection) Merge(params []byte) (Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&f.params)
	return f, err
}

func (h expediterRedirection) Exec(ctx context.Context, scope *scp) error {
	loc, err := url.ParseRequestURI(h.params.Location)

	if err != nil {
		return fmt.Errorf("could not redirect: %s", err)
	}

	status := h.params.HTTPStatus

	if !checkStatus("redirect", status) {
		return fmt.Errorf("could not redirect: wrong status %d", status)
	}

	http.Redirect(scope.Writer(), scope.Request(), loc.String(), status)

	return nil
}

func (pp errorHandler) Exec(ctx context.Context, scope *scp, err error) {
	type (
		responseHelper struct {
			ErrResponse struct {
				Msg string `json:"msg"`
			} `json:"error"`
		}
	)

	resp := responseHelper{
		ErrResponse: struct {
			Msg string "json:\"msg\""
		}{
			Msg: err.Error(),
		},
	}

	// set http status code
	scope.Writer().WriteHeader(http.StatusInternalServerError)

	// set body
	json.NewEncoder(scope.Writer()).Encode(resp)

}

func checkStatus(typ string, status int) bool {
	switch typ {
	case "redirect":
		return status >= 300 && status <= 399
	default:
		return true
	}
}
