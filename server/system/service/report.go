package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/minions"
	rep "github.com/cortezaproject/corteza-server/pkg/report"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cast"
)

type (
	report struct {
		ac        reportAccessController
		eventbus  eventDispatcher
		actionlog actionlog.Recorder
		store     store.Storer

		users UserService
	}

	reportAccessController interface {
		CanSearchReports(context.Context) bool
		CanCreateReport(context.Context) bool
		CanReadReport(context.Context, *types.Report) bool
		CanUpdateReport(context.Context, *types.Report) bool
		CanDeleteReport(context.Context, *types.Report) bool
		CanRunReport(context.Context, *types.Report) bool
	}
)

var (
	reporters = make(map[string]rep.DatasourceProvider)
)

// Report is a default report service initializer
func Report(s store.Storer, ac reportAccessController, al actionlog.Recorder, eb eventDispatcher) *report {
	return &report{
		store:     s,
		ac:        ac,
		actionlog: al,
		eventbus:  eb,

		users: DefaultUser,
	}
}

func (svc *report) RegisterReporter(key string, r rep.DatasourceProvider) {
	reporters[key] = r
}

func (svc *report) LookupByID(ctx context.Context, ID uint64) (report *types.Report, err error) {
	var (
		aaProps = &reportActionProps{report: &types.Report{ID: ID}}
	)

	err = func() error {
		if ID == 0 {
			return ReportErrInvalidID()
		}

		if report, err = store.LookupReportByID(ctx, svc.store, ID); err != nil {
			return ReportErrInvalidID().Wrap(err)
		}

		if !svc.ac.CanReadReport(ctx, report) {
			return ReportErrNotAllowedToRead()
		}

		return nil
	}()

	return report, svc.recordAction(ctx, aaProps, ReportActionLookup, err)
}

func (svc *report) Search(ctx context.Context, rf types.ReportFilter) (rr types.ReportSet, f types.ReportFilter, err error) {
	var (
		aaProps = &reportActionProps{filter: &rf}
	)

	// For each fetched item, store backend will check if it is valid or not
	rf.Check = func(res *types.Report) (bool, error) {
		if !svc.ac.CanReadReport(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchReports(ctx) {
			return ReportErrNotAllowedToSearch()
		}

		if len(rf.Labels) > 0 {
			rf.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Report{}.LabelResourceKind(),
				rf.Labels,
			)

			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(rf.LabeledIDs) == 0 {
				return nil
			}
		}

		if rr, f, err = store.SearchReports(ctx, svc.store, rf); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledReports(rr)...); err != nil {
			return err
		}

		return nil

	}()

	return rr, f, svc.recordAction(ctx, aaProps, ReportActionSearch, err)
}

func (svc *report) Create(ctx context.Context, new *types.Report) (report *types.Report, err error) {
	var (
		aaProps = &reportActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateReport(ctx) {
			return ReportErrNotAllowedToCreate()
		}

		// if err = svc.eventbus.WaitFor(ctx, event.ReportBeforeCreate(new, nil)); err != nil {
		// 	return
		// }

		// Set new values after beforeCreate events are emitted
		new.ID = nextID()
		new.CreatedAt = *now()

		if new.Meta == nil {
			new.Meta = &types.ReportMeta{}
		}

		new = svc.setIDs(new)

		if err = store.CreateReport(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		report = new

		// _ = svc.eventbus.WaitFor(ctx, event.ReportAfterCreate(new, nil))
		return nil
	}()

	return report, svc.recordAction(ctx, aaProps, ReportActionCreate, err)
}

func (svc *report) Update(ctx context.Context, upd *types.Report) (report *types.Report, err error) {
	var (
		aaProps = &reportActionProps{update: upd}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return ReportErrInvalidID()
		}

		if report, err = store.LookupReportByID(ctx, svc.store, upd.ID); err != nil {
			return
		}

		aaProps.setReport(report)

		if !svc.ac.CanUpdateReport(ctx, report) {
			return ReportErrNotAllowedToUpdate()
		}

		// if err = svc.eventbus.WaitFor(ctx, event.ReportBeforeUpdate(upd, report)); err != nil {
		// 	return
		// }

		// Assign changed values after afterUpdate events are emitted
		report.Handle = upd.Handle
		report.Meta = upd.Meta
		report.Scenarios = upd.Scenarios
		report.Sources = upd.Sources
		report.Blocks = upd.Blocks
		report.UpdatedAt = now()

		if upd.Meta != nil {
			report.Meta = upd.Meta
		}

		report = svc.setIDs(report)

		if err = store.UpdateReport(ctx, svc.store, report); err != nil {
			return err
		}

		if label.Changed(report.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}
			report.Labels = upd.Labels
		}

		// _ = svc.eventbus.WaitFor(ctx, event.ReportAfterUpdate(upd, report))
		return nil
	}()

	return report, svc.recordAction(ctx, aaProps, ReportActionUpdate, err)
}

