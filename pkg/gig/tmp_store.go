package gig

import (
	"context"
	"fmt"
)

func getGig(ctx context.Context, id uint64) (g Gig, err error) {
	var ok bool
	if g, ok = gigStore[id]; !ok {
		err = fmt.Errorf("gig not found: %d", id)
		return
	}

	return
}

func updateGig(ctx context.Context, old Gig) (g Gig, err error) {
	gigStore[old.ID] = old
	g = old
	return
}

func deleteGig(ctx context.Context, old *Gig) (err error) {
	delete(gigStore, old.ID)
	return
}
