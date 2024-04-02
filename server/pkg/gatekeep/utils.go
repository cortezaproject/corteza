package gatekeep

import (
	"context"
	"fmt"
	"sync"
)

func AwaitLocks(ctx context.Context, svc gksvc, cc ...Constraint) (locks []Lock, errs []error) {
	var (
		// state LockState
		wg = &sync.WaitGroup{}
	)

	locks = make([]Lock, len(cc))
	errs = make([]error, len(cc))

	for i, c := range cc {
		locks[i], errs[i] = svc.Lock(ctx, c)
		if errs[i] != nil {
			continue
		}

		if locks[i].State == lockStateLocked {
			continue
		}

		if locks[i].State == lockStateFailed {
			// @note should never happen
			errs[i] = fmt.Errorf("lock failed")
			continue
		}

		wg.Add(1)
		lID := svc.Subscribe(func(evt Event) {
			if !evt.Lock.matchesConstraints(c) {
				return
			}

			if evt.Lock.State == lockStateQueued {
				return
			}

			if evt.Lock.State == lockStateFailed {
				errs[i] = fmt.Errorf("lock failed")
			}

			wg.Done()
		})
		defer svc.Unsubscribe(lID)
	}

	wg.Wait()

	for i, err := range errs {
		if err != nil {
			locks[i] = Lock{}
		}
	}

	return
}
