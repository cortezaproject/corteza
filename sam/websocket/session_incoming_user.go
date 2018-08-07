package websocket

import (
	"context"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/websocket/incoming"
)

func (s *Session) userList(ctx context.Context, p *incoming.UserList) error {
	users, err := service.User().Find(ctx, nil)
	if err != nil {
		return err
	}
	return s.sendReply(payloadFromUsers(users))
}
