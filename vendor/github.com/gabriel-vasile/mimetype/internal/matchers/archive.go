package matchers

import "bytes"

// SevenZ matches a 7z archive.
func SevenZ(in []byte) bool {
	return bytes.HasPrefix(in, []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C})
}

// Gzip matched gzip files based on http://www.zlib.org/rfc-gzip.html#header-trailer.
func Gzip(in []byte) bool {
	return bytes.HasPrefix(in, []byte{0x1f, 0x8b})
}

// Crx matches a Chrome extension file: a zip archive prepended by "Cr24".
func Crx(in []byte) bool {
	return bytes.HasPrefix(in, []byte("Cr24"))
}

// Tar matches a (t)ape (ar)chive file.
func Tar(in []byte) bool {
	return len(in) > 262 && bytes.Equal(in[257:262], []byte("ustar"))
}

// Fits matches an Flexible Image Transport System file.
func Fits(in []byte) bool {
	return bytes.HasPrefix(in, []byte{
		0x53, 0x49, 0x4D, 0x50, 0x4C, 0x45, 0x20, 0x20, 0x3D, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x54,
	})
}

// Xar matches an eXtensible ARchive format file.
func Xar(in []byte) bool {
	return bytes.HasPrefix(in, []byte{0x78, 0x61, 0x72, 0x21})
}

// Bz2 matches a bzip2 file.
func Bz2(in []byte) bool {
	return bytes.HasPrefix(in, []byte{0x42, 0x5A, 0x68})
}

// Ar matches an ar (Unix) archive file.
func Ar(in []byte) bool {
	return bytes.HasPrefix(in, []byte{0x21, 0x3C, 0x61, 0x72, 0x63, 0x68, 0x3E})
}

// Deb matches a Debian package file.
func Deb(in []byte) bool {
	return len(in) > 8 && bytes.HasPrefix(in[8:], []byte{
		0x64, 0x65, 0x62, 0x69, 0x61, 0x6E, 0x2D,
		0x62, 0x69, 0x6E, 0x61, 0x72, 0x79,
	})
}

// Rar matches a RAR archive file.
func Rar(in []byte) bool {
	return bytes.HasPrefix(in, []byte("Rar!\x1A\x07\x00")) ||
		bytes.HasPrefix(in, []byte("Rar!\x1A\x07\x01\x00"))
}

// Warc matches a Web ARChive file.
func Warc(in []byte) bool {
	return bytes.HasPrefix(in, []byte("WARC/"))
}

// Zstd matches a Zstandard archive file.
func Zstd(in []byte) bool {
	return len(in) >= 4 &&
		(0x22 <= in[0] && in[0] <= 0x28 || in[0] == 0x1E) && // Different Zstandard versions.
		bytes.HasPrefix(in[1:], []byte{0xB5, 0x2F, 0xFD})
}

// Cab matches a Cabinet archive file.
func Cab(in []byte) bool {
	return bytes.HasPrefix(in, []byte("MSCF"))
}

// Rpm matches an RPM or Delta RPM package file.
func Rpm(in []byte) bool {
	return len(in) > 4 &&
		(bytes.HasPrefix(in, []byte{0xed, 0xab, 0xee, 0xdb}) ||
			bytes.HasPrefix(in, []byte("drpm")))
}

// Xz matches an xz compressed stream based on https://tukaani.org/xz/xz-file-format.txt.
func Xz(in []byte) bool {
	return bytes.HasPrefix(in, []byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00})
}

// Lzip matches an Lzip compressed file.
func Lzip(in []byte) bool {
	return bytes.HasPrefix(in, []byte{0x4c, 0x5a, 0x49, 0x50})
}

// Cpio matches a cpio archive file
func Cpio(in []byte) bool {
	return bytes.HasPrefix(in, []byte("070707")) ||
		bytes.HasPrefix(in, []byte("070701")) ||
		bytes.HasPrefix(in, []byte("070702"))
}
