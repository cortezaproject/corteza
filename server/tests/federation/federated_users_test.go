package federation

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/crusttech/human/server/federation/service"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	st "github.com/crusttech/human/server/system/types"
)

// Records written by federation are left out of what a node serves, so the
// users it writes as are found whatever the caller may see
func TestFederatedUserIDs(t *testing.T) {
	h := newHelper(t)

	var (
		ctx       = context.Background()
		federated = &st.User{ID: id.Next(), Email: fmt.Sprintf("%d@federation.corteza", id.Next()), CreatedAt: time.Now()}
		// a handle is the user's to change, so it decides nothing
		person = &st.User{ID: id.Next(), Handle: fmt.Sprintf("federation_%d", id.Next()), Email: "person@example.tld", CreatedAt: time.Now()}
	)

	h.noError(store.CreateUser(ctx, service.DefaultStore, federated, person))

	// a caller with no permissions at all
	ids, err := service.FederatedUserIDs(auth.SetIdentityToContext(ctx, &st.User{ID: id.Next()}), service.DefaultStore)
	h.noError(err)

	h.a.Contains(ids, federated.ID)
	h.a.Contains(ids, auth.FederationUser().ID)
	h.a.NotContains(ids, person.ID)
}
