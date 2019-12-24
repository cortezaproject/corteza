package corredor

type (
	Script struct {
		Name        string
		Label       string
		Description string
		Errors      []string
	}

	ScriptSet []*Script

	MatchedScriptSet []MatchedScript
	MatchedScript    struct {
		script *Script
	}
)
