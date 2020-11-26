package yaml

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

type (
	// decoder is a wrapper struct for yaml related methods
	decoder struct{}

	EnvoyMarshler interface {
		MarshalEnvoy() ([]resource.Interface, error)
	}

	nodeDecoder interface {
		DecodeNodes(ctx context.Context) ([]resource.Interface, error)
	}
)

func Decoder() *decoder {
	return &decoder{}
}

// CanDecodeFile
func (y *decoder) CanDecodeFile(i os.FileInfo) bool {
	switch strings.Trim(filepath.Ext(i.Name()), ".") {
	case "yaml", "yml":
		return true
	}

	return false
}

func (y *decoder) Decode(ctx context.Context, r io.Reader, dctx *envoy.DecoderOpts) ([]resource.Interface, error) {
	var (
		doc = &Document{}
	)

	if err := yaml.NewDecoder(r).Decode(doc); err != nil {
		return nil, err
	}

	return doc.Decode(ctx)
}
