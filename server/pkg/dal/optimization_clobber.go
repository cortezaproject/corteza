package dal

// pipelineClobberSteps tries to reduce the amount of pipeline steps by offloading
// higher operations to the lower levels
//
// As an example; the aggregation can be offloaded to the database for faster
// execution.
func pipelineClobberSteps(in Pipeline) (Pipeline, error) {
	if len(in) <= 1 {
		// Can't optimize further :upsidedownface:
		return in, nil
	}

	// Outline
	// - get a nicer pipeline representation
	// 	 @todo make the pipeline nicer in the first place
	// - traverse from the lief nodes up; try to clobber if the lief node allows it
	//
	// The clobbering for a branch ends when there is a node that can't be clobbered.
	// The progression ends because application level nodes can't be offloaded to.
	ll := wrapPpSteps(in)
	for _, l := range ll {
		for {
			// When there is no parent, we can't progress further
			if l.parent == nil {
				break
			}

			// if step can't clobber, skip
			cs, ok := l.step.(clobberableStep)
			if !ok {
				break
			}

			// if child fails to clobber parent, skip to the next child
			// @note for now we can end the clobbering if any of the steps
			//       can't be clobbered as all of the application defined steps
			//       are focused on the single op. and can't do anything else.
			if !cs.clobber(l.parent.step) {
				break
			}

			// if clobbered successfully, update references
			if l.parent != nil && l.parent.parent != nil {
				// - update child ref of the parent's parent
				for i, c := range l.parent.parent.child {
					if c == l.parent {
						l.parent.parent.child[i] = l
					}
				}
			}

			l.parent = l.parent.parent

			// @todo for now, clobbering ends after one successfull instance; this is due
			//       to the current DB implementation doesn't allow nested things.
			break
		}
	}

	// convert back to the pipeline representation; update deps in the process
	out := make(Pipeline, 0, len(in))
	seen := make(map[*ppStepWrap]bool, len(in))
	for _, l := range ll {
		out = append(out, unwrapPpSteps(l, seen)...)
	}

	return out, nil
}
