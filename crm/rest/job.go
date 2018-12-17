package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/service"
)

var _ = errors.Wrap

type (
	Job struct {
		jobSvc service.JobService
	}
)

func (Job) New() *Job {
	return &Job{jobSvc: service.DefaultJob}
}

func (ctrl *Job) List(ctx context.Context, r *request.JobList) (interface{}, error) {
	return nil, errors.New("Not implemented: Job.list")
}

func (ctrl *Job) Run(ctx context.Context, r *request.JobRun) (interface{}, error) {
	return nil, errors.New("Not implemented: Job.run")
}

func (ctrl *Job) Get(ctx context.Context, r *request.JobGet) (interface{}, error) {
	return nil, errors.New("Not implemented: Job.get")
}

func (ctrl *Job) Logs(ctx context.Context, r *request.JobLogs) (interface{}, error) {
	return nil, errors.New("Not implemented: Job.logs")
}

func (ctrl *Job) Update(ctx context.Context, r *request.JobUpdate) (interface{}, error) {
	return nil, errors.New("Not implemented: Job.update")
}

func (ctrl *Job) Delete(ctx context.Context, r *request.JobDelete) (interface{}, error) {
	return nil, errors.New("Not implemented: Job.delete")
}
