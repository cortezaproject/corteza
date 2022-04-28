package crs

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/crs/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"go.uber.org/zap"
)

type (
	// the core struct that outlines the compose record store facility
	composeRecordStore struct {
		stores map[uint64]StoreConnection

		// Indexed by corresponding storeID
		models map[uint64]data.ModelSet

		primary StoreConnection

		logger *zap.Logger
		inDev  bool
	}

	crsDefiner interface {
		ComposeRecordStoreID() uint64
		StoreDSN() string
		Capabilities() capabilities.Set
	}
)

const (
	defaultStoreID uint64 = 0
)

// ComposeRecordStore initializes a fresh record store where the given store serves as the default
func ComposeRecordStore(ctx context.Context, log *zap.Logger, inDev bool, primary crsDefiner, stores ...crsDefiner) (*composeRecordStore, error) {
	crs := &composeRecordStore{
		stores:  make(map[uint64]StoreConnection),
		models:  make(map[uint64]data.ModelSet),
		primary: nil,

		logger: log,
		inDev:  inDev,
	}

	var err error

	crs.primary, err = connect(ctx, log, primary, inDev)
	if err != nil {
		return nil, err
	}

	return crs, crs.AddStore(ctx, stores...)
}

// AddStore registers the given store definitions as compose record stores
func (crs *composeRecordStore) AddStore(ctx context.Context, definers ...crsDefiner) (err error) {
	for _, definer := range definers {
		crs.stores[definer.ComposeRecordStoreID()], err = connect(ctx, crs.logger, definer, crs.inDev)
		if err != nil {
			return
		}
	}

	return nil
}

// RemoveStore removes the given store definition as a compose record store
func (crs *composeRecordStore) RemoveStore(ctx context.Context, storeID uint64, storeIDs ...uint64) (err error) {
	for _, storeID := range append(storeIDs, storeID) {
		s := crs.stores[storeID]
		if s == nil {
			return fmt.Errorf("can not remove compose record store %d: store does not exist", storeID)
		}

		// Potential cleanups
		if err = s.Close(ctx); err != nil {
			return
		}

		// Remove from registry
		delete(crs.stores, storeID)
	}

	return nil
}

// ---
// Utilities

func (crs *composeRecordStore) getModel(store uint64, ident string) *data.Model {
	for _, model := range crs.models[store] {
		if model.Ident == ident {
			return model
		}
	}

	return nil
}

// getStore returns a store for the given identifier/capabilities combination
func (crs *composeRecordStore) getStore(ctx context.Context, storeID uint64, cc ...capabilities.Capability) (store StoreConnection, can capabilities.Set, err error) {
	err = func() error {
		// get the requested store
		if storeID == defaultStoreID {
			store = crs.primary
		} else {
			store = crs.stores[storeID]
		}
		if store == nil {
			return fmt.Errorf("could not get store %d: store does not exist", storeID)
		}

		// check if store supports requested capabilities
		if !store.Can(cc...) {
			return fmt.Errorf("store does not support requested capabilities: %v", capabilities.Set(cc).Diff(store.Capabilities()))
		}
		can = store.Capabilities()
		return nil
	}()

	if err != nil {
		err = fmt.Errorf("could not connect to store %d: %v", storeID, err)
		return
	}

	return
}
