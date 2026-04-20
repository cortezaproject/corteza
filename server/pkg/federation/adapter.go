package federation

import (
	"io"

	"github.com/crusttech/human/server/pkg/options"
)

type (
	EncoderAdapter interface {
		BuildStructure(io.Writer, options.FederationOpt, interface{}) (interface{}, error)
		BuildData(io.Writer, options.FederationOpt, interface{}) (interface{}, error)
	}
)
