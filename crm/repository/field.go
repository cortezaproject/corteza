package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crusttech/crust/crm/types"
	"os"
	"path"
	"path/filepath"
)

type (
	Field interface {
		With(ctx context.Context) Field

		FindByName(name string) (*types.Field, error)
		Find() ([]*types.Field, error)
	}

	field struct {
		*repository
	}
)

func NewField(ctx context.Context) Field {
	return &field{
		repository: &repository{
			ctx: ctx,
		},
	}
}

func (r *field) With(ctx context.Context) Field {
	return &field{
		repository: r.repository.With(ctx),
	}
}

// Finds field by it's name and returns it
func (f *field) FindByName(name string) (*types.Field, error) {
	return f.fieldDecode(fmt.Sprintf(fieldPath, name))
}

// Returns all known fields
func (f *field) Find() ([]*types.Field, error) {
	matches, err := filepath.Glob(fmt.Sprintf(fieldPath, "*"))
	if err != nil {
		return nil, err
	}

	res := make([]*types.Field, len(matches))
	for i, match := range matches {
		if res[i], err = f.fieldDecode(match); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (f *field) fieldDecode(filepath string) (*types.Field, error) {
	file, err := os.Open(filepath)
	if err != nil {
		// @todo wrap error
		return nil, err
	}

	defer file.Close()

	// Removes path and extension from full filename
	fieldTypeFromPath := func(filepath string) string {
		t := path.Base(filepath)
		return t[:len(t)-5]
	}

	// Preset field's type with name of the file (sans .json)
	// if type is explicitly set within the file, it will be overwritten
	field := &types.Field{Type: fieldTypeFromPath(filepath)}
	if err := json.NewDecoder(file).Decode(&field); err != nil {
		// @todo wrap error
		return nil, err
	}

	return field, nil
}
