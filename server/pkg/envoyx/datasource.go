package envoyx

type (
	Provider interface {
		Next(out map[string]string) (more bool, err error)
		Reset() error
		Ident() string
	}

	Datasource interface {
		Next(out map[string]string) (ident string, more bool, err error)
		Reset() error
		SetProvider(Provider) bool
	}
)

func SetDecoderSources(nn NodeSet, dd ...Provider) {
	for _, n := range nn {
		if n.Datasource == nil {
			continue
		}

		for _, d := range dd {
			if n.Datasource.SetProvider(d) {
				break
			}
		}
	}
}