func (svc *report) Delete(ctx context.Context, ID uint64) (err error) {
	var (
		aaProps = &reportActionProps{}
		report  *types.Report
	)

	err = func() (err error) {
		if ID == 0 {
			return ReportErrInvalidID()
		}

		if report, err = store.LookupReportByID(ctx, svc.store, ID); err != nil {
			return
		}

		aaProps.setReport(report)

		if !svc.ac.CanDeleteReport(ctx, report) {
			return ReportErrNotAllowedToDelete()
		}

		// if err = svc.eventbus.WaitFor(ctx, event.ReportBeforeDelete(nil, report)); err != nil {
		// 	return
		// }

		report.DeletedAt = now()
		if err = store.UpdateReport(ctx, svc.store, report); err != nil {
			return
		}

		// _ = svc.eventbus.WaitFor(ctx, event.ReportAfterDelete(nil, report))
		return nil
	}()

	return svc.recordAction(ctx, aaProps, ReportActionDelete, err)
}

func (svc *report) Undelete(ctx context.Context, ID uint64) (err error) {
	var (
		aaProps = &reportActionProps{}
		report  *types.Report
	)

	err = func() (err error) {
		if ID == 0 {
			return ReportErrInvalidID()
		}

		if report, err = store.LookupReportByID(ctx, svc.store, ID); err != nil {
			return
		}

		aaProps.setReport(report)

		if !svc.ac.CanDeleteReport(ctx, report) {
			return ReportErrNotAllowedToUndelete()
		}

		// if err = svc.eventbus.WaitFor(ctx, event.ReportBeforeUndelete(nil, app)); err != nil {
		// 	return
		// }

		report.DeletedAt = nil
		if err = store.UpdateReport(ctx, svc.store, report); err != nil {
			return
		}

		// _ = svc.eventbus.WaitFor(ctx, event.ReportAfterUndelete(nil, app))
		return nil
	}()

	return svc.recordAction(ctx, aaProps, ReportActionUndelete, err)
}

// actionlog?
func (svc *report) DescribeFresh(ctx context.Context, src types.ReportDataSourceSet, st rep.StepDefinitionSet, sources ...string) (out rep.FrameDescriptionSet, err error) {
	// var (
	// 	aaProps = &reportActionProps{}
	// )

	out = make(rep.FrameDescriptionSet, 0, len(sources)*2)

	err = func() (err error) {
		if !svc.ac.CanCreateReport(ctx) {
			return ReportErrNotAllowedToCreate()
		}

		ss := src.ModelSteps()
		ss = append(ss, st...)

		// Model the report
		model, err := rep.Model(ctx, reporters, ss...)
		if err != nil {
			return
		}

		err = model.Run(ctx)
		if err != nil {
			return
		}

		var auxOut rep.FrameDescriptionSet
		for _, s := range sources {
			auxOut, err = model.Describe(ctx, s)
			if err != nil {
				return err
			}

			out = append(out, auxOut...)
		}

		return nil
	}()

	return out, err
}

