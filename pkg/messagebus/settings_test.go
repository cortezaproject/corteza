package messagebus

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_settingsUnmarshal(t *testing.T) {
	var (
		tcc = []struct {
			name    string
			payload string
			err     error
			expect  QueueSettingsMeta
		}{
			{
				name:    "settings defaults",
				payload: `{}`,
				err:     nil,
				expect:  QueueSettingsMeta{PollDelay: nil, DispatchEvents: true},
			},
			{
				name:    "settings enabled dispatch events",
				payload: `{"poll_delay": "7s", "dispatch_events": true}`,
				err:     nil,
				expect:  QueueSettingsMeta{PollDelay: makeDelay(time.Second * 7), DispatchEvents: true},
			},
			{
				name:    "settings disabled dispatch events",
				payload: `{"poll_delay": "7s", "dispatch_events": false}`,
				err:     nil,
				expect:  QueueSettingsMeta{PollDelay: makeDelay(time.Second * 7), DispatchEvents: false},
			},
			{
				name:    "settings invalid poll delay",
				payload: `{"poll_delay": "7seconds", "dispatch_events": false}`,
				err:     nil,
				expect:  QueueSettingsMeta{PollDelay: nil, DispatchEvents: false},
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			s := &QueueSettingsMeta{}
			err := json.Unmarshal([]byte(tc.payload), s)

			req.Equal(tc.err, err)
			req.Equal(tc.expect, *s)
		})
	}
}
