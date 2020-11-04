package yaml

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

type (
	// decoder is a wrapper struct for yaml related methods
	decoder struct {
		loader loader
	}

	loader interface {
		LoadComposeNamespace()
	}

	EnvoyMarshler interface {
		MarshalEnvoy() ([]resource.Interface, error)
	}

	nodeDecoder interface {
		DecodeNodes(ctx context.Context, l loader) ([]resource.Interface, error)
	}
)

func Decoder(l loader) *decoder {
	return &decoder{l}
}

// CanDecodeFile
func (y *decoder) CanDecodeFile(i os.FileInfo) bool {
	switch filepath.Ext(i.Name()) {
	case "yaml", "yml":
		return true
	}

	return false
}

func (y *decoder) Decode(ctx context.Context, r io.Reader, i os.FileInfo) ([]resource.Interface, error) {
	var (
		doc = &Document{}
	)

	if err := yaml.NewDecoder(r).Decode(doc); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", i.Name(), err)
	}

	return doc.Decode(ctx, y.loader)
}
