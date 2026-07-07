package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/locale"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/reporting"
	"github.com/crusttech/human/server/system/types"
	"github.com/modern-go/reflect2"
	"github.com/spf13/cast"
)

// The CRUD skeleton (LookupByID, Search, Create, Update, Delete, Undelete,
// loadReport and toLabeledReports) is generated in report.gen.go from
// system/report.cue.
//
// This file owns the struct, access-controller interface, constructor, the
// before-create / before-update hooks the generated Create / Update call into,
// and the custom methods (Describe, Run and helpers).

type (
	report struct {
		ac        reportAccessController
		eventbus  eventDispatcher
		actionlog actionlog.Recorder
		store     store.Storer
		locale    locale.Locale

		users UserService

		pipelineRunner pipelineRunner
	}

	reportAccessController interface {
		CanSearchReports(context.Context) bool
		CanCreateReport(context.Context) bool
		CanReadReport(context.Context, *types.Report) bool
		CanUpdateReport(context.Context, *types.Report) bool
		CanDeleteReport(context.Context, *types.Report) bool
		CanRunReport(context.Context, *types.Report) bool
	}

	pipelineRunner interface {
		Run(context.Context, dal.Pipeline) (dal.Iterator, error)
		Dryrun(context.Context, dal.Pipeline) error
		FindModel(dal.ModelRef) *dal.Model
	}
)

// Report is a default report service initializer
func Report(s store.Storer, ac reportAccessController, al actionlog.Recorder, eb eventDispatcher) *report {
	return &report{
		store:     s,
		ac:        ac,
		actionlog: al,
		eventbus:  eb,
		locale:    locale.Global(),

		users: DefaultUser,

		pipelineRunner: dal.Service(),
	}
}

// beforeCreate runs after the access check and before the generated Create
// assigns the ID / timestamps and persists. It defaults the Meta and assigns
// IDs to nested scenarios / blocks / elements.
func (svc *report) beforeCreate(ctx context.Context, new *types.Report) error {
	if new.Meta == nil {
		new.Meta = &types.ReportMeta{}
	}

	svc.setIDs(new)
	return nil
}

// beforeUpdate runs after the stale-data guard and before the generated Update
// copies the mutable fields onto the loaded record. We assign IDs to nested
// scenarios / blocks / elements on the incoming resource so the subsequent
// field copy carries them onto the persisted record.
func (svc *report) beforeUpdate(ctx context.Context, upd, existing *types.Report) error {
	svc.setIDs(upd)
	return nil
}

func (svc *report) onDescribe(ctx context.Context, _ *reportActionProps, src types.ReportDataSourceSet, st types.ReportStepSet, sources ...string) (out []reporting.FrameDescription, err error) {
	out = make([]reporting.FrameDescription, 0, len(sources)*2)

	ss := src.ReportSteps()
	ss = append(ss, st...)

	out, err = reporting.Describe(ctx, svc.pipelineRunner, ss, sources)
	return out, err
}

func (svc *report) onRun(ctx context.Context, aProps *reportActionProps, reportID uint64, dd reporting.FrameDefinitionSet) (out []*reporting.Frame, err error) {
	var (
		iter dal.Iterator
		ff   []*reporting.Frame
	)
	out = make([]*reporting.Frame, 0, 4)

	r, err := loadReport(ctx, svc.store, reportID)
	if err != nil {
		return out, err
	}

	if !svc.ac.CanRunReport(ctx, r) {
		return out, ReportErrNotAllowedToRun()
	}

	ss := r.Sources.ReportSteps()
	ss = append(ss, r.Blocks.ReportSteps()...)

	runs, err := reporting.Runs(svc.pipelineRunner, ss, dd)
	if err != nil {
		return out, err
	}

	// @todo this can be ran in paralel
	for _, run := range runs {
		err = func() (err error) {
			iter, err = svc.pipelineRunner.Run(ctx, run.Pipeline)
			if err != nil {
				return
			}
			defer iter.Close()

			ff, err = reporting.Frames(ctx, iter, run)
			if err != nil {
				return
			}

			err = svc.enhance(ctx, ff)
			if err != nil {
				return
			}

			out = append(out, ff...)
			return
		}()

		if err != nil {
			return out, err
		}
	}

	return out, nil
}

// enhance is a temporary function that enriches the output to satisfy some current requirements.
// @todo extend core implementation to support such operatons
//
// - userID is replaced by the user name || username || email || handle || userID
func (svc *report) enhance(ctx context.Context, ff []*reporting.Frame) (err error) {
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

	var uID uint64
	for _, f := range ff {
		userCols := make([]int, 0, len(f.Columns))
		for i, c := range f.Columns {
			// Translate system columns
			if c.System {
				pp := strings.Split(c.Name, ".")
				c.Label = svc.locale.T(ctx, "compose", fmt.Sprintf("field.system.%s", pp[len(pp)-1]))
				f.Columns[i] = c
			}

			// Collect user columns to replace IDs with labels
			if c.Kind != "User" {
				continue
			}
			userCols = append(userCols, i)
		}

		for _, r := range f.Rows {
			for _, ci := range userCols {
				col := r[ci]
				if reflect2.IsNil(col) {
					continue
				}
				uID, err = cast.ToUint64E(col)
				if err != nil {
					continue
				}

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
					r[ci] = strconv.FormatUint(uID, 10)
					if usr.Name != "" {
						r[ci] = usr.Name
					} else if usr.Username != "" {
						r[ci] = usr.Username
					} else if usr.Email != "" {
						r[ci] = usr.Email
					} else if usr.Handle != "" {
						r[ci] = usr.Handle
					}
				}
			}
		}
	}

	return err
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
