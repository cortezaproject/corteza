package external

import (
	"os"
	"testing"

	"github.com/crusttech/human/server/pkg/logger"
)

func TestMain(m *testing.M) {
	logger.SetDefault(logger.MakeDebugLogger())
	os.Exit(m.Run())
}
