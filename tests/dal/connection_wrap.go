package dal

import "github.com/cortezaproject/corteza-server/pkg/dal/capabilities"

type (
	connection struct {
		id           uint64
		dsn          string
		capabilities capabilities.Set
	}
)

// Connection is a utility to define generic store connections
func Connection(id uint64, dsn string, cc ...capabilities.Capability) connection {
	return connection{
		id:           id,
		dsn:          dsn,
		capabilities: cc,
	}
}

// Receivers to conform to interface

func (s connection) ComposeRecordStoreID() uint64 {
	return s.id
}

func (s connection) StoreDSN() string {
	return s.dsn
}

func (s connection) Capabilities() capabilities.Set {
	return s.capabilities
}
