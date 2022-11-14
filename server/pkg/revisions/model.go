package revisions

import (
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/id"
)

const (
	RevisionResourceType = "corteza::system:revision"
)

// Model returns generic dal.Model for storing revisions
//
// Returns only basic
func Model() *dal.Model {
	// make revision model
	return &dal.Model{
		ResourceID: id.Next(),

		Ident: "revisions",

		ResourceType: RevisionResourceType,
		Attributes: dal.AttributeSet{
			&dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{Ident: "ts", Store: &dal.CodecPlain{}, Type: &dal.TypeTimestamp{}},
			&dal.Attribute{Ident: "revision", Store: &dal.CodecPlain{}, Type: &dal.TypeNumber{}},
			&dal.Attribute{Ident: "operation", Store: &dal.CodecPlain{}, Type: &dal.TypeNumber{}},
			&dal.Attribute{Ident: "rel_resource", Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{Ident: "rel_user", Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{Ident: "delta", Store: &dal.CodecPlain{}, Type: &dal.TypeJSON{}},
			&dal.Attribute{Ident: "comment", Store: &dal.CodecPlain{}, Type: &dal.TypeText{}},
		},
	}
}
