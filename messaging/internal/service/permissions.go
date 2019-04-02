package service

import (
	"context"

	"github.com/crusttech/crust/internal/auth"
	internalRules "github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/messaging/internal/repository"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	permissions struct {
		db  db
		ctx context.Context

		rules systemService.RulesService
	}

	resource interface {
		PermissionResource() internalRules.Resource
	}

	PermissionsService interface {
		With(context.Context) PermissionsService

		Effective() (ee []effectivePermission, err error)

		CanAccess() bool
		CanGrant() bool
		CanCreatePublicChannel() bool
		CanCreatePrivateChannel() bool
		CanCreateGroupChannel() bool

		CanUpdateChannel(ch *types.Channel) bool
		CanReadChannel(ch *types.Channel) bool
		CanJoinChannel(ch *types.Channel) bool
		CanLeaveChannel(ch *types.Channel) bool
		CanDeleteChannel(ch *types.Channel) bool
		CanUndeleteChannel(ch *types.Channel) bool
		CanArchiveChannel(ch *types.Channel) bool
		CanUnarchiveChannel(ch *types.Channel) bool

		CanManageChannelMembers(ch *types.Channel) bool
		CanManageChannelWebhooks(ch *types.Channel) bool
		CanManageChannelAttachments(ch *types.Channel) bool

		CanSendMessage(ch *types.Channel) bool
		CanReplyMessage(ch *types.Channel) bool
		CanEmbedMessage(ch *types.Channel) bool
		CanAttachMessage(ch *types.Channel) bool
		CanUpdateOwnMessages(ch *types.Channel) bool
		CanUpdateMessages(ch *types.Channel) bool
		CanDeleteOwnMessages(ch *types.Channel) bool
		CanDeleteMessages(ch *types.Channel) bool
		CanReactMessage(ch *types.Channel) bool
	}

	effectivePermission struct {
		Resource  internalRules.Resource `json:"resource"`
		Operation string                 `json:"operation"`
		Allow     bool                   `json:"allow"`
	}
)

func Permissions() PermissionsService {
	return (&permissions{
		rules: systemService.DefaultRules,
	}).With(context.Background())
}

func (p *permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:  db,
		ctx: ctx,

		rules: systemService.Rules(ctx),
	}
}

func (p *permissions) Effective() (ee []effectivePermission, err error) {
	ep := func(res internalRules.Resource, op string, allow bool) effectivePermission {
		return effectivePermission{
			Resource:  res,
			Operation: op,
			Allow:     allow,
		}
	}

	ee = append(ee, ep(types.PermissionResource, "access", p.CanAccess()))
	ee = append(ee, ep(types.PermissionResource, "grant", p.CanGrant()))
	ee = append(ee, ep(types.PermissionResource, "channel.public.create", p.CanCreatePublicChannel()))
	ee = append(ee, ep(types.PermissionResource, "channel.private.create", p.CanCreatePrivateChannel()))
	ee = append(ee, ep(types.PermissionResource, "channel.group.create", p.CanCreateGroupChannel()))

	return
}

func (p *permissions) CanAccess() bool {
	return p.checkAccess(types.PermissionResource, "access")
}

func (p *permissions) CanGrant() bool {
	return p.checkAccess(types.PermissionResource, "grant")
}

func (p *permissions) CanCreatePublicChannel() bool {
	return p.checkAccess(types.PermissionResource, "channel.public.create", p.allow())
}

func (p *permissions) CanCreatePrivateChannel() bool {
	return p.checkAccess(types.PermissionResource, "channel.private.create", p.allow())
}

func (p *permissions) CanCreateGroupChannel() bool {
	return p.checkAccess(types.PermissionResource, "channel.group.create", p.allow())
}

func (p *permissions) CanUpdateChannel(ch *types.Channel) bool {
	return p.checkAccess(ch, "update", p.isChannelOwnerFallback(ch))
}

func (p *permissions) CanReadChannel(ch *types.Channel) bool {
	return p.checkAccess(ch, "read", p.canReadFallback(ch))
}

func (p *permissions) CanJoinChannel(ch *types.Channel) bool {
	return p.checkAccess(ch, "join", p.canJoinFallback(ch))
}

func (p *permissions) CanLeaveChannel(ch *types.Channel) bool {
	return p.checkAccess(ch, "leave", p.allow())
}

func (p *permissions) CanArchiveChannel(ch *types.Channel) bool {
	return p.checkAccess(ch, "archive", p.isChannelOwnerFallback(ch))
}

