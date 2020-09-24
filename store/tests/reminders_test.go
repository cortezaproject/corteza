package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
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
	})

	t.Run("ordered search", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
