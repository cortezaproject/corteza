package gig

type (
	Preprocessor interface {
		Ref() string
		Params() map[string]interface{}
	}
	PreprocessorSet []Preprocessor
)

func PreprocessorNoop() (out preprocessorNoop) {
	return preprocessorNoop{}
}
