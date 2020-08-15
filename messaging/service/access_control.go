package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	accessControl struct {
		permissions accessControlPermissionServicer
		actionlog   actionlog.Recorder
	}

	accessControlPermissionServicer interface {
		Can([]uint64, permissions.Resource, permissions.Operation, ...permissions.CheckAccessFunc) bool
		Grant(context.Context, permissions.Whitelist, ...*permissions.Rule) error
		FindRulesByRoleID(roleID uint64) (rr permissions.RuleSet)
		ResourceFilter([]uint64, permissions.Resource, permissions.Operation, permissions.Access) *permissions.ResourceFilter
	}
)

func AccessControl(perm accessControlPermissionServicer) *accessControl {
	return &accessControl{
		permissions: perm,
		actionlog:   DefaultActionlog,
	}
}

// Effective returns a list of effective service-level permissions
func (svc accessControl) Effective(ctx context.Context) (ee permissions.EffectiveSet) {
	ee = permissions.EffectiveSet{}

	ee.Push(types.MessagingPermissionResource, "access", svc.CanAccess(ctx))
	ee.Push(types.MessagingPermissionResource, "grant", svc.CanGrant(ctx))
	ee.Push(types.MessagingPermissionResource, "settings.read", svc.CanReadSettings(ctx))
	ee.Push(types.MessagingPermissionResource, "settings.manage", svc.CanManageSettings(ctx))
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

func (svc accessControl) CanReadSettings(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingPermissionResource, "settings.read")
}

func (svc accessControl) CanManageSettings(ctx context.Context) bool {
	return svc.can(ctx, types.MessagingPermissionResource, "settings.manage")
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
	return svc.can(ctx, ch.PermissionResource(), "update", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanReadChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "read", svc.canReadFallback(ch))
}

func (svc accessControl) CanJoinChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "join", svc.canJoinFallback(ctx, ch))
}

func (svc accessControl) CanLeaveChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "leave", permissions.Allowed)
}

func (svc accessControl) CanArchiveChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "archive", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanUnarchiveChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "unarchive", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanDeleteChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "delete", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanUndeleteChannel(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "undelete", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanManageChannelMembers(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "members.manage", svc.isChannelOwnerFallback(ctx, ch))
}

func (svc accessControl) CanChangeChannelMembershipPolicy(ctx context.Context, ch *types.Channel) bool {
	// @todo introduce dedicated channel op. for this.
	return svc.can(ctx, types.MessagingPermissionResource, "grant")
}

func (svc accessControl) CanManageChannelAttachments(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "attachments.manage")
}

func (svc accessControl) CanSendMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "message.send", svc.canSendMessagesFallback(ch))
}

func (svc accessControl) CanReplyMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "message.reply", permissions.Allowed)
}

func (svc accessControl) CanEmbedMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "message.embed", permissions.Allowed)
}

func (svc accessControl) CanAttachMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "message.attach", svc.canSendMessagesFallback(ch))
}

func (svc accessControl) CanUpdateOwnMessages(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "message.update.own", permissions.Allowed)
}

func (svc accessControl) CanUpdateMessages(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "message.update.all")
}

func (svc accessControl) CanDeleteOwnMessages(ctx context.Context, ch *types.Channel) bool {
	// @todo implement
	return svc.can(ctx, ch.PermissionResource(), "message.delete.own", permissions.Allowed)
}

func (svc accessControl) CanDeleteMessages(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "message.delete.all")
}

func (svc accessControl) CanReactMessage(ctx context.Context, ch *types.Channel) bool {
	return svc.can(ctx, ch.PermissionResource(), "message.react", permissions.Allowed)
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

func (svc accessControl) can(ctx context.Context, res permissions.Resource, op permissions.Operation, ff ...permissions.CheckAccessFunc) bool {
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

func (svc accessControl) filter(ctx context.Context, res permissions.Resource, op permissions.Operation, a permissions.Access) *permissions.ResourceFilter {
	var (
		u     = auth.GetIdentityFromContext(ctx)
		roles = u.Roles()
	)

	if auth.IsSuperUser(u) {
		// Temp solution to allow migration from passing context to ResourceFilter
		// and checking "superuser" privileges there to more sustainable solution
		// (eg: creating super-role with allow-all)
		return permissions.NewSuperuserFilter()
	}

	return svc.permissions.ResourceFilter(roles, res, op, a)
}

func (svc accessControl) Grant(ctx context.Context, rr ...*permissions.Rule) error {
	if !svc.CanGrant(ctx) {
		return AccessControlErrNotAllowedToSetPermissions()
	}

	if err := svc.permissions.Grant(ctx, svc.Whitelist(), rr...); err != nil {
		return AccessControlErrGeneric().Wrap(err)
	}

	svc.logGrants(ctx, rr)

	return nil
}

func (svc accessControl) logGrants(ctx context.Context, rr []*permissions.Rule) {
	if svc.actionlog == nil {
		return
	}

	for _, r := range rr {
		g := AccessControlActionGrant(&accessControlActionProps{r})
		g.log = r.String()
		g.resource = r.Resource.String()

		svc.actionlog.Record(ctx, g)
	}
}

func (svc accessControl) FindRulesByRoleID(ctx context.Context, roleID uint64) (permissions.RuleSet, error) {
	if !svc.CanGrant(ctx) {
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	return svc.permissions.FindRulesByRoleID(roleID), nil
}

func (svc accessControl) Whitelist() permissions.Whitelist {
	var wl = permissions.Whitelist{}

	wl.Set(
		types.MessagingPermissionResource,
		"access",
		"grant",
		"settings.read",
		"settings.manage",
		"channel.public.create",
		"channel.private.create",
		"channel.group.create",
	)

	wl.Set(
		types.ChannelPermissionResource,
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
