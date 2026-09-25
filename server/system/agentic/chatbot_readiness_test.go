package agentic

import (
	"testing"

	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

// withNote is what carries the readiness warning to the caller; the note itself
// needs a live agent service, so what is pinned here is the part that does not.
func TestWithNote(t *testing.T) {
	t.Run("nothing to say leaves the extras alone", func(t *testing.T) {
		require.Nil(t, withNote(nil, ""))
		in := map[string]string{"url": "http://x"}
		require.Equal(t, in, withNote(in, ""))
		require.NotContains(t, withNote(in, ""), "note")
	})

	t.Run("a note is added alongside the links", func(t *testing.T) {
		out := withNote(map[string]string{"url": "http://x"}, "not ready")
		require.Equal(t, "not ready", out["note"])
		require.Equal(t, "http://x", out["url"])
	})

	t.Run("a note survives having no extras to join", func(t *testing.T) {
		require.Equal(t, "not ready", withNote(nil, "not ready")["note"])
	})
}

// TL;DR: a chatbot with nothing to warn about produces no note.
// Example: the note must not fire on form scenarios or on a chatbot that names
// no agent, or every write carries a warning nobody can act on.
func TestChatbotReadinessNote_NothingToWarnAbout(t *testing.T) {
	require.Empty(t, chatbotReadinessNote(nil))
	require.Empty(t, chatbotReadinessNote(&sysTypes.Chatbot{}))

	formOnly := &sysTypes.Chatbot{Scenarios: sysTypes.ChatbotScenarios{
		{ID: "f", Type: "form"},
	}}
	require.Empty(t, chatbotReadinessNote(formOnly))

	// A conversation scenario with no agent is the service's error to report,
	// not a readiness warning.
	noAgent := &sysTypes.Chatbot{Scenarios: sysTypes.ChatbotScenarios{
		{ID: "c", Type: "conversation"},
	}}
	require.Empty(t, chatbotReadinessNote(noAgent))
}

func TestChatbotReadinessNote_RunAs(t *testing.T) {
	conversation := func(runAs uint64) *sysTypes.Chatbot {
		return &sysTypes.Chatbot{Scenarios: sysTypes.ChatbotScenarios{
			{ID: "talk", Type: "conversation", AgentID: 7, RunAs: runAs},
		}}
	}

	note := chatbotReadinessNote(conversation(0))
	require.Contains(t, note, `"talk"`)
	require.Contains(t, note, "runAs")

	require.Empty(t, chatbotReadinessNote(conversation(42)))
}
