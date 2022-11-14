package resource

import "fmt"

type (
	shaper interface {
		Shape([]Interface) ([]Interface, error)
	}
)

// Shape shapes ResourceDatasets based on their correlated Template
//
// During the shaping step, no data is removed, nor added to the original resource slice.
// The shapers are ran in the same order they were provided.
func Shape(rr []Interface, ss ...shaper) ([]Interface, error) {
	ii := make([]Interface, 0, int(len(rr)/2))

	// Firstly handle any resource shaping
	for _, s := range ss {
		sdd, err := s.Shape(rr)
		if err != nil {
			return nil, err
		}
		ii = append(ii, sdd...)
	}

	// Cleanup the final shaped resource slice.
	// @todo make this cleaner!!1!'
	//
	// After this point, all of the data should be ready for further processing.
	for _, r := range rr {
		if _, ok := r.(*ComposeRecordTemplate); ok {
			continue
		}
		if _, ok := r.(*ResourceDataset); ok {
			continue
		}

		ii = append(ii, r)
	}

	return ii, nil
}

// findResourceDataset finds and returns the first dataset that matches the
// provided identifiers.
func findResourceDataset(rr []Interface, ii Identifiers) *ResourceDataset {
	for _, r := range rr {
		genR, ok := r.(*ResourceDataset)
		if !ok {
			continue
		}
		if !genR.Identifiers().HasAny(ii) {
			continue
		}

		return genR
	}
	return nil
}

func genericSourceErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("generic source unresolved %v", ii.StringSlice())
}
