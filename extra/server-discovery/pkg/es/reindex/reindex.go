package reindex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/extra/server-discovery/pkg/options"
	"github.com/davecgh/go-spew/spew"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/jmoiron/sqlx/types"
	"go.uber.org/zap"
)

type (
	docsSources struct {
		endpoint string
		index    string
		action   string
		params   map[string]string
		callback func(*document)
	}

	rspDiscoveryDocuments struct {
		Error *struct {
			Message string
		}
		Response *struct {
			Filter struct {
				NextPage string
			}

			Documents []*document
		}
	}

	// auxiliary struct for parsing indexable documents from Discovery API
	document struct {
		ID     string
		Index  string
		Source json.RawMessage
	}

	feedSources struct {
		endpoint string
		index    string
		action   string
		params   map[string]string
		callback func(*document)
	}

	feedResponse struct {
		Error *struct {
			Message string
		}
		Response *struct {
			Filter struct {
				NextPage string
			}

			ActivityLogs []ActivityLog `json:"activityLogs"`
		}
	}

	ActivityLog struct {
		ID             uint64         `json:"activityID,string"`
		ResourceID     uint64         `json:"resourceID,string"`
		ResourceType   string         `json:"resourceType"`
		ResourceAction string         `json:"resourceAction"`
		Timestamp      time.Time      `json:"timestamp"`
		Meta           types.JSONText `json:"meta"`
	}

	ActivityLogMeta struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		ModuleID    uint64 `json:"moduleID,string"`
	}

	esService interface {
		BulkIndexer() (esutil.BulkIndexer, error)
	}

	embedderService interface {
		GenerateEmbeddings(input string) ([]float64, error)
	}

	apiClientService interface {
		HttpClient() *http.Client
		Mappings() (*http.Request, error)
		Feed(url.Values) (*http.Request, error)
		Resources(string, url.Values) (*http.Request, error)
		Request(string) (*http.Request, error)
		Authenticate() error
	}

	reIndexer struct {
		log            *zap.Logger
		esOpt          options.EsOpt
		es             esService
		esClient       *elasticsearch.Client
		embedder       embedderService
		api            apiClientService
		assureMappings func(context.Context) error
	}
)

const (
	IndexTpl    = "corteza-%s-%s"
	MarkerIndex = "corteza_indexer_state"
	MarkerID    = "last_index_time"
)

func ReIndexer(log *zap.Logger, es esService, esc *elasticsearch.Client, api apiClientService, embedderSvc embedderService, esOpt options.EsOpt, assureMappings func(context.Context) error) *reIndexer {
	return &reIndexer{
		log:            log,
		esOpt:          esOpt,
		es:             es,
		esClient:       esc,
		embedder:       embedderSvc,
		api:            api,
		assureMappings: assureMappings,
	}
}

func (ri *reIndexer) ReindexAll(ctx context.Context, esb esutil.BulkIndexer, indexPrefix string) error {
	var (
		srcQueue = make(chan *docsSources, 100)
		bErr     = ri.reindexManager(ctx, esb, indexPrefix, srcQueue)
	)

	srcQueue <- &docsSources{
		endpoint: "/system/users",
		index:    "system-users",
	}

	postProcModules := func(namespaceID string) func(d *document) {
		return func(d *document) {
			srcQueue <- &docsSources{
				endpoint: fmt.Sprintf("/compose/namespaces/%s/modules/%s/records", namespaceID, d.ID),
				index:    fmt.Sprintf("compose-records-%s-%s", namespaceID, d.ID),
			}
		}
	}

	postProcNamespaces := func(d *document) {
		srcQueue <- &docsSources{
			endpoint: fmt.Sprintf("/compose/namespaces/%s/modules", d.ID),
			index:    "compose-modules",
			callback: postProcModules(d.ID),
		}
	}

	_ = postProcModules
	_ = postProcNamespaces

	srcQueue <- &docsSources{
		endpoint: "/compose/namespaces",
		index:    "compose-namespaces",
		callback: postProcNamespaces,
	}
	_ = fmt.Errorf("blocking error")
	return <-bErr
}

