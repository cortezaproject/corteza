package service

import (
	"context"

	internalRules "github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/messaging/repository"
	"github.com/crusttech/crust/messaging/types"
	systemTypes "github.com/crusttech/crust/system/types"
)

type (
	rules struct {
		db  db
		ctx context.Context

		// identity is passed with context
		resources internalRules.ResourcesInterface
		team      *systemTypes.Team
		org       *types.Organisation
	}

	RulesService interface {
		With(context.Context) RulesService

		// Applies mostly to admin panel
		isAdmin() bool

		// Individual rules for administration
		canManageOrganisation() bool
		canManageRoles() bool
		canManageChannels() bool

		// types.Team derived from identity uses a wildcard match
		canManageWebhooks(ch *types.Channel) bool

		// Messaging rules
		canSendMessages(ch *types.Channel) bool
		canEmbedLinks(ch *types.Channel) bool
		canAttachFiles(ch *types.Channel) bool
		canUpdateOwnMessages(ch *types.Channel) bool
		canUpdateMessages(ch *types.Channel) bool
		canReact(ch *types.Channel) bool
	}
)

func Rules() RulesService {
	return (&rules{
		team: &systemTypes.Team{},
	}).With(context.Background())
}

func (r *rules) With(ctx context.Context) RulesService {
	db := repository.DB(ctx)
	org := repository.Organization(ctx)
	return &rules{
		db:   db,
		ctx:  ctx,
		org:  org,
		team: r.team,

		resources: internalRules.NewResources(ctx, db),
	}
}

// @todo: honor defaults from (org/team/channel).Permissions()

func (r *rules) isAdmin() bool {
	op := "admin"
	return r.hasAccess(op, r.org.PermissionDefault(op), r.org.Resource().String())
}

func (r *rules) canManageOrganisation() bool {
	op := "manage.organisation"
	return r.hasAccess(op, r.org.PermissionDefault(op), r.org.Resource().String())
}

func (r *rules) canManageRoles() bool {
	op := "manage.roles"
	return r.hasAccess(op, r.org.PermissionDefault(op), r.org.Resource().String())
}

func (r *rules) canManageChannels() bool {
	op := "manage.channels"
	return r.hasAccess(op, r.org.PermissionDefault(op), r.org.Resource().String())
}

func (r *rules) canManageWebhooks(ch *types.Channel) bool {
	op := "manage.webhooks"
	return r.hasAccess(op, r.org.PermissionDefault(op), r.org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *rules) canSendMessages(ch *types.Channel) bool {
	op := "message.send"
	return r.hasAccess(op, r.org.PermissionDefault(op), r.org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *rules) canEmbedLinks(ch *types.Channel) bool {
	op := "message.embed"
	return r.hasAccess(op, r.org.PermissionDefault(op), r.org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *rules) canAttachFiles(ch *types.Channel) bool {
	op := "message.attach"
	return r.hasAccess(op, r.org.PermissionDefault(op), r.org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *rules) canUpdateOwnMessages(ch *types.Channel) bool {
	op := "message.update_own"
	return r.hasAccess(op, r.org.PermissionDefault(op), r.org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *rules) canUpdateMessages(ch *types.Channel) bool {
	op := "message.update_all"
	return r.hasAccess(op, r.org.PermissionDefault(op), r.org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *rules) canReact(ch *types.Channel) bool {
	op := "message.react"
	return r.hasAccess(op, r.org.PermissionDefault(op), r.org.Resource().String(), r.team.Resource().All(), ch.Resource().String())
}

func (r *rules) hasAccess(operation string, value internalRules.Access, scopes ...string) bool {
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
		case internalRules.Allow:
			return true
		case internalRules.Deny:
			return false
		default: // inherit
		}
	}
	return false
}
