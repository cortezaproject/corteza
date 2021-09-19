package reporter

import (
	"testing"
)

func Test_model_describe_basic(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, r, _   = loadScenario(ctx, s, t, h)
		ff        = describeNoErr(ctx, h, m, r.Describe...)
	)

	t.Run("load/users", func(t *testing.T) {
		ff := ff.FilterBySource("users")
		h.a.Len(ff, 1)

		h.a.Equal(
			"id<Record>, join_key<String>, first_name<String>, last_name<String>, email<String>, number_of_numbers<Number>, dob<DateTime>, createdAt<Date>, createdBy<User>, updatedAt<Date>, updatedBy<User>, deletedAt<Date>, deletedBy<User>",
			ff[0].Columns.String(),
		)
	})

	t.Run("load/jobs", func(t *testing.T) {
		ff := ff.FilterBySource("jobs")
		h.a.Len(ff, 1)

		h.a.Equal(
			"id<Record>, name<String>, type<Select>, cost<Number>, time_spent<Number>, usr<String>, createdAt<Date>, createdBy<User>, updatedAt<Date>, updatedBy<User>, deletedAt<Date>, deletedBy<User>",
			ff[0].Columns.String(),
		)
	})

	t.Run("join", func(t *testing.T) {
		ff := ff.FilterBySource("joined")
		h.a.Len(ff, 2)

		h.a.Equal(
			"id<Record>, join_key<String>, first_name<String>, last_name<String>, email<String>, number_of_numbers<Number>, dob<DateTime>, createdAt<Date>, createdBy<User>, updatedAt<Date>, updatedBy<User>, deletedAt<Date>, deletedBy<User>",
			ff.FilterByRef("users")[0].Columns.String(),
		)

		h.a.Equal(
			"id<Record>, name<String>, type<Select>, cost<Number>, time_spent<Number>, usr<String>, createdAt<Date>, createdBy<User>, updatedAt<Date>, updatedBy<User>, deletedAt<Date>, deletedBy<User>",
			ff.FilterByRef("jobs")[0].Columns.String(),
		)
	})
}
