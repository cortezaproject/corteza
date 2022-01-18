package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/gig"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	GigService interface {
		Create(context.Context, string, gig.UpdatePayload) (*gig.Gig, error)
		Read(context.Context, uint64) (*gig.Gig, error)
		Update(context.Context, uint64, gig.UpdatePayload) (*gig.Gig, error)
		AddSources(context.Context, uint64, gig.UpdatePayload) (*gig.Gig, error)
		RemoveSources(context.Context, uint64, uint64) (*gig.Gig, error)
		Delete(context.Context, uint64) error
		Undelete(context.Context, uint64) error

		Prepare(context.Context, uint64) error
		Exec(context.Context, uint64) error
		Complete(context.Context, uint64) error
		Cleanup(context.Context, uint64) error

		Output(context.Context, uint64) (gig.SourceSet, error)
		Status(context.Context, uint64) (gig.WorkerStatus, error)
		State(context.Context, uint64) (gig.WorkerState, error)

		Tasks(context.Context) gig.TaskDefSet
		Workers(context.Context) gig.WorkerDefSet
	}

	gigAccessController interface {
		CanSearchGigs(context.Context) bool
		CanCreateGig(context.Context) bool
		CanReadGig(context.Context, *types.Gig) bool
		CanUpdateGig(context.Context, *types.Gig) bool
		CanDeleteGig(context.Context, *types.Gig) bool
		CanUndeleteGig(context.Context, *types.Gig) bool
		CanExecGig(context.Context, *types.Gig) bool
	}

	gigService struct {
		actionlog actionlog.Recorder

		ac gigAccessController

		store   store.Storer
		manager gig.Service
	}
)

func Gig(m gig.Service) *gigService {
	return &gigService{
		ac:        DefaultAccessControl,
		store:     DefaultStore,
		actionlog: DefaultActionlog,
		manager:   m,
	}
}

func (svc gigService) findByID(ctx context.Context, gigID uint64) (*gig.Gig, error) {
	tmp, err := svc.manager.Read(ctx, gigID)
	if err != nil {
		return nil, err
	}

	return &tmp, err
}

func (svc gigService) wrapResponse(g gig.Gig, err error) (*gig.Gig, error) {
	if err != nil {
		return nil, err
	}

	return &g, nil
}

func (svc gigService) Create(ctx context.Context, worker string, pl gig.UpdatePayload) (g *gig.Gig, err error) {
	var (
		gProps = &gigServiceActionProps{}
	)

	err = func() error {
		if !svc.ac.CanCreateGig(ctx) {
			return GigServiceErrNotAllowedToCreate()
		}

		pl.Worker, err = svc.initWorker(ctx, worker)
		if err != nil {
			return err
		}

		g, err = svc.wrapResponse(svc.manager.Create(ctx, pl))
		if err != nil {
			return err
		}

		gProps.setNew(g.TySystemWrapper())
		return nil
	}()

	return g, svc.recordAction(ctx, gProps, GigServiceActionCreate, err)
}

func (svc gigService) Read(ctx context.Context, gigID uint64) (g *gig.Gig, err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
	)

	err = func() error {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		g, err = svc.wrapResponse(svc.manager.Read(ctx, gigID))
		if err != nil {
			return err
		}

		gProps.setGig(g.TySystemWrapper())

		if !svc.ac.CanReadGig(ctx, g.TySystemWrapper()) {
			return GigServiceErrNotAllowedToRead()
		}

		return err
	}()

	return g, svc.recordAction(ctx, gProps, GigServiceActionLookup, err)
}

func (svc gigService) Update(ctx context.Context, gigID uint64, pl gig.UpdatePayload) (g *gig.Gig, err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
	)

	err = func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}

		gProps.setGig(old.TySystemWrapper())

		if !svc.ac.CanUpdateGig(ctx, old.TySystemWrapper()) {
			return GigServiceErrNotAllowedToUpdate()
		}

		g, err = svc.wrapResponse(svc.manager.Update(ctx, *old, pl))
		if err != nil {
			return
		}

		gProps.setNew(g.TySystemWrapper())

		return
	}()

	return g, svc.recordAction(ctx, gProps, GigServiceActionUpdate, err)
}

