package settings

import (
	"context"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/jmoiron/sqlx/types"
)

type (
	Importer struct {
		settings ValueSet
	}

	ImportKeeper interface {
		BulkSet(vv ValueSet) (err error)
	}
)

func NewImporter() *Importer {
	return &Importer{}
}

// CastSet - resolves settings:
//   <ValueSet> [ <Value>, ... ]
func (imp *Importer) CastSet(settings interface{}) (err error) {
	if !deinterfacer.IsMap(settings) {
		return errors.New("expecting map of settings")
	}

	return deinterfacer.Each(settings, func(_ int, name string, value interface{}) error {
		return imp.addSetting(name, value)
	})
}

func (imp *Importer) addSetting(name string, value interface{}) (err error) {
	v, ok := value.(string)
	if !ok {
		return errors.New("value must be string")
	}

	setting := &Value{
		Name:  name,
		Value: types.JSONText(v),
	}

	imp.settings = append(imp.settings, setting)
	return nil
}

func (imp *Importer) Store(ctx context.Context, k ImportKeeper) (err error) {
	return k.BulkSet(imp.settings)
}
