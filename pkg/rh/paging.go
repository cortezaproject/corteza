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

func (pf *PageFilter) NormalizePerPage(min, max, def uint) {
	pf.PerPage = NormalizePerPage(pf.PerPage, min, max, def)
}

func (pf *PageFilter) NormalizePerPageWithDefaults() {
	pf.PerPage = NormalizePerPage(pf.PerPage, PER_PAGE_MIN, PER_PAGE_MAX, PER_PAGE_DEFAULT)
}

func (pf *PageFilter) NormalizePerPageNoMax() {
	pf.PerPage = NormalizePerPage(pf.PerPage, PER_PAGE_MIN, 0, PER_PAGE_DEFAULT)
}
