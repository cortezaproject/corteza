package report

type (
	stepLoad struct {
		ds  Datasource
		def *LoadStepDefinition
	}

	loadedDataset struct {
		def *LoadStepDefinition
		ds  Datasource
	}

	LoadStepDefinition struct {
		Name       string                 `json:"name"`
		Source     string                 `json:"source"`
		Definition map[string]interface{} `json:"definition"`
		Columns    FrameColumnSet         `json:"columns"`
		Rows       *RowDefinition         `json:"rows,omitempty"`
	}
)
