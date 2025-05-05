package documents

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/discovery/service"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/options"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	sysService "github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	systemResources struct {
		opt      options.DiscoveryOpt
		settings *types.AppSettings

		rbac interface {
			SignificantRoles(ctx context.Context, res rbac.Resource, op string) (aRR, dRR []uint64, err error)
		}

		ac interface {
			CanReadUser(ctx context.Context, r *types.User) bool
		}
		usr interface {
			Find(context.Context, types.UserFilter) (types.UserSet, types.UserFilter, error)
		}
	}
)

func SystemResources() *systemResources {
	return &systemResources{
		opt:      service.DefaultOption,
		settings: sysService.CurrentSettings,
		rbac:     rbac.Global(),
		ac:       sysService.DefaultAccessControl,
		usr:      sysService.DefaultUser,
	}
}

func (d systemResources) Users(ctx context.Context, limit uint, cur string, userID uint64, deleted uint) (rsp *Response, err error) {
	return rsp, func() (err error) {
		if !d.settings.Discovery.SystemUsers.Enabled {
			return errors.Internal("system user indexing disabled")
		}

		var (
			uu types.UserSet
			f  = types.UserFilter{
				Deleted: filter.State(deleted),
			}
		)

		if userID > 0 {
			f.UserID = append(f.UserID, id.String(userID))
		}

		if f.Paging, err = filter.NewPaging(limit, cur); err != nil {
			return err
		}

		if uu, f, err = d.usr.Find(ctx, f); err != nil {
			return err
		}

		rsp = &Response{
			Documents: make([]Document, len(uu)),
			Filter: Filter{
				Limit:    limit,
				NextPage: f.NextPage,
			},
		}

		for i, u := range uu {
			doc := &docUser{
				ResourceType: "system:user",
				UserID:       u.ID,
				Email:        u.Email,
				Name:         u.Name,
				Handle:       u.Handle,
				Suspended:    u.SuspendedAt,
				Created:      makePartialChange(&u.CreatedAt),
				Updated:      makePartialChange(u.UpdatedAt),
				Deleted:      makePartialChange(u.DeletedAt),
			}
			if len(d.opt.CortezaDomain) > 0 && u.ID > 0 {
				doc.Url = fmt.Sprintf("%s/admin/system/user/edit/%d", d.opt.CortezaDomain, u.ID)
			}

			doc.Security.AllowedRoles, doc.Security.DeniedRoles, err = d.rbac.SignificantRoles(ctx, u, "read")
			if err != nil {
				return
			}

			rsp.Documents[i].ID = u.ID
			rsp.Documents[i].Source = doc
		}

		return nil
	}()
}
