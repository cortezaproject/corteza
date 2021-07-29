package report

import (
	"context"
	"errors"
	"fmt"
)

type (
	stepLoad struct {
		dsp DatasourceProvider
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

func (j *stepLoad) Run(ctx context.Context, _ ...Datasource) (Datasource, error) {
	return j.dsp.Datasource(ctx, j.Def().Load)
}

func (j *stepLoad) Validate() error {
	pfx := "invalid load step: "

	// base things...
	switch {
	case j.def.Name == "":
		return errors.New(pfx + "dimension name not defined")

	case j.def.Source == "":
		return errors.New(pfx + "datasource not defined")
	case j.def.Definition == nil:
		return errors.New(pfx + "source definition not provided")
	}

	// provider
	switch {
	case j.dsp == nil:
		return errors.New(pfx + "datasource provider not defined")
	}

	// columns...
	for i, g := range j.def.Columns {
		if g.Name == "" {
			return fmt.Errorf("%scolumn key alias missing for column: %d", pfx, i)
		}
	}

	return nil
}

func (d *stepLoad) Name() string {
	return d.def.Name
}

func (d *stepLoad) Source() []string {
	return nil
}

func (d *stepLoad) Def() *StepDefinition {
	return &StepDefinition{Load: d.def}
}
