package repository

import (
	"context"
	"testing"
)

func TestEvents(t *testing.T) {
	repo := &repository{}
	repo.With(context.Background(), nil)
}
