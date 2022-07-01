package app

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"go.uber.org/zap"
)

func (app *CortezaApp) initDAL(ctx context.Context, log *zap.Logger) (err error) {
	// Verify that primary store is connected
	// or return error
	if app.Store == nil {
		return fmt.Errorf("primary store not connected")
	}

	// Init DAL and prepare default connection
	svc, err := dal.New(log.Named("dal"), app.Opt.Environment.IsDevelopment())
	if err != nil {
		return err
	}

	dal.SetGlobal(svc)

	return
}