func (p *permissions) CanUnarchiveChannel(ch *types.Channel) bool {
	return p.checkAccess(ch, "unarchive", p.isChannelOwnerFallback(ch))
}

func (p *permissions) CanDeleteChannel(ch *types.Channel) bool {
	return p.checkAccess(ch, "delete", p.isChannelOwnerFallback(ch))
}

func (p *permissions) CanUndeleteChannel(ch *types.Channel) bool {
	return p.checkAccess(ch, "undelete", p.isChannelOwnerFallback(ch))
}

func (p *permissions) CanManageChannelMembers(ch *types.Channel) bool {
	return p.checkAccess(ch, "members.manage", p.isChannelOwnerFallback(ch))
}

func (p *permissions) CanManageChannelWebhooks(ch *types.Channel) bool {
	return p.checkAccess(ch, "webhooks.manage")
}

func (p *permissions) CanManageChannelAttachments(ch *types.Channel) bool {
	return p.checkAccess(ch, "attachments.manage")
}

func (p *permissions) CanSendMessage(ch *types.Channel) bool {
	return p.checkAccess(ch, "message.send", p.canSendMessagesFallback(ch))
}

func (p *permissions) CanReplyMessage(ch *types.Channel) bool {
	return p.checkAccess(ch, "message.reply", p.allow())
}

func (p *permissions) CanEmbedMessage(ch *types.Channel) bool {
	return p.checkAccess(ch, "message.embed", p.allow())
}

func (p *permissions) CanAttachMessage(ch *types.Channel) bool {
	return p.checkAccess(ch, "message.attach", p.allow())
}

func (p *permissions) CanUpdateOwnMessages(ch *types.Channel) bool {
	return p.checkAccess(ch, "message.update.own", p.allow())
}

func (p *permissions) CanUpdateMessages(ch *types.Channel) bool {
	return p.checkAccess(ch, "message.update.all")
}

func (p *permissions) CanDeleteOwnMessages(ch *types.Channel) bool {
	// @todo implement
	return p.checkAccess(ch, "message.delete.own", p.allow())
}

func (p *permissions) CanDeleteMessages(ch *types.Channel) bool {
	return p.checkAccess(ch, "message.delete.all")
}

func (p *permissions) CanReactMessage(ch *types.Channel) bool {
	return p.checkAccess(ch, "message.react", p.allow())
}

func (p permissions) canJoinFallback(ch *types.Channel) func() internalRules.Access {
	return func() internalRules.Access {
		userID := auth.GetIdentityFromContext(p.ctx).Identity()

		isMember := ch.Member != nil
		isCreator := ch.CreatorID == userID
		isOwner := isCreator || (isMember && ch.Member.Type == types.ChannelMembershipTypeOwner)
		isPublic := ch.Type == types.ChannelTypePublic

		if (ch.IsValid() && isPublic) || isOwner {
			return internalRules.Allow
		}
		return internalRules.Deny
	}
}

func (p permissions) canReadFallback(ch *types.Channel) func() internalRules.Access {
	return func() internalRules.Access {
		if ch.IsValid() && (ch.Type == types.ChannelTypePublic || ch.Member != nil) {
			return internalRules.Allow
		}
		return internalRules.Deny
	}
}

func (p permissions) canSendMessagesFallback(ch *types.Channel) func() internalRules.Access {
	return func() internalRules.Access {
		if ch.IsValid() && (ch.Type == types.ChannelTypePublic || ch.Member != nil) {
			return internalRules.Allow
		}
		return internalRules.Deny
	}
}

func (p permissions) allow() func() internalRules.Access {
	return func() internalRules.Access {
		return internalRules.Allow
	}
}

func (p permissions) isChannelOwnerFallback(ch *types.Channel) func() internalRules.Access {
	return func() internalRules.Access {
		userID := auth.GetIdentityFromContext(p.ctx).Identity()

		isMember := ch.Member != nil
		isCreator := ch.CreatorID == userID
		isOwner := isCreator || (isMember && ch.Member.Type == types.ChannelMembershipTypeOwner)

		if isOwner {
			return internalRules.Allow
		}
		return internalRules.Deny
	}
}

func (p *permissions) checkAccess(res resource, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	access := p.rules.Check(res.PermissionResource(), operation, fallbacks...)
	if access == internalRules.Allow {
		return true
	}
	return false
}
