package settings

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Export transforms a given ValueSet into a yaml exportable structure
func Export(ss ValueSet) (o yaml.MapSlice) {
	o = yaml.MapSlice{}
	for _, s := range ss {
		var v interface{}
		json.Unmarshal(s.Value, &v)
		setting := yaml.MapItem{
			Key:   s.Name,
			Value: v,
		}

		o = append(o, setting)
	}

	return o
}
