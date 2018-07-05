package crm

import (
	"os"
	"path"
	"path/filepath"
	"encoding/json"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

func (*Types) List(r *typesListRequest) (interface{}, error) {
	matches, err := filepath.Glob("../crm/types/*.json")
	if err != nil {
		return nil, err
	}

	decodeFile := func(filename string) (map[string]interface{}, error) {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		result := make(map[string]interface{})
		return result, json.NewDecoder(file).Decode(&result)
	}

	res := make([]interface{}, 0)
	for _, match := range matches {
		t := path.Base(match)
		t = t[:len(t)-5]
		params, err := decodeFile(match)
		if err != nil {
			return nil, errors.Wrap(err, "Error when parsing "+match)
		}
		params["type"] = t
		res = append(res, params)
	}
	return res, nil
}

func (*Types) Type(r *typesTypeRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Types.type")
}
