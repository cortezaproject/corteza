package service

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// attachmentDenyAllAC denies everything
type attachmentDenyAllAC struct{}

func (attachmentDenyAllAC) CanReadNamespace(context.Context, *types.Namespace) bool   { return false }
func (attachmentDenyAllAC) CanUpdateNamespace(context.Context, *types.Namespace) bool { return false }
func (attachmentDenyAllAC) CanCreateNamespace(context.Context) bool                   { return false }
func (attachmentDenyAllAC) CanReadModule(context.Context, *types.Module) bool         { return false }
func (attachmentDenyAllAC) CanReadPage(context.Context, *types.Page) bool             { return false }
func (attachmentDenyAllAC) CanUpdatePage(context.Context, *types.Page) bool           { return false }
func (attachmentDenyAllAC) CanReadRecord(context.Context, *types.Record) bool         { return false }
func (attachmentDenyAllAC) CanUpdateRecord(context.Context, *types.Record) bool       { return false }
func (attachmentDenyAllAC) CanCreateRecordOnModule(context.Context, *types.Module) bool {
	return false
}

func TestAttachmentAccess(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = sqlite.ConnectInMemory(ctx)

		ownerID    = nextID()
		strangerID = nextID()

		ownerCtx    = auth.SetIdentityToContext(ctx, auth.Authenticated(ownerID))
		strangerCtx = auth.SetIdentityToContext(ctx, auth.Authenticated(strangerID))

		nsA = &types.Namespace{Name: "a", Slug: "a", ID: nextID(), CreatedAt: *now()}
		nsB = &types.Namespace{Name: "b", Slug: "b", ID: nextID(), CreatedAt: *now()}
	)

	if err != nil {
		t.Fatalf("failed to init sqlite in-memory db: %v", err)
	}

	if err = store.Upgrade(ctx, zap.NewNop(), s); err != nil {
		t.Fatalf("failed to upgrade store: %v", err)
	}

	if err = s.TruncateComposeNamespaces(ctx); err != nil {
		t.Fatalf("failed to truncate compose namespaces: %v", err)
	}

	if err = s.TruncateComposeAttachments(ctx); err != nil {
		t.Fatalf("failed to truncate compose attachments: %v", err)
	}

	if err = s.TruncateComposePages(ctx); err != nil {
		t.Fatalf("failed to truncate compose pages: %v", err)
	}

	if err = store.CreateComposeNamespace(ctx, s, nsA, nsB); err != nil {
		t.Fatalf("failed to seed namespaces: %v", err)
	}

	var (
		allowAll = &attachment{store: s, ac: &accessControl{rbac: &rbac.ServiceAllowAll{}}}
		denyAll  = &attachment{store: s, ac: attachmentDenyAllAC{}}

		makeAtt = func(t *testing.T, kind string, namespaceID uint64) *types.Attachment {
			att := &types.Attachment{
				ID:          nextID(),
				Kind:        kind,
				Name:        "secret.pdf",
				NamespaceID: namespaceID,
				OwnerID:     ownerID,
				CreatedAt:   *now(),
			}

			require.NoError(t, store.CreateComposeAttachment(ctx, s, att))
			return att
		}
	)

	t.Run("find by id", func(t *testing.T) {
		att := makeAtt(t, types.RecordAttachment, nsA.ID)

		t.Run("own namespace", func(t *testing.T) {
			res, err := allowAll.FindByID(strangerCtx, nsA.ID, att.ID)
			require.NoError(t, err)
			require.Equal(t, att.ID, res.ID)
		})

		t.Run("foreign namespace", func(t *testing.T) {
			_, err := allowAll.FindByID(strangerCtx, nsB.ID, att.ID)
			require.EqualError(t, err, AttachmentErrNotFound().Error())
		})

		t.Run("no namespace", func(t *testing.T) {
			// used by automation, access is still checked on attachment's namespace
			res, err := allowAll.FindByID(strangerCtx, 0, att.ID)
			require.NoError(t, err)
			require.Equal(t, att.ID, res.ID)

			_, err = denyAll.FindByID(strangerCtx, 0, att.ID)
			require.EqualError(t, err, AttachmentErrNotAllowedToReadNamespace().Error())
		})

		t.Run("not allowed to read namespace", func(t *testing.T) {
			_, err := denyAll.FindByID(strangerCtx, nsA.ID, att.ID)
			require.Error(t, err)
			require.Equal(t, AttachmentErrNotAllowedToReadNamespace().Error(), err.Error())
		})

		t.Run("deleted", func(t *testing.T) {
			del := makeAtt(t, types.RecordAttachment, nsA.ID)
			del.DeletedAt = now()
			require.NoError(t, store.UpdateComposeAttachment(ctx, s, del))

			_, err := allowAll.FindByID(strangerCtx, nsA.ID, del.ID)
			require.EqualError(t, err, AttachmentErrNotFound().Error())
		})
	})

	t.Run("delete by id", func(t *testing.T) {
		t.Run("foreign namespace", func(t *testing.T) {
			att := makeAtt(t, types.RecordAttachment, nsA.ID)

			err := allowAll.DeleteByID(strangerCtx, nsB.ID, att.ID)
			require.EqualError(t, err, AttachmentErrNotFound().Error())

			res, err := store.LookupComposeAttachmentByID(ctx, s, att.ID)
			require.NoError(t, err)
			require.Nil(t, res.DeletedAt)
		})

		t.Run("stranger without permissions", func(t *testing.T) {
			att := makeAtt(t, types.RecordAttachment, nsA.ID)

			err := denyAll.DeleteByID(strangerCtx, nsA.ID, att.ID)
			require.Error(t, err)

			// same without the namespace
			err = denyAll.DeleteByID(strangerCtx, 0, att.ID)
			require.Error(t, err)

			res, err := store.LookupComposeAttachmentByID(ctx, s, att.ID)
			require.NoError(t, err)
			require.Nil(t, res.DeletedAt)
		})

		t.Run("owner", func(t *testing.T) {
			att := makeAtt(t, types.RecordAttachment, nsA.ID)

			require.NoError(t, denyAll.DeleteByID(ownerCtx, nsA.ID, att.ID))

			res, err := store.LookupComposeAttachmentByID(ctx, s, att.ID)
			require.NoError(t, err)
			require.NotNil(t, res.DeletedAt)
		})

		t.Run("namespace manager", func(t *testing.T) {
			att := makeAtt(t, types.RecordAttachment, nsA.ID)

			require.NoError(t, allowAll.DeleteByID(strangerCtx, nsA.ID, att.ID))

			res, err := store.LookupComposeAttachmentByID(ctx, s, att.ID)
			require.NoError(t, err)
			require.NotNil(t, res.DeletedAt)
		})
	})

	t.Run("find by page", func(t *testing.T) {
		pg := &types.Page{ID: nextID(), NamespaceID: nsA.ID, Title: "p", CreatedAt: *now()}
		require.NoError(t, store.CreateComposePage(ctx, s, pg))

		f := types.AttachmentFilter{NamespaceID: nsA.ID, Kind: types.PageAttachment, PageID: pg.ID}

		t.Run("allowed to read page", func(t *testing.T) {
			// store does not implement filtering by page, only the access check matters here
			_, _, err := allowAll.Search(strangerCtx, f)
			if err != nil {
				require.NotEqual(t, AttachmentErrNotAllowedToReadPage().Error(), err.Error())
			}
		})

		t.Run("not allowed to read page", func(t *testing.T) {
			_, _, err := denyAll.Search(strangerCtx, f)
			require.Error(t, err)
			require.Equal(t, AttachmentErrNotAllowedToReadPage().Error(), err.Error())
		})
	})
}

