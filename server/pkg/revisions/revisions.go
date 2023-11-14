package revisions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/cast2"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

type (
	// generic revision struct
	Revision struct {
		ID        uint64    `json:"changeID,string"`
		Timestamp time.Time `json:"timestamp"`

		ResourceID uint64 `json:"resourceID,string"`

		Revision  int       `json:"revision"`
		Operation Operation `json:"operation"`

		UserID uint64 `json:"userID,string"`

		Changes []*Change `json:"changes"`

		Comment string `json:"comment"`
	}

	Change struct {
		// changed field
		Key string `json:"key"`

		Old []any `json:"old,omitempty"`
		New []any `json:"new,omitempty"`
	}

	Filter struct {
		ResourceID uint64 `json:"resourceID,string"`
	}
)

var (
	now = func() time.Time {
		return time.Now().Round(time.Second)
	}
)

func Make(op Operation, revision int, resourceID, userID uint64) (rev *Revision) {
	return &Revision{
		ID:         id.Next(),
		Timestamp:  now(),
		ResourceID: resourceID,
		Revision:   revision,
		UserID:     userID,
		Operation:  op,
	}
}

// untested
func (r *Revision) CollectChanges(new, old dal.ValueGetter, skip ...string) (err error) {
	var (
		cc = make([]*Change, 0)
		ch *Change

		// old values count
		ovc map[string]uint

		val any
	)

	if new == nil {
		return
	}

	if old != nil {
		// old values count
		ovc = old.CountValues()

	}

keys:
	for key, count := range new.CountValues() {
		if count == 0 && old == nil {
			continue keys
		}

		for _, s := range skip {
			if key == s {
				continue keys
			}
		}

		ch = &Change{
			Key: key,
			New: make([]any, 0, count),
			Old: make([]any, 0, count),
		}

		for pos := uint(0); pos < count; pos++ {
			val, err = new.GetValue(key, pos)
			if err != nil {
				return
			}

			ch.New = append(ch.New, val)
		}

		if old == nil {
			cc = append(cc, ch)
			continue
		}

		for pos := uint(0); pos < ovc[key]; pos++ {
			val, err = old.GetValue(key, pos)
			if err != nil {
				return
			}

			ch.Old = append(ch.Old, val)
		}

		if len(ch.New) != len(ch.Old) {
			// different sizes, means that something has changed
			cc = append(cc, ch)
			continue
		}

		for i := range ch.New {
			// go over all and append on first found difference
			if ch.New[i] != ch.Old[i] {
				cc = append(cc, ch)
				break
			}
		}
	}

	r.Changes = cc
	return
}

// CountValues satisfies dal.ValueGetter interface
func (r *Revision) CountValues() map[string]uint {
	// signaling DAL that each attribute has exactly one value!
	return nil
}

// GetValue satisfies dal.ValueGetter interface
func (r *Revision) GetValue(ident string, _ uint) (any, error) {
	switch ident {
	case "id":
		return r.ID, nil

	case "ts":
		return r.Timestamp, nil

	case "rel_resource":
		return r.ResourceID, nil

	case "revision":
		return r.Revision, nil

	case "operation":
		return r.Operation, nil

	case "rel_user":
		return r.UserID, nil

	case "delta":
		return r.Changes, nil

	case "comment":
		return r.Comment, nil

	}

	return nil, nil
}

func (r *Revision) SetValue(name string, _ uint, value any) error {
	switch name {
	case "id":
		return cast2.Uint64(value, &r.ID)

	case "ts":
		return cast2.Time(value, &r.Timestamp)

	case "rel_resource":
		return cast2.Uint64(value, &r.ResourceID)

	case "revision":
		return cast2.Int(value, &r.Revision)

	case "operation":
		return cast2.String(value, &r.Operation)

	case "rel_user":
		return cast2.Uint64(value, &r.UserID)

	case "delta":
		if bb, is := value.([]byte); is {
			return json.Unmarshal(bb, &r.Changes)
		}

		return fmt.Errorf("unexpected type for delta: %T", value)

	case "comment":
		return cast2.String(value, &r.Comment)
	}

	return nil
}

func (f Filter) Constraints() map[string][]any {
	return map[string][]any{
		"rel_resource": {f.ResourceID},
	}
}

func (f Filter) Expression() string                        { return "" }
func (f Filter) OrderBy() filter.SortExprSet               { return nil }
func (f Filter) Limit() uint                               { return 0 }
func (f Filter) Cursor() *filter.PagingCursor              { return nil }
func (f Filter) StateConstraints() map[string]filter.State { return nil }
