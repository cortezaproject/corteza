package service

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/crusttech/human/server/system/types"
)

// decode mirrors runPollAction: numbers stay json.Number to preserve large cursors.
func decodePoll(t *testing.T, s string) any {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader([]byte(s)))
	dec.UseNumber()
	var v any
	require.NoError(t, dec.Decode(&v))
	return v
}

func TestPollEngine_GmailShape(t *testing.T) {
	// A Gmail users.history.list response: items are nested two arrays deep.
	list := decodePoll(t, `{
		"history": [
			{"messagesAdded": [{"message": {"id": "m1", "threadId": "t1", "labelIds": ["INBOX","UNREAD"]}}]},
			{"messagesAdded": [
				{"message": {"id": "m2", "threadId": "t2", "labelIds": ["SENT"]}},
				{"message": {"id": "m3", "threadId": "t3", "labelIds": ["INBOX"]}}
			]}
		],
		"historyId": 9876543210
	}`)

	// Flatten history.*.messagesAdded.*.message → three message objects.
	items := collectByPath(list, []string{"history", "*", "messagesAdded", "*", "message"})
	require.Len(t, items, 3)

	// Cursor keeps full precision (json.Number, not float64).
	require.Equal(t, "9876543210", extractString(list, []string{"historyId"}))

	// INBOX filter: m1 and m3 pass, m2 (SENT) is excluded.
	inbox := &types.ConnectionPollFilter{Field: []string{"labelIds"}, Contains: "INBOX"}
	var passed []string
	for _, it := range items {
		if pollItemPasses(it, inbox) {
			passed = append(passed, extractString(it, []string{"id"}))
		}
	}
	require.Equal(t, []string{"m1", "m3"}, passed)
}

func TestPollItemVars_HeaderSelectors(t *testing.T) {
	// An enriched Gmail message (format=metadata): from/subject live in headers.
	msg := decodePoll(t, `{
		"id": "m1",
		"threadId": "t1",
		"snippet": "hello there",
		"payload": {"headers": [
			{"name": "From", "value": "alice@example.com"},
			{"name": "Subject", "value": "Weekly sync"}
		]}
	}`)

	payload := []types.ConnectionWebhookField{
		{Name: "messageId", Selector: []string{"id"}},
		{Name: "threadId", Selector: []string{"threadId"}},
		{Name: "snippet", Selector: []string{"snippet"}},
		{Name: "from", Selector: []string{"payload", "headers[name=From]", "value"}},
		{Name: "subject", Selector: []string{"payload", "headers[name=Subject]", "value"}},
	}

	vars := pollItemVars(msg, 42, payload)
	require.Equal(t, "m1", vars["messageId"])
	require.Equal(t, "t1", vars["threadId"])
	require.Equal(t, "hello there", vars["snippet"])
	require.Equal(t, "alice@example.com", vars["from"])
	require.Equal(t, "Weekly sync", vars["subject"])
	require.Equal(t, uint64(42), vars["configuredConnectionID"])
}
