package yaml

import (
	. "github.com/cortezaproject/corteza-server/pkg/y7s"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

func decodeTimestamps(n *yaml.Node) (*resource.Timestamps, error) {
	var (
		st = &resource.Timestamps{}
	)

	return st, EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "createdat":
			return DecodeScalar(v, "created at", &st.CreatedAt)
		case "updatedat":
			return DecodeScalar(v, "updated at", &st.UpdatedAt)
		case "deletedat":
			return DecodeScalar(v, "deleted at", &st.DeletedAt)
		case "suspendedat":
			return DecodeScalar(v, "suspended at", &st.SuspendedAt)
		case "archivedat":
			return DecodeScalar(v, "archived at", &st.ArchivedAt)
		}
		return nil
	})
}

func decodeUserstamps(n *yaml.Node) (*resource.Userstamps, error) {
	var (
		us = &resource.Userstamps{}
	)

	return us, EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "createdby",
			"creatorid":
			return DecodeScalar(v, "created by", &us.CreatedBy)
		case "updatedby":
			return DecodeScalar(v, "updated by", &us.UpdatedBy)
		case "deletedby":
			return DecodeScalar(v, "deleted by", &us.DeletedBy)
		case "ownedby":
			return DecodeScalar(v, "owned by", &us.OwnedBy)
		}
		return nil
	})
}
