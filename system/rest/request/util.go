package request

//lint:file-ignore U1000 Ignore unused code, part of request pkg toolset

// is checks if string s is contained in matches
func is(s string, matches ...string) bool {
	for _, v := range matches {
		if s == v {
			return true
		}
	}
	return false
}

func void(...interface{}) (interface{}, error) {
	return nil, nil
}
