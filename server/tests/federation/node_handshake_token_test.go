package federation

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/federation/service"
	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
)

// Nodes without a generated pair token must not accept handshakes with an empty token
func TestNodeHandshakeInitEmptyPairToken(t *testing.T) {
	h := newHelper(t)
	h.clearNodes()

	n := &types.Node{
		ID:        id.Next(),
		Name:      "node",
		BaseURL:   "https://example.tld/federation",
		Status:    types.NodeStatusPending,
		CreatedAt: time.Now(),
	}
	h.noError(store.CreateFederationNode(context.Background(), service.DefaultStore, n))

	err := service.DefaultNode.HandshakeInit(context.Background(), n.ID, "", n.ID, "attacker-auth-token")
	h.a.Error(err)

	stored := h.lookupNodeByID(n.ID)
	h.a.Empty(stored.AuthToken)
	h.a.Equal(types.NodeStatusPending, stored.Status)
}
