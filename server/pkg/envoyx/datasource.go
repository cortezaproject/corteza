package envoyx

import "context"

type (
	Provider interface {
		Next(ctx context.Context, out map[string]string) (more bool, err error)
		Reset(ctx context.Context) error
		SetIdent(string)
		Ident() string
	}

	Datasource interface {
		Next(ctx context.Context, out map[string]string) (ident []string, more bool, err error)
		Reset(ctx context.Context) error
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