func (ri *reIndexer) reindexManager(ctx context.Context, esb esutil.BulkIndexer, indexPrefix string, srcQueue chan *docsSources) chan error {
	var qErr = make(chan error)
	const maxQueueLen = 3

	go func(esb esutil.BulkIndexer) {
		var (
			pQueueLen        = -1
			pQueueStaleCount int

			ticker = time.NewTicker(time.Second)
		)

		defer ticker.Stop()
		defer func() {
			qErr <- nil
		}()

		for {
			select {
			case <-ctx.Done():
				if ctx.Err() != context.Canceled {
					ri.log.Error(ctx.Err().Error())
				} else {
					ri.log.Info("stopped")
				}
				return

			case ds := <-srcQueue:
				if ds == nil {
					// graceful termination
					ri.log.Info("done")
					return
				}

				err := ri.reindex(ctx, esb, indexPrefix, ds)
				if err != nil {
					ri.log.Error("failed to reindex", zap.Error(err), zap.String("endpoint", ds.endpoint))
					return
				}

			case <-ticker.C:
				if pQueueLen != len(srcQueue) {
					pQueueStaleCount = maxQueueLen
				} else {
					pQueueStaleCount--
				}

				if pQueueStaleCount <= 0 {
					ri.log.Info("idle")
					return
				}

				pQueueLen = len(srcQueue)

				//esb, err := ri.es.BulkIndexer()
				//if err != nil {
				//	qErr <- err
				//}

				s := esb.Stats()
				ri.log.Debug("batch indexing stats",
					zap.Uint64("added", s.NumAdded),
					zap.Uint64("created", s.NumCreated),
					zap.Uint64("updated", s.NumUpdated),
					zap.Uint64("deleted", s.NumDeleted),
					zap.Uint64("flushed", s.NumFlushed),
					zap.Uint64("failed", s.NumFailed),
					zap.Uint64("indexed", s.NumIndexed),
					zap.Uint64("requests", s.NumRequests),
					zap.Int("queue length", pQueueLen),
				)
			}
		}
	}(esb)

	println("returning")
	return qErr
}

