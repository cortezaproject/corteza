package indexer

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server-discovery/pkg/api"
	"github.com/cortezaproject/corteza-server-discovery/pkg/es"
	"github.com/cortezaproject/corteza-server-discovery/pkg/es/mapping"
	"github.com/cortezaproject/corteza-server-discovery/pkg/es/reindex"
	"github.com/cortezaproject/corteza-server-discovery/pkg/options"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type (
	Config struct {
		Corteza options.CortezaOpt
		ES      options.EsOpt
		Indexer options.IndexerOpt
	}

	esService interface {
		Client() (*elasticsearch.Client, error)
		BulkIndexer() (esutil.BulkIndexer, error)
	}

	apiClientService interface {
		HttpClient() *http.Client
		Mappings() (*http.Request, error)
		Feed(url.Values) (*http.Request, error)
		Resources(string, url.Values) (*http.Request, error)
		Request(string) (*http.Request, error)
		Authenticate() error
	}

	mappingService interface {
		Mappings(ctx context.Context, esc *elasticsearch.Client, indexPrefix string) (err error)
		ConfigurationMapping(ctx context.Context) (err error)
	}

	reIndexService interface {
		ReindexAll(ctx context.Context, esb esutil.BulkIndexer, indexPrefix string) error
		Watch(ctx context.Context)
	}
)

var (
	DefaultEs        esService
	DefaultApiClient apiClientService
	DefaultMapper    mappingService
	DefaultReIndexer reIndexService
)

func Initialize(ctx context.Context, log *zap.Logger, c Config) (err error) {
	DefaultEs, err = es.ES(log, c.ES)
	if err != nil {
		return
	}

	schema := c.Indexer.Schemas[0]
	if len(schema.ClientKey) == 0 || len(schema.ClientSecret) == 0 {
		return fmt.Errorf("client key and secret is missing")
	}
	DefaultApiClient, err = api.Client(c.Corteza, schema.ClientKey, schema.ClientSecret)
	if err != nil {
		return
	}

	// Map indexing for resources
	DefaultMapper = mapping.Mapper(log, DefaultEs, DefaultApiClient)

	err = DefaultMapper.ConfigurationMapping(ctx)
	if err != nil {
		return err
	}

	esc, err := DefaultEs.Client()
	if err != nil {
		return fmt.Errorf("failed to prepare es client: %w", err)
	}

	// @todo: private/public/protected indexing
	err = DefaultMapper.Mappings(ctx, esc, "private")
	if err != nil {
		return err
	}

	// Reindexing existing mapping if needed
	DefaultReIndexer = reindex.ReIndexer(log, DefaultEs, DefaultApiClient, c.ES)

	esb, err := DefaultEs.BulkIndexer()
	if err != nil {
		return fmt.Errorf("failed to prepare bulk indexer: %w", err)
	}

	err = DefaultReIndexer.ReindexAll(ctx, esb, "private")
	if err != nil {
		return err
	}

	if err = esb.Close(ctx); err != nil {
		return fmt.Errorf("failed to close bulk indexer: %w", err)
	}

	return
}

func Watchers(ctx context.Context) {
	// Initiate watcher for reindexing resource
	DefaultReIndexer.Watch(ctx)

	return
}
