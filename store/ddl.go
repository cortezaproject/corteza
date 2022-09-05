package store

import "github.com/cortezaproject/corteza-server/pkg/dal"

type (
	DDL interface {
		// EnsureModel ensures model exists in the store in form of a collection or a database table
		//
		// This includes all model's attributes, constraints and indexes.
		// All must exist and be of the right type.
		EnsureModel(model *dal.Model) error
	}
)
