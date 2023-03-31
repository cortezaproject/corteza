package envoy

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

func (e StoreEncoder) prepare(ctx context.Context, p envoyx.EncodeParams, s store.Storer, rt string, nn envoyx.NodeSet) (err error) {
	switch rt {
	case rbac.RuleResourceType:
		return e.prepareRbacRule(ctx, p, s, nn)
	case types.ResourceTranslationResourceType:
		return e.prepareResourceTranslation(ctx, p, s, nn)
	case types.SettingValueResourceType:
		return e.prepareSetting(ctx, p, s, nn)
	}

	return
}

func (e StoreEncoder) encode(ctx context.Context, p envoyx.EncodeParams, s store.Storer, rt string, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	switch rt {
	case rbac.RuleResourceType:
		return e.encodeRbacRules(ctx, p, s, nn, tree)
	case types.ResourceTranslationResourceType:
		return e.encodeResourceTranslations(ctx, p, s, nn, tree)
	case types.SettingValueResourceType:
		return e.encodeSettings(ctx, p, s, nn, tree)
	}

	return
}

func (e StoreEncoder) setApplicationDefaults(res *types.Application) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	if res.Unify == nil {
		res.Unify = &types.ApplicationUnify{}
	}

	return
}

func (e StoreEncoder) validateApplication(res *types.Application) (err error) {
	return
}

func (e StoreEncoder) setApigwRouteDefaults(res *types.ApigwRoute) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return
}

func (e StoreEncoder) validateApigwRoute(res *types.ApigwRoute) (err error) {
	return
}

func (e StoreEncoder) setApigwFilterDefaults(res *types.ApigwFilter) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return
}

func (e StoreEncoder) validateApigwFilter(res *types.ApigwFilter) (err error) {
	return
}

func (e StoreEncoder) setAuthClientDefaults(res *types.AuthClient) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return
}

func (e StoreEncoder) validateAuthClient(res *types.AuthClient) (err error) {
	return
}

func (e StoreEncoder) setQueueDefaults(res *types.Queue) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return
}

func (e StoreEncoder) validateQueue(res *types.Queue) (err error) {
	return
}

func (e StoreEncoder) setReportDefaults(res *types.Report) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return
}

func (e StoreEncoder) validateReport(res *types.Report) (err error) {
	return
}

func (e StoreEncoder) setRoleDefaults(res *types.Role) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return
}

func (e StoreEncoder) validateRole(res *types.Role) (err error) {
	return
}

func (e StoreEncoder) setTemplateDefaults(res *types.Template) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return
}

func (e StoreEncoder) validateTemplate(res *types.Template) (err error) {
	return
}

func (e StoreEncoder) setUserDefaults(res *types.User) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return
}

func (e StoreEncoder) validateUser(res *types.User) (err error) {
	return
}

func (e StoreEncoder) setDalConnectionDefaults(res *types.DalConnection) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return
}

func (e StoreEncoder) validateDalConnection(res *types.DalConnection) (err error) {
	return
}

func (e StoreEncoder) setDalSensitivityLevelDefaults(res *types.DalSensitivityLevel) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	return
}

func (e StoreEncoder) validateDalSensitivityLevel(res *types.DalSensitivityLevel) (err error) {
	return
}
