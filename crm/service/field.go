package service

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"

	_ "github.com/crusttech/crust/crm/types"
)

var _ = errors.Wrap

type Field struct{}

func (Field) New() *Field {
	return &Field{}
}

func (*Field) List() (interface{}, error) {
	matches, err := filepath.Glob("../crm/data/*.json")
	if err != nil {
		return nil, err
	}

	res := make([]interface{}, 0)
	for _, match := range matches {
		t := path.Base(match)
		t = t[:len(t)-5]
		params, err := decodeJSON(match)
		if err != nil {
			return nil, errors.Wrap(err, "Error when parsing "+match)
		}
		params["type"] = t
		res = append(res, params)
	}
	return res, nil
}

func (*Field) Type(id string) (interface{}, error) {
	if id == "" {
		return nil, errors.New("Missing id parameter")
	}
	params, err := decodeJSON("../crm/data/" + id + ".json")
	if err != nil {
		return nil, errors.Wrap(err, "Error reading field type: "+id)
	}
	params["type"] = id
	return params, nil
}

func decodeJSON(filename string) (map[string]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result := make(map[string]interface{})
	return result, json.NewDecoder(file).Decode(&result)
}
