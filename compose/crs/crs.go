package crs

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/crs/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/data"
)

type (
	// the core struct that outlines the compose record store facility
	composeRecordStore struct {
		drivers []driver
		stores  map[uint64]*storeWrap

		// Indexed by corresponding storeID
		models map[uint64]data.ModelSet

		primary *storeWrap
	}

	storeWrap struct {
		store        connFunc
		driver       driver
		dsn          string
		capabilities capabilities.Set
	}

	connFunc func(ctx context.Context) (Store, error)

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
func ComposeRecordStore(ctx context.Context, primary crsDefiner, d driver, drivers ...driver) (*composeRecordStore, error) {
	drivers = append([]driver{d}, drivers...)

	crs := &composeRecordStore{
		drivers: drivers,
		stores:  make(map[uint64]*storeWrap),
		models:  make(map[uint64]data.ModelSet),
		primary: nil,
	}

	d = crs.getDriver(primary)
	if d == nil {
		return nil, fmt.Errorf("could not add default store: no supported driver found")
	}

	crs.primary = &storeWrap{
		store:        storeConnWrap(d, primary),
		driver:       d,
		dsn:          primary.StoreDSN(),
		capabilities: primary.Capabilities(),
	}

	return crs, nil
}

// AddStore registers the given store definitions as compose record stores
func (crs *composeRecordStore) AddStore(ctx context.Context, definers ...crsDefiner) error {
	for _, definer := range definers {
		if crs.stores[definer.ComposeRecordStoreID()] != nil {
			return fmt.Errorf("can not add compose record store %d: already defined", definer.ComposeRecordStoreID())
		}

		d := crs.getDriver(definer)
		if d == nil {
			return fmt.Errorf("could not add store %d: no supported driver found", definer.ComposeRecordStoreID())
		}

		crs.stores[definer.ComposeRecordStoreID()] = &storeWrap{
			store:        storeConnWrap(d, definer),
			driver:       d,
			dsn:          definer.StoreDSN(),
			capabilities: definer.Capabilities(),
		}
	}
	return nil
}

// RemoveStore removes the given store definition as a compose record store
func (crs *composeRecordStore) RemoveStore(ctx context.Context, storeIDs ...uint64) (err error) {
	for _, storeID := range storeIDs {
		s := crs.stores[storeID]
		if s == nil {
			return fmt.Errorf("can not remove compose record store %d: store does not exist", storeID)
		}

		// Any potential driver cleanup
		err = s.driver.Close(ctx, s.dsn)
		if err != nil {
			return
		}

		// Remove from registry
		delete(crs.stores, storeID)
	}

	return nil
}

// ---

func (crs *composeRecordStore) getModel(store uint64, ident string) *data.Model {
	for _, model := range crs.models[store] {
		if model.Ident == ident {
			return model
		}
	}

	return nil
}

// getStore returns a store for the given identifier/capabilities combination
func (crs *composeRecordStore) getStore(ctx context.Context, storeID uint64, cc ...capabilities.Capability) (store Store, can capabilities.Set, err error) {
	err = func() error {
		// get the requested store
		var wrap *storeWrap
		if storeID == defaultStoreID {
			wrap = crs.primary
		} else {
			wrap = crs.stores[storeID]
		}
		if wrap == nil {
			return fmt.Errorf("could not get store %d: store does not exist", storeID)
		}

		// check if store supports requested capabilities
		if !wrap.capabilities.IsSuperset(cc...) {
			return fmt.Errorf("store does not support requested capabilities: %v", capabilities.Set(cc).Diff(wrap.capabilities))
		}

		store, err = wrap.store(ctx)
		can = wrap.capabilities
		return nil
	}()

	if err != nil {
		err = fmt.Errorf("could not connect to store %d: %v", storeID, err)
		return
	}

	return
}

// getDriver returns a driver which can be used with the given store
func (crs *composeRecordStore) getDriver(def crsDefiner) driver {
	for _, d := range crs.drivers {
		if !d.Can(def.StoreDSN(), def.Capabilities()...) {
			continue
		}

		return d
	}

	return nil
}

func storeConnWrap(d driver, def crsDefiner) connFunc {
	return func(ctx context.Context) (Store, error) {
		return d.Store(ctx, def.StoreDSN())
	}
}
