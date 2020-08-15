package settings

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Importer struct {
		settings types.SettingValueSet
	}

	ImportKeeper interface {
		BulkSet(ctx context.Context, vv types.SettingValueSet) (err error)
	}
)

func NewImporter() *Importer {
	return &Importer{}
}

// CastSet - resolves settings:
//   <ValueSet> [ <SettingValue>, ... ]
func (imp *Importer) CastSet(settings interface{}) (err error) {
	if !deinterfacer.IsMap(settings) {
		return fmt.Errorf("expecting map of settings")
	}

	return deinterfacer.Each(settings, func(_ int, name string, value interface{}) error {
		return imp.addSetting(name, value)
	})
}

func (imp *Importer) addSetting(name string, value interface{}) (err error) {
	// Convert to interface{}, since json.Marshal cant handle map[interface{}]interface{}
	setting := &types.SettingValue{Name: name}

	if err = setting.SetValue(deinterfacer.Simplify(value)); err != nil {
		return err
	}

	imp.settings = append(imp.settings, setting)
	return nil
}

func (imp *Importer) Store(ctx context.Context, k ImportKeeper) (err error) {
	return k.BulkSet(ctx, imp.settings)
}

func (imp *Importer) GetValues() types.SettingValueSet {
	return imp.settings
}
