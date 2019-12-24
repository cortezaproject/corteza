package messaging

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/app"
)

func TestConfigure(t *testing.T) {
	var _ app.Runnable = &App{}
}
