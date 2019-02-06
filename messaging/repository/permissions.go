package repository

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/messaging/types"
	systemTypes "github.com/crusttech/crust/system/types"
)

type PermissionsRepository interface {
	With(context.Context, *factory.DB) PermissionsRepository

	// Applies mostly to admin panel
	isAdmin(org *types.Organisation) bool

	// Individual rules for administration
	canManageOrganisation(org *types.Organisation) bool
	canManageRoles(org *types.Organisation) bool
	canManageChannels(org *types.Organisation) bool

	// types.Team derived from identity uses a wildcard match
	canManageWebhooks(org *types.Organisation, ch *types.Channel) bool

	// Messaging rules
	canSendMessages(org *types.Organisation, ch *types.Channel) bool
	canEmbedLinks(org *types.Organisation, ch *types.Channel) bool
	canAttachFiles(org *types.Organisation, ch *types.Channel) bool
	canEditOwnMessages(org *types.Organisation, ch *types.Channel) bool
	canEditMessages(org *types.Organisation, ch *types.Channel) bool
	canReact(org *types.Organisation, ch *types.Channel) bool
}

type permissions struct {
	*repository

	// identity is passed with context
	resources rules.ResourcesInterface
	team      *systemTypes.Team
}

func Permissions(ctx context.Context, db *factory.DB) PermissionsRepository {
	return (&permissions{
		team: &systemTypes.Team{},
	}).With(ctx, db)
}

func (r *permissions) With(ctx context.Context, db *factory.DB) PermissionsRepository {
	return &permissions{
		team:       r.team,
		resources:  rules.NewResources(ctx, db),
		repository: r.repository.With(ctx, db),
	}
}

var (
	ErrPermissionsNotLoggedIn = repositoryError("PermissionsNotLoggedIn")
)

// @todo: honor defaults from (org/team/channel).Permissions()

func (r *permissions) isAdmin(org *types.Organisation) bool {
	op := "admin"
	return r.hasAccess(op, org.PermissionDefault(op), org.Resource().String())
}

func (r *permissions) canManageOrganisation(org *types.Organisation) bool {
	op := "manage.organisation"
	return r.hasAccess(op, org.PermissionDefault(op), org.Resource().String())
}

func (r *permissions) canManageRoles(org *types.Organisation) bool {
	op := "manage.roles"
	return r.hasAccess(op, org.PermissionDefault(op), org.Resource().String())
}

func (r *permissions) canManageChannels(org *types.Organisation) bool {
	op := "manage.channels"
	return r.hasAccess(op, org.PermissionDefault(op), org.Resource().String())
}

func (r *permissions) canManageWebhooks(org *types.Organisation, ch *types.Channel) bool {
	op := "manage.webhooks"
	return r.hasAccess(op, org.PermissionDefault(op), org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *permissions) canSendMessages(org *types.Organisation, ch *types.Channel) bool {
	op := "text.send"
	return r.hasAccess(op, org.PermissionDefault(op), org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *permissions) canEmbedLinks(org *types.Organisation, ch *types.Channel) bool {
	op := "text.embed"
	return r.hasAccess(op, org.PermissionDefault(op), org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *permissions) canAttachFiles(org *types.Organisation, ch *types.Channel) bool {
	op := "text.attach"
	return r.hasAccess(op, org.PermissionDefault(op), org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *permissions) canEditOwnMessages(org *types.Organisation, ch *types.Channel) bool {
	op := "text.edit_own"
	return r.hasAccess(op, org.PermissionDefault(op), org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *permissions) canEditMessages(org *types.Organisation, ch *types.Channel) bool {
	op := "text.edit_all"
	return r.hasAccess(op, org.PermissionDefault(op), org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *permissions) canReact(org *types.Organisation, ch *types.Channel) bool {
	op := "text.react"
	return r.hasAccess(op, org.PermissionDefault(op), org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *permissions) hasAccess(operation string, value rules.Access, scopes ...string) bool {
	// reverse scopes from to order it from most-least significant
	// aka: [0]channel [1]teams [2]org
	last := len(scopes) - 1
	for i := 0; i < len(scopes)/2; i++ {
		scopes[i], scopes[last-i] = scopes[last-i], scopes[i]
	}

	for _, scope := range scopes {
		if scope == "" {
			continue
		}
		switch r.resources.IsAllowed(scope, operation) {
		case rules.Allow:
			return true
		case rules.Deny:
			return false
		default: // inherit
		}
	}
	return false
}
