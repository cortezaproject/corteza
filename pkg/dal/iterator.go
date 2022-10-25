package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	// Iterator provides an interface for loading data from the underlying store
	Iterator interface {
		Next(ctx context.Context) bool
		More(uint, ValueGetter) error
		Err() error
		Scan(ValueSetter) error
		Close() error

		BackCursor(ValueGetter) (*filter.PagingCursor, error)
		ForwardCursor(ValueGetter) (*filter.PagingCursor, error)

		// // -1 means unknown
		// Total() int
		// Cursor() any
		// // ... do we need anything else here?
	}
)

// IteratorEncodeJSON helper function that encodes each item from the iterator as JSON
// and writes it to th given io.Writer.
//
// target initialization function is intentionally used to avoid use of reflection
func IteratorEncodeJSON(ctx context.Context, w io.Writer, iter Iterator, initTarget func() ValueSetter) (err error) {
	var (
		target   ValueSetter
		firstOut = false
	)

	for iter.Next(ctx) {
		if err = iter.Err(); err != nil {
			return
		}

		if firstOut {
			if _, err = w.Write([]byte(`,`)); err != nil {
				return
			}
		}

		firstOut = true

		target = initTarget()

		if err = iter.Scan(target); err != nil {
			return
		}

		err = json.NewEncoder(w).Encode(target)
		if err != nil {
			return
		}
	}

	return
}

// IteratorPaging helper function for record paging cursor and total
func IteratorPaging(ctx context.Context, iter Iterator, pp *filter.Paging, ss filter.Sorting, fn func(i Iterator) (ValueGetter, bool)) (err error) {
	if pp == nil {
		return
	}

	if pp.PageCursor != nil {
		if pp.IncPageNavigation || pp.IncTotal {
			return fmt.Errorf("not allowed to fetch page navigation or total item count with page cursor")
		}
	}

	pp.Total = 0
	pp.PrevPage = nil
	pp.NextPage = nil
	pp.PageNavigation = []*filter.Page{}

	const howMuchMore = 1000

	var (
		counter uint
		total   uint

		overfetch = pp.IncTotal || pp.IncPageNavigation

		cur  *filter.PagingCursor
		page = filter.Page{
			Page:   1,
			Count:  0,
			Cursor: nil,
		}
	)

	for iter.Next(ctx) {
		if err = iter.Err(); err != nil {
			return
		}

		r, ok := fn(iter)
		if !ok {
			continue
		}

		counter++
		total++
		page.Count++
		nextPage := pp.Limit > 0 && total%pp.Limit == 0

		if pp.PrevPage == nil {
			pp.PrevPage, err = iter.BackCursor(r)
			if err != nil {
				return
			}
		}

		// cursor for each page
		cur, err = iter.ForwardCursor(r)
		if err != nil {
			return
		}

		// We fetched enough and we don't need count/all pages; end because anything
		// extra would be useless processing
		if nextPage && !overfetch {
			break
		}

		// Additional processing only happens when we get to the next page so we can
		// safely skip it
		if !nextPage {
			continue
		}

		// Update paging params for the initial filtering
		if pp.NextPage == nil {
			pp.NextPage = cur
		}

		// Paging params for the current chunk
		pp.PageNavigation = append(pp.PageNavigation, &filter.Page{
			Page:   page.Page,
			Count:  page.Count,
			Cursor: page.Cursor,
		})

		// Prepare paging params for the next chunk
		page = filter.Page{
			Page:   uint(len(pp.PageNavigation) + 1),
			Count:  0,
			Cursor: cur,
		}

		// Request more
		// If this was the first page, request more because the limit was exceeded
		// If this wasn't the first page, request more after we've reached the refetch count
		if len(pp.PageNavigation) == 1 || counter == howMuchMore {
			counter = 0

			// request more items
			if err = iter.More(howMuchMore, r); err != nil {
				return
			}
		}
	}

	// push the last page to page navigation
	if page.Count > 0 {
		pp.PageNavigation = append(pp.PageNavigation, &filter.Page{
			Page:   page.Page,
			Count:  page.Count,
			Cursor: page.Cursor,
		})
	}

	if pp.PageCursor == nil {
		pp.PrevPage = nil
	}

	if pp.NextPage != nil && len(pp.PageNavigation) == 1 {
		pp.NextPage = nil
	}

	if pp.IncTotal {
		pp.Total = total
	}

	if !pp.IncPageNavigation {
		pp.PageNavigation = nil
	}

	return
}
