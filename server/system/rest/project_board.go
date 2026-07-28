package rest

import (
	"context"

	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ProjectBoard struct {
		boardSvc projectBoardService
	}

	projectBoardService interface {
		Board(ctx context.Context, rr *types.ProjectBoardRequest) (*types.ProjectBoardResult, error)
	}
)

func (ProjectBoard) New() *ProjectBoard {
	return &ProjectBoard{
		boardSvc: service.DefaultProjectBoard,
	}
}

func (ctrl *ProjectBoard) Board(ctx context.Context, r *request.ProjectBoardBoard) (interface{}, error) {
	return ctrl.boardSvc.Board(ctx, &types.ProjectBoardRequest{
		ProjectID:  r.ProjectID,
		RevisionID: r.RevisionID,
		Status:     r.Status,
		PageCursor: r.PageCursor,
		Limit:      r.Limit,
	})
}
