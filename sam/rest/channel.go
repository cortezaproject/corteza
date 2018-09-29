package rest

import (
	"context"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/outgoing"
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

	return ctrl.svc.ch.With(ctx).Create(channel)
}

func (ctrl *Channel) Edit(ctx context.Context, r *request.ChannelEdit) (interface{}, error) {
	channel := &types.Channel{
		ID:    r.ChannelID,
		Name:  r.Name,
		Topic: r.Topic,
	}

	return ctrl.svc.ch.With(ctx).Update(channel)

}

func (ctrl *Channel) Delete(ctx context.Context, r *request.ChannelDelete) (interface{}, error) {
	return nil, ctrl.svc.ch.With(ctx).Delete(r.ChannelID)
}

func (ctrl *Channel) Read(ctx context.Context, r *request.ChannelRead) (interface{}, error) {
	return ctrl.svc.ch.With(ctx).FindByID(r.ChannelID)
}

func (ctrl *Channel) List(ctx context.Context, r *request.ChannelList) (interface{}, error) {
	return ctrl.svc.ch.With(ctx).Find(&types.ChannelFilter{Query: r.Query})
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

	return ctrl.wrapAttachment(ctrl.svc.att.With(ctx).Create(
		r.ChannelID,
		r.Upload.Filename,
		r.Upload.Size,
		file))
}

func (ctrl *Channel) wrapAttachment(attachment *types.Attachment, err error) (*outgoing.Attachment, error) {
	if err != nil {
		return nil, err
	} else {
		return payload.Attachment(attachment), nil
	}
}
