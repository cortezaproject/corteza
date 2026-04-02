package searcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cortezaproject/corteza/extra/server-discovery/pkg/api"
	"github.com/cortezaproject/corteza/extra/server-discovery/pkg/options"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"go.uber.org/zap"
)

type (
	Config struct {
		Corteza    options.CortezaOpt
		ES         options.EsOpt
		HttpServer options.HttpServerOpt
		Searcher   options.SearcherOpt
	}

	apiClientService interface {
		HttpClient() *http.Client
		Namespaces() (*http.Request, error)
		Modules(uint64) (*http.Request, error)
		Request(string) (*http.Request, error)
		Authenticate() error
	}
)

var (
	DefaultLogger *zap.Logger
	DefaultConfig Config

	DefaultEsClient  *elasticsearch.Client
	DefaultApiClient apiClientService
)

func Initialize(_ context.Context, log *zap.Logger, c Config, esClient *elasticsearch.Client) (err error) {
	DefaultLogger = log.Named("service")
	DefaultConfig = c

	DefaultEsClient = esClient

	DefaultApiClient, err = api.Client(c.Corteza, c.Searcher.ClientKey, c.Searcher.ClientSecret)
	if err != nil {
		return
	}

	return
}

// @todo move this to es service
func validElasticResponse(res *esapi.Response, err error) error {
	if err != nil {
		return fmt.Errorf("failed to get response from search backend: %w", err)
	}

	if res.IsError() {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		var rsp struct {
			Error struct {
				Type   string
				Reason string
			}
		}

		if err := json.NewDecoder(res.Body).Decode(&rsp); err != nil {
			return fmt.Errorf("could not parse response body: %w", err)
		} else {
			return fmt.Errorf("search backend responded with an error: %s (type: %s, status: %s)", rsp.Error.Reason, rsp.Error.Type, res.Status())
		}
	}

	return nil
}
