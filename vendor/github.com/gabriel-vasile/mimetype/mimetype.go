// Package mimetype uses magic number signatures to detect the MIME type of a file.
//
// mimetype stores the list of MIME types in a tree structure with
// "application/octet-stream" at the root of the hierarchy. The hierarchy
// approach minimizes the number of checks that need to be done on the input
// and allows for more precise results once the base type of file has been
// identified.
package mimetype

import (
	"io"
	"mime"
	"os"

	"github.com/gabriel-vasile/mimetype/internal/matchers"
)

// Detect returns the MIME type found from the provided byte slice.
//
// The result is always a valid MIME type, with application/octet-stream
// returned when identification failed.
func Detect(in []byte) *MIME {
	if len(in) > matchers.ReadLimit {
		in = in[:matchers.ReadLimit]
	}
	return root.match(in)
}

// DetectReader returns the MIME type of the provided reader.
//
// The result is always a valid MIME type, with application/octet-stream
// returned when identification failed with or without an error.
// Any error returned is related to the reading from the input reader.
//
// DetectReader assumes the reader offset is at the start. If the input
// is a ReadSeeker you read from before, it should be rewinded before detection:
//  reader.Seek(0, io.SeekStart)
//
// To prevent loading entire files into memory, DetectReader reads at most
// matchers.ReadLimit bytes from the reader.
func DetectReader(r io.Reader) (*MIME, error) {
	in := make([]byte, matchers.ReadLimit)
	n, err := io.ReadFull(r, in)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return root, err
	}
	in = in[:n]

	return Detect(in), nil
}

// DetectFile returns the MIME type of the provided file.
//
// The result is always a valid MIME type, with application/octet-stream
// returned when identification failed with or without an error.
// Any error returned is related to the opening and reading from the input file.
//
// To prevent loading entire files into memory, DetectFile reads at most
// matchers.ReadLimit bytes from the input file.
func DetectFile(file string) (*MIME, error) {
	f, err := os.Open(file)
	if err != nil {
		return root, err
	}
	defer f.Close()

	return DetectReader(f)
}

// EqualsAny reports whether s MIME type is equal to any MIME type in mimes.
// MIME type equality test is done on the "type/subtype" section, ignores
// any optional MIME parameters, ignores any leading and trailing whitespace,
// and is case insensitive.
func EqualsAny(s string, mimes ...string) bool {
	s, _, _ = mime.ParseMediaType(s)
	for _, m := range mimes {
		m, _, _ = mime.ParseMediaType(m)
		if s == m {
			return true
		}
	}

	return false
}
