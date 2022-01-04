package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/gig"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	GigService interface {
		Create(context.Context, string, gig.UpdatePayload) (*gig.Gig, error)
		Read(context.Context, uint64) (*gig.Gig, error)
		Update(context.Context, uint64, gig.UpdatePayload) (*gig.Gig, error)
		AddSources(context.Context, uint64, gig.UpdatePayload) (*gig.Gig, error)
		RemoveSources(context.Context, uint64, uint64) (*gig.Gig, error)

		Prepare(context.Context, uint64) error
		Exec(context.Context, uint64) error
		Complete(context.Context, uint64) error
		Cleanup(context.Context, uint64) error

		Output(context.Context, uint64) (gig.SourceSet, error)
		Status(context.Context, uint64) (gig.WorkerStatus, error)

		Tasks(context.Context) gig.TaskDefSet
	}

	gigService struct {
		actionlog actionlog.Recorder

		ac userAccessController

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
	// var (
	// 	uaProps = &userActionProps{user: &types.User{ID: userID}}
	// )

	err = func() error {
		// if !svc.ac.CanReadUser(ctx, u) {
		// 	return UserErrNotAllowedToRead()
		// }

		// @todo access control
		// @todo actionlog?

		pl.Worker, err = svc.initWorker(ctx, worker)
		if err != nil {
			return err
		}

		g, err = svc.wrapResponse(svc.manager.Create(ctx, pl))
		return err
	}()

	return g, err // svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc gigService) Read(ctx context.Context, gigID uint64) (g *gig.Gig, err error) {
	// var (
	// 	uaProps = &userActionProps{user: &types.User{ID: userID}}
	// )

	err = func() error {
		// if !svc.ac.CanReadUser(ctx, u) {
		// 	return UserErrNotAllowedToRead()
		// }

		// @todo access control
		// @todo actionlog?

		g, err = svc.wrapResponse(svc.manager.Read(ctx, gigID))
		return err
	}()

	return g, err // svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc gigService) Update(ctx context.Context, gigID uint64, pl gig.UpdatePayload) (g *gig.Gig, err error) {
	// var (
	// 	uaProps = &userActionProps{user: &types.User{ID: userID}}
	// )

	err = func() (err error) {
		// if !svc.ac.CanUpdateGig(ctx, u) {
		// 	return UserErrNotAllowedToUpdate()
		// }

		// @todo access control
		// @todo actionlog?
		old, err := svc.findByID(ctx, gigID)
		g, err = svc.wrapResponse(svc.manager.Update(ctx, *old, pl))
		return
	}()

	return g, err // svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc gigService) AddSources(ctx context.Context, gigID uint64, pl gig.UpdatePayload) (g *gig.Gig, err error) {
	// var (
	// 	uaProps = &userActionProps{user: &types.User{ID: userID}}
	// )

	err = func() (err error) {
		// if !svc.ac.CanUpdateGig(ctx, u) {
		// 	return UserErrNotAllowedToUpdate()
		// }

		// @todo access control
		// @todo actionlog?

		old, err := svc.findByID(ctx, gigID)
		g, err = svc.wrapResponse(svc.manager.AddSources(ctx, *old, pl.Sources, pl.Decode...))
		return
	}()

	return g, err // svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc gigService) RemoveSources(ctx context.Context, gigID, sourceID uint64) (g *gig.Gig, err error) {
	// var (
	// 	uaProps = &userActionProps{user: &types.User{ID: userID}}
	// )

	err = func() (err error) {
		// if !svc.ac.CanUpdateGig(ctx, u) {
		// 	return UserErrNotAllowedToUpdate()
		// }

		// @todo access control
		// @todo actionlog?

		old, err := svc.findByID(ctx, gigID)
		g, err = svc.wrapResponse(svc.manager.RemoveSources(ctx, *old, sourceID))
		return
	}()

	return g, err // svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc gigService) Prepare(ctx context.Context, gigID uint64) (err error) {
	// var (
	// 	uaProps = &userActionProps{user: &types.User{ID: userID}}
	// )

	err = func() (err error) {
		// if !svc.ac.CanUpdateGig(ctx, u) {
		// 	return UserErrNotAllowedToUpdate()
		// }

		// @todo access control
		// @todo actionlog?

		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}
		_, err = svc.wrapResponse(svc.manager.Prepare(ctx, *old))
		return
	}()

	return err // svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc gigService) Exec(ctx context.Context, gigID uint64) (err error) {
	// var (
	// 	uaProps = &userActionProps{user: &types.User{ID: userID}}
	// )

	err = func() (err error) {
		// if !svc.ac.CanUpdateGig(ctx, u) {
		// 	return UserErrNotAllowedToUpdate()
		// }

		// @todo access control
		// @todo actionlog?

		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}
		_, err = svc.wrapResponse(svc.manager.Exec(ctx, *old))
		return
	}()

	return err // svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc gigService) Complete(ctx context.Context, gigID uint64) (err error) {
	// var (
	// 	uaProps = &userActionProps{user: &types.User{ID: userID}}
	// )

	err = func() (err error) {
		// if !svc.ac.CanUpdateGig(ctx, u) {
		// 	return UserErrNotAllowedToUpdate()
		// }

		// @todo access control
		// @todo actionlog?

		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}
		_, err = svc.wrapResponse(svc.manager.Complete(ctx, *old))
		return
	}()

	return err // svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc gigService) Cleanup(ctx context.Context, gigID uint64) (err error) {
	// var (
	// 	uaProps = &userActionProps{user: &types.User{ID: userID}}
	// )

	err = func() (err error) {
		// if !svc.ac.CanUpdateGig(ctx, u) {
		// 	return UserErrNotAllowedToUpdate()
		// }

		// @todo access control
		// @todo actionlog?

		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}
		_, err = svc.wrapResponse(svc.manager.Cleanup(ctx, *old))
		return
	}()

	return err // svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc gigService) Output(ctx context.Context, gigID uint64) (out gig.SourceSet, err error) {
	// var (
	// 	uaProps = &userActionProps{user: &types.User{ID: userID}}
	// )

	var g *gig.Gig
	err = func() (err error) {
		// if !svc.ac.CanUpdateGig(ctx, u) {
		// 	return UserErrNotAllowedToUpdate()
		// }

		// @todo access control
		// @todo actionlog?
		old, err := svc.findByID(ctx, gigID)
		if err != nil {
			return err
		}
		out, err = svc.manager.Output(ctx, *old)
		return
	}()

	return g.Output, err // svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc gigService) Status(ctx context.Context, gigID uint64) (out gig.WorkerStatus, err error) {
	// var (
	// 	uaProps = &userActionProps{user: &types.User{ID: userID}}
	// )

	var g *gig.Gig
	err = func() (err error) {
		// if !svc.ac.CanUpdateGig(ctx, u) {
		// 	return UserErrNotAllowedToUpdate()
		// }

		// @todo access control
		// @todo actionlog?

		g, err = svc.findByID(ctx, gigID)
		return
	}()

	return g.Status, err // svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc gigService) Tasks(ctx context.Context) (defs gig.TaskDefSet) {
	return svc.manager.Tasks(ctx)
}

func (svc gigService) initWorker(ctx context.Context, name string) (w gig.Worker, err error) {
	switch strings.ToLower(name) {
	case gig.WorkerHandleEnvoy:
		return gig.WorkerImport(svc.store), nil
	case gig.WorkerHandleAttachment:
		return nil, fmt.Errorf("worker not implemented: %s", name)
	case gig.WorkerHandleNoop:
		return gig.WorkerNoop(), nil
	}

	err = fmt.Errorf("unknown worker: %s", name)
	return
}
