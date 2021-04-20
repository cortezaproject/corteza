package websocket

import (
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/pkg/errors"
)

func (s *Session) dispatch(raw []byte) error {
	var p, err = payload.Unmarshal(raw)
	if err != nil {
		return errors.Wrap(err, "Session.incoming: payload malformed")
	}

	ctx := s.Context()

	switch {
	// Access token
	case p.Token != nil:
		return s.authenticate(ctx, p.Token)
	}

	return nil
}
