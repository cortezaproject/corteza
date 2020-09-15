package store

import (
	"context"
	"go.uber.org/zap"
)

type (
	// Storer interface combines interfaces of all supported store interfaces
	Storer interface {
		// SetLogger sets new logging facility
		//
		// Store facility should fallback to logger.Default when no logging facility is set
		//
		// Intentionally closely coupled with Zap logger since this is not some public lib
		// and it's highly unlikely we'll support different/multiple logging "backend"
		SetLogger(*zap.Logger)

		// Tx is a transaction handler
		Tx(context.Context, func(context.Context, Storer) error) error

		// All generated store interfaces
		storerGenerated
	}
)
