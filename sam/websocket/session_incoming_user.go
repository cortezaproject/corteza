package websocket

import (
	"context"
	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/incoming"
)

func (s *Session) userList(ctx context.Context, p *incoming.Users) error {
	users, err := s.svc.user.With(ctx).Find(nil)
	if err != nil {
		return err
	}
	return s.sendReply(payload.Users(users))
}
