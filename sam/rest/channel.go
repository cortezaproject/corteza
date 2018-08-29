package rest

import (
	"context"

	"github.com/crusttech/crust/sam/rest/request"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Channel struct {
		svc channelService
	}

	channelService interface {
		FindByID(ctx context.Context, channelID uint64) (*types.Channel, error)
		Find(ctx context.Context, filter *types.ChannelFilter) ([]*types.Channel, error)

		Create(ctx context.Context, channel *types.Channel) (*types.Channel, error)
		Update(ctx context.Context, channel *types.Channel) (*types.Channel, error)

		deleter
		archiver
	}
)

func (Channel) New(channelSvc channelService) *Channel {
	var ctrl = &Channel{}
	ctrl.svc = channelSvc
	return ctrl
}

func (ctrl *Channel) Create(ctx context.Context, r *request.ChannelCreate) (interface{}, error) {
	channel := &types.Channel{
		Name:  r.Name,
		Topic: r.Topic,
	}

	return ctrl.svc.Create(ctx, channel)
}

func (ctrl *Channel) Edit(ctx context.Context, r *request.ChannelEdit) (interface{}, error) {
	channel := &types.Channel{
		Name:  r.Name,
		Topic: r.Topic,
	}

	return ctrl.svc.Update(ctx, channel)

}

func (ctrl *Channel) Delete(ctx context.Context, r *request.ChannelDelete) (interface{}, error) {
	return nil, ctrl.svc.Delete(ctx, r.ChannelID)
}

func (ctrl *Channel) Read(ctx context.Context, r *request.ChannelRead) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, r.ChannelID)
}

func (ctrl *Channel) List(ctx context.Context, r *request.ChannelList) (interface{}, error) {
	return ctrl.svc.Find(ctx, &types.ChannelFilter{Query: r.Query})
}

func (ctrl *Channel) Members(ctx context.Context, r *request.ChannelMembers) (interface{}, error) {
	return nil, nil
}

func (ctrl *Channel) Join(ctx context.Context, r *request.ChannelJoin) (interface{}, error) {
	return nil, nil
}

func (ctrl *Channel) Part(ctx context.Context, r *request.ChannelPart) (interface{}, error) {
	return nil, nil
}

func (ctrl *Channel) Invite(ctx context.Context, r *request.ChannelInvite) (interface{}, error) {
	return nil, nil
}
