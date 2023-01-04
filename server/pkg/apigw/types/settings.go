package types

type (
	Settings struct {
		Enabled         bool
		ProfilerGlobal  bool
		ProfilerEnabled bool
	}

	AddSettings func(s Setter)

	Setter interface {
		UpdateSettings(st ...AddSettings)
	}
)
