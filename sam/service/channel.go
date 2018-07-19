package service

import (
	"context"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	channel struct {
		repository channelRepository
	}

	channelRepository interface {
		FindById(context.Context, uint64) (*types.Channel, error)
		Find(context.Context, *types.ChannelFilter) ([]*types.Channel, error)

		Create(context.Context, *types.Channel) (*types.Channel, error)
		Update(context.Context, *types.Channel) (*types.Channel, error)

		deleter
		archiver
	}
)

func Channel() *channel {
	return &channel{repository: repository.Channel()}
}

func (svc channel) FindById(ctx context.Context, id uint64) (*types.Channel, error) {
	// @todo: permission check if current user can read channel
	return svc.repository.FindById(ctx, id)
}

func (svc channel) Find(ctx context.Context, filter *types.ChannelFilter) ([]*types.Channel, error) {
	// @todo: permission check to return only channels that channel has access to
	// @todo: actual searching not just a full select
	return svc.repository.Find(ctx, filter)
}

func (svc channel) Create(ctx context.Context, mod *types.Channel) (*types.Channel, error) {
	// @todo: topic channelEvent/log entry
	// @todo: channel name cmessage/log entry
	// @todo: permission check if channel can add channel

	return svc.repository.Create(ctx, mod)
}

func (svc channel) Update(ctx context.Context, mod *types.Channel) (*types.Channel, error) {
	// @todo: topic change channelEvent/log entry
	// @todo: channel name change channelEvent/log entry
	// @todo: permission check if current user can edit channel
	// @todo: make sure archived & deleted entries can not be edited
	// @todo: handle channel movinga
	// @todo: handle channel archiving

	return svc.repository.Update(ctx, mod)
}

func (svc channel) Delete(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that channel has been removed (remove from web UI)
	// @todo: permissions check if current user can remove channel
	return svc.repository.Delete(ctx, id)
}

func (svc channel) Archive(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that channel has been removed (remove from web UI)
	// @todo: permissions check if current user can remove channel
	return svc.repository.Archive(ctx, id)
}

func (svc channel) Unarchive(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that channel has been removed (remove from web UI)
	// @todo: permissions check if current user can remove channel
	return svc.repository.Unarchive(ctx, id)
}
