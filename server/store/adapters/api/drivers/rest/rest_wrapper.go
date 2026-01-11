package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type (
	// restAPIWrapper is a simple wrapper that provides a simpler API to the API adapter
	restAPIWrapper struct {
		client *httpClient
	}
)

func (svc *restAPIWrapper) Run(ctx context.Context, method string, path string, payload []byte, headers map[string][]string) (rsp []byte, err error) {
	switch method {
	case "GET":
		return svc.procOut(svc.client.Get(ctx, path, headers))

	case "POST":
		return svc.procOut(svc.client.Post(ctx, path, payload, headers))

	case "PUT":
		return svc.procOut(svc.client.Put(ctx, path, payload, headers))

	case "PATCH":
		return svc.procOut(svc.client.Patch(ctx, path, payload, headers))

	case "DELETE":
		return svc.procOut(svc.client.Delete(ctx, path, headers))

	// @todo HEAD and OPTIONS but those don't really do anything at the moment

	default:
		panic(fmt.Sprintf("not supported %s", method))

	}
}

// @todo would make sense to stream the output
func (svc *restAPIWrapper) procOut(resp *http.Response, err error) (rsp []byte, _ error) {
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
