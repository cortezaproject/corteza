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

func TestActionPlaceholders_ConnectionVsScoped(t *testing.T) {
	// Gmail-style list: only {{cursor}} — connection-level, no scope needed.
	gmail := &types.ConnectionHTTPAction{
		Path: types.ConnectionTemplate{Value: "users/me/history"},
		QueryParams: map[string]types.ConnectionTemplate{
			"startHistoryId": {Value: "{{cursor}}"},
		},
	}
	ph := actionPlaceholders(gmail)
	delete(ph, "cursor")
	require.Empty(t, ph, "gmail poll needs no trigger scope")

	// Sheets-style list: references {{spreadsheetId}} + {{sheetName}} — scoped.
	sheets := &types.ConnectionHTTPAction{
		Path: types.ConnectionTemplate{Value: "{{spreadsheetId}}/values/{{sheetName}}"},
	}
	ph = actionPlaceholders(sheets)
	delete(ph, "cursor")
	require.True(t, ph["spreadsheetId"] && ph["sheetName"])
	require.Len(t, ph, 2)
}

func TestScopeSig_StableAndDistinct(t *testing.T) {
	params := map[string]bool{"spreadsheetId": true, "sheetName": true}
	a := map[string]string{"configurationID": "9", "spreadsheetId": "A", "sheetName": "Tab1"}
	b := map[string]string{"configurationID": "9", "spreadsheetId": "A", "sheetName": "Tab2"}

	// Deterministic: same scope -> same signature, ignoring non-param keys.
	require.Equal(t, scopeSig(a, params), scopeSig(a, params))
	require.Equal(t, "sheetName=Tab1|spreadsheetId=A", scopeSig(a, params))
	// Different tab -> different cursor key.
	require.NotEqual(t, scopeSig(a, params), scopeSig(b, params))
}

func TestPollVars_MergesScopeAndCursor(t *testing.T) {
	v := pollVars(map[string]string{"spreadsheetId": "X"}, "42")
	require.Equal(t, "42", v["cursor"])
	require.Equal(t, "X", v["spreadsheetId"])
}

func TestBuildSnapshot_KeyedAndPositional(t *testing.T) {
	sheets := collectByPath(decodePoll(t, `{"sheets":[
		{"properties":{"sheetId":"10","title":"A","index":0}},
		{"properties":{"sheetId":"20","title":"B","index":1}}
	]}`), []string{"sheets", "*"})
	snap, bodies := buildSnapshot(sheets, []string{"properties", "sheetId"})
	require.Len(t, snap, 2)
	require.Contains(t, snap, "10")
	require.Contains(t, snap, "20")
	require.NotNil(t, bodies["10"])

	// Rows have no stable key — fall back to position.
	rows := collectByPath(decodePoll(t, `{"values":[["a","b"],["c","d"]]}`), []string{"values", "*"})
	rsnap, _ := buildSnapshot(rows, nil)
	require.Len(t, rsnap, 2)
	require.Contains(t, rsnap, "0")
	require.Contains(t, rsnap, "1")
}

func TestDiffSnapshots_AddedRemovedChanged(t *testing.T) {
	old := map[string]string{"10": `{"t":"A"}`, "20": `{"t":"B"}`, "99": `{"t":"Z"}`}
	cur := map[string]string{"10": `{"t":"A"}`, "20": `{"t":"B2"}`, "30": `{"t":"C"}`}
	added, removed, changed := diffSnapshots(old, cur)
	require.ElementsMatch(t, []string{"30"}, added)   // new sheet -> sheet.created
	require.ElementsMatch(t, []string{"99"}, removed) // vanished sheet -> sheet.deleted
	require.ElementsMatch(t, []string{"20"}, changed) // renamed sheet -> sheet.updated
}

func TestBuildSnapshot_StableJSON_NoFalseChange(t *testing.T) {
	// Same content, different key order in the source must hash identically,
	// so a re-poll does not report a spurious "changed".
	a := collectByPath(decodePoll(t, `{"s":[{"a":1,"b":2}]}`), []string{"s", "*"})
	b := collectByPath(decodePoll(t, `{"s":[{"b":2,"a":1}]}`), []string{"s", "*"})
	sa, _ := buildSnapshot(a, nil)
	sb, _ := buildSnapshot(b, nil)
	require.Equal(t, sa["0"], sb["0"])
}
