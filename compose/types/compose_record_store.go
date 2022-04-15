package types

import "github.com/cortezaproject/corteza-server/compose/crs/capabilities"

type (
	ComposeRecordStore struct {
		ID       uint64
		Handle   string
		DSN      string
		Location string
		// ...

		// @todo IMO having it like so (instead of in a struct) allows for more
		// flexibility with less data
		Enforced    capabilities.Set
		Supported   capabilities.Set
		Unsupported capabilities.Set
		Enabled     capabilities.Set
	}
)

// @note this conforms to the crs.crsDefiner interface

func (crs ComposeRecordStore) ComposeRecordStoreID() uint64 {
	return crs.ID
}

func (crs ComposeRecordStore) StoreDSN() string {
	return crs.DSN
}

func (crs ComposeRecordStore) Capabilities() capabilities.Set {
	return append(crs.Enforced, crs.Supported...)
}

// ---
