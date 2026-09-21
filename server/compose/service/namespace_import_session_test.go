package service

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/stretchr/testify/require"
)

// Import sessions hold uploaded namespaces, they can only be run by the user that created them
func TestNamespaceImportRunForeignSession(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = sqlite.ConnectInMemory(ctx)

		ownerID    = nextID()
		strangerID = nextID()
		sessionID  = nextID()
	)

	require.NoError(t, err)

	svc := &namespace{store: s, ac: &accessControl{rbac: &rbac.ServiceAllowAll{}}}

	namespaceSessionStore[sessionID] = namespaceImportSession{SessionID: sessionID, UserID: ownerID}
	defer delete(namespaceSessionStore, sessionID)

	_, err = svc.ImportRun(
		auth.SetIdentityToContext(ctx, auth.Authenticated(strangerID)),
		sessionID,
		&types.Namespace{Name: "stolen"},
	)

	require.EqualError(t, err, NamespaceErrImportSessionNotFound().Error())
	require.Contains(t, namespaceSessionStore, sessionID, "session of another user must stay intact")
}
