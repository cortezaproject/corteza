package system

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	discoveryType "github.com/cortezaproject/corteza-server/pkg/discovery/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"testing"
)

func (h helper) clearActivityLog() {
	h.noError(store.TruncateResourceActivityLogs(context.Background(), service.DefaultStore))
}

func (h helper) repoMakeActivityLog() *discoveryType.ResourceActivity {
	var res = &discoveryType.ResourceActivity{
		ID:             id.Next(),
		ResourceID:     id.Next(),
		ResourceType:   (types.Record{}).LabelResourceKind(),
		ResourceAction: "create",
	}

	h.a.NoError(store.CreateResourceActivityLog(context.Background(), service.DefaultStore, res))

	return res
}

func TestCreateActivityLog(t *testing.T) {
	h := newHelper(t)
	h.clearActionLog()

	h.repoMakeActivityLog()
}