func (svc *report) Run(ctx context.Context, reportID uint64, dd rep.FrameDefinitionSet) (out []*rep.Frame, err error) {
	var (
		aaProps = &reportActionProps{}
	)
	out = make([]*rep.Frame, 0, 4)

	err = func() (err error) {
		// @todo evt bus?
		// if err = svc.eventbus.WaitFor(ctx, event.ReportBeforeUpdate(upd, report)); err != nil {
		// 	return
		// }

		if reportID == 0 {
			return ReportErrInvalidID()
		}

		r, err := store.LookupReportByID(ctx, svc.store, reportID)
		if err != nil {
			return err
		}

		// - ac
		if !svc.ac.CanRunReport(ctx, r) {
			return ReportErrNotAllowedToRun()
		}

		ss := r.Sources.ModelSteps()
		ss = append(ss, r.Blocks.ModelSteps()...)

		if err = ss.Validate(); err != nil {
			return ReportErrInvalidConfiguration().Wrap(err)
		}

		for _, d := range dd {
			if err = d.Validate(); err != nil {
				return ReportErrInvalidConfiguration().Wrap(err)
			}
		}

		// Model the report
		model, err := rep.Model(ctx, reporters, ss...)
		if err != nil {
			return
		}

		err = model.Run(ctx)
		if err != nil {
			return
		}

		auxdd := make([]*rep.FrameDefinition, 0, len(dd))
		for i, d := range dd {
			// first one; nothing special needed
			if i == 0 {
				auxdd = append(auxdd, d)
				continue
			}

			stp := model.GetStep(d.Source)
			if stp == nil {
				return fmt.Errorf("unknown source: %s", d.Source)
			}

			// if the current source matches the prev. source, and they both define references,
			// they fall into the same chunk.
			if stp.Def().Join != nil && (d.Source == dd[i-1].Source) && (d.Ref != "" && dd[i-1].Ref != "") && (d.Name == dd[i-1].Name) {
				auxdd = append(auxdd, d)
				continue
			}

			// if the current one doesn't fall into the current chunk, process
			// the chunk and reset it
			ff, err := model.Load(ctx, auxdd...)
			if err != nil {
				return err
			}
			out = append(out, ff...)

			auxdd = make([]*rep.FrameDefinition, 0, len(dd))
			auxdd = append(auxdd, d)
		}

		if len(auxdd) > 0 {
			ff, err := model.Load(ctx, auxdd...)
			if err != nil {
				return err
			}

			ff, err = svc.enhance(ctx, ff)

			out = append(out, ff...)
		}

		// _ = svc.eventbus.WaitFor(ctx, event.ReportAfterUpdate(upd, report))
		// return nil
		return nil
	}()

	return out, svc.recordAction(ctx, aaProps, ReportActionRun, err)
}

// enhance is a temporary function that enriches the output to satisfy some current requirements.
// @todo extend core implementation to support such operatons
//
// - userID is replaced by the user name || username || email || handle || userID
func (svc *report) enhance(ctx context.Context, ff []*rep.Frame) (_ []*rep.Frame, err error) {
	// Preload sys users
	uIndex := make(map[uint64]*types.User)
	uu, uf, err := svc.users.Find(ctx, types.UserFilter{Paging: filter.Paging{Limit: 1024}})
	if err != nil {
		return
	}
	hasMore := uf.NextPage != nil
	for i := range uu {
		uIndex[uu[i].ID] = uu[i]
	}

	for _, f := range ff {
		userCols := make([]int, 0, len(f.Columns))
		for i, c := range f.Columns {
			if c.Kind != "User" {
				continue
			}
			userCols = append(userCols, i)
		}

		for _, r := range f.Rows {
			for _, ci := range userCols {
				col := r[ci]
				if minions.IsNil(col) {
					continue
				}
				switch col.Type() {
				case "ID":
					uID := col.Get().(uint64)
					user, ok := uIndex[uID]
					if !ok && hasMore {
						user, err = svc.users.FindByID(ctx, uID)
						if err != nil && err != store.ErrNotFound {
							return
						}
					}

					if user == nil {
						continue
					} else if _, ok := uIndex[uID]; !ok {
						uIndex[uID] = user
					}

					if usr, ok := uIndex[uID]; ok {
						label := strconv.FormatUint(uID, 10)
						if usr.Name != "" {
							label = usr.Name
						} else if usr.Username != "" {
							label = usr.Username
						} else if usr.Email != "" {
							label = usr.Email
						} else if usr.Handle != "" {
							label = usr.Handle
						}

						r[ci], _ = expr.NewString(label)
					}
				}
			}
		}
	}

	return ff, err
}

func (svc *report) setIDs(r *types.Report) *types.Report {
	// scenarios
	for _, s := range r.Scenarios {
		if s.ScenarioID == 0 {
			s.ScenarioID = nextID()
		}
	}

	// blocks
	for _, b := range r.Blocks {
		if b.BlockID == 0 {
			b.BlockID = nextID()
		}

		// elements
		for _, elRaw := range b.Elements {
			el, ok := elRaw.(map[string]interface{})
			if !ok {
				continue
			}

			elID, ok := el["elementID"]
			sElID := cast.ToString(elID)
			if sElID != "" && sElID != "0" {
				continue
			}
			if cast.ToUint64(elID) != 0 {
				continue
			}

			el["elementID"] = strconv.FormatUint(nextID(), 10)
		}
	}

	return r
}

// toLabeledReports converts to []label.LabeledResource
//
// This function is auto-generated.
func toLabeledReports(set []*types.Report) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
