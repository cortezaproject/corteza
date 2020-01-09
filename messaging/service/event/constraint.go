package event

type (
	constraint interface {
		Name() string
		Values() []string
		Match(value string) bool
	}
)
