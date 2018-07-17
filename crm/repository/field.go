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

const (
	// @todo root should be configurable
	// @todo move this to db or stack it inside the binary or container
	fieldPath = "crm/data/%s.json"
)

type (
	field struct{}
)

func Field() field {
	return field{}
}

// Finds field by it's name and returns it
func (repo field) FindByName(ctx context.Context, name string) (*types.Field, error) {
	return repo.decode(fmt.Sprintf(fieldPath, name))
}

// Returns all known fields
func (repo field) Find(ctx context.Context) ([]*types.Field, error) {
	matches, err := filepath.Glob(fmt.Sprintf(fieldPath, "*"))
	if err != nil {
		return nil, err
	}

	res := make([]*types.Field, len(matches))
	for i, match := range matches {
		if res[i], err = repo.decode(match); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (repo field) decode(filepath string) (*types.Field, error) {
	file, err := os.Open(filepath)
	if err != nil {
		// @todo wrap error
		return nil, err
	}

	defer file.Close()

	// Preset field's type with name of the file (sans .json)
	// if type is explicitly set within the file, it will be overwritten
	field := &types.Field{Type: repo.typeFromPath(filepath)}
	if err := json.NewDecoder(file).Decode(&field); err != nil {
		// @todo wrap error
		return nil, err
	}

	return field, nil
}

// Removes path and extension from full filename
func (repo field) typeFromPath(filepath string) string {
	t := path.Base(filepath)
	return t[:len(t)-5]
}
