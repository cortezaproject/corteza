package es

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cortezaproject/corteza/extra/server-discovery/pkg/options"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"go.uber.org/zap"
)

type (
	es struct {
		log *zap.Logger
		opt options.EsOpt
		//esc *elasticsearch.Client
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

func ES(log *zap.Logger, opt options.EsOpt) (out *es, err error) {
	out = &es{log: log, opt: opt}
	return
}

func (es *es) Client() (*elasticsearch.Client, error) {
	config := elasticsearch.Config{
		Addresses:            es.opt.Addresses,
		EnableRetryOnTimeout: es.opt.EnableRetryOnTimeout,
		MaxRetries:           es.opt.MaxRetries,
	}
	if len(es.opt.Username) > 0 {
		config.Username = es.opt.Username
		config.Password = es.opt.Password
	}

	// if the user provided a CA certificate file, we need to load it
	// and set it in the TLS config
	if es.opt.CertFile != "" {
		rootCACert, err := os.ReadFile(es.opt.CertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate file: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(rootCACert) {
			return nil, fmt.Errorf("failed parse root certificate: %w", err)
		}

		config.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		}
	} else if !es.opt.Secure {
		config.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	return elasticsearch.NewClient(config)
}

func (es *es) BulkIndexer() (esutil.BulkIndexer, error) {
	client, err := es.Client()
	if err != nil {
		return nil, err
	}

	return esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     client,
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
