package rest

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
)

var _ = errors.Wrap

type (
	Channel struct {
		service channelService
	}

	channelService interface {
		FindById(context.Context, uint64) (*types.Channel, error)
		Find(context.Context, *types.ChannelFilter) ([]*types.Channel, error)

		Create(context.Context, *types.Channel) (*types.Channel, error)
		Update(context.Context, *types.Channel) (*types.Channel, error)

		deleter
		archiver
	}
)

func (Channel) New() *Channel {
	return &Channel{service: service.Channel()}
}

func (ctrl *Channel) Create(ctx context.Context, r *server.ChannelCreateRequest) (interface{}, error) {
	channel := types.Channel{}.
		New().
		SetName(r.Name).
		SetTopic(r.Topic).
		SetMeta([]byte("{}")).
		SetID(factory.Sonyflake.NextID())

	return ctrl.service.Create(ctx, channel)
}

func (ctrl *Channel) Edit(ctx context.Context, r *server.ChannelEditRequest) (interface{}, error) {
	channel := types.Channel{}.
		New().
		SetName(r.Name).
		SetTopic(r.Topic)

	return ctrl.service.Update(ctx, channel)

}

func (ctrl *Channel) Delete(ctx context.Context, r *server.ChannelDeleteRequest) (interface{}, error) {
	return nil, ctrl.service.Delete(ctx, r.ChannelId)
}

func (ctrl *Channel) Read(ctx context.Context, r *server.ChannelReadRequest) (interface{}, error) {
	return ctrl.service.FindById(ctx, r.ChannelId)
}

func (ctrl *Channel) List(ctx context.Context, r *server.ChannelListRequest) (interface{}, error) {
	return ctrl.service.Find(ctx, &types.ChannelFilter{Query: r.Query})
}
