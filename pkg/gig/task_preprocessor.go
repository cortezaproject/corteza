package gig

type (
	Preprocessor interface {
		Worker() []string
		Ref() string
		Params() map[string]interface{}
	}
	PreprocessorSet []Preprocessor
)

func PreprocessorNoopParams(_ map[string]interface{}) (out preprocessorNoop) {
	return PreprocessorNoop()
}

func PreprocessorNoop() (out preprocessorNoop) {
	return preprocessorNoop{}
}

func (p preprocessorNoop) Ref() string {
	return PreprocessorHandleNoop
}

func (p preprocessorNoop) Worker() []string {
	return nil
}

func (p preprocessorNoop) Params() map[string]interface{} {
	return nil
}
