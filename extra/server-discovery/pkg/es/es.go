package es

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"go.uber.org/zap"
)

type (
	es struct {
		log *zap.Logger
		esc *elasticsearch.Client
		//esb esutil.BulkIndexer
	}

	Service interface {
		EsClient() *elasticsearch.Client
		EsBulk() esutil.BulkIndexer
		Watch(ctx context.Context)
	}

	//apiClientService interface {
	//	HttpClient() *http.Client
	//	Mappings() (*http.Request, error)
	//	Resources(string, url.Values) (*http.Request, error)
	//	Request(string) (*http.Request, error)
	//	Authenticate() error
	//}
)

func ES(log *zap.Logger, esc *elasticsearch.Client) (out *es, err error) {
	out = &es{log: log, esc: esc}
	return
}

func (es *es) BulkIndexer() (esutil.BulkIndexer, error) {
	return esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     es.esc,
		FlushBytes: 5e+5,
	})
}

func ValidElasticResponse(res *esapi.Response, err error) error {
	if err != nil {
		return fmt.Errorf("failed to get response from search backend: %w", err)
	}

	if res.IsError() {
		defer res.Body.Close()
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
