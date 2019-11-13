package rh

const (
	PER_PAGE_DEFAULT uint = 50
	PER_PAGE_MAX          = 200
	PER_PAGE_MIN          = 10
)

type (
	PageFilter struct {
		Page    uint `json:"page"`
		PerPage uint `json:"perPage"`
		Count   uint `json:"count"`
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
