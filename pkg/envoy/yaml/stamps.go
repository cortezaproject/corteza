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
		}
		return err
	})
}
