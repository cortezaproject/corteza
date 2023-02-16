package types

import (
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	Icon struct {
	}

	IconFilter struct {
		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(icon *Icon) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	IconType string
)

const (
	// IconTypeLink empty or "link" (default):
	// Indicate that src will contain an absolute or relative link to an icon.
	// Can also be used for inline images (storing "base64:" prefixed string in source).
	// This type and reference is not validated by the backend.
	IconTypeLink IconType = "link" // type, source
	// IconTypeLibrary "library"
	// Source references an icon from a library. Ref's value should be in the following
	// notation: "font-awesome://<icon-identifier>".
	// This type and source is not validated by the backend.
	IconTypeLibrary IconType = "library" // type, library(font-awesome), icon name
	// IconTypeInlineSvg "svg"
	// SRC contains raw SVG document
	IconTypeInlineSvg IconType = "inline-svg" // source
	// IconTypeAttachment "attachment"
	// Reference (ID) to an existing attachment in local Corteza instance is expected
	// This type and reference must be validated by the backend.
	IconTypeAttachment IconType = "attachment" // type, Icon, name
)
