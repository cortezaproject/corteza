package yaml

import (
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

func decodeTimestamps(n *yaml.Node) (*resource.Timestamps, error) {
	var (
		st = &resource.Timestamps{}
	)

	// Little helper to make timestamps
	f := func(v *yaml.Node) (*resource.Timestamp, error) {
		aux := ""
		err := y7s.DecodeScalar(v, "decode "+v.Value, &aux)
		if err != nil {
			return nil, err
		}
		return resource.MakeTimestamp(aux), nil
	}

	return st, y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "createdat":
			st.CreatedAt, err = f(v)
		case "updatedat":
			st.UpdatedAt, err = f(v)
		case "deletedat":
			st.DeletedAt, err = f(v)
		case "suspendedat":
			st.SuspendedAt, err = f(v)
		case "archivedat":
			st.ArchivedAt, err = f(v)
		}
		return err
	})
}

func decodeUserstamps(n *yaml.Node) (*resource.Userstamps, error) {
	var (
		us = &resource.Userstamps{}
	)

	// Little helper to make userstamps
	f := func(v *yaml.Node) (*resource.Userstamp, error) {
		aux := ""
		err := y7s.DecodeScalar(v, "decode "+v.Value, &aux)
		if err != nil {
			return nil, err
		}
		return resource.MakeUserstampFromRef(aux), nil
	}

	return us, y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "createdby",
			"creatorid",
			"creator":
			us.CreatedBy, err = f(v)
		case "updatedby":
			us.UpdatedBy, err = f(v)
		case "deletedby":
			us.DeletedBy, err = f(v)
		case "ownedby",
			"ownerid",
			"owner":
			us.OwnedBy, err = f(v)
		case "runas",
			"runasid",
			"runner":
			us.RunAs, err = f(v)
		}
		return err
	})
}

// mapTimestamps helper encodes Timestamps into the mapping node
func encodeTimestamps(n *yaml.Node, ts *resource.Timestamps) (*yaml.Node, error) {
	if ts == nil {
		return n, nil
	}

	return addMap(n,
		"createdAt", ts.CreatedAt,
		"updatedAt", ts.UpdatedAt,
		"deletedAt", ts.DeletedAt,
		"archivedAt", ts.ArchivedAt,
		"suspendedAt", ts.SuspendedAt,
	)
}

// mapUserstamps helper encodes Userstamps into the mapping node
func encodeUserstamps(n *yaml.Node, us *resource.Userstamps) (*yaml.Node, error) {
	if us == nil {
		return n, nil
	}

	return addMap(n,
		"createdBy", us.CreatedBy,
		"updatedBy", us.UpdatedBy,
		"deletedBy", us.DeletedBy,
		"ownedBy", us.OwnedBy,
		"runAs", us.RunAs,
	)
}
