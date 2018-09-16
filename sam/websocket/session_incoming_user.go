package websocket

import (
	"context"
	"github.com/crusttech/crust/sam/websocket/incoming"
)

func (s *Session) userList(ctx context.Context, p *incoming.Users) error {
	users, err := s.svc.user.With(ctx).Find(nil)
	if err != nil {
		return err
	}
	return s.sendReply(payloadFromUsers(users))
}
