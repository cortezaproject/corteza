package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/messaging/websocket"
)

var _ = errors.Wrap

type Status struct {
	// xxx service.XXXService
}

func (Status) New() *Status {
	return &Status{}
}

func (ctrl *Status) List(ctx context.Context, r *request.StatusList) (interface{}, error) {
	type userStatus struct {
		UserID  uint64 `json:"userID,string"`
		Status  string `json:"present"`
		Icon    string `json:"icon"`
		Message string `json:"message"`
	}

	out := []userStatus{}
	for _, userID := range websocket.GetConnectedUsers() {
		out = append(out, userStatus{
			UserID: userID,
			Status: "online",
		})
	}

	return out, nil
}

func (ctrl *Status) Set(ctx context.Context, r *request.StatusSet) (interface{}, error) {
	return nil, errors.New("Not implemented: Status.set")
}

func (ctrl *Status) Delete(ctx context.Context, r *request.StatusDelete) (interface{}, error) {
	return nil, errors.New("Not implemented: Status.delete")
}
