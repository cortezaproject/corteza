package rh

import "strconv"

const (
	PER_PAGE_DEFAULT uint = 50
	PER_PAGE_MAX          = 200
	PER_PAGE_MIN          = 10
)

type (
	// @todo this needs to be refactored to support
	//       limit/offset params alongside page/perPage
	PageFilter struct {
		Page    uint `json:"page"`
		PerPage uint `json:"perPage"`
		Count   uint `json:"count"`

		// Limit  uint `json:"limit"`
		// Offset uint `json:"offset"`
	}
)

func Paging(page, perPage uint) PageFilter {
	if page == 0 {
		page = 1
	}

	return PageFilter{
		Page:    page,
		PerPage: perPage,
	}
}

func (pf *PageFilter) ParsePagination(input interface{}) {
	switch i := input.(type) {
	case map[string]string:
		if len(i["limit"]+i["offset"]) > 0 {
			// @todo to properly & fully support limit & offset
			//       we need to refactor pagination handling
			limit, _ := strconv.ParseUint(i["limit"], 10, 32)
			//offset, _ := strconv.ParseUint(i["offset"], 10, 32)

			// only basic support for now due to
			// limitation of PageFilter implementation
			pf.PerPage = uint(limit)
			//pf.Page = offset / limit
		} else if len(i["page"]+i["perPage"]) > 0 {
			page, _ := strconv.ParseUint(i["page"], 10, 32)
			perPage, _ := strconv.ParseUint(i["perPage"], 10, 32)
			pf.Page = uint(page)
			pf.PerPage = uint(perPage)
		}
	}
}
