package service

import (
	"context"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/permissions"
	"github.com/crusttech/crust/messaging/types"
)

type (
	accessControl struct {
		permissions accessControlPermissionServicer
	}

	accessControlPermissionServicer interface {
		Can(context.Context, permissions.Resource, permissions.Operation, ...permissions.CheckAccessFunc) bool
		Grant(context.Context, ...*permissions.Rule) error
	}

	permissionResource interface {
		PermissionResource() permissions.Resource
	}
)

func AccessControl(perm accessControlPermissionServicer) *accessControl {
	return &accessControl{
		permissions: perm,
	}
}

// Effective returns a list of effective service-level permissions
func (svc accessControl) Effective(ctx context.Context) (ee permissions.EffectiveSet) {
	ee = permissions.EffectiveSet{}

	ee.Push(types.MessagingPermissionResource, "access", svc.CanAccess(ctx))
	ee.Push(types.MessagingPermissionResource, "grant", svc.CanGrant(ctx))
	ee.Push(types.MessagingPermissionResource, "channel.public.create", svc.CanCreatePublicChannel(ctx))
	ee.Push(types.MessagingPermissionResource, "channel.private.create", svc.CanCreatePrivateChannel(ctx))
	ee.Push(types.MessagingPermissionResource, "channel.group.create", svc.CanCreateGroupChannel(ctx))

	return
}

func (svc accessControl) CanAccess(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingPermissionResource, "access")
}

func (svc accessControl) CanGrant(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingPermissionResource, "grant")
}

func (svc accessControl) CanCreatePublicChannel(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingPermissionResource, "channel.public.create", permissions.Allowed)
}

func (svc accessControl) CanCreatePrivateChannel(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingPermissionResource, "channel.private.create", permissions.Allowed)
}

func (svc accessControl) CanCreateGroupChannel(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingPermissionResource, "channel.group.create", permissions.Allowed)
}

func (svc accessControl) CanUpdateChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "update", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanReadChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "read", svc.canReadFallback(ch))
}

func (svc accessControl) CanJoinChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "join", svc.canJoinFallback(ctx, ch))
}

func (svc accessControl) CanLeaveChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "leave", permissions.Allowed)
}

func (svc accessControl) CanArchiveChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "archive", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanUnarchiveChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "unarchive", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanDeleteChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "delete", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanUndeleteChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "undelete", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanManageChannelMembers(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "members.manage", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanManageWebhooks(ctx context.Context, webhook *types.Webhook) bool {
	return svc.can(ctx, webhook, "webhook.manage.all")
}

func (svc accessControl) CanManageOwnWebhooks(ctx context.Context, webhook *types.Webhook) bool {
	return svc.can(ctx, webhook, "webhook.manage.own")
}

func (svc accessControl) CanManageChannelAttachments(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "attachments.manage")
}

func (svc accessControl) CanSendMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "message.send", svc.canSendMessagesFallback(ch))
}

func (svc accessControl) CanReplyMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "message.reply", permissions.Allowed)
}

func (svc accessControl) CanEmbedMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "message.embed", permissions.Allowed)
}

func (svc accessControl) CanAttachMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "message.attach", permissions.Allowed)
}

func (svc accessControl) CanUpdateOwnMessages(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "message.update.own", permissions.Allowed)
}

func (svc accessControl) CanUpdateMessages(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "message.update.all")
}

func (svc accessControl) CanDeleteOwnMessages(ctx context.Context, ch *types.Channel) bool {
	// @todo implement
	return svc.can(ctx, ch, "message.delete.own", permissions.Allowed)
}

func (svc accessControl) CanDeleteMessages(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "message.delete.all")
}

func (svc accessControl) CanReactMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch, "message.react", permissions.Allowed)
}

func (svc accessControl) canJoinFallback(ctx context.Context, ch *types.Channel) func() permissions.Access {
	return func() permissions.Access {
		userID := auth.GetIdentityFromContext(ctx).Identity()

		isMember := ch.Member != nil
		isCreator := ch.CreatorID == userID
		isOwner := isCreator || (isMember && ch.Member.Type == types.ChannelMembershipTypeOwner)
		isPublic := ch.Type == types.ChannelTypePublic

		if (ch.IsValid() && isPublic) || isOwner {
			return permissions.Allow
		}
		return permissions.Deny
	}
}

func (svc accessControl) canReadFallback(ch *types.Channel) func() permissions.Access {
	return func() permissions.Access {
		if ch.IsValid() && (ch.Type == types.ChannelTypePublic || ch.Member != nil) {
			return permissions.Allow
		}
		return permissions.Deny
	}
}

func (svc accessControl) canSendMessagesFallback(ch *types.Channel) func() permissions.Access {
	return func() permissions.Access {
		if ch.IsValid() && (ch.Type == types.ChannelTypePublic || ch.Member != nil) {
			return permissions.Allow
		}
		return permissions.Deny
	}
}

func (svc accessControl) isChannelOwnerFallback(ctx context.Context, ch *types.Channel) func() permissions.Access {
	return func() permissions.Access {
		userID := auth.GetIdentityFromContext(ctx).Identity()

		isMember := ch.Member != nil
		isCreator := ch.CreatorID == userID
		isOwner := isCreator || (isMember && ch.Member.Type == types.ChannelMembershipTypeOwner)

		if isOwner {
			return permissions.Allow
		}
		return permissions.Deny
	}
}

func (svc accessControl) can(ctx context.Context, res permissionResource, op permissions.Operation, ff ...permissions.CheckAccessFunc) bool {
	return svc.permissions.Can(ctx, res.PermissionResource(), op, ff...)
}

func (svc accessControl) Grant(ctx context.Context, rr ...*permissions.Rule) error {
	return svc.permissions.Grant(ctx, rr...)
}

// DefaultRules returns list of default rules for this compose service
func (svc accessControl) DefaultRules() permissions.RuleSet {
	var (
		messaging = types.MessagingPermissionResource
		channels  = types.ChannelPermissionResource.AppendWildcard()
		webhooks  = types.WebhookPermissionResource.AppendWildcard()

		allowAdm = func(res permissions.Resource, op permissions.Operation) *permissions.Rule {
			return &permissions.Rule{
				permissions.AdminRoleID,
				res,
				op,
				permissions.Allow}
		}
	)

	return permissions.RuleSet{
		{permissions.EveryoneRoleID, messaging, "access", permissions.Allow},

		allowAdm(messaging, "access"),
		allowAdm(messaging, "grant"),
		allowAdm(messaging, "channel.public.create"),
		allowAdm(messaging, "channel.private.create"),
		allowAdm(messaging, "channel.group.create"),

		allowAdm(channels, "update"),
		allowAdm(channels, "leave"),
		allowAdm(channels, "read"),
		allowAdm(channels, "join"),
		allowAdm(channels, "delete"),
		allowAdm(channels, "undelete"),
		allowAdm(channels, "archive"),
		allowAdm(channels, "unarchive"),
		allowAdm(channels, "members.manage"),
		allowAdm(channels, "webhooks.manage"),
		allowAdm(channels, "attachments.manage"),
		allowAdm(channels, "message.attach"),
		allowAdm(channels, "message.update.all"),
		allowAdm(channels, "message.update.own"),
		allowAdm(channels, "message.delete.all"),
		allowAdm(channels, "message.delete.own"),
		allowAdm(channels, "message.embed"),
		allowAdm(channels, "message.send"),
		allowAdm(channels, "message.reply"),
		allowAdm(channels, "message.react"),

		allowAdm(webhooks, "webhook.manage.all"),
		allowAdm(webhooks, "webhook.manage.own"),
	}
}