func (ri *reIndexer) reindex(ctx context.Context, esb esutil.BulkIndexer, indexPrefix string, ds *docsSources) (err error) {
	var (
		qs     = url.Values{"limit": []string{"500"}}
		req    *http.Request
		rsp    *http.Response
		cursor string
	)

	//esb, err := ri.es.BulkIndexer()
	//if err != nil {
	//	return fmt.Errorf("failed to prepare bulk indexer: %w", err)
	//}

	for {
		rspPayload := &rspDiscoveryDocuments{}

		if cursor != "" {
			// set new cursor and update source URL
			qs.Set("pageCursor", cursor)
		}

		if req, err = ri.api.Resources(ds.endpoint, qs); err != nil {
			return fmt.Errorf("failed to prepare resource request: %w", err)
		}

		if rsp, err = ri.api.HttpClient().Do(req.WithContext(ctx)); err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}

		if rsp.StatusCode != http.StatusOK {
			return fmt.Errorf("request resulted in an unexpected status '%s' for url '%s'", rsp.Status, req.URL)
		}

		//{
		//	d, err := httputil.DumpRequestOut(req, true)
		//	println(string(d))
		//	spew.Dump(err)
		//}
		//{
		//	d, err := httputil.DumpResponse(rsp, true)
		//	println(string(d))
		//	spew.Dump(err)
		//}

		if err = json.NewDecoder(rsp.Body).Decode(rspPayload); err != nil {
			return fmt.Errorf("failed to decode reindexing response: %w", err)
		}

		if err = rsp.Body.Close(); err != nil {
			return fmt.Errorf("failed to close reindexing response body: %w", err)
		}

		var docs int
		if rspPayload.Error != nil {
			ri.log.Debug("skipping",
				zap.String("index", fmt.Sprintf(IndexTpl, indexPrefix, ds.index)),
				zap.String("error", rspPayload.Error.Message),
			)
			return
		} else if rspPayload.Response != nil {
			docs = len(rspPayload.Response.Documents)
		}

		ri.log.Debug("reindexing",
			zap.Int("docs", docs),
			zap.String("index", fmt.Sprintf(IndexTpl, indexPrefix, ds.index)),
		)

		if docs == 0 {
			return
		}

		for _, doc := range rspPayload.Response.Documents {
			body := bytes.NewBuffer(doc.Source)
			err = ri.processResource(body)
			if err != nil {
				ri.log.Error("Failed to process record's for embeddings: ", zap.Error(err))
			}

			err = esb.Add(ctx, esutil.BulkIndexerItem{
				Index:      fmt.Sprintf(IndexTpl, indexPrefix, ds.index),
				Action:     "index",
				DocumentID: doc.ID,
				Body:       body,
				OnFailure: func(ctx context.Context, req esutil.BulkIndexerItem, rsp esutil.BulkIndexerResponseItem, err error) {
					spew.Dump(req)
					spew.Dump(rsp)
					spew.Dump(err)
				},
			})

			if err != nil {
				return err
			}

			if ds.callback != nil {
				go ds.callback(doc)
			}
		}

		cursor = rspPayload.Response.Filter.NextPage
		if rspPayload.Response.Filter.NextPage == "" {
			break
		}
	}

	//if err = esb.Close(ctx); err != nil {
	//	return fmt.Errorf("failed to close bulk indexer: %w", err)
	//}

	return nil
}

func (ri *reIndexer) feedReindex(ctx context.Context, esb esutil.BulkIndexer, indexPrefix string, ds *feedSources) (err error) {
	var (
		qs     = url.Values{"limit": []string{"500"}}
		req    *http.Request
		rsp    *http.Response
		cursor string
	)

	if ds == nil {
		ri.log.Debug("invalid resource for feed update")
		return
	}

	err = ri.assureMappings(ctx)
	if err != nil {
		return fmt.Errorf("failed to assure mappings: %w", err)
	}

	for {
		rspPayload := &rspDiscoveryDocuments{}

		if cursor != "" {
			// set new cursor and update source URL
			qs.Set("pageCursor", cursor)
		}

		if ds.params != nil {
			var (
				val string
				ok  bool
			)

			val, ok = ds.params["userID"]
			if ok {
				qs.Set("userID", val)
			}

			val, ok = ds.params["namespaceID"]
			if ok {
				qs.Set("namespaceID", val)
			}

			val, ok = ds.params["moduleID"]
			if ok {
				qs.Set("moduleID", val)
			}

			val, ok = ds.params["recordID"]
			if ok {
				qs.Set("recordID", val)
			}

			val, ok = ds.params["deleted"]
			if ok {
				qs.Set("deleted", val)
			}
		}

		if req, err = ri.api.Resources(ds.endpoint, qs); err != nil {
			return fmt.Errorf("failed to prepare resource request: %w", err)
		}

		if rsp, err = ri.api.HttpClient().Do(req.WithContext(ctx)); err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}

		if rsp.StatusCode != http.StatusOK {
			return fmt.Errorf("request resulted in an unexpected status '%s' for url '%s'", rsp.Status, req.URL)
		}

		//{
		//	d, err := httputil.DumpRequestOut(req, true)
		//	println(string(d))
		//	spew.Dump(err)
		//}
		//{
		//	d, err := httputil.DumpResponse(rsp, true)
		//	println(string(d))
		//	spew.Dump(err)
		//}

		if err = json.NewDecoder(rsp.Body).Decode(rspPayload); err != nil {
			return fmt.Errorf("failed to decode reindexing response: %w", err)
		}

		if err = rsp.Body.Close(); err != nil {
			return fmt.Errorf("failed to close reindexing response body: %w", err)
		}

		var (
			docs     int
			docIndex = fmt.Sprintf(IndexTpl, indexPrefix, ds.index)
		)
		if rspPayload.Error != nil {
			ri.log.Debug("skipping",
				zap.String("index", docIndex),
				zap.String("error", rspPayload.Error.Message),
			)
			return
		} else if rspPayload.Response != nil {
			docs = len(rspPayload.Response.Documents)
		}

		ri.log.Debug("feed reindexing",
			zap.Int("docs", docs),
			zap.String("index", docIndex),
		)

		if docs == 0 {
			return
		}

		action := "index"
		if len(ds.action) > 0 {
			action = ds.action
		}

		for _, doc := range rspPayload.Response.Documents {
			esbItem := esutil.BulkIndexerItem{
				Index:      docIndex,
				Action:     action,
				DocumentID: doc.ID,
				OnFailure: func(ctx context.Context, req esutil.BulkIndexerItem, rsp esutil.BulkIndexerResponseItem, err error) {
					spew.Dump(req)
					spew.Dump(rsp)
					spew.Dump(err)
				},
			}
			if action != "delete" {
				esbItem.Action = "index"
				body := bytes.NewBuffer(doc.Source)
				err = ri.processResource(body)

				esbItem.Body = body
			}

			err = esb.Add(ctx, esbItem)
			if err != nil {
				return err
			}

			if ds.callback != nil {
				go ds.callback(doc)
			}
		}

		cursor = rspPayload.Response.Filter.NextPage
		if rspPayload.Response.Filter.NextPage == "" {
			break
		}
	}

	return nil
}

