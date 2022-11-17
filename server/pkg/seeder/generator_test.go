package seeder

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerate(t *testing.T) {
	var (
		req = require.New(t)
	)

	req.Error(Generator(
		context.Background(),
		nil,
		0,
	), "generator is nil")

	count := 0
	req.NoError(Generator(
		context.Background(),
		func(ctx context.Context) error {
			count++
			return nil
		},
		5,
	))
	req.Equal(5, count)

	count = 0
	req.NoError(Generator(
		context.Background(),
		func(ctx context.Context) error {
			count++
			if count%2 == 0 {
				return nil
			}
			return errors.DuplicateData("foo")
		},
		5,
	))
	req.Equal(10, count)
}
