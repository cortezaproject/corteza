package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testReminders(t *testing.T, s remindersStore) {
	var (
		ctx = context.Background()
		req = require.New(t)

		//err      error
		reminder *types.Reminder
	)

	t.Run("create", func(t *testing.T) {
		reminder = &types.Reminder{
			ID:         42,
			CreatedAt:  time.Now(),
			AssignedAt: time.Now(),
		}
		req.NoError(s.CreateReminder(ctx, reminder))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		fetched, err := s.LookupReminderByID(ctx, reminder.ID)
		req.NoError(err)
		req.Equal(reminder.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		reminder = &types.Reminder{
			ID:         42,
			CreatedAt:  time.Now(),
			AssignedAt: time.Now(),
		}
		req.NoError(s.UpdateReminder(ctx, reminder))
	})

	//t.Run("delete/undelete", func(t *testing.T) {
	//	ID := reminder.ID
	//	reminder, err = s.LookupReminderByID(ctx, ID)
	//	req.NoError(err)
	//
	//	req.NoError(s.DeleteReminderByID(ctx, ID))
	//	reminder, err = s.LookupReminderByID(ctx, ID)
	//	req.NoError(err)
	//	req.NotNil(reminder.DeletedAt)
	//
	//	req.NoError(s.UndeleteReminderByID(ctx, ID))
	//	reminder, err = s.LookupReminderByID(ctx, ID)
	//	req.NoError(err)
	//	req.Nil(reminder.DeletedAt)
	//})

	t.Run("search", func(t *testing.T) {
		set, f, err := s.SearchReminders(ctx, types.ReminderFilter{})
		req.NoError(err)
		req.Len(set, 1)
		req.Equal(uint(1), f.Count)
	})

	t.Run("search by *", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("ordered search", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