func (ri *reIndexer) feedReindexManager(ctx context.Context, esb esutil.BulkIndexer, indexPrefix string, feedQueue chan *feedSources) chan error {
	var qErr = make(chan error)
	const maxQueueLen = 3

	go func(esb esutil.BulkIndexer) {
		var (
			pQueueLen        = -1
			pQueueStaleCount int

			ticker = time.NewTicker(time.Second)
		)

		defer ticker.Stop()
		defer func() {
			qErr <- nil
		}()

		for {
			select {
			case <-ctx.Done():
				if ctx.Err() != context.Canceled {
					ri.log.Error(ctx.Err().Error())
				} else {
					ri.log.Info("feed changes stopped")
				}
				return

			case ds := <-feedQueue:
				if ds == nil {
					// graceful termination
					ri.log.Info("feed changes done")
					return
				}

				err := ri.feedReindex(ctx, esb, indexPrefix, ds)
				if err != nil {
					ri.log.Error("failed to reindex", zap.Error(err), zap.String("endpoint", ds.endpoint))
					return
				}

			case <-ticker.C:
				if pQueueLen != len(feedQueue) {
					pQueueStaleCount = maxQueueLen
				} else {
					pQueueStaleCount--
				}

				if pQueueStaleCount <= 0 {
					ri.log.Info("feed changes idle")
					return
				}

				pQueueLen = len(feedQueue)

				s := esb.Stats()
				ri.log.Debug("feed batch indexing stats",
					zap.Uint64("added", s.NumAdded),
					zap.Uint64("created", s.NumCreated),
					zap.Uint64("updated", s.NumUpdated),
					zap.Uint64("deleted", s.NumDeleted),
					zap.Uint64("flushed", s.NumFlushed),
					zap.Uint64("failed", s.NumFailed),
					zap.Uint64("indexed", s.NumIndexed),
					zap.Uint64("requests", s.NumRequests),
					zap.Int("queue length", pQueueLen),
				)
			}
		}
	}(esb)

	println("feed changes returning")
	return qErr
}

