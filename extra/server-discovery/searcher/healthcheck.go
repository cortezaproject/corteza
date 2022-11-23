package searcher

import (
	"context"
	"fmt"
	"io"
)

// Healthcheck for searcher
func Healthcheck(ctx context.Context) error {
	if DefaultEs == nil {
		return nil
	}

	esc, err := DefaultEs.Client()
	if esc == nil || err != nil {
		return fmt.Errorf("stopped")
	}

	res, err := esc.Ping(esc.Ping.WithContext(ctx))
	if err = validElasticResponse(res, err); err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	return nil
}
