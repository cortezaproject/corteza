package envoy

import (
	"context"
	"sync"
)

type (
	Provider interface {
		Next(ctx context.Context) (*ResourceState, error)
	}
)

// @todo errors!
func Encode(ctx context.Context, p Provider, ee ...Encoder) error {
	var err error
	rcc := make([]Rc, len(ee))

	var wg sync.WaitGroup
	wg.Add(len(ee))

	for i, e := range ee {
		rcc[i] = make(Rc)

		go func() {
			defer wg.Done()

			eErr := e.Encode(ctx, rcc[i])
			if err != nil {
				err = eErr
			}
		}()
	}

	for {
		rs, err := p.Next(ctx)
		if err != nil {
			return err
		}

		for _, rc := range rcc {
			rc <- rs
		}
		if rs == nil {
			break
		}
	}

	wg.Wait()
	return nil
}
