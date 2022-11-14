package types

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"golang.org/x/text/language"
)

type (
	ResourceTranslation struct {
		ID uint64 `json:"translationID,string"`

		Lang     Lang   `json:"lang"`
		Resource string `json:"resource"`
		K        string `json:"key"`
		Message  string `json:"message"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`

		OwnedBy   uint64 `json:"ownedBy,string"`
		CreatedBy uint64 `json:"createdBy,string" `
		UpdatedBy uint64 `json:"updatedBy,string,omitempty" `
		DeletedBy uint64 `json:"deletedBy,string,omitempty" `
	}

	ResourceTranslationFilter struct {
		TranslationID []uint `json:"translationID"`
		Lang          string `json:"lang"`
		Resource      string `json:"resource"`
		ResourceType  string `json:"resourceType"`
		OwnerID       uint64 `json:"ownerID,string"`

		Deleted filter.State `json:"deleted"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	Lang struct {
		language.Tag
	}
)

// Value scan on ReportDataSourceSet gracefully handles conversion from NULL
func (l Lang) Value() (driver.Value, error) {
	return l.String(), nil
}

func (l *Lang) Scan(value interface{}) error {
	v := ""

	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch aux := value.(type) {
	case string:
		v = aux
	case []uint8:
		v = string(value.([]byte))
	}

	l.Tag = language.Make(v)
	return nil
}

func (a *ResourceTranslation) Compare(b *locale.ResourceTranslation) bool {
	return strings.EqualFold(a.K, b.Key) && a.Lang.Tag.String() == b.Lang
}

func (set ResourceTranslationSet) New(bb locale.ResourceTranslationSet) (out ResourceTranslationSet) {
outer:
	for _, b := range bb {
		for _, a := range set {
			// It's not new
			if a.Compare(b) {
				continue outer
			}
		}

		out = append(out, &ResourceTranslation{
			Lang:     Lang{Tag: language.Make(b.Lang)},
			Resource: b.Resource,
			K:        b.Key,
			Message:  b.Msg,
		})
	}
	return
}

func (set ResourceTranslationSet) Old(bb locale.ResourceTranslationSet) (out [][2]*ResourceTranslation) {
	for _, b := range bb {
		for _, a := range set {
			// It's not new
			if a.Compare(b) {
				aux := *a
				aux.Message = b.Msg

				// old, new
				out = append(out, [2]*ResourceTranslation{a, &aux})
				break
			}
		}
	}
	return
}

func FromLocale(ll locale.ResourceTranslationSet) (out ResourceTranslationSet) {
	for _, l := range ll {
		out = append(out, &ResourceTranslation{
			Lang:     Lang{Tag: language.Make(l.Lang)},
			Resource: l.Resource,
			K:        l.Key,
			Message:  l.Msg,
		})
	}

	return out
}