func (ri *reIndexer) feedReindexChanges(ctx context.Context, esb esutil.BulkIndexer, indexPrefix string, als []ActivityLog) error {
	var (
		updateQueue = make(chan *feedSources, 100)
		bErr        = ri.feedReindexManager(ctx, esb, indexPrefix, updateQueue)

		duplicateMap = make(map[string]ActivityLog)
	)

	for _, al := range als {
		activityType := fmt.Sprintf("%d-%s", al.ResourceID, al.ResourceAction)
		if val, ok := duplicateMap[activityType]; ok {
			continue
		} else {
			duplicateMap[activityType] = val
		}

		action := al.ResourceAction
		switch al.ResourceType {
		case "system:user":
			updateQueue <- &feedSources{
				endpoint: "/system/users",
				index:    "system-users",
				action:   action,
				params: map[string]string{
					"userID":  fmt.Sprintf("%d", al.ResourceID),
					"deleted": "1",
				},
			}
			break

		case "compose:namespace":
			updateQueue <- &feedSources{
				endpoint: "/compose/namespaces",
				index:    "compose-namespaces",
				action:   action,
				params: map[string]string{
					"namespaceID": fmt.Sprintf("%d", al.ResourceID),
					"deleted":     "1",
				},
			}
			break

		case "compose:module":
			var meta ActivityLogMeta
			err := al.Meta.Unmarshal(&meta)
			if err != nil {
				return err
			}

			updateQueue <- &feedSources{
				endpoint: fmt.Sprintf("/compose/namespaces/%d/modules", meta.NamespaceID),
				index:    "compose-modules",
				action:   action,
				params: map[string]string{
					"moduleID": fmt.Sprintf("%d", al.ResourceID),
					"deleted":  "1",
				},
			}
			break

		case "compose:record":
			var meta ActivityLogMeta
			err := al.Meta.Unmarshal(&meta)
			if err != nil {
				return err
			}

			updateQueue <- &feedSources{
				endpoint: fmt.Sprintf("/compose/namespaces/%d/modules/%d/records", meta.NamespaceID, meta.ModuleID),
				index:    fmt.Sprintf("compose-records-%d-%d", meta.NamespaceID, meta.ModuleID),
				action:   action,
				params: map[string]string{
					"recordID": fmt.Sprintf("%d", al.ResourceID),
					"deleted":  "1",
				},
			}
			break

		default:
			break
		}

	}

	return <-bErr
}

func (ri *reIndexer) Watch(ctx context.Context) {
	var now time.Time

	isFirst, err := ri.IsFirstRunWithMarker(ctx)
	if err != nil {
		ri.log.Warn(fmt.Sprintf("failed to check first run indexes marker: %s", err))
		isFirst = false
	}

	if isFirst {
		now = time.Now().AddDate(0, -ri.esOpt.IndexBackFillMonths, 0)
		ri.log.Info(fmt.Sprintf("first run detected: starting from %d months ago (%s)",
			ri.esOpt.IndexBackFillMonths, now.UTC().Format(time.RFC3339)))

		defer func() {
			if err := ri.SaveLastIndexTime(ctx); err != nil {
				ri.log.Error(fmt.Sprintf("failed to save indexer state to [corteza_indexer_state]: %s", err))
			}
		}()
	} else {
		now = time.Now()
		ri.log.Info("continuing reindexing from current time")
	}

	timeOut := ri.esOpt.IndexInterval
	ticker := time.NewTicker(time.Second * time.Duration(timeOut))

	go func() {
		defer ticker.Stop()
		processing := false

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if processing {
					ri.log.Warn("skipping feed changes reindexing: already processing")
					continue
				}

				processing = true
				go func() {
					defer func() {
						processing = false
					}()

					esb, err := ri.es.BulkIndexer()
					if err != nil {
						ri.log.Error(fmt.Sprintf("failed to prepare bulk indexer for feed changes: %s", err))
						return
					}

					var (
						req *http.Request
						rsp *http.Response

						qs = url.Values{"from": []string{now.UTC().Format(time.RFC3339)}}
					)

					feeds := &feedResponse{}

					// store time before making request
					tmpTime := time.Now()

					if req, err = ri.api.Feed(qs); err != nil {
						ri.log.Error(fmt.Sprintf("failed to prepare feed request: %s", err))
						return
					}

					if rsp, err = ri.api.HttpClient().Do(req.WithContext(ctx)); err != nil {
						ri.log.Error(fmt.Sprintf("failed to send feed request: %s", err))
						return
					}

					if rsp.StatusCode != http.StatusOK {
						ri.log.Error(fmt.Sprintf("request resulted in an unexpected status '%s' for feed", rsp.Status))
						return
					}

					if err = json.NewDecoder(rsp.Body).Decode(feeds); err != nil {
						ri.log.Error(fmt.Sprintf("failed to decode feed resources: %s", err))
						return
					}

					if err = rsp.Body.Close(); err != nil {
						ri.log.Error(fmt.Sprintf("failed to close feed response body: %s", err))
						return
					}

					// Update time after successful request
					now = tmpTime

					if feeds != nil && feeds.Response != nil && len(feeds.Response.ActivityLogs) > 0 {
						err = ri.feedReindexChanges(ctx, esb, "private", feeds.Response.ActivityLogs)
						if err != nil {
							ri.log.Error(fmt.Sprintf("failed to update indexes for feed changes: %s", err))
							return
						}
					} else {
						ri.log.Debug(fmt.Sprintf("No feed changes since last %d seconds; current time: %s", timeOut, now.UTC().String()))
					}

					_ = esb.Close(ctx)
				}()
			}
		}
	}()
}

