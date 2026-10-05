package wfexec

import (
	"encoding/json"

	"github.com/crusttech/human/server/pkg/id"
)

type (
	frameAlias Frame

	// frameJSON is a Frame as JSON carries it: IDs written as strings, read
	// from strings or numbers.
	frameJSON struct {
		*frameAlias
		SessionID id.Uint64 `json:"sessionID"`
		StateID   id.Uint64 `json:"stateID"`
		ParentID  id.Uint64 `json:"parentID"`
		StepID    id.Uint64 `json:"stepID"`
	}
)

func (f Frame) MarshalJSON() ([]byte, error) {
	return json.Marshal(frameJSON{
		frameAlias: (*frameAlias)(&f),
		SessionID:  id.Uint64(f.SessionID),
		StateID:    id.Uint64(f.StateID),
		ParentID:   id.Uint64(f.ParentID),
		StepID:     id.Uint64(f.StepID),
	})
}

func (f *Frame) UnmarshalJSON(data []byte) error {
	aux := frameJSON{frameAlias: (*frameAlias)(f)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	f.SessionID = uint64(aux.SessionID)
	f.StateID = uint64(aux.StateID)
	f.ParentID = uint64(aux.ParentID)
	f.StepID = uint64(aux.StepID)
	return nil
}
