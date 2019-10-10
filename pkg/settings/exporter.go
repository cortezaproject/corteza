package settings

import "gopkg.in/yaml.v2"

// Export transforms a given ValueSet into a yaml exportable structure
func Export(ss ValueSet) (o yaml.MapSlice) {
	o = yaml.MapSlice{}
	for _, s := range ss {
		setting := yaml.MapItem{
			Key:   s.Name,
			Value: s.Value.String(),
		}
		o = append(o, setting)
	}

	return o
}
