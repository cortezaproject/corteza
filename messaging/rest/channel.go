package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/api"

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/payload/outgoing"
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
		Type:  types.ChannelType(r.Type),
		// Due to golang's inability do decode uint64 slice from string slice, we're expecting
		// string input for members (for now)
		// https://github.com/golang/go/issues/15624
		Members: payload.ParseUint64s(r.Members),

		MembershipPolicy: r.MembershipPolicy,
	}

	return ctrl.wrap(ctrl.svc.ch.Create(ctx, channel))
}

func (ctrl *Channel) Update(ctx context.Context, r *request.ChannelUpdate) (interface{}, error) {
	channel := &types.Channel{
		ID:    r.ChannelID,
		Name:  r.Name,
		Topic: r.Topic,
		Type:  types.ChannelType(r.Type),

		MembershipPolicy: r.MembershipPolicy,
	}

	return ctrl.wrap(ctrl.svc.ch.Update(ctx, channel))
}

func (ctrl *Channel) State(ctx context.Context, r *request.ChannelState) (interface{}, error) {
	switch r.State {
	case "delete":
		return ctrl.wrap(ctrl.svc.ch.Delete(ctx, r.ChannelID))
	case "undelete":
		return ctrl.wrap(ctrl.svc.ch.Undelete(ctx, r.ChannelID))
	case "archive":
		return ctrl.wrap(ctrl.svc.ch.Archive(ctx, r.ChannelID))
	case "unarchive":
		return ctrl.wrap(ctrl.svc.ch.Unarchive(ctx, r.ChannelID))
	}

	return nil, nil
}

func (ctrl *Channel) SetFlag(ctx context.Context, r *request.ChannelSetFlag) (interface{}, error) {
	switch r.Flag {
	case "pinned", "hidden", "ignored":
		return ctrl.wrap(ctrl.svc.ch.SetFlag(ctx, r.ChannelID, types.ChannelMembershipFlag(r.Flag)))
	}

	return nil, nil
}

func (ctrl *Channel) RemoveFlag(ctx context.Context, r *request.ChannelRemoveFlag) (interface{}, error) {
	return ctrl.wrap(ctrl.svc.ch.SetFlag(ctx, r.ChannelID, types.ChannelMembershipFlagNone))
}

func (ctrl *Channel) Read(ctx context.Context, r *request.ChannelRead) (interface{}, error) {
	return ctrl.wrap(ctrl.svc.ch.FindByID(ctx, r.ChannelID))
}

func (ctrl *Channel) List(ctx context.Context, r *request.ChannelList) (interface{}, error) {
	return ctrl.wrapSet(ctrl.svc.ch.Find(ctx, types.ChannelFilter{Query: r.Query}))
}

func (ctrl *Channel) Members(ctx context.Context, r *request.ChannelMembers) (interface{}, error) {
	return ctrl.wrapMemberSet(ctrl.svc.ch.FindMembers(ctx, r.ChannelID))
}

func (ctrl *Channel) Invite(ctx context.Context, r *request.ChannelInvite) (interface{}, error) {
	// Due to golang's inability do decode uint64 slice from string slice, we're expecting
	// string input for members (for now)
	// https://github.com/golang/go/issues/15624
	return ctrl.wrapMemberSet(ctrl.svc.ch.InviteUser(ctx, r.ChannelID, payload.ParseUint64s(r.UserID)...))
}

func (ctrl *Channel) Join(ctx context.Context, r *request.ChannelJoin) (interface{}, error) {
	return ctrl.wrapMemberSet(ctrl.svc.ch.AddMember(ctx, r.ChannelID, r.UserID))
}

func (ctrl *Channel) Part(ctx context.Context, r *request.ChannelPart) (interface{}, error) {
	return api.OK(), ctrl.svc.ch.DeleteMember(ctx, r.ChannelID, r.UserID)
}

func (ctrl *Channel) Attach(ctx context.Context, r *request.ChannelAttach) (interface{}, error) {
	file, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	att, err := ctrl.svc.att.CreateMessageAttachment(
		ctx,
		r.Upload.Filename,
		r.Upload.Size,
		file,
		r.ChannelID,
		r.ReplyTo,
	)

	if err != nil {
		return nil, err
	}

	return payload.Attachment(att, auth.GetIdentityFromContext(ctx).Identity()), nil
}

func (ctrl *Channel) wrap(channel *types.Channel, err error) (*outgoing.Channel, error) {
	if err != nil {
		return nil, err
	} else {
		return payload.Channel(channel), nil
	}
}

func (ctrl *Channel) wrapSet(cc types.ChannelSet, f types.ChannelFilter, err error) (*outgoing.ChannelSet, error) {
	if err != nil {
		return nil, err
	} else {
		return payload.Channels(cc), nil
	}
}

func (ctrl *Channel) wrapMemberSet(mm types.ChannelMemberSet, err error) (*outgoing.ChannelMemberSet, error) {
	if err != nil {
		return nil, err
	} else {
		return payload.ChannelMembers(mm), nil
	}
}
