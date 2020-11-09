package service

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/stretchr/testify/require"
)

type (
	testProcesser struct{}
)

func (p *testProcesser) Process(ctx context.Context, payload []byte) (int, error) {
	return 0, nil
}

func TestSyncer_process(t *testing.T) {
	var (
		req    = require.New(t)
		syncer = &Syncer{}

		ctx = context.Background()
	)

	c := make(chan Url, 2)
	u := types.SyncerURI{
		Limit:    10,
		NextPage: "123",
	}
	tp := &testProcesser{}

	go syncer.Process(ctx, []byte(`{"response":{"filter":{"limit":1, "nextPage":"456"}}}`), c, u, tp)

	select {
	case url := <-c:
		req.Equal(url.Url.NextPage, "456")
		break
	}
}

func TestSyncer_parseHeader(t *testing.T) {
	var (
		req    = require.New(t)
		syncer = &Syncer{}

		ctx = context.Background()
	)

	c := make(chan Url, 2)
	u := types.SyncerURI{
		Limit:    10,
		NextPage: "123",
	}
	tp := &testProcesser{}

	n, err := syncer.Process(ctx, []byte(`{"response":{"filt`), c, u, tp)

	req.EqualError(err, "unexpected end of JSON input")
	req.Equal(n, 0)
}
