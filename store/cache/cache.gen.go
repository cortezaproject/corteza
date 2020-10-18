package cache

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_cache_struct.gen.go.tpl
// Definitions:
//  - store/actionlog.yaml
//  - store/applications.yaml
//  - store/attachments.yaml
//  - store/compose_attachments.yaml
//  - store/compose_charts.yaml
//  - store/compose_module_fields.yaml
//  - store/compose_modules.yaml
//  - store/compose_namespaces.yaml
//  - store/compose_pages.yaml
//  - store/credentials.yaml
//  - store/messaging_attachments.yaml
//  - store/messaging_channel_members.yaml
//  - store/messaging_channels.yaml
//  - store/messaging_flags.yaml
//  - store/messaging_mentions.yaml
//  - store/messaging_message_attachments.yaml
//  - store/messaging_messages.yaml
//  - store/messaging_unread.yaml
//  - store/rbac_rules.yaml
//  - store/reminders.yaml
//  - store/role_members.yaml
//  - store/roles.yaml
//  - store/settings.yaml
//  - store/users.yaml

//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"github.com/cortezaproject/corteza-server/store"
	"github.com/dgraph-io/ristretto"
)

type (
	Cache struct {
		store.Storer

		composeCharts     *ristretto.Cache
		composeModules    *ristretto.Cache
		composeNamespaces *ristretto.Cache
		composePages      *ristretto.Cache
		roleMembers       *ristretto.Cache
		roles             *ristretto.Cache
		users             *ristretto.Cache
	}
)

var _ store.Users = &Cache{}

func Connect(s store.Storer) (store.Storer, error) {
	var (
		err error
		c   = &Cache{Storer: s}
	)

	c.composeCharts, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 20, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return nil, err
	}

	c.composeModules, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 20, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return nil, err
	}

	c.composeNamespaces, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 20, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return nil, err
	}

	c.composePages, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 20, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return nil, err
	}

	c.roleMembers, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 20, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return nil, err
	}

	c.roles, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 20, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return nil, err
	}

	c.users, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 20, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}
