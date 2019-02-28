package service

import (
	"context"

	internalRules "github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/messaging/repository"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	permissions struct {
		db  db
		ctx context.Context

		prm systemService.PermissionsService
	}

	PermissionsService interface {
		With(context.Context) PermissionsService

		CanAccessMessaging() bool
		CanGrantMessaging() bool
		CanCreatePublicChannel() bool
		CanCreatePrivateChannel() bool
		CanCreateDirectChannel() bool

		CanUpdate(ch *types.Channel) bool
		CanRead(ch *types.Channel) bool
		CanJoin(ch *types.Channel) bool
		CanLeave(ch *types.Channel) bool

		CanManageMembers(ch *types.Channel) bool
		CanManageWebhooks(ch *types.Channel) bool
		CanManageAttachments(ch *types.Channel) bool

		CanSendMessage(ch *types.Channel) bool
		CanReplyMessage(ch *types.Channel) bool
		CanEmbedMessage(ch *types.Channel) bool
		CanAttachMessage(ch *types.Channel) bool
		CanUpdateOwnMessages(ch *types.Channel) bool
		CanUpdateMessages(ch *types.Channel) bool
		CanReactMessage(ch *types.Channel) bool
	}
)

func Permissions() PermissionsService {
	return (&permissions{
		prm: systemService.Permissions(),
	}).With(context.Background())
}

func (p *permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:  db,
		ctx: ctx,

		prm: p.prm.With(ctx),
	}
}

func (p *permissions) CanAccessMessaging() bool {
	return p.checkAccess("messaging", "access")
}

func (p *permissions) CanGrantMessaging() bool {
	return p.checkAccess("messaging", "grant")
}

func (p *permissions) CanCreatePublicChannel() bool {
	return p.checkAccess("messaging", "channel.public.create")
}

func (p *permissions) CanCreatePrivateChannel() bool {
	return p.checkAccess("messaging", "channel.private.create")
}

func (p *permissions) CanCreateDirectChannel() bool {
	return p.checkAccess("messaging", "channel.direct.create")
}

func (p *permissions) CanUpdate(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "update")
}

func (p *permissions) CanRead(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "read")
}

func (p *permissions) CanJoin(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "join")
}

func (p *permissions) CanLeave(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "leave")
}

func (p *permissions) CanManageMembers(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "members.manage")
}

func (p *permissions) CanManageWebhooks(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "webhooks.manage")
}

func (p *permissions) CanManageAttachments(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "attachments.manage")
}

func (p *permissions) CanSendMessage(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "message.send")
}

func (p *permissions) CanReplyMessage(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "message.reply")
}

func (p *permissions) CanEmbedMessage(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "message.embed")
}

func (p *permissions) CanAttachMessage(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "message.attach")
}

func (p *permissions) CanUpdateOwnMessages(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "message.update.own")
}

func (p *permissions) CanUpdateMessages(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "message.update.all")
}

func (p *permissions) CanReactMessage(ch *types.Channel) bool {
	return p.checkAccess(ch.Resource().String(), "message.react")
}

func (p *permissions) checkAccess(resource string, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	access := p.prm.Check(resource, operation, fallbacks...)
	if access == internalRules.Allow {
		return true
	}
	return false
}
