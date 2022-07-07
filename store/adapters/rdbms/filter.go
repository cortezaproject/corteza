package rdbms

import (
	"fmt"
	"strings"
	"time"

	automationType "github.com/cortezaproject/corteza-server/automation/types"
	composeType "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	discoveryType "github.com/cortezaproject/corteza-server/pkg/discovery/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	labelsType "github.com/cortezaproject/corteza-server/pkg/label/types"
	systemType "github.com/cortezaproject/corteza-server/system/types"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func DefaultFilters() (f *extendedFilters) {
	f = &extendedFilters{}

	f.Actionlog = func(s *Store, f actionlog.Filter) (ee []goqu.Expression, _ actionlog.Filter, err error) {
		if ee, f, err = ActionlogFilter(f); err != nil {
			return
		}

		// make sure we always sort ID, descending
		if f.Sorting, err = filter.NewSorting("id DESC"); err != nil {
			return
		}

		if f.BeforeActionID > 0 {
			ee = append(ee, goqu.C("id").Lt(f.BeforeActionID))
		}

		if f.FromTimestamp != nil {
			ee = append(ee, goqu.C("ts").Gte(f.FromTimestamp))
		}

		if f.ToTimestamp != nil {
			ee = append(ee, goqu.C("ts").Lte(f.FromTimestamp))
		}

		if f.Limit == 0 || f.Limit > MaxLimit {
			f.Limit = MaxLimit
		}

		return ee, f, err
	}

	f.Application = func(s *Store, f systemType.ApplicationFilter) (ee []goqu.Expression, _ systemType.ApplicationFilter, err error) {
		if ee, f, err = ApplicationFilter(f); err != nil {
			return
		}

		if len(f.FlaggedIDs) > 0 {
			ee = append(ee, goqu.C("id").In(f.FlaggedIDs))
		}

		return ee, f, err
	}

	f.AutomationSession = func(s *Store, f automationType.SessionFilter) (ee []goqu.Expression, _ automationType.SessionFilter, err error) {
		if ee, f, err = AutomationSessionFilter(f); err != nil {
			return
		}

		if len(f.Status) > 0 {
			ee = append(ee, goqu.C("status").In(f.Status))
		}

		return ee, f, err
	}

	f.ComposeAttachment = func(s *Store, f composeType.AttachmentFilter) (ee []goqu.Expression, _ composeType.AttachmentFilter, err error) {
		if ee, f, err = ComposeAttachmentFilter(f); err != nil {
			return
		}

		switch f.Kind {
		case composeType.PageAttachment:
			// @todo implement filtering by page
			if f.PageID > 0 {
				err = fmt.Errorf("filtering by pageID not implemented")
				return
			}

		case composeType.RecordAttachment:
			panic("@todo pending implementation")
			//query = query.
			//	Join("compose_record_value AS v ON (v.ref = a.id)")
			//
			//if f.ModuleID > 0 {
			//	query = query.
			//		Join("compose_record AS r ON (r.id = v.record_id)").
			//		Where(squirrel.Eq{"r.module_id": f.ModuleID})
			//}
			//
			//if f.RecordID > 0 {
			//	query = query.Where(squirrel.Eq{"v.record_id": f.RecordID})
			//}
			//
			//if f.FieldName != "" {
			//	query = query.Where(squirrel.Eq{"v.name": f.FieldName})
			//}

		default:
			err = fmt.Errorf("unsupported kind value")
			return
		}

		if f.Filter != "" {
			err = fmt.Errorf("filtering by filter not implemented")
			return
		}

		return ee, f, nil
	}

	f.ComposePage = func(s *Store, f composeType.PageFilter) (ee []goqu.Expression, _ composeType.PageFilter, err error) {
		if ee, f, err = ComposePageFilter(f); err != nil {
			return
		}

		if f.ParentID > 0 {
			ee = append(ee, goqu.C("self_id").Eq(f.ParentID))
		} else if f.Root {
			ee = append(ee, goqu.C("self_id").Eq(0))
		}

		return ee, f, nil
	}

	f.Label = func(store *Store, f labelsType.LabelFilter) (ee []goqu.Expression, _ labelsType.LabelFilter, err error) {
		if ee, f, err = LabelFilter(f); err != nil {
			return
		}

		if len(f.Filter) > 0 {
			values := make([]goqu.Expression, 0, len(f.Filter))

			for k, v := range f.Filter {
				values = append(values, exp.Ex{"name": k, "value": v})
			}

			ee = append(ee, goqu.Or(values...))
		}

		return ee, f, nil
	}

	f.Reminder = func(s *Store, f systemType.ReminderFilter) (ee []goqu.Expression, _ systemType.ReminderFilter, err error) {
		if ee, f, err = ReminderFilter(f); err != nil {
			return
		}

		if f.ExcludeDismissed {
			ee = append(ee, goqu.C("dismissed_at").IsNull())
		}

		if !f.IncludeDeleted {
			ee = append(ee, goqu.C("deleted_at").IsNull())
		}

		if f.ScheduledOnly {
			ee = append(ee, goqu.C("remind_at").IsNotNull())
		}

		if f.Resource != "" {
			ee = append(ee, goqu.C("resource").Like(f.Resource+"%"))
		}

		if f.ScheduledFrom != nil {
			ee = append(ee, goqu.C("remind_at").Gte(f.ScheduledFrom.Format(time.RFC3339)))
		}
		if f.ScheduledUntil != nil {
			ee = append(ee, goqu.C("remind_at").Lte(f.ScheduledUntil.Format(time.RFC3339)))
		}

		return ee, f, nil
	}

	f.ResourceTranslation = func(s *Store, f systemType.ResourceTranslationFilter) (ee []goqu.Expression, _ systemType.ResourceTranslationFilter, err error) {
		if ee, f, err = ResourceTranslationFilter(f); err != nil {
			return
		}

		if f.ResourceType != "" {
			ee = append(ee, goqu.C("resource").Like(f.ResourceType+"%"))
		}

		return ee, f, nil
	}

	f.Role = func(s *Store, f systemType.RoleFilter) (ee []goqu.Expression, _ systemType.RoleFilter, err error) {
		if ee, f, err = RoleFilter(f); err != nil {
			return
		}

		if f.MemberID > 0 {
			memberships := roleMemberSelectQuery(s.Dialect).
				Select("rel_role").
				Where(goqu.C("rel_user").In(f.MemberID))

			ee = append(ee, goqu.C("id").In(memberships))
		}

		return ee, f, nil
	}

	f.User = func(s *Store, f systemType.UserFilter) (ee []goqu.Expression, _ systemType.UserFilter, err error) {
		if ee, f, err = UserFilter(f); err != nil {
			return
		}

		if !f.AllKinds {
			ee = append(ee, goqu.C("kind").Eq(f.Kind))
		}

		if len(f.RoleID) > 0 {
			members := roleMemberSelectQuery(s.Dialect).
				Select("rel_user").
				Where(goqu.C("rel_role").In(f.RoleID))

			ee = append(ee, goqu.C("id").In(members))
		}

		return ee, f, nil
	}

	f.SettingValue = func(s *Store, f systemType.SettingsFilter) (ee []goqu.Expression, _ systemType.SettingsFilter, err error) {
		if ee, f, err = SettingValueFilter(f); err != nil {
			return
		}

		if len(f.Prefix) > 0 {
			ee = append(ee, goqu.C("name").Like(f.Prefix+"%"))
		}

		return ee, f, nil
	}

	f.FederationExposedModule = func(s *Store, f types.ExposedModuleFilter) (ee []goqu.Expression, _ types.ExposedModuleFilter, err error) {
		if ee, f, err = FederationExposedModuleFilter(f); err != nil {
			return
		}

		if f.LastSync > 0 {
			t := time.Unix(int64(f.LastSync), 0)

			if !t.IsZero() {
				ts := t.UTC().Format(time.RFC3339)
				ee = append(ee, goqu.Or(
					goqu.C("updated_at").Gte(ts),
					goqu.C("created_at").Gte(ts),
				))
			}
		}

		return ee, f, nil
	}

	f.ResourceActivity = func(s *Store, f discoveryType.ResourceActivityFilter) (ee []goqu.Expression, _ discoveryType.ResourceActivityFilter, err error) {
		if ee, f, err = ResourceActivityFilter(f); err != nil {
			return
		}

		// Always sort by ID descending
		//query = query.OrderBy("id DESC")

		if f.FromTimestamp != nil {
			ee = append(ee, goqu.C("ts").Gte(f.ToTimestamp))
		}

		if f.ToTimestamp != nil {
			ee = append(ee, goqu.C("ts").Lte(f.ToTimestamp))
		}

		if f.Limit == 0 || f.Limit > MaxLimit {
			f.Limit = MaxLimit
		}

		return ee, f, err
	}

	f.DataPrivacyRequest = func(s *Store, f systemType.DataPrivacyRequestFilter) (ee []goqu.Expression, _ systemType.DataPrivacyRequestFilter, err error) {
		if ee, f, err = DataPrivacyRequestFilter(f); err != nil {
			return
		}

		if len(f.Kind) > 0 {
			ee = append(ee, goqu.C("kind").In(f.Kind))
		}

		if len(f.Status) > 0 {
			ee = append(ee, goqu.C("status").In(f.Status))
		}

		if f.Limit == 0 || f.Limit > MaxLimit {
			f.Limit = MaxLimit
		}

		return ee, f, err
	}

	return
}

func Order(sort filter.SortExprSet, sortables map[string]string) (oo []exp.OrderedExpression, err error) {
	return order(sort, sortables)
}

func order(sort filter.SortExprSet, sortables map[string]string) (oo []exp.OrderedExpression, err error) {
	var (
		has bool
		col string
	)

	for _, s := range sort {
		if len(sortables) > 0 {
			if col, has = sortables[strings.ToLower(s.Column)]; !has {
				return nil, fmt.Errorf("column %q is not sortable", s.Column)
			}
			s.Column = col
		}

		if s.Descending {
			oo = append(oo, exp.NewOrderedExpression(goqu.I(s.Column), exp.DescSortDir, exp.NoNullsSortType))
		} else {
			oo = append(oo, exp.NewOrderedExpression(goqu.I(s.Column), exp.AscDir, exp.NoNullsSortType))
		}
	}

	return
}

func stateNilComparison(lit string, fs filter.State) goqu.Expression {
	switch fs {
	case filter.StateExclusive:
		// only not-null values
		return goqu.Literal(lit).IsNotNull()

	case filter.StateInclusive:
		// no filter
		return nil

	default:
		// exclude all non-null values
		return goqu.Literal(lit).IsNull()
	}
}

func stateFalseComparison(lit string, fs filter.State) goqu.Expression {
	switch fs {
	case filter.StateExcluded:
		// only true
		return goqu.Literal(lit).IsTrue()

	case filter.StateExclusive:
		// only false
		return goqu.Literal(lit).IsFalse()

	default:
		return nil
	}
}
