package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

type (
	accessControl struct {
		permissions accessControlRBACServicer
		actionlog   actionlog.Recorder
	}

	accessControlRBACServicer interface {
		Can([]uint64, rbac.Resource, rbac.Operation, ...rbac.CheckAccessFunc) bool
		Grant(context.Context, rbac.Whitelist, ...*rbac.Rule) error
		FindRulesByRoleID(roleID uint64) (rr rbac.RuleSet)
	}
)

func AccessControl(perm accessControlRBACServicer) *accessControl {
	return &accessControl{
		permissions: perm,
		actionlog:   DefaultActionlog,
	}
}

// Effective returns a list of effective service-level permissions
func (svc accessControl) Effective(ctx context.Context) (ee rbac.EffectiveSet) {
	ee = rbac.EffectiveSet{}

	ee.Push(types.MessagingRBACResource, "grant", svc.CanGrant(ctx))
	ee.Push(types.MessagingRBACResource, "settings.read", svc.CanReadSettings(ctx))
	ee.Push(types.MessagingRBACResource, "settings.manage", svc.CanManageSettings(ctx))
	ee.Push(types.MessagingRBACResource, "channel.public.create", svc.CanCreatePublicChannel(ctx))
	ee.Push(types.MessagingRBACResource, "channel.private.create", svc.CanCreatePrivateChannel(ctx))
	ee.Push(types.MessagingRBACResource, "channel.group.create", svc.CanCreateGroupChannel(ctx))

	return
}

func (svc accessControl) CanGrant(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingRBACResource, "grant")
}

func (svc accessControl) CanReadSettings(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingRBACResource, "settings.read")
}

func (svc accessControl) CanManageSettings(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingRBACResource, "settings.manage")
}

func (svc accessControl) CanCreatePublicChannel(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingRBACResource, "channel.public.create", rbac.Allowed)
}

func (svc accessControl) CanCreatePrivateChannel(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingRBACResource, "channel.private.create", rbac.Allowed)
}

func (svc accessControl) CanCreateGroupChannel(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingRBACResource, "channel.group.create", rbac.Allowed)
}

func (svc accessControl) CanUpdateChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "update", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanReadChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "read", svc.canReadFallback(ch))
}

func (svc accessControl) CanJoinChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "join", svc.canJoinFallback(ctx, ch))
}

func (svc accessControl) CanLeaveChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "leave", rbac.Allowed)
}

func (svc accessControl) CanArchiveChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "archive", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanUnarchiveChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "unarchive", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanDeleteChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "delete", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanUndeleteChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "undelete", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanManageChannelMembers(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "members.manage", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanChangeChannelMembershipPolicy(ctx context.Context, ch *types.Channel) bool {
	// @todo introduce dedicated channel op. for this.
	return svc.can(ctx, types.MessagingRBACResource, "grant")
}

func (svc accessControl) CanManageChannelAttachments(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "attachments.manage")
}

func (svc accessControl) CanSendMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "message.send", svc.canSendMessagesFallback(ch))
}

func (svc accessControl) CanReplyMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "message.reply", rbac.Allowed)
}

func (svc accessControl) CanEmbedMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "message.embed", rbac.Allowed)
}

func (svc accessControl) CanAttachMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "message.attach", svc.canSendMessagesFallback(ch))
}

func (svc accessControl) CanUpdateOwnMessages(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "message.update.own", rbac.Allowed)
}

func (svc accessControl) CanUpdateMessages(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "message.update.all")
}

func (svc accessControl) CanDeleteOwnMessages(ctx context.Context, ch *types.Channel) bool {
	// @todo implement
	return svc.can(ctx, ch.RBACResource(), "message.delete.own", rbac.Allowed)
}

func (svc accessControl) CanDeleteMessages(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "message.delete.all")
}

func (svc accessControl) CanReactMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.RBACResource(), "message.react", rbac.Allowed)
}

