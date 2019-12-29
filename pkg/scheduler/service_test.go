package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
)

func TestMainServiceFunctions(t *testing.T) {
	r := require.New(t)

	const (
		loopInterval = time.Millisecond
		actionWait   = loopInterval * 10
	)

	r.Nil(gScheduler)
	Setup(zap.NewNop(), eventbus.New(), time.Second)
	r.NotNil(gScheduler)
	r.False(gScheduler.Started())
	r.Equal(gScheduler, Service())
	Service().Start(context.Background())
	time.Sleep(actionWait)
	r.True(gScheduler.Started())
	gScheduler.Stop()
	time.Sleep(actionWait)
	r.False(gScheduler.Started())
}
