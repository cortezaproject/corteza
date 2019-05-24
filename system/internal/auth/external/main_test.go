package external

import (
	"os"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/internal/logger"
)

func TestMain(m *testing.M) {
	logger.Init(zapcore.DebugLevel)
	os.Exit(m.Run())
}
