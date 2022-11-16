package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cortezaproject/corteza/server/federation/types"
)

type (
	Syncer struct {
		client http.Client
	}

	Url struct {
		Url  types.SyncerURI
		Meta Processer
	}
	Payload struct {
		Payload io.Reader
		Meta    Processer
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

const FederationUserToken string = "authToken"

func (h *Syncer) Queue(url Url, out chan Url) {
	out <- url
}

func (h *Syncer) Fetch(ctx context.Context, url string) (io.Reader, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if authToken := ctx.Value(FederationUserToken); authToken != nil {
		req.Header.Add("Authorization", `Bearer `+authToken.(string))
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("invalid return status: %d", resp.StatusCode))
	}

	return resp.Body, nil
}

func (h *Syncer) Process(ctx context.Context, payload []byte, out chan Url, url types.SyncerURI, processer Processer) (ProcesserResponse, error) {
	aux, err := h.ParseHeader(ctx, payload)

	if err != nil {
		return 0, err
	}

	if aux.Response.Filter.NextPage != "" {
		url.NextPage = aux.Response.Filter.NextPage

		out <- Url{
			Url:  url,
			Meta: processer,
		}
	}

	return processer.Process(ctx, payload)
}

func (h *Syncer) ParseHeader(ctx context.Context, payload []byte) (aux AuxResponseSet, err error) {
	err = json.Unmarshal(payload, &aux)
	return
}

func NewSyncer() *Syncer {
	return &Syncer{
		client: http.Client{
			Timeout: time.Duration(3) * time.Second,
		},
	}
}
