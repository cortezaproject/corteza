package rest

import (
	"context"
	"strconv"
	"strings"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/filter"
	pkgid "github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
	"github.com/go-oauth2/oauth2/v4/errors"
)

type (
	Actionlog struct {
		actionSvc actionlog.Recorder
		reportSvc actionlogReportService
		store     store.Users
		ac        actionlogAccessController
	}

	actionlogReportService interface {
		Run(ctx context.Context, rr *actionlog.ReportRequest) (*actionlog.ReportResult, error)
	}

	actionlogAccessController interface {
		CanReadActionLog(ctx context.Context) bool
	}

	// Extend actionlog.Action so we can
	// provide user's email
	actionlogActionPayload struct {
		*actionlog.Action
		Actor string                 `json:"actor,omitempty"`
		Meta  map[string]interface{} `json:"meta,omitempty"`
	}

	actionlogPayload struct {
		Filter actionlog.Filter          `json:"filter"`
		Set    []*actionlogActionPayload `json:"set"`
	}
)

func (Actionlog) New() *Actionlog {
	return &Actionlog{
		actionSvc: service.DefaultActionlog,
		reportSvc: service.DefaultActionlogReport,
		ac:        service.DefaultAccessControl,
		store:     service.DefaultStore,
	}
}

func (ctrl *Actionlog) List(ctx context.Context, r *request.ActionlogList) (interface{}, error) {
	if !ctrl.ac.CanReadActionLog(ctx) {
		return nil, errors.ErrAccessDenied
	}

	var (
		err error
		f   = actionlog.Filter{
			FromTimestamp:  r.From,
			ToTimestamp:    r.To,
			BeforeActionID: r.BeforeActionID,
			ActorID:        r.ActorID,
			Resource:       r.Resource,
			Action:         r.Action,
			Origin:         r.Origin,
			Limit:          r.Limit,
		}
	)

	ee, f, err := ctrl.actionSvc.Find(ctx, f)

	return ctrl.makeFilterPayload(ctx, ee, f, err)
}

func (ctrl *Actionlog) Report(ctx context.Context, r *request.ActionlogReport) (interface{}, error) {
	rr := &actionlog.ReportRequest{
		Dimensions: r.Dimensions,
		Metrics:    r.Metrics,
		Filter: actionlog.Filter{
			FromTimestamp: r.From,
			ToTimestamp:   r.To,
			ActorID:       r.ActorID,
			Resource:      r.Resource,
			Action:        r.Action,
			Origin:        r.Origin,
			Limit:         r.Limit,
		},
	}

	result, err := ctrl.reportSvc.Run(ctx, rr)
	if err != nil {
		return nil, err
	}

	mapReportSeverityNames(result)
	return result, nil
}

// severity comes from DB as a raw int; map to name so the payload is human-readable (Decision J)
func mapReportSeverityNames(result *actionlog.ReportResult) {
	if result == nil {
		return
	}
	for _, row := range result.Set {
		v, ok := row.Dimensions["severity"]
		if !ok {
			continue
		}
		var n int64
		switch sv := v.(type) {
		case int64:
			n = sv
		case int:
			n = int64(sv)
		case float64:
			n = int64(sv)
		default:
			continue
		}
		row.Dimensions["severity"] = actionlog.Severity(n).String()
	}
}

func (ctrl Actionlog) makeFilterPayload(ctx context.Context, ee []*actionlog.Action, f actionlog.Filter, err error) (*actionlogPayload, error) {
	if err != nil {
		return nil, err
	}

	var (
		pp = make([]*actionlogActionPayload, len(ee))
	)

	// Remap events to payload structs
	for e := range ee {
		pp[e] = &actionlogActionPayload{Action: ee[e]}
		pp[e].Meta = ee[e].Meta

		sanitizeMapStringInterface(pp[e].Meta)
	}

	err = userPreloader(
		ctx,
		ctrl.store,
		func(c chan uint64) {
			for e := range ee {
				c <- ee[e].ActorID
			}

			close(c)
		},
		types.UserFilter{
			Deleted:   filter.StateInclusive,
			Suspended: filter.StateInclusive,
		},
		func(u *types.User) error {
			for p := range pp {
				if pp[p].ActorID == u.ID {
					pp[p].Actor = u.Name
					if pp[p].Actor == "" {
						pp[p].Actor = u.Email
					}
				}
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &actionlogPayload{Filter: f, Set: pp}, nil
}

// Making sure JS can read old entries that were encoded
// as numeric. When decoded from the JSON stored in the DB
// they are of type float64. We'll cast them and store as uint64
//
// We do not want to do that to any numeric values just for the keys ending with
// "ID" or "By"
func sanitizeMapStringInterface(m map[string]interface{}) {
	for k := range m {
		switch v := m[k].(type) {
		case float64:
			if strings.HasSuffix(k, "ID") || strings.HasSuffix(k, "By") {
				// make sure uint64 values on fields ending with ID
				// are properly encoded as strings
				m[k] = strconv.FormatUint(uint64(v), 10)
			}

		case map[string]interface{}:
			sanitizeMapStringInterface(v)
		}
	}
}

// Preloader collects all ids of users, loads them and sets them back
//
// We'll be accessing the store directly since this is protected with action-log.read operation check.
//
// @todo move action log collection and user merging to dedicated service
func userPreloader(ctx context.Context, s store.Users, g func(chan uint64), f types.UserFilter, collect func(*types.User) error) error {
	var (
		// channel that will collect the IDs in the getter
		ch = make(chan uint64, 0)

		// unique index for IDs
		unq = make(map[uint64]bool)
	)

	// Reset the collection of the IDs
	f.UserID = make([]string, 0)

	// Call getter and collect the IDs
	go g(ch)

rangeLoop:
	for {
		select {
		case <-ctx.Done():
			close(ch)
			break rangeLoop
		case id, ok := <-ch:
			if !ok {
				// Channel closed
				break rangeLoop
			}

			if !unq[id] {
				unq[id] = true
				f.UserID = append(f.UserID, pkgid.String(id))
			}
		}

	}

	// Load all users (even if deleted, suspended) from the given list of IDs
	uu, _, err := store.SearchUsers(ctx, s, f)

	if err != nil {
		return err
	}

	return uu.Walk(collect)
}
