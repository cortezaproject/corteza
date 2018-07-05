package crm

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

func typesDecodeJSON(filename string) (map[string]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result := make(map[string]interface{})
	return result, json.NewDecoder(file).Decode(&result)
}

func (*Types) List(r *typesListRequest) (interface{}, error) {
	matches, err := filepath.Glob("../crm/types/*.json")
	if err != nil {
		return nil, err
	}

	res := make([]interface{}, 0)
	for _, match := range matches {
		t := path.Base(match)
		t = t[:len(t)-5]
		params, err := typesDecodeJSON(match)
		if err != nil {
			return nil, errors.Wrap(err, "Error when parsing "+match)
		}
		params["type"] = t
		res = append(res, params)
	}
	return res, nil
}

func (*Types) Type(r *typesTypeRequest) (interface{}, error) {
	if r.id == "" {
		return nil, errors.New("Missing id parameter")
	}
	params, err := typesDecodeJSON("../crm/types/" + r.id + ".json")
	if err != nil {
		return nil, errors.Wrap(err, "Error reading type: "+r.id)
	}
	params["type"] = r.id
	return params, nil
}
