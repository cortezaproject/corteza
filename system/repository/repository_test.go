package repository

import (
	"context"

	"testing"
)

func tx(t *testing.T, f func() error) {
	db := DB(context.Background())

	if err := db.Begin(); err != nil {
		t.Errorf("Could not begin transaction: %v", err)

	}

	if err := f(); err != nil {
		t.Errorf("Test transaction resulted in an error: %v", err)
	}

	if err := db.Rollback(); err != nil {
		t.Errorf("Could not rollback transaction: %v", err)
	}
}
