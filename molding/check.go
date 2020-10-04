package molding

import "fmt"

func CheckNodes(nn ...Node) error {
	return checkNodes(Nodes{}, nn...)
}

func checkNodes(callstack Nodes, nn ...Node) error {
	// indexing callstack
	var (
		revStack = map[Node]bool{}
		csString = callstack.String()
	)
	for _, n := range callstack {
		revStack[n] = true
	}

	for _, n := range nn {
		if revStack[n] {
			// checking (indexed) callstack
			// to see if node is already checked
			//
			// this is needed to prevent inf-loop checks in case
			// of cyclic graphs (loop scenario)
			return nil
		}

		switch c := n.(type) {
		case Joiner:
			if len(c.Paths()) == 0 {
				return fmt.Errorf("expecting at least one input path for join gateway (%s)", csString)
			}

			for _, p := range c.Paths() {
				nxt, is := p.(Iterator)
				if !is {
					return fmt.Errorf("invalid parent node used, does not implement Next() (%s)", csString)
				}

				if nxt.Next() != n {
					return fmt.Errorf("invalid parent node used should point back ti join gateway node (%s)", csString)
				}
			}

		case Finalizer:
			return nil

		case Tester:
			if len(c.Paths()) == 0 {
				return fmt.Errorf("expecting at least one output path for gateway (%s)", csString)
			}

			return checkNodes(append(callstack, c), c.Paths()...)

		case Executor:
			if c.Next() == nil {
				return fmt.Errorf("next node of executor node not set (%s)", csString)
			}

			return checkNodes(append(callstack, c), c.Next())

		default:
			return fmt.Errorf("unknown node (%s)", csString)
		}
	}

	return nil
}