func (svc accessControl) canJoinFallback(ctx context.Context, ch *types.Channel) func() rbac.Access {
	return func() rbac.Access {
		userID := auth.GetIdentityFromContext(ctx).Identity()

		isMember := ch.Member != nil
		isCreator := ch.CreatorID == userID
		isOwner := isCreator || (isMember && ch.Member.Type == types.ChannelMembershipTypeOwner)
		isPublic := ch.Type == types.ChannelTypePublic

		if (ch.IsValid() && isPublic) || isOwner {
			return rbac.Allow
		}
		return rbac.Deny
	}
}

func (svc accessControl) canReadFallback(ch *types.Channel) func() rbac.Access {
	return func() rbac.Access {
		if ch.IsValid() && (ch.Type == types.ChannelTypePublic || ch.Member != nil) {
			return rbac.Allow
		}
		return rbac.Deny
	}
}

func (svc accessControl) canSendMessagesFallback(ch *types.Channel) func() rbac.Access {
	return func() rbac.Access {
		if ch.IsValid() && (ch.Type == types.ChannelTypePublic || ch.Member != nil) {
			return rbac.Allow
		}
		return rbac.Deny
	}
}

func (svc accessControl) isChannelOwnerFallback(ctx context.Context, ch *types.Channel) func() rbac.Access {
	return func() rbac.Access {
		userID := auth.GetIdentityFromContext(ctx).Identity()

		isMember := ch.Member != nil
		isCreator := ch.CreatorID == userID
		isOwner := isCreator || (isMember && ch.Member.Type == types.ChannelMembershipTypeOwner)

		if isOwner {
			return rbac.Allow
		}
		return rbac.Deny
	}
}

func (svc accessControl) can(ctx context.Context, res rbac.Resource, op rbac.Operation, ff ...rbac.CheckAccessFunc) bool {
	var (
		u     = auth.GetIdentityFromContext(ctx)
		roles = u.Roles()
	)

	if auth.IsSuperUser(u) {
		// Temp solution to allow migration from passing context to ResourceFilter
		// and checking "superuser" privileges there to more sustainable solution
		// (eg: creating super-role with allow-all)
		return true
	}

	return svc.permissions.Can(roles, res, op, ff...)
}

func (svc accessControl) Grant(ctx context.Context, rr ...*rbac.Rule) error {
	if !svc.CanGrant(ctx) {
		return AccessControlErrNotAllowedToSetPermissions()
	}

	if err := svc.permissions.Grant(ctx, svc.Whitelist(), rr...); err != nil {
		return AccessControlErrGeneric().Wrap(err)
	}

	svc.logGrants(ctx, rr)

	return nil
}

func (svc accessControl) logGrants(ctx context.Context, rr []*rbac.Rule) {
	if svc.actionlog == nil {
		return
	}

	for _, r := range rr {
		g := AccessControlActionGrant(&accessControlActionProps{r})
		g.log = r.String()
		g.resource = r.Resource.String()

		svc.actionlog.Record(ctx, g.ToAction())
	}
}

func (svc accessControl) FindRulesByRoleID(ctx context.Context, roleID uint64) (rbac.RuleSet, error) {
	if !svc.CanGrant(ctx) {
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	return svc.permissions.FindRulesByRoleID(roleID), nil
}

func (svc accessControl) Whitelist() rbac.Whitelist {
	var wl = rbac.Whitelist{}

	wl.Set(
		types.MessagingRBACResource,
		"grant",
		"settings.read",
		"settings.manage",
		"channel.public.create",
		"channel.private.create",
		"channel.group.create",
	)

	wl.Set(
		types.ChannelRBACResource,
		"update",
		"read",
		"join",
		"leave",
		"delete",
		"undelete",
		"archive",
		"unarchive",
		"members.manage",
		"attachments.manage",
		"message.send",
		"message.reply",
		"message.embed",
		"message.attach",
		"message.update.own",
		"message.update.all",
		"message.delete.own",
		"message.delete.all",
		"message.react",
	)

	return wl
}
