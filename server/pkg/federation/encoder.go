package federation

import (
	"encoding/json"
	"io"

	"github.com/crusttech/human/server/pkg/options"
)

const (
	ActivityStreamsStructure EncodingFormat = 0
	HumanInternalStructure EncodingFormat = 1
	ActivityStreamsData      EncodingFormat = 2
	HumanInternalData      EncodingFormat = 3
)

type (
	EncodingFormat int

	Encoder struct {
		w io.Writer
		o options.FederationOpt
	}
)

func NewEncoder(w io.Writer, o options.FederationOpt) *Encoder {
	return &Encoder{w: w, o: o}
}

// Encode the specific format per payload
func (e Encoder) Encode(payload interface{}, t EncodingFormat) error {
	var (
		ea   EncoderAdapter
		enc  = json.NewEncoder(e.w)
		resp interface{}
		err  error
	)

	switch t {
	case ActivityStreamsStructure:
		ea = &EncoderAdapterActivityStreams{}
		resp, err = ea.BuildStructure(e.w, e.o, payload)
		break
	case HumanInternalStructure:
		ea = &EncoderAdapterHumanInternal{}
		resp, err = ea.BuildStructure(e.w, e.o, payload)
		break
	case ActivityStreamsData:
		ea = &EncoderAdapterActivityStreams{}
		resp, err = ea.BuildData(e.w, e.o, payload)
		break
	case HumanInternalData:
		ea = &EncoderAdapterHumanInternal{}
		resp, err = ea.BuildData(e.w, e.o, payload)
		break
	}

	if err != nil {
		return err
	}

	return enc.Encode(resp)
}
