package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/compose_charts.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	composeChartCreate struct {
		Done chan struct{}
		res  *types.Chart
		err  error
	}

	composeChartUpdate struct {
		Done chan struct{}
		res  *types.Chart
		err  error
	}

	composeChartRemove struct {
		Done chan struct{}
		res  *types.Chart
		err  error
	}
)

// CreateComposeChart creates a new ComposeChart
// create job that can be pushed to store's transaction handler
func CreateComposeChart(res *types.Chart) *composeChartCreate {
	return &composeChartCreate{res: res}
}

// Do Executes composeChartCreate job
func (j *composeChartCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateComposeChart(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateComposeChart creates a new ComposeChart
// update job that can be pushed to store's transaction handler
func UpdateComposeChart(res *types.Chart) *composeChartUpdate {
	return &composeChartUpdate{res: res}
}

// Do Executes composeChartUpdate job
func (j *composeChartUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateComposeChart(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveComposeChart creates a new ComposeChart
// remove job that can be pushed to store's transaction handler
func RemoveComposeChart(res *types.Chart) *composeChartRemove {
	return &composeChartRemove{res: res}
}

// Do Executes composeChartRemove job
func (j *composeChartRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveComposeChart(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}
