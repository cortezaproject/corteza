package rest

import (
	"context"

	"github.com/crusttech/crust/sam/rest/server"
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

func (ctrl *Channel) Create(ctx context.Context, r *server.ChannelCreateRequest) (interface{}, error) {
	channel := &types.Channel{
		Name:  r.Name,
		Topic: r.Topic,
	}

	return ctrl.svc.Create(ctx, channel)
}

func (ctrl *Channel) Edit(ctx context.Context, r *server.ChannelEditRequest) (interface{}, error) {
	channel := &types.Channel{
		Name:  r.Name,
		Topic: r.Topic,
	}

	return ctrl.svc.Update(ctx, channel)

}

func (ctrl *Channel) Delete(ctx context.Context, r *server.ChannelDeleteRequest) (interface{}, error) {
	return nil, ctrl.svc.Delete(ctx, r.ChannelID)
}

func (ctrl *Channel) Read(ctx context.Context, r *server.ChannelReadRequest) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, r.ChannelID)
}

func (ctrl *Channel) List(ctx context.Context, r *server.ChannelListRequest) (interface{}, error) {
	return ctrl.svc.Find(ctx, &types.ChannelFilter{Query: r.Query})
}

func (ctrl *Channel) Members(ctx context.Context, r *server.ChannelMembersRequest) (interface{}, error) {
	return nil, nil
}

func (ctrl *Channel) Join(ctx context.Context, r *server.ChannelJoinRequest) (interface{}, error) {
	return nil, nil
}

func (ctrl *Channel) Part(ctx context.Context, r *server.ChannelPartRequest) (interface{}, error) {
	return nil, nil
}

func (ctrl *Channel) Invite(ctx context.Context, r *server.ChannelInviteRequest) (interface{}, error) {
	return nil, nil
}