func TestAttachmentFindForServing(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = sqlite.ConnectInMemory(ctx)

		nsA = &types.Namespace{Name: "a", Slug: "a", ID: nextID(), CreatedAt: *now()}
		nsB = &types.Namespace{Name: "b", Slug: "b", ID: nextID(), CreatedAt: *now()}
	)

	if err != nil {
		t.Fatalf("failed to init sqlite in-memory db: %v", err)
	}

	if err = store.Upgrade(ctx, zap.NewNop(), s); err != nil {
		t.Fatalf("failed to upgrade store: %v", err)
	}

	if err = s.TruncateComposeNamespaces(ctx); err != nil {
		t.Fatalf("failed to truncate compose namespaces: %v", err)
	}

	if err = s.TruncateComposeAttachments(ctx); err != nil {
		t.Fatalf("failed to truncate compose attachments: %v", err)
	}

	if err = store.CreateComposeNamespace(ctx, s, nsA, nsB); err != nil {
		t.Fatalf("failed to seed namespaces: %v", err)
	}

	var (
		svc = &attachment{store: s, ac: attachmentDenyAllAC{}}
		att = &types.Attachment{ID: nextID(), Kind: types.RecordAttachment, Name: "secret.pdf", NamespaceID: nsA.ID, CreatedAt: *now()}
	)

	require.NoError(t, store.CreateComposeAttachment(ctx, s, att))

	t.Run("matching kind and namespace", func(t *testing.T) {
		res, err := svc.FindForServing(ctx, nsA.ID, types.RecordAttachment, att.ID)
		require.NoError(t, err)
		require.Equal(t, att.ID, res.ID)
	})

	t.Run("spoofed public kind", func(t *testing.T) {
		for _, kind := range []string{types.PageAttachment, types.IconAttachment, types.NamespaceAttachment} {
			_, err := svc.FindForServing(ctx, nsA.ID, kind, att.ID)
			require.EqualError(t, err, AttachmentErrNotFound().Error(), "kind %q", kind)
		}
	})

	t.Run("foreign namespace", func(t *testing.T) {
		_, err := svc.FindForServing(ctx, nsB.ID, types.RecordAttachment, att.ID)
		require.EqualError(t, err, AttachmentErrNotFound().Error())
	})
}
