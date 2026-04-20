package federation

import (
	"io"

	"github.com/crusttech/human/server/pkg/options"
)

type (
	EncoderAdapterHumanInternal struct{}

	ResponseWrapper struct {
		Response interface{} `json:"response"`
	}
)

// Build a default Human response
func (a EncoderAdapterHumanInternal) BuildStructure(w io.Writer, o options.FederationOpt, p interface{}) (interface{}, error) {
	return ResponseWrapper{
		Response: listModuleResponseHumanInternal{
			Filter: p.(ListStructurePayload).Filter,
			Set:    p.(ListStructurePayload).Set,
		}}, nil
}

// Build a default Human response
func (a EncoderAdapterHumanInternal) BuildData(w io.Writer, o options.FederationOpt, p interface{}) (interface{}, error) {
	return ResponseWrapper{
		Response: listRecordResponseHumanInternal{
			Filter: p.(ListDataPayload).Filter,
			Set:    p.(ListDataPayload).Set,
		}}, nil
}
