package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crusttech/human/server/system/agentic/observability"
	"github.com/crusttech/human/server/system/types"
)

func newTestPreview(t *testing.T, cap int) *chatbotPreview {
	t.Helper()
	return &chatbotPreview{
		byID: make(map[string]*ChatbotPreviewSession),
		cap:  cap,
		bus:  observability.NewBus(),
	}
}

func minimalChatbot() types.Chatbot {
	return types.Chatbot{
		Scenarios: types.ChatbotScenarios{
			{ID: "intro", Type: "form"},
			{ID: "chat", Type: "conversation", AgentID: 1},
			{ID: "outro", Type: "static_message"},
		},
	}
}

func TestChatbotPreview_OpenEvictsOldest(t *testing.T) {
	s := newTestPreview(t, 3)

	a := s.Open(minimalChatbot(), 100)
	b := s.Open(minimalChatbot(), 101)
	c := s.Open(minimalChatbot(), 102)

	assert.Equal(t, 3, len(s.order))
	assert.NotNil(t, s.Find(a.ID))
	assert.NotNil(t, s.Find(b.ID))
	assert.NotNil(t, s.Find(c.ID))

	d := s.Open(minimalChatbot(), 103)

	assert.Equal(t, 3, len(s.order))
	assert.Nil(t, s.Find(a.ID), "oldest should be evicted")
	assert.NotNil(t, s.Find(b.ID))
	assert.NotNil(t, s.Find(c.ID))
	assert.NotNil(t, s.Find(d.ID))
}

func TestChatbotPreview_StepProgression(t *testing.T) {
	s := newTestPreview(t, 8)
	ps := s.Open(minimalChatbot(), 200)

	s.StartStep(ps, 0)
	step := s.CurrentStep(ps)
	assert.NotNil(t, step)
	assert.Equal(t, "active", step.Status)
	assert.Equal(t, 0, step.ScenarioIndex)

	s.FinalizeStep(ps)
	assert.Equal(t, "complete", step.Status)

	s.StartStep(ps, 1)
	step = s.CurrentStep(ps)
	assert.Equal(t, 1, step.ScenarioIndex)
	assert.Equal(t, "active", step.Status)

	s.FinalizeStep(ps)
	s.StartStep(ps, 2)
	step = s.CurrentStep(ps)
	assert.Equal(t, 2, step.ScenarioIndex)

	s.FinalizeStep(ps)
	// Past last scenario → session closes.
	s.StartStep(ps, 3)
	assert.Equal(t, "closed", ps.Status)
}

func TestChatbotPreview_HandoffResumes(t *testing.T) {
	s := newTestPreview(t, 8)
	ps := s.Open(minimalChatbot(), 300)
	s.StartStep(ps, 1)

	h := s.RequestHandoff(ps)
	assert.Equal(t, "handoff_requested", ps.Status)
	assert.Equal(t, "requested", h.Status)

	s.ActivateHandoff(ps, "alice")
	assert.Equal(t, "handoff_active", ps.Status)
	assert.Equal(t, "active", ps.Handoff.Status)

	s.CloseHandoff(ps)
	assert.Equal(t, "active", ps.Status, "session must resume after handoff close")
	assert.Equal(t, "closed", ps.Handoff.Status)
	assert.NotNil(t, ps.Handoff.Closed)
}

func TestChatbotPreview_HooksLogged(t *testing.T) {
	s := newTestPreview(t, 8)
	cb := minimalChatbot()
	cb.Scenarios[1].Automation = types.ChatbotScenarioAutomation{
		Before: types.ChatbotAutomationHook{Automation: "corteza::automation:ng-automation/1", Async: false},
		After:  types.ChatbotAutomationHook{Automation: "corteza::automation:ng-automation/2", Async: true},
	}
	ps := s.Open(cb, 400)

	s.StartStep(ps, 1)
	step := s.CurrentStep(ps)
	assert.Len(t, step.HookLog, 1, "before hook should be logged but not executed")

	s.FinalizeStep(ps)
	assert.Len(t, step.HookLog, 2, "after hook should be logged but not executed")
}
