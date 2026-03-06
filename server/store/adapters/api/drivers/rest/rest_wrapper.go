package rest

import (
	"context"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/api/cred_registry"
)

type (
	// restAPIWrapper is a simple wrapper that provides a simpler API to the API adapter
	restAPIWrapper struct {
		client       *httpClient
		dsn          dal.DSN
		connectionID uint64
	}
)

func (svc *restAPIWrapper) Run(ctx context.Context, method string, path string, payload []byte, headers map[string][]string) (statusCode int, outHeaders map[string][]string, rsp []byte, err error) {
	client, err := svc.appendAuth(ctx, svc.cloneClient())
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

	case "HEAD":
		return svc.procOut(client.Head(ctx, path, headers))

	// @todo OPTIONS

	default:
		panic(fmt.Sprintf("not supported %s", method))

	}
}

func (svc *restAPIWrapper) appendAuth(ctx context.Context, client *httpClient) (_ *httpClient, err error) {
	switch strings.ToLower(svc.dsn.AuthType) {
	case "basic":
		// noop, handled by URL construction
		break

	case "bearer":
		return svc.appendAuthBearer(client, svc.dsn.Token)

	case "apikey":
		return svc.appendAuthApiKey(client, svc.dsn.APIKey)

	case "oauth2_client_credentials":
		return svc.appendAuthOAuth2(ctx, client)

	default:
		err = fmt.Errorf("unknown auth type: %s", svc.dsn.AuthType)
		return
	}

	return client, nil
}

func (svc *restAPIWrapper) appendAuthBearer(client *httpClient, token string) (_ *httpClient, err error) {
	client.SetHeader("Authorization", "Bearer "+token)
	return client, nil
}

func (svc *restAPIWrapper) appendAuthApiKey(client *httpClient, apiKey string) (_ *httpClient, err error) {
	client.SetHeader(svc.dsn.APIKeyHeader, "token "+apiKey)
	return client, nil
}

func (svc *restAPIWrapper) appendAuthOAuth2(ctx context.Context, client *httpClient) (_ *httpClient, err error) {
	token, err := cred_registry.Default().GetAccessToken(ctx, svc.connectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth2 access token: %w", err)
	}

	client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	return client, nil
}

func (svc *restAPIWrapper) cloneClient() *httpClient {
	c := *svc.client
	c.headers = make(map[string][]string, len(svc.client.headers))
	maps.Copy(c.headers, svc.client.headers)
	return &c
}

// @todo would make sense to stream the output
func (svc *restAPIWrapper) procOut(resp *http.Response, err error) (statusCode int, outHeaders map[string][]string, rsp []byte, _ error) {
	if err != nil {
		return 0, nil, nil, err
	}

	if resp == nil {
		return 0, nil, nil, fmt.Errorf("nil response")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, resp.Header, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp.StatusCode, resp.Header, body, nil
}