func (ri *reIndexer) processResource(source *bytes.Buffer) error {
	var data map[string]any
	if err := json.Unmarshal(source.Bytes(), &data); err != nil {
		return fmt.Errorf("Error unmarshaling: %v\n", err)
	}

	resourceType, ok := data["resourceType"].(string)
	if !ok {
		return fmt.Errorf("resourceType not found or not a string")
	}

	// only generate embeddings for compose:record for now
	if resourceType == "compose:record" {
		catchAll, ok := data["catch_all"].([]interface{})
		if !ok {
			return fmt.Errorf("records catch_all not found or not an array")
		}

		embeddings, err := ri.embedder.GenerateEmbeddings(convertValuesToString(catchAll))
		if err != nil {
			return err
		}

		data["vectorsValue"] = embeddings

		updatedJSON, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("Error marshaling: %v\n", err)
		}

		source.Reset()
		source.Write(updatedJSON)
	}

	return nil
}

func (ri *reIndexer) IsFirstRunWithMarker(ctx context.Context) (bool, error) {
	res, err := ri.esClient.Get(
		MarkerIndex,
		MarkerID,
		ri.esClient.Get.WithContext(ctx),
	)

	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	if res.IsError() {
		// marker index does'nt exist
		if res.StatusCode == http.StatusNotFound {
			return true, err
		}
	}

	return false, nil
}

func (ri *reIndexer) SaveLastIndexTime(ctx context.Context) error {
	doc := map[string]interface{}{
		"last_run":    time.Now().UTC().Format(time.RFC3339),
		"initialized": true,
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal marker document: %w", err)
	}

	// create/update the marker document
	res, err := ri.esClient.Index(
		MarkerIndex,
		bytes.NewReader(docBytes),
		ri.esClient.Index.WithContext(ctx),
		ri.esClient.Index.WithDocumentID(MarkerID),
		ri.esClient.Index.WithRefresh("true"),
	)

	if err != nil {
		return fmt.Errorf("failed to save marker index document: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to save marker index document: %s", res.Status())
	}

	return nil
}

func convertValuesToString(values []interface{}) string {
	var parts []string

	for _, item := range values {
		switch val := item.(type) {
		case string:
			parts = append(parts, val)
		case float64:
			parts = append(parts, fmt.Sprintf("%.0f", val))
		default:
			parts = append(parts, fmt.Sprintf("%v", val))
		}
	}

	return strings.Join(parts, " ")
}
