package model

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/dal"
)

type (
	modelReplacer interface {
		ReplaceModel(ctx context.Context, model *dal.Model) (err error)
	}
)

var (
	models []*dal.Model
)

func Models() dal.ModelSet {
	return models
}
