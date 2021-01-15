package matchers

import "bytes"

// Pdf matches a Portable Document Format file.
func Pdf(in []byte) bool {
	return bytes.HasPrefix(in, []byte{0x25, 0x50, 0x44, 0x46})
}

// Fdf matches a Forms Data Format file.
func Fdf(in []byte) bool {
	return bytes.HasPrefix(in, []byte("%FDF"))
}

// DjVu matches a DjVu file.
func DjVu(in []byte) bool {
	if len(in) < 12 {
		return false
	}
	if !bytes.HasPrefix(in, []byte{0x41, 0x54, 0x26, 0x54, 0x46, 0x4F, 0x52, 0x4D}) {
		return false
	}
	return bytes.HasPrefix(in[12:], []byte("DJVM")) ||
		bytes.HasPrefix(in[12:], []byte("DJVU")) ||
		bytes.HasPrefix(in[12:], []byte("DJVI")) ||
		bytes.HasPrefix(in[12:], []byte("THUM"))
}

// Mobi matches a Mobi file.
func Mobi(in []byte) bool {
	return len(in) > 67 && bytes.Equal(in[60:68], []byte("BOOKMOBI"))
}

// Lit matches a Microsoft Lit file.
func Lit(in []byte) bool {
	return bytes.HasPrefix(in, []byte("ITOLITLS"))
}

// P7s matches an .p7s signature File (PEM, Base64).
func P7s(in []byte) bool {
	// Check for PEM Encoding.
	if bytes.HasPrefix(in, []byte("-----BEGIN PKCS7")) {
		return true
	}
	// Check if DER Encoding is long enough.
	if len(in) < 20 {
		return false
	}
	// Magic Bytes for the signedData ASN.1 encoding.
	startHeader := [][]byte{{0x30, 0x80}, {0x30, 0x81}, {0x30, 0x82}, {0x30, 0x83}, {0x30, 0x84}}
	signedDataMatch := []byte{0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01, 0x07}
	// Check if Header is correct. There are multiple valid headers.
	for i, match := range startHeader {
		// If first bytes match, then check for ASN.1 Object Type.
		if bytes.HasPrefix(in, match) {
			if bytes.HasPrefix(in[i+2:], signedDataMatch) {
				return true
			}
		}
	}

	return false
}
