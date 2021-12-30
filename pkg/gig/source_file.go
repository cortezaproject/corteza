package gig

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

type (
	fileSource struct {
		id       uint64
		path     string
		name     string
		size     int64
		mime     string
		ext      string
		checksum string

		decoders []Decoder
	}
)

func prepFileSource() *fileSource {
	return &fileSource{
		id: nextID(),
	}
}

func FileSourceFromURI(ctx context.Context, uri string) (src Source, err error) {
	var (
		tmp                   = prepFileSource()
		path, name, mime, ext string
		size                  int64
	)

	if isRemote(uri) {
		path, name, mime, ext, size, err = loadRemote(ctx, uri)
		if err != nil {
			return
		}
	} else {
		path, name, mime, ext, size, err = loadLocal(uri)
		if err != nil {
			return
		}
	}

	tmp.path = path
	tmp.name = name
	tmp.size = size
	tmp.mime = mime
	tmp.ext = ext
	tmp.checksum, err = getChecksum(tmp)

	return tmp, err
}

func FileSourceFromBlob(ctx context.Context, name string, r io.Reader) (src Source, err error) {
	tmp := prepFileSource()

	tmpF, err := createTempFile(r)
	if err != nil {
		return
	}
	defer tmpF.Close()

	stats, err := tmpF.Stat()
	if err != nil {
		return
	}

	tmp.name = name
	tmp.path = tmpF.Name()
	tmp.size = stats.Size()
	tmp.ext = getFileExt(name)
	tmp.mime, err = detectMime(tmpF, name)
	tmp.checksum, err = getChecksum(tmp)

	return tmp, err
}

func PrepareSourceFromDirectory(ctx context.Context, path string) (sources SourceSet, err error) {
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		src, err := FileSourceFromURI(ctx, p)
		if err != nil {
			return err
		}

		sources = append(sources, src)

		return nil
	})
	return
}

func (s fileSource) ID() uint64 {
	return s.id
}

func (s fileSource) Checksum() string {
	return s.checksum
}

func (f fileSource) FileName() string {
	return f.name
}

func (f fileSource) Name() string {
	return f.path
}

func (f fileSource) MimeType() string {
	return f.mime
}

func (f fileSource) Ext() string {
	return f.ext
}

func (f fileSource) Size() int64 {
	return f.size
}

func (f fileSource) Read() (r io.Reader, err error) {
	return os.Open(f.path)
}

func (f fileSource) ReadSafe() (r io.Reader) {
	r, err := f.Read()
	if err != nil {
		panic(err)
	}

	return
}

func (f fileSource) Cleanup() error {
	return deleteTempFile(f.Name())
}

func (f *fileSource) SetDecoders(dd ...Decoder) {
	f.decoders = dd
}

func (f fileSource) Decoders() []Decoder {
	return f.decoders
}

func (f fileSource) HasDecoders() bool {
	return len(f.decoders) > 0
}

func (f fileSource) Mode() os.FileMode {
	// These are only readable sources so this is fine
	// -r--r--r--
	return 0444
}

func (f fileSource) ModTime() time.Time {
	return time.Now()
}

func (f fileSource) IsDir() bool {
	return false
}

func (f fileSource) Sys() interface{} {
	// Interface states we can return nil here
	return nil
}

func isRemote(uri string) bool {
	if uri == "" {
		return false
	}

	if strings.Contains(uri, "127") {
		return true
	}

	url, err := url.Parse(uri)
	if err != nil {
		return false
	}

	switch strings.ToLower(url.Scheme) {
	case "http", "https":
		return true
	}

	return false
}

func detectMime(in io.Reader, names ...string) (string, error) {
	// @todo move this somewhere else; the package doesn't have yaml definitions
	mime.AddExtensionType(".yaml", "application/x-yaml")
	mime.AddExtensionType(".yaml", "text/yaml")
	mime.AddExtensionType(".yml", "application/x-yaml")
	mime.AddExtensionType(".yml", "text/yaml")

	dft, err := mimetype.DetectReader(in)
	if err != nil {
		return "", err
	}
	out := dft.String()

	for _, n := range names {
		ext := getFileExt(n)
		tmp := mime.TypeByExtension(ext)
		if tmp != "" {
			out = tmp
		}
	}

	return out, nil
}

func getName(path string) string {
	pp := strings.Split(path, string(os.PathSeparator))
	return pp[len(pp)-1]
}

func getFileExt(n string) string {
	pp := strings.Split(n, ".")
	return "." + strings.Join(pp[1:], ".")
}

func loadRemote(ctx context.Context, uri string) (path, name, mime, ext string, size int64, err error) {
	// Download
	rsp, err := http.Get(uri)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	// Temporarily store
	tmpF, err := createTempFile(rsp.Body)
	if err != nil {
		return
	}
	defer tmpF.Close()

	// Update source
	stats, err := os.Lstat(tmpF.Name())
	if err != nil {
		return
	}

	// Try to detect name
	url, err := url.Parse(uri)
	if err != nil {
		return
	}
	name = getName(url.Path)

	path = tmpF.Name()
	size = stats.Size()
	ext = getFileExt(name)
	mime, err = detectMime(tmpF, name)
	return
}

func loadLocal(p string) (path, name, mime, ext string, size int64, err error) {
	file, err := os.Open(p)
	if err != nil {
		return
	}
	defer file.Close()

	tmpF, err := createTempFile(file)
	if err != nil {
		return
	}
	defer tmpF.Close()

	// Try to detect name
	name = getName(file.Name())

	stats, err := file.Stat()
	if err != nil {
		return
	}

	path = tmpF.Name()
	size = stats.Size()
	ext = getFileExt(name)
	mime, err = detectMime(file, name)
	return
}

func getChecksum(src Source) (cs string, err error) {
	r, err := src.Read()
	if err != nil {
		return
	}

	hasher := sha256.New()
	if _, err = io.Copy(hasher, r); err != nil {
		return
	}

	cs = fmt.Sprintf("%x", hasher.Sum(nil))
	return
}
