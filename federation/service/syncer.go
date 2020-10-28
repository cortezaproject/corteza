package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/federation/types"
)

type (
	Syncer struct {
		Service  SharedModuleService
		CService service.RecordService
	}

	AuxResponseSet struct {
		Response struct {
			Filter struct {
				NodeID          string
				ComposeModuleID string
				Query           string
				Limit           int
				NextPage        string
				PrevPage        string
			} `json:"filter"`
		} `json:"response"`
	}
)

func (h *Syncer) Queue(url types.SyncerURI, out chan types.SyncerURI) {
	out <- url
}

func (h *Syncer) Fetch(ctx context.Context, url string) (io.Reader, error) {
	client := http.Client{
		Timeout: time.Duration(3) * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("404")
	}

	return resp.Body, nil
}

func (h *Syncer) Process(ctx context.Context, payload []byte, nodeID uint64, out chan types.SyncerURI, url types.SyncerURI, handler *Syncer, fn func(ctx context.Context, payload []byte, nodeID uint64, handler *Syncer) error) error {
	aux := AuxResponseSet{}
	err := json.Unmarshal(payload, &aux)

	if err != nil {
		return err
	}

	if aux.Response.Filter.NextPage != "" {
		url.NextPage = aux.Response.Filter.NextPage

		// out <- url

		// sync data
		// out <- fmt.Sprintf("%s/federation/exposed/modules/196342359342989002/records/?limit=%d&pageCursor=%s", "http://localhost:8084", aux.Response.Filter.Limit, aux.Response.Filter.NextPage)
	}

	return fn(ctx, payload, nodeID, handler)
}

func NewSyncer() *Syncer {
	return &Syncer{
		Service:  DefaultSharedModule,
		CService: service.DefaultRecord,
	}
}
