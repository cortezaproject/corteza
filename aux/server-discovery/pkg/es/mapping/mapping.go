package mapping

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server-discovery/pkg/es/reindex"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type (
	esIndex struct {
		Name      string `json:"index"`
		Health    string `json:"health"`
		Status    string `json:"status"`
		DocsCount string `json:"docs.count"`
		StoreSize string `json:"store.size"`
	}

	rspDiscoveryMappings struct {
		Response []*mapping
	}

	reqMapping struct {
		// @todo settings
		Mappings struct {
			Properties map[string]*property `json:"properties,omitempty"`
		} `json:"mappings,omitempty"`
	}

	mapping struct {
		Index        string               `json:"index"`
		Properties   map[string]*property `json:"properties"`
		Mapping      map[string]*property `json:"mapping"`
		DocumentsURL string               `json:"documentsURL"`
	}

	property struct {
		// https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping-types.html
		Type string `json:"type,omitempty"`

		// Boost factor
		// https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping-boost.html
		Boost float32 `json:"boost,omitempty"`

		Properties map[string]*property `json:"properties,omitempty"`
	}

	esService interface {
		Client() (*elasticsearch.Client, error)
		BulkIndexer() (esutil.BulkIndexer, error)
	}

	apiClientService interface {
		HttpClient() *http.Client
		Mappings() (*http.Request, error)
		Resources(string, url.Values) (*http.Request, error)
		Request(string) (*http.Request, error)
		Authenticate() error
	}

	mapper struct {
		log *zap.Logger
		es  esService
		api apiClientService
	}
)

func Mapper(log *zap.Logger, esc esService, api apiClientService) *mapper {
	return &mapper{
		log: log,
		es:  esc,
		api: api,
	}
}

// Mappings fetches mappings from discovery server and update elastic search indexes
func (m *mapper) Mappings(ctx context.Context, esc *elasticsearch.Client, indexPrefix string) (err error) {
	var (
		req             *http.Request
		rsp             *http.Response
		rspPayload      = &rspDiscoveryMappings{}
		buf             = &bytes.Buffer{}
		esRsp           *esapi.Response
		existingIndexes []*esIndex
		index           string
	)

	if req, err = m.api.Mappings(); err != nil {
		return fmt.Errorf("failed to prepare mapping request: %w", err)
	}

	//d, _ := httputil.DumpRequest(req, true)
	//println(string(d))

	if rsp, err = m.api.HttpClient().Do(req.WithContext(ctx)); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("request resulted in an unexpected status: %s", rsp.Status)
	}

	//d, _ := httputil.DumpResponse(rsp, true)
	//println(string(d))

	if err = json.NewDecoder(rsp.Body).Decode(rspPayload); err != nil {
		return fmt.Errorf("failed to decode mapping response: %w", err)
	}
	if err = rsp.Body.Close(); err != nil {
		return fmt.Errorf("failed to close mapping response body: %w", err)
	}

	if existingIndexes, err = m.getExistingIndexes(ctx); err != nil {
		return fmt.Errorf("failed to fetch existing indexes: %w", err)
	}

	indexMap := m.mapExistingIndexes(existingIndexes)

	//esc, err := m.es.Client()
	//if err != nil {
	//	return fmt.Errorf("failed to prepare es client: %w", err)
	//}

	for _, im := range rspPayload.Response {
		buf.Reset()
		esReq := reqMapping{}
		//esReq.Mappings.Properties = im.Properties
		esReq.Mappings.Properties = im.Mapping

		if err = json.NewEncoder(buf).Encode(esReq); err != nil {
			return
		}

		index = fmt.Sprintf(reindex.IndexTpl, indexPrefix, im.Index)
		iLog := m.log.With(zap.String("name", index))

		if e := indexMap[index]; e != nil {
			iLog.Info("index exists",
				zap.String("health", e.Health),
				zap.String("status", e.Status),
				zap.String("size", e.StoreSize),
				zap.String("documents", e.DocsCount),
			)

			continue
		}

		if esRsp, err = esc.Indices.Create(index, esc.Indices.Create.WithBody(buf)); esRsp.IsError() || err != nil {
			if err != nil {
				iLog.Error("index creation failed", zap.Error(err))
			}
			if len(esRsp.String()) > 0 {
				iLog.Error(fmt.Sprintf("index creation failed due to %s", esRsp.String()))
			}
			continue
		}

		if err = esRsp.Body.Close(); err != nil {
			return
		}

		iLog.Info("index created")
	}

	return
}

func (m *mapper) mapExistingIndexes(ii []*esIndex) (im map[string]*esIndex) {
	im = make(map[string]*esIndex)
	for _, i := range ii {
		im[i.Name] = i
	}

	return
}

func (m *mapper) getExistingIndexes(ctx context.Context) (ii []*esIndex, err error) {
	var (
		esRsp *esapi.Response
	)

	ii = make([]*esIndex, 100)

	esc, err := m.es.Client()
	if err != nil {
		return
	}

	esRsp, err = esc.Cat.Indices(
		esc.Cat.Indices.WithContext(ctx),
		esc.Cat.Indices.WithFormat("json"),
	)
	if err != nil {
		return
	}
	defer esRsp.Body.Close()

	return ii, json.NewDecoder(esRsp.Body).Decode(&ii)
}

// ConfigurationMapping ...
func (m *mapper) ConfigurationMapping(ctx context.Context) (err error) {
	var (
		buf = &bytes.Buffer{}

		// @todo better naming and ref please parent structs
		configMap = reqMapping{
			Mappings: struct {
				Properties map[string]*property `json:"properties,omitempty"`
			}(struct{ Properties map[string]*property }{Properties: map[string]*property{
				"incrementalIndex": {
					//Type: "object",
					Properties: map[string]*property{
						"updatedAt": {Type: "date"},
					},
				},
			}}),
		}
	)
	const (
		CortezaConfigurationIndex = "corteza-configuration-test"
	)

	index := fmt.Sprintf(CortezaConfigurationIndex)
	esc, err := m.es.Client()
	if err != nil {
		return
	}

	if err = json.NewEncoder(buf).Encode(configMap); err != nil {
		return
	}

	// Check whether a particular configuration index exists
	res, err := esc.Indices.Exists([]string{index}, esc.Indices.Exists.WithContext(ctx))
	if res != nil && err != nil && !res.IsError() {
		if err != nil {
			m.log.Error("Configuration index exist check failed", zap.Error(err))
		} else {
			// Put mapping if exist otherwise create it
			if res, err := esc.Indices.Create(index, esc.Indices.Create.WithBody(buf)); res.IsError() || err != nil {
				if err != nil {
					m.log.Error("Configuration index creation failed", zap.Error(err))
				}
				if len(res.String()) > 0 {
					m.log.Error(fmt.Sprintf("Configuration index creation failed due to %s", res.String()))
				}
			}

			m.log.Info("Configuration index created")
			return
		}
	}

	m.log.Info("Configuration index exist")

	return
}

// @todo: maybe add const for all es mapping property type
// @todo: add helper methods for all Indices & use it instead of esc.Indices.
