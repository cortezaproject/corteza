package service

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/logger"
	internalRules "github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/messaging/internal/repository"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	permissions struct {
		db     db
		ctx    context.Context
		logger *zap.Logger

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

		CanUpdateChannel(*types.Channel) bool
		CanReadChannel(*types.Channel) bool
		CanJoinChannel(*types.Channel) bool
		CanLeaveChannel(*types.Channel) bool
		CanDeleteChannel(*types.Channel) bool
		CanUndeleteChannel(*types.Channel) bool
		CanArchiveChannel(*types.Channel) bool
		CanUnarchiveChannel(*types.Channel) bool

		CanManageChannelMembers(*types.Channel) bool
		CanManageChannelAttachments(*types.Channel) bool

		CanManageWebhooks(*types.Webhook) bool
		CanManageOwnWebhooks(*types.Webhook) bool

		CanSendMessage(*types.Channel) bool
		CanReplyMessage(*types.Channel) bool
		CanEmbedMessage(*types.Channel) bool
		CanAttachMessage(*types.Channel) bool
		CanUpdateOwnMessages(*types.Channel) bool
		CanUpdateMessages(*types.Channel) bool
		CanDeleteOwnMessages(*types.Channel) bool
		CanDeleteMessages(*types.Channel) bool
		CanReactMessage(*types.Channel) bool
	}

	effectivePermission struct {
		Resource  internalRules.Resource `json:"resource"`
		Operation string                 `json:"operation"`
		Allow     bool                   `json:"allow"`
	}
)

func Permissions(ctx context.Context) PermissionsService {
	return (&permissions{
		logger: DefaultLogger.Named("permissions"),
	}).With(ctx)
}

func (svc permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		rules: systemService.Rules(ctx),
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc permissions) log(fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
}

func (svc permissions) Effective() (ee []effectivePermission, err error) {
	ep := func(res internalRules.Resource, op string, allow bool) effectivePermission {
		return effectivePermission{
			Resource:  res,
			Operation: op,
			Allow:     allow,
		}
	}

	ee = append(ee, ep(types.PermissionResource, "access", svc.CanAccess()))
	ee = append(ee, ep(types.PermissionResource, "grant", svc.CanGrant()))
	ee = append(ee, ep(types.PermissionResource, "channel.public.create", svc.CanCreatePublicChannel()))
	ee = append(ee, ep(types.PermissionResource, "channel.private.create", svc.CanCreatePrivateChannel()))
	ee = append(ee, ep(types.PermissionResource, "channel.group.create", svc.CanCreateGroupChannel()))

	return
}

func (svc permissions) CanAccess() bool {
	return svc.checkAccess(types.PermissionResource, "access")
}

func (svc permissions) CanGrant() bool {
	return svc.checkAccess(types.PermissionResource, "grant")
}

func (svc permissions) CanCreatePublicChannel() bool {
	return svc.checkAccess(types.PermissionResource, "channel.public.create", svc.allow())
}

func (svc permissions) CanCreatePrivateChannel() bool {
	return svc.checkAccess(types.PermissionResource, "channel.private.create", svc.allow())
}

func (svc permissions) CanCreateGroupChannel() bool {
	return svc.checkAccess(types.PermissionResource, "channel.group.create", svc.allow())
}

func (svc permissions) CanUpdateChannel(ch *types.Channel) bool {
	return svc.checkAccess(ch, "update", svc.isChannelOwnerFallback(ch))
}

func (svc permissions) CanReadChannel(ch *types.Channel) bool {
	return svc.checkAccess(ch, "read", svc.canReadFallback(ch))
}

func (svc permissions) CanJoinChannel(ch *types.Channel) bool {
	return svc.checkAccess(ch, "join", svc.canJoinFallback(ch))
}

func (svc permissions) CanLeaveChannel(ch *types.Channel) bool {
	return svc.checkAccess(ch, "leave", svc.allow())
}

func (svc permissions) CanArchiveChannel(ch *types.Channel) bool {
	return svc.checkAccess(ch, "archive", svc.isChannelOwnerFallback(ch))
}

