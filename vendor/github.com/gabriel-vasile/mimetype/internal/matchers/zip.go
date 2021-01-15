package matchers

import "bytes"

// Zip matches a zip archive.
func Zip(in []byte) bool {
	return len(in) > 3 &&
		in[0] == 0x50 && in[1] == 0x4B &&
		(in[2] == 0x3 || in[2] == 0x5 || in[2] == 0x7) &&
		(in[3] == 0x4 || in[3] == 0x6 || in[3] == 0x8)
}

// Odt matches an OpenDocument Text file.
func Odt(in []byte) bool {
	return len(in) > 30 && bytes.HasPrefix(in[30:], []byte("mimetypeapplication/vnd.oasis.opendocument.text"))
}

// Ott matches an OpenDocument Text Template file.
func Ott(in []byte) bool {
	return len(in) > 30 && bytes.HasPrefix(in[30:], []byte("mimetypeapplication/vnd.oasis.opendocument.text-template"))
}

// Ods matches an OpenDocument Spreadsheet file.
func Ods(in []byte) bool {
	return len(in) > 30 && bytes.HasPrefix(in[30:], []byte("mimetypeapplication/vnd.oasis.opendocument.spreadsheet"))
}

// Ots matches an OpenDocument Spreadsheet Template file.
func Ots(in []byte) bool {
	return len(in) > 30 && bytes.HasPrefix(in[30:], []byte("mimetypeapplication/vnd.oasis.opendocument.spreadsheet-template"))
}

// Odp matches an OpenDocument Presentation file.
func Odp(in []byte) bool {
	return len(in) > 30 && bytes.HasPrefix(in[30:], []byte("mimetypeapplication/vnd.oasis.opendocument.presentation"))
}

// Otp matches an OpenDocument Presentation Template file.
func Otp(in []byte) bool {
	return len(in) > 30 && bytes.HasPrefix(in[30:], []byte("mimetypeapplication/vnd.oasis.opendocument.presentation-template"))
}

// Odg matches an OpenDocument Drawing file.
func Odg(in []byte) bool {
	return len(in) > 30 && bytes.HasPrefix(in[30:], []byte("mimetypeapplication/vnd.oasis.opendocument.graphics"))
}

// Otg matches an OpenDocument Drawing Template file.
func Otg(in []byte) bool {
	return len(in) > 30 && bytes.HasPrefix(in[30:], []byte("mimetypeapplication/vnd.oasis.opendocument.graphics-template"))
}

// Odf matches an OpenDocument Formula file.
func Odf(in []byte) bool {
	return len(in) > 30 && bytes.HasPrefix(in[30:], []byte("mimetypeapplication/vnd.oasis.opendocument.formula"))
}

// Epub matches an EPUB file.
func Epub(in []byte) bool {
	return len(in) > 30 && bytes.HasPrefix(in[30:], []byte("mimetypeapplication/epub+zip"))
}

// Jar matches a Java archive file.
func Jar(in []byte) bool {
	t := zipTokenizer{in: in}
	for i, tok := 0, t.next(); i < 10 && tok != ""; i, tok = i+1, t.next() {
		if tok == "META-INF/MANIFEST.MF" {
			return true
		}
	}

	return false
}
