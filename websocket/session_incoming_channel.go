package websocket

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/payload/incoming"
)

func (s *Session) authenticate(ctx context.Context, p *incoming.Token) error {
	// Get JWT claims
	claims, err := p.ParseWithClaims()
	if err != nil {
		s.Close()
		return err
	}

	// Get identity using JWT claims
	identity := auth.ClaimsToIdentity(claims)

	// Update the existing ws sessions if exists or create new one
	if store.CountConnections(identity.Identity()) > 0 {
		store.Walk(func(session *Session) {
			session.user = identity
		})
	} else {
		s.user = identity
		store.Save(s)
	}

	return nil
}
