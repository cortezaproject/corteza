package corredor

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/options"
	"go.uber.org/zap"
)

// tested with
// go test -count 10 -race -run TestDataRace ./pkg/corredor/...
func TestDataRace(t *testing.T) {
	go Setup(zap.NewNop(), options.CorredorOpt{})
	go Healthcheck(context.Background())
}
