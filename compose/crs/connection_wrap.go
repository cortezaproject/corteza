package crs

import "github.com/cortezaproject/corteza-server/compose/crs/capabilities"

type (
	crsConnectionWrap struct {
		id           uint64
		dsn          string
		capabilities capabilities.Set
	}
)

// CRSConnectionWrap is a utility to define generic store connections
func CRSConnectionWrap(id uint64, dsn string, cc ...capabilities.Capability) crsConnectionWrap {
	return crsConnectionWrap{
		id:           id,
		dsn:          dsn,
		capabilities: cc,
	}
}

// Receivers to conform to interface

func (s crsConnectionWrap) ComposeRecordStoreID() uint64 {
	return s.id
}

func (s crsConnectionWrap) StoreDSN() string {
	return s.dsn
}

func (s crsConnectionWrap) Capabilities() capabilities.Set {
	return s.capabilities
}
