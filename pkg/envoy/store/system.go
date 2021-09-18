package store

import (
	"context"
	"strings"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cast"
)

type (
	roleFilter        types.RoleFilter
	userFilter        types.UserFilter
	templateFilter    types.TemplateFilter
	applicationFilter types.ApplicationFilter
	apiGwRouteFilter  types.ApigwRouteFilter
	reportFilter      types.ReportFilter
	settingFilter     types.SettingsFilter
	rbacFilter        struct {
		rbac.RuleFilter
		// This will help us determine what rules for what resources we are able to export
		resourceID map[uint64]bool
		strict     bool
	}
	resourceTranslationFilter types.ResourceTranslationFilter

	systemDecoder struct {
		resourceID []uint64
		ux         *userIndex
	}
)

func newSystemDecoder(ux *userIndex) *systemDecoder {
	return &systemDecoder{
		resourceID: make([]uint64, 0, 200),
		ux:         ux,
	}
}

func (d *systemDecoder) decodeRoles(ctx context.Context, s store.Storer, ff []*roleFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.RoleSet
	var fn types.RoleFilter
	var err error

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchRoles(ctx, types.RoleFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				mm = append(mm, newRole(n))
				d.resourceID = append(d.resourceID, n.ID)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

func (d *systemDecoder) decodeUsers(ctx context.Context, s store.Storer, ff []*userFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.UserSet
	var fn types.UserFilter
	var err error

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchUsers(ctx, types.UserFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				mm = append(mm, newUser(n))
				d.resourceID = append(d.resourceID, n.ID)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

func (d *systemDecoder) decodeTemplates(ctx context.Context, s store.Storer, ff []*templateFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.TemplateSet
	var fn types.TemplateFilter
	var err error

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchTemplates(ctx, types.TemplateFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				mm = append(mm, newTemplate(n))
				d.resourceID = append(d.resourceID, n.ID)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

func (d *systemDecoder) decodeAPIGWRoutes(ctx context.Context, s store.Storer, ff []*apiGwRouteFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.ApigwRouteSet
	var fn types.ApigwRouteFilter
	var err error

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchApigwRoutes(ctx, types.ApigwRouteFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			// filters
			for _, n := range nn {
				gwf, _, err := s.SearchApigwFilters(ctx, types.ApigwFilterFilter{RouteID: n.ID})
				if err != nil {
					return &auxRsp{
						err: err,
					}
				}

				mm = append(mm, newAPIGateway(n, gwf, d.ux))
				d.resourceID = append(d.resourceID, n.ID)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

func (d *systemDecoder) decodeReports(ctx context.Context, s store.Storer, ff []*reportFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.ReportSet
	var fn types.ReportFilter
	var err error

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchReports(ctx, types.ReportFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				mm = append(mm, newReport(n, d.ux))
				d.resourceID = append(d.resourceID, n.ID)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

func (d *systemDecoder) decodeApplications(ctx context.Context, s store.Storer, ff []*applicationFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.ApplicationSet
	var fn types.ApplicationFilter
	var err error

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchApplications(ctx, types.ApplicationFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				// Index users
				err = d.ux.add(
					ctx,
					n.OwnerID,
				)

				mm = append(mm, newApplication(n, d.ux))
				d.resourceID = append(d.resourceID, n.ID)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}
func (d *systemDecoder) decodeSettings(ctx context.Context, s store.Storer, ff []*settingFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.SettingValueSet
	var err error

	for _, f := range ff {
		aux := *f

		for {
			nn, _, err = s.SearchSettings(ctx, types.SettingsFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				// Index users
				err = d.ux.add(
					ctx,
					n.OwnedBy,
					n.UpdatedBy,
				)

				mm = append(mm, newSetting(n, d.ux))
			}
			// mm = append(mm, NewSettings(nn))

			break
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

func (d *systemDecoder) decodeRbac(ctx context.Context, s store.Storer, ff []*rbacFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn rbac.RuleSet
	var err error

	c := func(r *resource.Ref, f *rbacFilter, path ...*resource.Ref) (bool, error) {
		if r == nil {
			return true, nil
		}

		// the first identifier is the most specific .. the ID
		id, err := cast.ToUint64E(r.Identifiers.First())
		if err != nil {
			return false, err
		}

		if r.ResourceType == composeTypes.ModuleFieldResourceType {
			id, err = cast.ToUint64E(path[1].Identifiers.First())
			if err != nil {
				return false, err
			}
		}

		return f.resourceID[id], nil
	}

	for _, f := range ff {
		aux := *f

		for {
			nn, _, err = s.SearchRbacRules(ctx, rbac.RuleFilter(aux.RuleFilter))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				r, err := newRbacRule(n)
				rt := strings.Split(r.refRbacResource, "/")[0]
				if err != nil {
					return &auxRsp{
						err: err,
					}
				}

				// somesort of a generic rule; no specifc resource
				// @todo check for pp inclusion!!
				if r.refRbacRes == nil && len(r.refPathRes) == 0 {
					mm = append(mm, r)
				} else {
					// Check the resource ref and the path refs for validity.
					// ComposeRecords are fetched in chunks so this check is not valid here.
					if rt != composeTypes.RecordResourceType {
						if ok, err := c(r.refRbacRes, f, r.refPathRes...); err != nil {
							return &auxRsp{
								err: err,
							}
						} else if !ok {
							continue
						}
					}

					for _, p := range r.refPathRes {
						if ok, err := c(p, f); err != nil {
							return &auxRsp{
								err: err,
							}
						} else if !ok {
							continue
						}
					}

					mm = append(mm, r)
				}
			}

			break
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

func (d *systemDecoder) decodeResourceTranslation(ctx context.Context, s store.Storer, ff []*resourceTranslationFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var (
		resMap = make(map[string]*resourceTranslation)

		nn  types.ResourceTranslationSet
		fn  types.ResourceTranslationFilter
		err error
	)

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchResourceTranslations(ctx, types.ResourceTranslationFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				if _, ok := resMap[n.Resource]; !ok {
					resMap[n.Resource], err = newResourceTranslation(types.ResourceTranslationSet{n})
					if err != nil {
						return &auxRsp{
							err: err,
						}
					}

					// parse the resource wo we can define relations/check if we can unmarshal it
					_, ref, pp, err := resource.ParseResourceTranslation(n.Resource)
					resMap[n.Resource].refLocaleRes = ref
					resMap[n.Resource].refPathRes = pp
					if err != nil {
						return &auxRsp{
							err: err,
						}
					}

					continue
				}

				resMap[n.Resource].locales = append(resMap[n.Resource].locales, n)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}

			break
		}
	}

	for _, lr := range resMap {
		mm = append(mm, lr)
	}

	return &auxRsp{
		mm: mm,
	}
}

func (df *DecodeFilter) systemFromResource(rr ...string) *DecodeFilter {
	for _, r := range rr {
		if !strings.HasPrefix(r, "system") {
			continue
		}

		id := ""
		if strings.Count(r, ":") == 2 && !strings.HasSuffix(r, "*") {
			// There is an identifier
			aux := strings.Split(r, ":")

			id = aux[len(aux)-1]
			r = strings.Join(aux[:len(aux)-1], ":")
		}

		switch strings.ToLower(r) {
		case "system:role":
			df = df.Roles(&types.RoleFilter{
				Query: id,
			})
		case "system:user":
			df = df.Users(&types.UserFilter{
				Query:    id,
				AllKinds: true,
			})
		case "system:template":
			df = df.Templates(&types.TemplateFilter{
				Handle: id,
			})
			templateID, err := cast.ToUint64E(id)
			if err == nil && templateID > 0 {
				df = df.Templates(&types.TemplateFilter{
					TemplateID: []uint64{templateID},
				})
			}
		case "system:apigw-route":
			df = df.APIGWRoutes(&types.ApigwRouteFilter{
				Route: id,
			})
		case "system:report":
			df = df.Reports(&types.ReportFilter{
				Handle: id,
			})
			reportID, err := cast.ToUint64E(id)
			if err == nil && reportID > 0 {
				df = df.Reports(&types.ReportFilter{
					ReportID: []uint64{reportID},
				})
			}

		case "system:application":
			df = df.Applications(&types.ApplicationFilter{
				Query: id,
			})
		case "system:setting":
			df = df.Settings(&types.SettingsFilter{})
		case "system:rbac":
			df = df.Rbac(&rbac.RuleFilter{})
		case "system:resource-translation":
			df = df.Rbac(&rbac.RuleFilter{})
		}
	}

	return df
}

func (df *DecodeFilter) systemFromRef(rr ...*resource.Ref) *DecodeFilter {
	for _, r := range rr {
		if strings.Index(r.ResourceType, "system") < 0 {
			continue
		}

		switch r.ResourceType {
		case types.RoleResourceType:
			for _, i := range r.Identifiers.StringSlice() {
				df = df.Roles(&types.RoleFilter{
					Query: i,
				})
			}
		case types.UserResourceType:
			for _, i := range r.Identifiers.StringSlice() {
				df = df.Users(&types.UserFilter{
					Query:    i,
					AllKinds: true,
				})
			}
		case types.TemplateResourceType:
			for _, i := range r.Identifiers.StringSlice() {
				df = df.Templates(&types.TemplateFilter{
					Handle: i,
				})
				templateID, err := cast.ToUint64E(i)
				if err == nil && templateID > 0 {
					df = df.Templates(&types.TemplateFilter{
						TemplateID: []uint64{templateID},
					})
				}
			}
		case types.ApplicationResourceType:
			for _, i := range r.Identifiers.StringSlice() {
				df = df.Applications(&types.ApplicationFilter{
					Query: i,
				})
			}
		}
	}

	return df
}

// Roles adds a new RoleFilter
func (df *DecodeFilter) Roles(f *types.RoleFilter) *DecodeFilter {
	if df.roles == nil {
		df.roles = make([]*roleFilter, 0, 1)
	}
	df.roles = append(df.roles, (*roleFilter)(f))
	return df
}

// Users adds a new UserFilter
func (df *DecodeFilter) Users(f *types.UserFilter) *DecodeFilter {
	if df.users == nil {
		df.users = make([]*userFilter, 0, 1)
	}
	df.users = append(df.users, (*userFilter)(f))
	return df
}

// Templates adds a new TemplateFilter
func (df *DecodeFilter) Templates(f *types.TemplateFilter) *DecodeFilter {
	if df.templates == nil {
		df.templates = make([]*templateFilter, 0, 1)
	}
	df.templates = append(df.templates, (*templateFilter)(f))
	return df
}

func (df *DecodeFilter) APIGWRoutes(f *types.ApigwRouteFilter) *DecodeFilter {
	if df.apiGwRoutes == nil {
		df.apiGwRoutes = make([]*apiGwRouteFilter, 0, 1)
	}
	df.apiGwRoutes = append(df.apiGwRoutes, (*apiGwRouteFilter)(f))
	return df
}

func (df *DecodeFilter) Reports(f *types.ReportFilter) *DecodeFilter {
	if df.reports == nil {
		df.reports = make([]*reportFilter, 0, 1)
	}
	df.reports = append(df.reports, (*reportFilter)(f))
	return df
}

// Applications adds a new ApplicationFilter
func (df *DecodeFilter) Applications(f *types.ApplicationFilter) *DecodeFilter {
	if df.applications == nil {
		df.applications = make([]*applicationFilter, 0, 1)
	}
	df.applications = append(df.applications, (*applicationFilter)(f))
	return df
}

// Settings adds a new SettingsFilter
func (df *DecodeFilter) Settings(f *types.SettingsFilter) *DecodeFilter {
	if df.settings == nil {
		df.settings = make([]*settingFilter, 0, 1)
	}
	df.settings = append(df.settings, (*settingFilter)(f))
	return df
}

// Rbac adds a new RuleFilter
func (df *DecodeFilter) Rbac(f *rbac.RuleFilter) *DecodeFilter {
	if df.rbac == nil {
		df.rbac = make([]*rbacFilter, 0, 1)
	} else {
		// There can only be a single rbac filter
		// since it makes no sense to have multiple of
		return df
	}

	df.rbac = append(df.rbac, &rbacFilter{RuleFilter: *f})
	return df
}

// ResourceTranslation adds a new ResourceTranslationFilter
func (df *DecodeFilter) ResourceTranslation(f *types.ResourceTranslationFilter) *DecodeFilter {
	if df.resourceTranslations == nil {
		df.resourceTranslations = make([]*resourceTranslationFilter, 0, 1)
	} else {
		// There can only be a single resourceTranslations filter
		// since it makes no sense to have multiple of
		return df
	}

	df.resourceTranslations = append(df.resourceTranslations, (*resourceTranslationFilter)(f))
	return df
}

func (df *DecodeFilter) RbacStrict(f *rbac.RuleFilter) *DecodeFilter {
	if df.rbac == nil {
		df.rbac = make([]*rbacFilter, 0, 1)
	} else {
		// There can only be a single rbac filter
		// since it makes no sense to have multiple of
		return df
	}

	df.rbac = append(df.rbac, &rbacFilter{RuleFilter: *f, strict: true})
	return df
}

// allowRbacResource adds a new resource identifier to supported resource rules
func (df *DecodeFilter) allowRbacResource(id ...uint64) {
	if df.rbac == nil || len(df.rbac) == 0 {
		return
	}
	rf := df.rbac[0]

	if rf.resourceID == nil {
		rf.resourceID = make(map[uint64]bool)
	}
	for _, i := range id {
		rf.resourceID[i] = true
	}
}
