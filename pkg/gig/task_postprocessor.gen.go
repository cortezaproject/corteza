package gig

type (
	postprocessorNoop    struct{}
	postprocessorDiscard struct{}
	postprocessorSave    struct{}
	postprocessorArchive struct {
		encoding archive
		name     string
	}
)

var (
	PostprocessorHandleArchive = "archive"
	PostprocessorHandleDiscard = "discard"
	PostprocessorHandleNoop    = "noop"
	PostprocessorHandleSave    = "save"
)

func postprocessorDefinitions() TaskDefSet {
	return TaskDefSet{
		{
			Ref:  string(PostprocessorHandleDiscard),
			Kind: TaskDecoder,
		},
		{
			Ref:  string(PostprocessorHandleNoop),
			Kind: TaskDecoder,
		},
		{
			Ref:  string(PostprocessorHandleSave),
			Kind: TaskDecoder,
		},
		{
			Ref:  string(PostprocessorHandleArchive),
			Kind: TaskDecoder,
			Params: []taskDefParam{
				{
					Name:        "format",
					Kind:        "String",
					Required:    false,
					Description: "What archive format to encode into",
				},
				{
					Name:        "name",
					Kind:        "String",
					Required:    false,
					Description: "The output name of the archive",
				},
			},
		},
	}
}
