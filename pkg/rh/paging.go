package rh

import "strconv"

const (
	PER_PAGE_DEFAULT uint = 50
	PER_PAGE_MAX          = 200
	PER_PAGE_MIN          = 10
)

type (
	// PageFilter supports page/perPage (one based) and limit/offset
	// pagination.
	//
	// Limit/offset is prioritised over page/perPage
	//
	PageFilter struct {
		// If limit is set to a positive number,
		// paging mechanisms will use limit/offset
		// Otherwise page/perPage is used
		Limit  uint `json:"limit,omitempty"`
		Offset uint `json:"offset,omitempty"`

		Page    uint `json:"page,omitempty"`
		PerPage uint `json:"perPage,omitempty"`

		// Count is used when filter and pagination are send back
		// with the response
		Count uint `json:"count"`
	}

	paginationParams interface {
		GetLimit() uint
		GetOffset() uint
		GetPage() uint
		GetPerPage() uint
	}
)

func Paging(p paginationParams) PageFilter {
	return PageFilter{
		Limit:   p.GetLimit(),
		Offset:  p.GetOffset(),
		Page:    p.GetPage(),
		PerPage: p.GetPerPage(),
	}
}

// Limit creates PageFilter struct from limit and, optionally offset
func Limit(a ...uint) PageFilter {
	switch len(a) {
	case 1:
		return PageFilter{Limit: a[0]}
	case 2:
		return PageFilter{Limit: a[0], Offset: a[1]}
	}

	return PageFilter{}
}

func (pf *PageFilter) ParsePagination(input interface{}) error {
	return parsePagination(pf, input)
}

func parsePagination(pf *PageFilter, input interface{}) (err error) {
	switch i := input.(type) {
	case map[string]string:
		conv := func(v *uint, name string) error {
			if _, has := i[name]; has {
				pv, err := strconv.ParseUint(i[name], 10, 32)
				if err != nil {
					return err
				}

				*v = uint(pv)
			}

			return nil
		}

		if len(i["limit"]+i["offset"]) > 0 {
			if err = conv(&pf.Limit, "limit"); err != nil {
				return
			}

			if err = conv(&pf.Offset, "offset"); err != nil {
				return
			}

			return
		}

		if len(i["page"]+i["perPage"]) > 0 {
			if err = conv(&pf.Page, "page"); err != nil {
				return
			}

			if err = conv(&pf.PerPage, "perPage"); err != nil {
				return
			}

			return
		}
	}

	return nil
}