func (svc permissions) CanUnarchiveChannel(ch *types.Channel) bool {
	return svc.checkAccess(ch, "unarchive", svc.isChannelOwnerFallback(ch))
}

func (svc permissions) CanDeleteChannel(ch *types.Channel) bool {
	return svc.checkAccess(ch, "delete", svc.isChannelOwnerFallback(ch))
}

func (svc permissions) CanUndeleteChannel(ch *types.Channel) bool {
	return svc.checkAccess(ch, "undelete", svc.isChannelOwnerFallback(ch))
}

func (svc permissions) CanManageChannelMembers(ch *types.Channel) bool {
	return svc.checkAccess(ch, "members.manage", svc.isChannelOwnerFallback(ch))
}

func (svc permissions) CanManageWebhooks(webhook *types.Webhook) bool {
	return svc.checkAccess(webhook, "webhook.manage.all")
}

func (svc permissions) CanManageOwnWebhooks(webhook *types.Webhook) bool {
	return svc.checkAccess(webhook, "webhook.manage.own")
}

func (svc permissions) CanManageChannelAttachments(ch *types.Channel) bool {
	return svc.checkAccess(ch, "attachments.manage")
}

func (svc permissions) CanSendMessage(ch *types.Channel) bool {
	return svc.checkAccess(ch, "message.send", svc.canSendMessagesFallback(ch))
}

func (svc permissions) CanReplyMessage(ch *types.Channel) bool {
	return svc.checkAccess(ch, "message.reply", svc.allow())
}

func (svc permissions) CanEmbedMessage(ch *types.Channel) bool {
	return svc.checkAccess(ch, "message.embed", svc.allow())
}

func (svc permissions) CanAttachMessage(ch *types.Channel) bool {
	return svc.checkAccess(ch, "message.attach", svc.allow())
}

func (svc permissions) CanUpdateOwnMessages(ch *types.Channel) bool {
	return svc.checkAccess(ch, "message.update.own", svc.allow())
}

func (svc permissions) CanUpdateMessages(ch *types.Channel) bool {
	return svc.checkAccess(ch, "message.update.all")
}

func (svc permissions) CanDeleteOwnMessages(ch *types.Channel) bool {
	// @todo implement
	return svc.checkAccess(ch, "message.delete.own", svc.allow())
}

func (svc permissions) CanDeleteMessages(ch *types.Channel) bool {
	return svc.checkAccess(ch, "message.delete.all")
}

func (svc permissions) CanReactMessage(ch *types.Channel) bool {
	return svc.checkAccess(ch, "message.react", svc.allow())
}

func (svc permissions) canJoinFallback(ch *types.Channel) func() internalRules.Access {
	return func() internalRules.Access {
		userID := auth.GetIdentityFromContext(svc.ctx).Identity()

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

func (svc permissions) canReadFallback(ch *types.Channel) func() internalRules.Access {
	return func() internalRules.Access {
		if ch.IsValid() && (ch.Type == types.ChannelTypePublic || ch.Member != nil) {
			return internalRules.Allow
		}
		return internalRules.Deny
	}
}

func (svc permissions) canSendMessagesFallback(ch *types.Channel) func() internalRules.Access {
	return func() internalRules.Access {
		if ch.IsValid() && (ch.Type == types.ChannelTypePublic || ch.Member != nil) {
			return internalRules.Allow
		}
		return internalRules.Deny
	}
}

func (svc permissions) allow() func() internalRules.Access {
	return func() internalRules.Access {
		return internalRules.Allow
	}
}

func (svc permissions) isChannelOwnerFallback(ch *types.Channel) func() internalRules.Access {
	return func() internalRules.Access {
		userID := auth.GetIdentityFromContext(svc.ctx).Identity()

		isMember := ch.Member != nil
		isCreator := ch.CreatorID == userID
		isOwner := isCreator || (isMember && ch.Member.Type == types.ChannelMembershipTypeOwner)

		if isOwner {
			return internalRules.Allow
		}
		return internalRules.Deny
	}
}

func (svc permissions) checkAccess(res resource, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	access := svc.rules.Check(res.PermissionResource(), operation, fallbacks...)
	if access == internalRules.Allow {
		return true
	}
	return false
}
