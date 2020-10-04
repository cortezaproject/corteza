package molding

import "strconv"

type (
	//Type interface {
	//	Kind() string
	//}

	Variables map[string]interface{}
)

// Assign takes base variables and assigns all new variables
func (vv Variables) Merge(nn ...Variables) Variables {
	var (
		out = Variables{}
	)

	nn = append([]Variables{vv}, nn...)
	for i := range nn {
		for k, v := range nn[i] {
			out[k] = v
		}
	}

	return out
}

func (vv Variables) Bool(key string, def bool) bool {
	if v, has := vv[key]; has {
		// @todo handle other types than string...
		o, err := strconv.ParseBool(v.(string))
		if err != nil {
			return o
		}
	}

	return def
}

func (vv Variables) Int64(key string, def int64) int64 {
	if v, has := vv[key]; has {
		// @todo handle other types than string...
		o, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return o
		}
	}

	return def
}

func (vv Variables) Uint64(key string, def uint64) uint64 {
	if v, has := vv[key]; has {
		// @todo handle other types than string...
		o, err := strconv.ParseUint(v.(string), 10, 64)
		if err != nil {
			return o
		}
	}

	return def
}

func (vv Variables) Float64(key string, def float64) float64 {
	if v, has := vv[key]; has {
		// @todo handle other types than string...
		o, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return o
		}
	}

	return def
}