func (svc gigService) AddSources(ctx context.Context, gigID uint64, pl gig.UpdatePayload) (g *gig.Gig, err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
	)

	err = func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}

		gProps.setGig(old.TySystemWrapper())

		if !svc.ac.CanUpdateGig(ctx, old.TySystemWrapper()) {
			return GigServiceErrNotAllowedToUpdate()
		}

		g, err = svc.wrapResponse(svc.manager.AddSources(ctx, *old, pl.Sources, pl.Decode...))
		if err != nil {
			return
		}

		gProps.setNew(g.TySystemWrapper())
		return
	}()

	return g, svc.recordAction(ctx, gProps, GigServiceActionUpdate, err)
}

func (svc gigService) RemoveSources(ctx context.Context, gigID, sourceID uint64) (g *gig.Gig, err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
	)

	err = func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}

		gProps.setGig(old.TySystemWrapper())

		if !svc.ac.CanUpdateGig(ctx, old.TySystemWrapper()) {
			return GigServiceErrNotAllowedToUpdate()
		}

		g, err = svc.wrapResponse(svc.manager.RemoveSources(ctx, *old, sourceID))
		if err != nil {
			return
		}

		gProps.setNew(g.TySystemWrapper())
		return
	}()

	return g, svc.recordAction(ctx, gProps, GigServiceActionUpdate, err)
}

func (svc gigService) Delete(ctx context.Context, gigID uint64) (err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
		g      *gig.Gig
	)

	err = (func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		g, err = svc.findByID(ctx, gigID)
		if err != nil {
			return
		}

		gProps.setGig(g.TySystemWrapper())

		if !svc.ac.CanDeleteGig(ctx, g.TySystemWrapper()) {
			return GigServiceErrNotAllowedToDelete()
		}

		g, err = svc.wrapResponse(svc.manager.Delete(ctx, *g))
		if err != nil {
			return
		}

		gProps.setNew(g.TySystemWrapper())
		return nil
	})()

	return svc.recordAction(ctx, gProps, GigServiceActionDelete, err)
}
func (svc gigService) Undelete(ctx context.Context, gigID uint64) (err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
		g      *gig.Gig
	)

	err = (func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		g, err = svc.findByID(ctx, gigID)
		if err != nil {
			return
		}

		gProps.setGig(g.TySystemWrapper())

		if !svc.ac.CanUndeleteGig(ctx, g.TySystemWrapper()) {
			return GigServiceErrNotAllowedToUndelete()
		}

		g, err = svc.wrapResponse(svc.manager.Undelete(ctx, *g))
		if err != nil {
			return
		}

		gProps.setNew(g.TySystemWrapper())
		return nil
	})()

	return svc.recordAction(ctx, gProps, GigServiceActionUndelete, err)
}

func (svc gigService) Prepare(ctx context.Context, gigID uint64) (err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
		g      *gig.Gig
	)

	err = func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}

		gProps.setGig(old.TySystemWrapper())

		if !svc.ac.CanExecGig(ctx, old.TySystemWrapper()) {
			return GigServiceErrNotAllowedToExec()
		}

		g, err = svc.wrapResponse(svc.manager.Prepare(ctx, *old))
		if err != nil {
			return
		}

		gProps.setNew(g.TySystemWrapper())
		return
	}()

	return svc.recordAction(ctx, gProps, GigServiceActionExec, err)
}

func (svc gigService) Exec(ctx context.Context, gigID uint64) (err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
		g      *gig.Gig
	)

	err = func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}

		gProps.setGig(old.TySystemWrapper())

		if !svc.ac.CanExecGig(ctx, old.TySystemWrapper()) {
			return GigServiceErrNotAllowedToExec()
		}

		g, err = svc.wrapResponse(svc.manager.Exec(ctx, *old))
		if err != nil {
			return
		}

		gProps.setNew(g.TySystemWrapper())
		return
	}()

	return svc.recordAction(ctx, gProps, GigServiceActionExec, err)
}

