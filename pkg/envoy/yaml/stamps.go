package yaml

import (
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

func decodeTimestamps(n *yaml.Node) (*resource.Timestamps, error) {
	var (
		st = &resource.Timestamps{}
	)

	return st, eachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "createdat":
			return decodeScalar(v, "created at", &st.CreatedAt)
		case "updatedat":
			return decodeScalar(v, "updated at", &st.UpdatedAt)
		case "deletedat":
			return decodeScalar(v, "deleted at", &st.DeletedAt)
		case "suspendedat":
			return decodeScalar(v, "suspended at", &st.SuspendedAt)
		case "archivedat":
			return decodeScalar(v, "archived at", &st.ArchivedAt)
		}
		return nil
	})
}

func decodeUserstamps(n *yaml.Node) (*resource.Userstamps, error) {
	var (
		us = &resource.Userstamps{}
	)

	return us, eachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "createdby",
			"creatorid":
			return decodeScalar(v, "created by", &us.CreatedBy)
		case "updatedby":
			return decodeScalar(v, "updated by", &us.UpdatedBy)
		case "deletedby":
			return decodeScalar(v, "deleted by", &us.DeletedBy)
		case "ownedby":
			return decodeScalar(v, "owned by", &us.OwnedBy)
		}
		return nil
	})
}
