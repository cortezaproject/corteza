package repository

import (
	"github.com/pkg/errors"
	"github.com/titpetric/statik/fs"

	"github.com/crusttech/crust/crm/repository/files"
	"github.com/crusttech/crust/crm/rest/request"
)

// Module chart request parameters
type ModuleChart struct {
	Name        string
	Description string
	XAxis       string
	XMin        string
	XMax        string
	YAxis       string
	GroupBy     string
	Sum         string
	Count       string
	Kind        string
	ModuleID    uint64 `json:",string"`
}

func (m *module) Chart(r *request.ModuleChart) (interface{}, error) {
	statikFS, err := fs.New(files.Data())
	if err != nil {
		return nil, errors.Wrap(err, "Error creating statik filesystem")
	}

	filename := "/chart-" + r.Kind + ".json"
	body, err := fs.ReadFile(statikFS, filename)
	if err != nil {
		return nil, errors.Wrap(err, "Error when reading "+filename)
	}
	return body, err
}
