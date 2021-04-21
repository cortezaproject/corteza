package websocket

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

func (s *session) authenticate(ctx context.Context, p *Auth) error {
	// Get JWT claims
	claims, err := p.ParseWithClaims()
	if err != nil {
		s.Close()
		return err
	}

	// Get identity using JWT claims
	identity := auth.ClaimsToIdentity(claims)

	// Update the existing ws sessions if exists or create new one
	if existingSession := s.Get(identity.Identity()); existingSession == nil {
		s.user = identity
		s.Save()
	} else {
		existingSession.user = identity
	}

	return nil
}
