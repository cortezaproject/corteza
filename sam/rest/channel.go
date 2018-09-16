package rest

import (
	"context"

	"github.com/crusttech/crust/sam/rest/request"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Channel struct {
		svc struct {
			ch  service.ChannelService
			att service.AttachmentService
		}
	}
)

func (Channel) New() *Channel {
	ctrl := &Channel{}
	ctrl.svc.ch = service.DefaultChannel
	ctrl.svc.att = service.DefaultAttachment

	return ctrl
}

func (ctrl *Channel) Create(ctx context.Context, r *request.ChannelCreate) (interface{}, error) {
	channel := &types.Channel{
		Name:  r.Name,
		Topic: r.Topic,
	}

	return ctrl.svc.ch.Create(ctx, channel)
}

func (ctrl *Channel) Edit(ctx context.Context, r *request.ChannelEdit) (interface{}, error) {
	channel := &types.Channel{
		ID:    r.ChannelID,
		Name:  r.Name,
		Topic: r.Topic,
	}

	return ctrl.svc.ch.Update(ctx, channel)

}

func (ctrl *Channel) Delete(ctx context.Context, r *request.ChannelDelete) (interface{}, error) {
	return nil, ctrl.svc.ch.Delete(ctx, r.ChannelID)
}

func (ctrl *Channel) Read(ctx context.Context, r *request.ChannelRead) (interface{}, error) {
	return ctrl.svc.ch.FindByID(ctx, r.ChannelID)
}

func (ctrl *Channel) List(ctx context.Context, r *request.ChannelList) (interface{}, error) {
	return ctrl.svc.ch.Find(ctx, &types.ChannelFilter{Query: r.Query})
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

func (ctrl *Channel) Attach(ctx context.Context, r *request.ChannelAttach) (interface{}, error) {
	file, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return ctrl.svc.att.Create(
		ctx,
		r.ChannelID,
		r.Upload.Filename,
		r.Upload.Size,
		file)
}
