package websocket

import (
	"context"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/incoming"
	systemService "github.com/crusttech/crust/system/service"
)

func (s *Session) userList(ctx context.Context, p *incoming.Users) error {
	users, err := systemService.User(ctx).Find(nil)
	if err != nil {
		return err
	}
	return s.sendReply(payload.Users(users))
}
