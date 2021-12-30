package gig

import "encoding/json"

type (
	Preprocessor interface {
		Ref() preprocessor
		Worker() []string
		Params() interface{}
	}

	PreprocessorWrap struct {
		Ref    preprocessor    `json:"ref"`
		Params json.RawMessage `json:"params"`
	}
	PreprocessorWrapSet []PreprocessorWrap

	preprocessor string
)

func PreprocessorNoop() (out preprocessorNoop, err error) {
	return preprocessorNoop{}, nil
}

func (p preprocessorNoop) Ref() preprocessor {
	return PreprocessorHandleNoop
}

func (p preprocessorNoop) Worker() []string {
	return nil
}

func (p preprocessorNoop) Params() interface{} {
	return nil
}

func ParsePreprocessorWrap(ss []string) (out PreprocessorWrapSet, err error) {
	for _, s := range ss {
		aux := make(PreprocessorWrapSet, 0, 2)
		err = json.Unmarshal([]byte(s), &aux)
		if err != nil {
			return
		}

		out = append(out, aux...)
	}
	return
}
