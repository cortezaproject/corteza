package settings

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
)

type (
	Importer struct {
		settings ValueSet
	}

	ImportKeeper interface {
		BulkSet(ctx context.Context, vv ValueSet) (err error)
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
	// Convert to interface{}, since json.Marshal cant handle map[interface{}]interface{}
	v, err := json.Marshal(deinterfacer.Simplify(value))
	if err != nil {
		return err
	}

	setting := &Value{
		Name:  name,
		Value: v,
	}

	imp.settings = append(imp.settings, setting)
	return nil
}

func (imp *Importer) Store(ctx context.Context, k ImportKeeper) (err error) {
	return k.BulkSet(ctx, imp.settings)
}

func (imp *Importer) GetValues() ValueSet {
	return imp.settings
}
