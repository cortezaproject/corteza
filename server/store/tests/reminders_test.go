package tests

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func testReminders(t *testing.T, s store.Reminders) {
	var (
		ctx = context.Background()

		makeNew = func(nn ...string) *types.Reminder {
			// minimum data set for new user
			name := strings.Join(nn, "")
			thisID := id.Next()
			return &types.Reminder{
				ID:         thisID,
				Resource:   "resource+" + name,
				AssignedTo: thisID,
				CreatedAt:  time.Now(),
				AssignedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Reminder) {
			req := require.New(t)
			req.NoError(s.TruncateReminders(ctx))
			reminder := makeNew()
			req.NoError(s.CreateReminder(ctx, reminder))
			return req, reminder
		}

		truncAndFill = func(t *testing.T, l int) (*require.Assertions, types.ReminderSet) {
			req := require.New(t)
			req.NoError(s.TruncateReminders(ctx))

			set := make([]*types.Reminder, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew(string(rand.Bytes(10)))
			}

			req.NoError(s.CreateReminder(ctx, set...))
			return req, set
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.CreateReminder(ctx, makeNew()))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, reminder := truncAndCreate(t)
		fetched, err := s.LookupReminderByID(ctx, reminder.ID)
		req.NoError(err)
		req.Equal(reminder.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, reminder := truncAndCreate(t)
		req.NoError(s.UpdateReminder(ctx, reminder))
	})

	t.Run("search", func(t *testing.T) {
		t.Run("by ID", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, f, err := s.SearchReminders(ctx, types.ReminderFilter{ReminderID: []uint64{prefill[0].ID}})
			req.NoError(err)
			req.Equal([]uint64{prefill[0].ID}, f.ReminderID)
			req.Len(set, 1)
		})

		t.Run("by resource", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, _, err := s.SearchReminders(ctx, types.ReminderFilter{Resource: prefill[0].Resource})
			req.NoError(err)
			req.Len(set, 1)
		})

		t.Run("by assigned to", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)
			set, f, err := s.SearchReminders(ctx, types.ReminderFilter{AssignedTo: prefill[0].AssignedTo})
			req.NoError(err)
			req.Equal(prefill[0].AssignedTo, f.AssignedTo)
			req.Len(set, 1)
		})

		t.Run("by dismissed", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			prefill[0].DismissedAt = &(prefill[0].CreatedAt)
			s.UpdateReminder(ctx, prefill[0])

			set, _, err := s.SearchReminders(ctx, types.ReminderFilter{ExcludeDismissed: true})
			req.NoError(err)
			req.Len(set, 4)
		})

		t.Run("by scheduled", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			prefill[0].RemindAt = &(prefill[0].CreatedAt)
			s.UpdateReminder(ctx, prefill[0])

			set, _, err := s.SearchReminders(ctx, types.ReminderFilter{ScheduledOnly: true})
			req.NoError(err)
			req.Len(set, 1)
		})

		t.Run("with check", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)
			set, _, err := s.SearchReminders(ctx, types.ReminderFilter{
				Check: func(user *types.Reminder) (bool, error) {
					// simple check that matches with the first user from prefill
					return user.ID == prefill[0].ID, nil
				},
			})
			req.NoError(err)
			req.Len(set, 1)
			req.Equal(prefill[0].ID, set[0].ID)
		})

		t.Run("paging with check", func(t *testing.T) {
			// hide every third reminder from the check fn
			hide := func(prefill types.ReminderSet) (map[uint64]bool, int) {
				hidden, visible := make(map[uint64]bool), 0
				for i, r := range prefill {
					if i%3 == 0 {
						hidden[r.ID] = true
					} else {
						visible++
					}
				}
				return hidden, visible
			}

			t.Run("page holds no more than the limit and no duplicates", func(t *testing.T) {
				req, prefill := truncAndFill(t, 20)
				hidden, _ := hide(prefill)

				set, _, err := s.SearchReminders(ctx, types.ReminderFilter{
					Paging: filter.Paging{Limit: 5},
					Check:  func(r *types.Reminder) (bool, error) { return !hidden[r.ID], nil },
				})
				req.NoError(err)
				req.Len(set, 5)

				seen := make(map[uint64]bool)
				for _, r := range set {
					req.False(seen[r.ID], "reminder %d returned twice", r.ID)
					seen[r.ID] = true
				}
			})

			t.Run("walking every page yields each visible item once", func(t *testing.T) {
				req, prefill := truncAndFill(t, 20)
				hidden, visible := hide(prefill)

				var (
					cursor *filter.PagingCursor
					seen   = make(map[uint64]int)
				)

				for page := 0; page < 30; page++ {
					set, f, err := s.SearchReminders(ctx, types.ReminderFilter{
						Paging: filter.Paging{Limit: 3, PageCursor: cursor},
						Check:  func(r *types.Reminder) (bool, error) { return !hidden[r.ID], nil },
					})
					req.NoError(err)

					for _, r := range set {
						req.False(hidden[r.ID], "hidden reminder %d was returned", r.ID)
						seen[r.ID]++
					}

					if f.NextPage == nil {
						break
					}
					cursor = f.NextPage
				}

				req.Len(seen, visible)
				for id, times := range seen {
					req.Equal(1, times, "reminder %d returned %d times", id, times)
				}
			})

			t.Run("fills a page past a long run of rejected items", func(t *testing.T) {
				req, prefill := truncAndFill(t, 200)

				// hide everything but the last 5
				hidden := make(map[uint64]bool)
				for _, r := range prefill[:len(prefill)-5] {
					hidden[r.ID] = true
				}

				set, _, err := s.SearchReminders(ctx, types.ReminderFilter{
					Paging: filter.Paging{Limit: 3},
					Check:  func(r *types.Reminder) (bool, error) { return !hidden[r.ID], nil },
				})
				req.NoError(err)
				req.Len(set, 3)
			})
		})

		t.Run("total", func(t *testing.T) {
			t.Run("omitted when not requested", func(t *testing.T) {
				req, _ := truncAndFill(t, 5)

				set, f, err := s.SearchReminders(ctx, types.ReminderFilter{
					Paging: filter.Paging{Limit: 2},
				})
				req.NoError(err)
				req.Len(set, 2)
				req.Equal(uint(0), f.Total)
			})

			t.Run("without limit", func(t *testing.T) {
				req, _ := truncAndFill(t, 5)

				set, f, err := s.SearchReminders(ctx, types.ReminderFilter{
					Paging: filter.Paging{IncTotal: true},
				})
				req.NoError(err)
				req.Len(set, 5)
				req.Equal(uint(5), f.Total)
			})

			t.Run("with limit above count", func(t *testing.T) {
				req, _ := truncAndFill(t, 5)

				set, f, err := s.SearchReminders(ctx, types.ReminderFilter{
					Paging: filter.Paging{Limit: 10, IncTotal: true},
				})
				req.NoError(err)
				req.Len(set, 5)
				req.Equal(uint(5), f.Total)
			})

			t.Run("with limit equal to count", func(t *testing.T) {
				req, _ := truncAndFill(t, 5)

				set, f, err := s.SearchReminders(ctx, types.ReminderFilter{
					Paging: filter.Paging{Limit: 5, IncTotal: true},
				})
				req.NoError(err)
				req.Len(set, 5)
				req.Equal(uint(5), f.Total)
			})

			t.Run("with limit below count", func(t *testing.T) {
				req, _ := truncAndFill(t, 5)

				set, f, err := s.SearchReminders(ctx, types.ReminderFilter{
					Paging: filter.Paging{Limit: 2, IncTotal: true},
				})
				req.NoError(err)
				req.Len(set, 2)
				req.Equal(uint(5), f.Total)
			})

			t.Run("counts only what check keeps", func(t *testing.T) {
				req, prefill := truncAndFill(t, 5)

				set, f, err := s.SearchReminders(ctx, types.ReminderFilter{
					Paging: filter.Paging{Limit: 2, IncTotal: true},
					Check: func(r *types.Reminder) (bool, error) {
						return r.ID != prefill[0].ID && r.ID != prefill[1].ID, nil
					},
				})
				req.NoError(err)
				req.Len(set, 2)
				req.Equal(uint(3), f.Total)
			})

			t.Run("refused with a page cursor", func(t *testing.T) {
				req, _ := truncAndFill(t, 5)

				_, f, err := s.SearchReminders(ctx, types.ReminderFilter{
					Paging: filter.Paging{Limit: 2},
				})
				req.NoError(err)
				req.NotNil(f.NextPage)

				_, _, err = s.SearchReminders(ctx, types.ReminderFilter{
					Paging: filter.Paging{Limit: 2, IncTotal: true, PageCursor: f.NextPage},
				})
				req.Error(err)
			})
		})
	})

	t.Run("ordered search", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
