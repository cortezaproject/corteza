package options

type (
	Options struct {
		Corteza     CortezaOpt
		ES          EsOpt
		Indexer     IndexerOpt
		Searcher    SearcherOpt
		Environment EnvironmentOpt
		HTTPServer  HttpServerOpt
		WaitFor     WaitForOpt
	}
)

// @todo: add codegen for options

func Init() (opt *Options, err error) {
	indexer, err := Indexer()
	if err != nil {
		return
	}

	searcher, err := Searcher()
	if err != nil {
		return
	}

	es, err := ES()
	if err != nil {
		return
	}

	corteza, err := Corteza()
	if err != nil {
		return
	}

	return &Options{
		Corteza:     *corteza,
		ES:          *es,
		Indexer:     *indexer,
		Searcher:    *searcher,
		Environment: *Environment(),
		HTTPServer:  *HttpServer(),
		WaitFor:     *WaitFor(),
	}, nil
}
