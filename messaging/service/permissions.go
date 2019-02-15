package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/organization"
	internalRules "github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/messaging/repository"
	"github.com/crusttech/crust/messaging/types"
	systemRepository "github.com/crusttech/crust/system/repository"
)

type (
	permissions struct {
		db  db
		ctx context.Context

		team    systemRepository.TeamRepository
		channel repository.ChannelRepository

		scopes    internalRules.ScopeInterface
		resources internalRules.ResourcesInterface
	}

	PermissionsService interface {
		With(ctx context.Context) PermissionsService

		List() (interface{}, error)
		Get(teamID uint64, resource string) (interface{}, error)
		Set(teamID uint64, rules []internalRules.Rules) (interface{}, error)

		Scopes(scope string) (interface{}, error)
	}
)

func Permissions(scopes internalRules.ScopeInterface) PermissionsService {
	return (&permissions{
		scopes: scopes,
	}).With(context.Background())
}

func (p *permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:  db,
		ctx: ctx,

		team:    systemRepository.Team(ctx, db),
		channel: repository.Channel(ctx, db),

		scopes:    p.scopes,
		resources: internalRules.NewResources(ctx, db),
	}
}

func (p *permissions) List() (interface{}, error) {
	return p.scopes.List(), nil
}

func (p *permissions) Get(teamID uint64, resource string) (interface{}, error) {
	return p.resources.ListGrants(teamID, resource)
}

func (p *permissions) Set(teamID uint64, rules []internalRules.Rules) (interface{}, error) {
	var err error
	for _, rule := range rules {
		err = p.resources.Grant(
			teamID,
			rule.Resource,
			[]string{rule.Operation},
			rule.Value,
		)
		if err != nil {
			break
		}
	}
	return nil, err
}

func (p *permissions) Scopes(scope string) (interface{}, error) {
	switch scope {
	case "organization":
		// @todo organizations from DB once multi-org
		// return p.organizaion.Find(nil)
		orgs := []types.Organisation{
			types.Organisation{
				organization.Crust(),
			},
		}
		return orgs, nil
	case "team":
		return p.team.Find(nil)
	case "channel":
		return p.channel.FindChannels(nil)
	}
	return nil, errors.New("no scope defined")
}
