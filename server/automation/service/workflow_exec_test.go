package service

import (
	"context"
	"sync"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/stretchr/testify/require"
)

func TestWorkflowExecNotLoaded(t *testing.T) {
	var (
		req = require.New(t)
		svc = &workflow{cache: map[uint64]*wfCacheItem{}, muxCache: &sync.RWMutex{}}
	)

	_, _, err := svc.exec(context.Background(), &types.Workflow{ID: 1}, types.WorkflowExecParams{})
	req.EqualError(err, "workflow is disabled or deleted and cannot run")
}