func (svc gigService) Complete(ctx context.Context, gigID uint64) (err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
		g      *gig.Gig
	)

	err = func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}

		gProps.setGig(old.TySystemWrapper())

		if !svc.ac.CanUpdateGig(ctx, old.TySystemWrapper()) {
			return GigServiceErrNotAllowedToUpdate()
		}

		g, err = svc.wrapResponse(svc.manager.Complete(ctx, *old))
		if err != nil {
			return
		}

		gProps.setNew(g.TySystemWrapper())

		return
	}()

	return svc.recordAction(ctx, gProps, GigServiceActionUpdate, err)
}

func (svc gigService) Cleanup(ctx context.Context, gigID uint64) (err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
		g      *gig.Gig
	)

	err = func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}

		gProps.setGig(old.TySystemWrapper())

		if !svc.ac.CanUpdateGig(ctx, old.TySystemWrapper()) {
			return GigServiceErrNotAllowedToUpdate()
		}

		g, err = svc.wrapResponse(svc.manager.Cleanup(ctx, *old))
		if err != nil {
			return
		}

		gProps.setNew(g.TySystemWrapper())

		return
	}()

	return svc.recordAction(ctx, gProps, GigServiceActionUpdate, err)
}

func (svc gigService) Output(ctx context.Context, gigID uint64) (out gig.SourceSet, err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
		g      gig.Gig
	)

	err = func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		g, err = svc.manager.Read(ctx, gigID)
		if err != nil {
			return err
		}

		gProps.setGig(g.TySystemWrapper())

		if !svc.ac.CanReadGig(ctx, g.TySystemWrapper()) {
			return GigServiceErrNotAllowedToRead()
		}

		out, err = svc.manager.Output(ctx, g)
		return
	}()

	return g.Output, svc.recordAction(ctx, gProps, GigServiceActionLookup, err)
}

func (svc gigService) Status(ctx context.Context, gigID uint64) (out gig.WorkerStatus, err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
		g      gig.Gig
	)

	err = func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		g, err = svc.manager.Read(ctx, gigID)
		if err != nil {
			return err
		}

		gProps.setGig(g.TySystemWrapper())

		if !svc.ac.CanReadGig(ctx, g.TySystemWrapper()) {
			return GigServiceErrNotAllowedToRead()
		}

		return
	}()

	return g.Status, svc.recordAction(ctx, gProps, GigServiceActionLookup, err)
}

func (svc gigService) State(ctx context.Context, gigID uint64) (out gig.WorkerState, err error) {
	var (
		gProps = &gigServiceActionProps{gig: &types.Gig{ID: gigID}}
		g      gig.Gig
	)

	err = func() (err error) {
		if gigID == 0 {
			return GigServiceErrInvalidID()
		}

		g, err = svc.manager.Read(ctx, gigID)
		if err != nil {
			return err
		}

		gProps.setGig(g.TySystemWrapper())

		if !svc.ac.CanReadGig(ctx, g.TySystemWrapper()) {
			return GigServiceErrNotAllowedToRead()
		}

		out, err = svc.manager.State(ctx, g)
		if err != nil {
			return err
		}

		return
	}()

	return out, svc.recordAction(ctx, gProps, GigServiceActionLookup, err)
}

func (svc gigService) Tasks(ctx context.Context) (defs gig.TaskDefSet) {
	return svc.manager.Tasks(ctx)
}

func (svc gigService) Workers(ctx context.Context) (defs gig.WorkerDefSet) {
	return svc.manager.Workers(ctx)
}

func (svc gigService) initWorker(ctx context.Context, name string) (w gig.Worker, err error) {
	switch strings.ToLower(name) {
	case gig.WorkerHandleExport:
		return gig.WorkerExport(svc.store), nil
	case gig.WorkerHandleImport:
		return gig.WorkerImport(svc.store), nil
	case gig.WorkerHandleAttachment:
		return nil, fmt.Errorf("worker not implemented: %s", name)
	case gig.WorkerHandleNoop:
		return gig.WorkerNoop(), nil
	}

	err = fmt.Errorf("unknown worker: %s", name)
	return
}
