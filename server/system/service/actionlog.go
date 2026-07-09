package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	actionlogReport struct {
		ac        actionlogReportAccessController
		actionlog actionlog.Recorder
		store     store.Storer
	}

	actionlogReportAccessController interface {
		CanReadActionLog(context.Context) bool
	}
)

func ActionlogReport() *actionlogReport {
	return &actionlogReport{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc actionlogReport) Run(ctx context.Context, rr *actionlog.ReportRequest) (*actionlog.ReportResult, error) {
	if err := svc.checkScope(ctx, scope.CapRead); err != nil {
		return nil, err
	}

	if !svc.ac.CanReadActionLog(ctx) {
		return nil, errors.Unauthorized("not allowed to read action log")
	}

	if rr.FromTimestamp == nil || rr.ToTimestamp == nil {
		return nil, errors.InvalidData("report requires both from and to timestamps")
	}

	if rr.ToTimestamp.Before(*rr.FromTimestamp) {
		return nil, errors.InvalidData("report time range is inverted")
	}

	if err := rr.Normalize(); err != nil {
		return nil, errors.InvalidData("invalid report request: %v", err)
	}

	set, err := svc.actionlog.Report(ctx, *rr)
	if err != nil {
		return nil, err
	}

	labels, err := svc.actorLabels(ctx, set)
	if err != nil {
		return nil, err
	}

	return &actionlog.ReportResult{
		Dimensions:  rr.Dimensions,
		Metrics:     rr.Metrics,
		ActorLabels: labels,
		Set:         set,
	}, nil
}

func (svc actionlogReport) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

// Rows grouped by actor carry only the ID; resolve them to a human-readable
// label. Direct store access is safe here — protected by action-log.read.
func (svc actionlogReport) actorLabels(ctx context.Context, set actionlog.ReportRowSet) (map[string]string, error) {
	var actorIDs []string

	for _, row := range set {
		if actorID, ok := row.Dimensions["actor"].(string); ok && actorID != "" && actorID != "0" {
			actorIDs = append(actorIDs, actorID)
		}
	}

	if len(actorIDs) == 0 {
		return nil, nil
	}

	uu, _, err := store.SearchUsers(ctx, svc.store, types.UserFilter{
		UserID:    actorIDs,
		Deleted:   filter.StateInclusive,
		Suspended: filter.StateInclusive,
	})

	if err != nil {
		return nil, err
	}

	labels := make(map[string]string, len(uu))
	for _, u := range uu {
		label := u.Name
		if label == "" {
			label = u.Email
		}

		labels[id.String(u.ID)] = label
	}

	return labels, nil
}
