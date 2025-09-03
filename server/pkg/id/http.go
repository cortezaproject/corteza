package id

func ParseID(ss string) (p ID, err error) {
	if len(ss) == 0 {
		return
	}

	return ByteID([]byte(ss))
}

func ParseIDs(ss []string) (out []ID, err error) {
	if len(ss) == 0 {
		return
	}

	var aux ID
	for _, s := range ss {
		aux, err = ByteID([]byte(s))
		if err != nil {
			return
		}

		out = append(out, aux)
	}

	return
}
