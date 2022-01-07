package gig

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cast"
)

type (
	archive int
)

const (
	ArchiveZIP archive = iota
	ArchiveTar
)

func isTarGz(in Source) bool {
	return in.Ext() == ".tar.gz"
}

func extractTarGz(ctx context.Context, in Source) (out SourceSet, err error) {
	arch, err := in.Read()
	if err != nil {
		return nil, err
	}

	uncompressedStream, err := gzip.NewReader(arch)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		source := &fileSource{}
		switch header.Typeflag {
		case tar.TypeReg:
			file, err := createTempFile(tarReader)
			if err != nil {
				return nil, err
			}

			// Try to detect name
			pp := strings.Split(header.Name, string(os.PathSeparator))
			source.name = pp[len(pp)-1]
			source.size = header.Size
			source.ext = getFileExt(header.Name)
			source.mime, err = detectMime(file, header.Name)

			out = append(out, source)
		}
	}

	return
}

func compressTarGz(ctx context.Context, sources SourceSet, baseName string) (archive, name string, err error) {
	var f *os.File
	f, err = createTempFile(nil)
	if err != nil {
		return
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, src := range sources {
		var header *tar.Header

		header, err = tar.FileInfoHeader(src, src.FileName())
		if err != nil {
			return
		}

		header.Name = src.FileName()

		if err = tw.WriteHeader(header); err != nil {
			return
		}

		var r io.Reader
		if r, err = src.Read(); err != nil {
			return
		}

		if _, err = io.Copy(tw, r); err != nil {
			return
		}

	}

	name = baseName
	if name == "" {
		name = "archive"
	}
	name = strings.Split(name, ".")[0]

	return f.Name(), fmt.Sprintf("%s.tar.gz", name), nil
}

func archiveFromParams(p interface{}) archive {
	c := cast.ToString(p)
	if c == "" {
		return ArchiveZIP
	}

	// Ignoring the error and defaulting
	out, _ := archiveFromString(c)
	return out
}

func archiveFromString(arch string) (archive, error) {
	switch strings.ToLower(arch) {
	case "zip":
		return ArchiveZIP, nil
	case "tar":
		return ArchiveTar, nil
	default:
		return 0, fmt.Errorf("unknown archive: %s", arch)
	}
}

func (f archive) String() string {
	switch f {
	case ArchiveZIP:
		return ".zip"
	case ArchiveTar:
		return ".tar.gz"
	}

	return "N/A"
}
