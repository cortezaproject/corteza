package yaml

import (
	"context"
	"io"
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

// CanDecodeFile checks if the file can be handled by this decoder
//
// @todo Add support for this; current library is unable to detect this.
func (y *decoder) CanDecodeFile(f io.Reader) bool {
	return false
}

func (y *decoder) CanDecodeExt(ext string) bool {
	pt := strings.Split(ext, ".")
	return strings.TrimSpace(pt[len(pt)-1]) == "yaml"
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
