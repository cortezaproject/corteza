package service

import (
	"context"

	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	organisation struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		rpo repository.OrganisationRepository
	}

	OrganisationService interface {
		With(ctx context.Context) OrganisationService

		FindByID(organisationID uint64) (*types.Organisation, error)
		Find(filter *types.OrganisationFilter) ([]*types.Organisation, error)

		Create(organisation *types.Organisation) (*types.Organisation, error)
		Update(organisation *types.Organisation) (*types.Organisation, error)

		Archive(ID uint64) error
		Unarchive(ID uint64) error
		Delete(ID uint64) error
	}
)

func Organisation(ctx context.Context) OrganisationService {
	return (&organisation{
		logger: DefaultLogger.Named("organisation"),
	}).With(ctx)
}

func (svc organisation) With(ctx context.Context) OrganisationService {
	db := repository.DB(ctx)
	return &organisation{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		rpo: repository.Organisation(ctx, db),
	}
}

func (svc organisation) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}
func (svc organisation) FindByID(id uint64) (*types.Organisation, error) {
	// @todo: permission check if current user can read organisation
	return svc.rpo.FindByID(id)
}

func (svc organisation) Find(filter *types.OrganisationFilter) ([]*types.Organisation, error) {
	// @todo: permission check to return only organisations that organisation has access to
	// @todo: actual searching not just a full select
	return svc.rpo.Find(filter)
}

func (svc organisation) Create(mod *types.Organisation) (*types.Organisation, error) {
	// @todo: permission check if current user can add/edit organisation
	// @todo: make sure archived & deleted entries can not be edited

	return svc.rpo.Create(mod)
}

func (svc organisation) Update(mod *types.Organisation) (*types.Organisation, error) {
	// @todo: permission check if current user can add/edit organisation
	// @todo: make sure archived & deleted entries can not be edited

	return svc.rpo.Update(mod)
}

func (svc organisation) Delete(id uint64) error {
	// @todo: permissions check if current user can remove organisation
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	return svc.rpo.DeleteByID(id)
}

func (svc organisation) Archive(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	// @todo: permissions check if current user can archive organisation
	return svc.rpo.ArchiveByID(id)
}

func (svc organisation) Unarchive(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	// @todo: permissions check if current user can unarchive organisation
	return svc.rpo.UnarchiveByID(id)
}

var _ OrganisationService = &organisation{}
