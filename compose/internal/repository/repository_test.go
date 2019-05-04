// +build unit integration

package repository

import (
	"context"
	"testing"

	"github.com/crusttech/crust/internal/test"
)

func TestRepository(t *testing.T) {
	repo := &repository{}
	repo.With(context.Background(), nil)
}

func tx(t *testing.T, f func() error) {
	var err error
	db := DB(context.Background())

	err = db.Begin()
	test.Assert(t, err == nil, "Could not begin transaction: %+v", err)

	err = f()
	test.Assert(t, err == nil, "Test transaction resulted in an error: %+v", err)

	err = db.Rollback()
	test.Assert(t, err == nil, "Could not rollback transaction: %+v", err)
}
