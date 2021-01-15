package mimetype

import (
	"mime"
)

// MIME struct holds information about a file format: the string representation
// of the MIME type, the extension and the parent file format.
type MIME struct {
	mime      string
	aliases   []string
	extension string
	matchFunc func([]byte) bool
	children  []*MIME
	parent    *MIME
}

// String returns the string representation of the MIME type, e.g., "application/zip".
func (m *MIME) String() string {
	return m.mime
}

// Extension returns the file extension associated with the MIME type.
// It includes the leading dot, as in ".html". When the file format does not
// have an extension, the empty string is returned.
func (m *MIME) Extension() string {
	return m.extension
}

// Parent returns the parent MIME type from the hierarchy.
// Each MIME type has a non-nil parent, except for the root MIME type.
//
// For example, the application/json and text/html MIME types have text/plain as
// their parent because they are text files who happen to contain JSON or HTML.
// Another example is the ZIP format, which is used as container
// for Microsoft Office files, EPUB files, JAR files and others.
func (m *MIME) Parent() *MIME {
	return m.parent
}

// Is checks whether this MIME type, or any of its aliases, is equal to the
// expected MIME type. MIME type equality test is done on the "type/subtype"
// section, ignores any optional MIME parameters, ignores any leading and
// trailing whitespace, and is case insensitive.
func (m *MIME) Is(expectedMIME string) bool {
	// Parsing is needed because some detected MIME types contain parameters
	// that need to be stripped for the comparison.
	expectedMIME, _, _ = mime.ParseMediaType(expectedMIME)
	found, _, _ := mime.ParseMediaType(m.mime)

	if expectedMIME == found {
		return true
	}
	for _, alias := range m.aliases {
		if alias == expectedMIME {
			return true
		}
	}

	return false
}

func newMIME(mime, extension string, matchFunc func([]byte) bool, children ...*MIME) *MIME {
	m := &MIME{
		mime:      mime,
		extension: extension,
		matchFunc: matchFunc,
		children:  children,
	}

	for _, c := range children {
		c.parent = m
	}

	return m
}

func (m *MIME) alias(aliases ...string) *MIME {
	m.aliases = aliases
	return m
}

// match does a depth-first search on the matchers tree.
// It returns the deepest successful matcher for which all the children fail.
func (m *MIME) match(in []byte) *MIME {
	for _, c := range m.children {
		if c.matchFunc(in) {
			return c.match(in)
		}
	}

	return m
}

func (m *MIME) flatten() []*MIME {
	out := []*MIME{m}
	for _, c := range m.children {
		out = append(out, c.flatten()...)
	}

	return out
}
