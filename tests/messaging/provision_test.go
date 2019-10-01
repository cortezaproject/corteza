package messaging

import (
	"testing"

	"github.com/cortezaproject/corteza-server/messaging/importer"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	provision "github.com/cortezaproject/corteza-server/provision/messaging"
)

func TestProvisioning(t *testing.T) {
	h := newHelper(t)
	ctx := auth.SetSuperUserContext(h.secCtx())

	readers, err := impAux.ReadStatic(provision.Asset)
	h.a.NoError(err)
	h.a.NoError(importer.Import(ctx, readers...))
}
