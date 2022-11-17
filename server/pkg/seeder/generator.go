package seeder

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/errors"
)

type (
	fn func(ctx context.Context) error
)

const (
	maxRetry = 10
)

// Generator function generates synthetic data by calling generator function
//
// It will retry on error up to maxRetry times
func Generator(ctx context.Context, generator fn, total uint) (err error) {
	var (
		retry uint
	)

	if generator == nil {
		return fmt.Errorf("generator is nil")
	}

	for total > 0 {
		if retry > maxRetry {
			return fmt.Errorf("max retry count (%d) reached", maxRetry)
		}

		err = generator(ctx)

		if err == nil {
			retry = 0
			total--
			continue
		}

		if errors.IsDuplicateData(err) {
			retry++
			continue
		}

		return err
	}

	return nil
}
