package service

import (
	"fmt"
	"io"

	"github.com/cespare/xxhash/v2"
)

// hashFileContent hashes the full content of r using xxHash64 and returns
// a 16-character lowercase hex string.
//
// The caller must seek r back to 0 before calling this function (if it is a
// ReadSeeker), and must seek back to 0 again after this call before passing
// r to any subsequent read operation.
func hashFileContent(r io.Reader) (string, error) {
	h := xxhash.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return fmt.Sprintf("%016x", h.Sum64()), nil
}

// DuplicateAttachmentError is returned by CreateRecordAttachment when a file
// with the same content hash already exists in the target module and
// enforceUniqueness is enabled on the field.
type DuplicateAttachmentError struct {
	ExistingAttachmentID uint64
	ExistingRecordID     uint64
	ModuleID             uint64
	NamespaceID          uint64
	ConflictAction       string // "modal", "newTab", "sameTab", "alert"
	RecordPageID         uint64 // resolved page ID for the target module
	NamespaceSlug        string // resolved namespace slug
}

func (e *DuplicateAttachmentError) Error() string {
	return "duplicate file: this file already exists in the target module"
}
