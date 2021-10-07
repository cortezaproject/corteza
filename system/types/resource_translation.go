package types

import (
	"database/sql/driver"
	"strings"
	"time"
	"unicode/utf8"

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
			Message:  sanitizeHtml(b.Msg),
		})
	}
	return
}

func (set ResourceTranslationSet) Old(bb locale.ResourceTranslationSet) (out ResourceTranslationSet) {
	for _, b := range bb {
		for _, a := range set {
			// It's not new
			if a.Compare(b) {
				a.Message = sanitizeHtml(b.Msg)
				out = append(out, a)
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

// sanitizeHtml Aggressively strips HTML tags from a string.
// It will only keep anything between `>` and `<`.
func sanitizeHtml(s string) string {
	const (
		htmlTagStart = 60 // Unicode `<`
		htmlTagEnd   = 62 // Unicode `>`
	)

	// Set up a string builder and allocate enough memory for the new string.
	var builder strings.Builder
	builder.Grow(len(s) + utf8.UTFMax)

	in := false // True if we are inside an HTML tag.
	start := 0  // The index of the previous start tag character `<`
	end := 0    // The index of the previous end tag character `>`

	for i, c := range s {
		// If this is the last character, and we are not in an HTML tag, save it.
		if (i+1) == len(s) && end >= start {
			builder.WriteString(s[end:])
		}

		// Keep going if the character is not `<` or `>`
		if c != htmlTagStart && c != htmlTagEnd {
			continue
		}

		if c == htmlTagStart {
			// Only update the start if we are not in a tag.
			// This make sure we strip `<<br>` not just `<br>`
			if !in {
				start = i
			}
			in = true

			// Write the valid string between the close and start of the two tags.
			builder.WriteString(s[end:start])
			continue
		}
		// else c == htmlTagEnd
		in = false
		end = i + 1
	}
	s = builder.String()
	return s
}
