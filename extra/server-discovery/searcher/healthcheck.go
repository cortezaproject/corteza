package searcher

import (
	"context"
	"io"
)

// Healthcheck for searcher
func Healthcheck(ctx context.Context) error {
	if DefaultApiClient == nil {
		return nil
	}

	esc := DefaultEsClient

	res, err := esc.Ping(esc.Ping.WithContext(ctx))
	if err = validElasticResponse(res, err); err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	return nil
}
