package importer

import (
	"context"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

type (
	Importer struct {
		roleFinder roleFinder

		roles *RoleImport

		permissions importer.PermissionImporter
	}
)

func NewImporter(rf roleFinder, p importer.PermissionImporter) *Importer {
	return &Importer{
		roleFinder:  rf,
		roles:       NewRoleImporter(rf, p),
		permissions: p,
	}
}

func (imp *Importer) YAML(r io.Reader) (err error) {
	var aux interface{}

	if err = yaml.NewDecoder(r).Decode(&aux); err != nil {
		return
	}

	return imp.Cast(aux)
}

func (imp *Importer) Cast(in interface{}) (err error) {
	return deinterfacer.Each(in, func(index int, key string, val interface{}) (err error) {
		switch key {
		case "roles":
			return imp.roles.CastSet(val)
		case "role":
			return imp.roles.CastSet([]interface{}{val})

		case "allow", "deny":
			return imp.permissions.CastResourcesSet(key, val)

		default:
			err = fmt.Errorf("unexpected key %q", key)
		}

		return err
	})
}

func (imp *Importer) Store(ctx context.Context, rk roleKeeper, pk permissions.ImportKeeper) (err error) {
	err = imp.roles.Store(ctx, rk)
	if err != nil {
		return
	}

	err = imp.permissions.Store(ctx, pk)
	if err != nil {
		return
	}

	return nil
}
