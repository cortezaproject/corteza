package expr

type (
	exprPath struct {
		path   string
		i      int
		isLast bool

		start, end int
	}
)

// Path initializes a new exprPath helper to efficiently traverse the path
func Path(p string) (out exprPath) {
	return exprPath{path: p}
}

func (p exprPath) More() bool {
	return p.start < len(p.path)
}

func (p exprPath) Get() string {
	return p.path[p.start:p.end]
}

func (p exprPath) Rest() string {
	var rest string
	if p.end+1 < len(p.path) {
		rest = p.path[p.end:]
	}

	// @todo this is fugly but it'll do the trick for now
	// Clean it up please :)
	if len(rest) > 0 && (rest[0] == '.' || rest[0] == ']') {
		rest = rest[1:]
	}
	if len(rest) > 0 && (rest[0] == '.' || rest[0] == ']') {
		rest = rest[1:]
	}

	return rest
}

func (p exprPath) Next() (out exprPath, err error) {
	if !p.More() {
		return p, nil
	}

	if p.end > 0 {
		p.start = p.end + 1
	}

	var ()

	p.start, p.end, p.isLast, err = nxtRange(p.path, p.start)
	if err != nil {
		return p, err
	}

	p.i++

	return p, nil
}

func nxtRange(path string, start int) (startOut, end int, isLast bool, err error) {
	startOut = start
	for i := start; i < len(path); i++ {
		switch path[i] {
		// This thing concludes the prev ident
		case '.', '[':
			if i == len(path)-1 {
				return startOut, -1, false, invalidPathErr
			}
			if i > 0 {
				if path[i-1] == ']' {
					startOut++
					continue
				}
			}
			return startOut, i, false, nil

		case ']':
			// If we're at the end, that's that
			if i == len(path)-1 {
				return startOut, i, true, nil
			}

			if path[i+1] != '.' && path[i+1] != '[' {
				return startOut, -1, false, invalidPathErr
			} else {
				return startOut, i, false, nil
			}
		}
	}

	return startOut, len(path), true, nil
}
