package feed

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza/server/discovery/service"
	"github.com/cortezaproject/corteza/server/discovery/types"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/options"

	"github.com/cortezaproject/corteza/server/pkg/rbac"
)

type (
	resourceActivity struct {
		opt options.DiscoveryOpt

		rbac interface {
			SignificantRoles(ctx context.Context, res rbac.Resource, op string) (aRR, dRR []uint64, err error)
		}

		ac interface {
		}

		resActivity service.ResourceActivityService
	}
)

func ResourceActivity() *resourceActivity {
	return &resourceActivity{
		opt:         service.DefaultOption,
		rbac:        rbac.Global(),
		resActivity: service.DefaultResourceActivity,
	}
}

func (a resourceActivity) ResourceActivities(ctx context.Context, limit uint, cur string, from *time.Time, to *time.Time) (rsp *Response, err error) {
	return rsp, func() (err error) {
		var (
			aa types.ResourceActivitySet

			f = types.ResourceActivityFilter{}
		)

		if from == nil || from.IsZero() {
			return errors.Internal("invalid or missing from timestamp")
		}

		if to == nil || to.IsZero() {
			now := time.Now()
			to = &now
		}

		if from.After(*to) {
			return errors.Internal("invalid from timestamp, it must be before to timestamp")
		}

		f.FromTimestamp = from
		f.ToTimestamp = to

		if f.Paging, err = filter.NewPaging(limit, cur); err != nil {
			return err
		}

		if aa, f, err = a.resActivity.Find(ctx, f); err != nil {
			return err
		}

		rsp = &Response{
			ActivityLogs: make([]ActivityLog, 0),
			Filter: Filter{
				Limit:    limit,
				NextPage: f.NextPage,
			},
		}

		err = aa.Walk(func(a *types.ResourceActivity) error {
			rsp.ActivityLogs = append(rsp.ActivityLogs, ActivityLog{
				ID:             a.ID,
				ResourceID:     a.ResourceID,
				ResourceType:   a.ResourceType,
				ResourceAction: a.ResourceAction,
				Timestamp:      a.Timestamp,
				Meta:           a.Meta,
			})

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	}()
}
