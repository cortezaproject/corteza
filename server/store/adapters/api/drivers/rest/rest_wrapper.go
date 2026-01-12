package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/dal"
)

type (
	// restAPIWrapper is a simple wrapper that provides a simpler API to the API adapter
	restAPIWrapper struct {
		client *httpClient
		dsn    dal.DSN
	}
)

func (svc *restAPIWrapper) Run(ctx context.Context, method string, path string, payload []byte, headers map[string][]string) (rsp []byte, err error) {
	client := svc.client
	client, err = svc.appendAuth(client)
	if err != nil {
		return
	}

	switch method {
	case "GET":
		return svc.procOut(client.Get(ctx, path, headers))

	case "POST":
		return svc.procOut(client.Post(ctx, path, payload, headers))

	case "PUT":
		return svc.procOut(client.Put(ctx, path, payload, headers))

	case "PATCH":
		return svc.procOut(client.Patch(ctx, path, payload, headers))

	case "DELETE":
		return svc.procOut(client.Delete(ctx, path, headers))

	// @todo HEAD and OPTIONS but those don't really do anything at the moment

	default:
		panic(fmt.Sprintf("not supported %s", method))

	}
}

func (svc *restAPIWrapper) appendAuth(client *httpClient) (_ *httpClient, err error) {
	switch strings.ToLower(svc.dsn.AuthType) {
	case "basic":
		// noop, handled by URL construction
		break

	case "bearer":
		return svc.appendAuthBearer(client, svc.dsn.Token)

	case "apikey":
		return svc.appendAuthApiKey(client, svc.dsn.APIKey)

	default:
		err = fmt.Errorf("unknown auth type: %s", svc.dsn.AuthType)
		return
	}

	return client, nil
}

func (svc *restAPIWrapper) appendAuthBearer(client *httpClient, token string) (_ *httpClient, err error) {
	client.SetHeader("Bearer", token)
	return client, nil
}

func (svc *restAPIWrapper) appendAuthApiKey(client *httpClient, apiKey string) (_ *httpClient, err error) {
	// @todo this header could be custom
	client.SetHeader("X-API-Key", apiKey)
	return client, nil
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
