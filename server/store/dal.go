package store

import "github.com/cortezaproject/corteza/server/pkg/dal"

// DAL uses given store as DAL connection
//
// This is mainly used to wrap primary store connection with DAL connection wrap
// and use it to interact with records in a primary DB
func DAL(s Storer) dal.Connection {
	return s.ToDalConn()
}
