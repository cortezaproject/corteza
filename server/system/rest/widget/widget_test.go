package widget

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

// sessionSvcStub records the steps started and what has been stored.
type sessionSvcStub struct {
	ChatbotSessionService

	mu      sync.Mutex
	started []int
	steps   types.ChatbotSessionStepSet
}

func (s *sessionSvcStub) Open(context.Context, *types.Chatbot) (*types.ChatbotSession, *types.AiConversation, error) {
	return &types.ChatbotSession{ID: 11}, &types.AiConversation{ID: 22}, nil
}

func (s *sessionSvcStub) Start(_ context.Context, _ *types.Chatbot, _, _ uint64, idx int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.started = append(s.started, idx)
	s.steps = append(s.steps, &types.ChatbotSessionStep{ScenarioIndex: idx})
}

func (s *sessionSvcStub) FindStepsBySession(context.Context, uint64) (types.ChatbotSessionStepSet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.steps, nil
}

func (s *sessionSvcStub) startedCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.started)
}

func TestFirstStepStartsWhenTheStreamListens(t *testing.T) {
	svc := &sessionSvcStub{}
	c := &Controller{sessionSvc: svc, serverSecret: "secret"}
	cb := &types.Chatbot{ID: 1, WidgetKey: "key", Scenarios: types.ChatbotScenarios{{ID: "hello"}}}

	// Opening a session starts nothing: nobody is listening yet.
	req := httptest.NewRequest(http.MethodPost, "/session", nil)
	req = req.WithContext(context.WithValue(req.Context(), ctxKeyChatbot, cb))
	rec := httptest.NewRecorder()
	c.createSession(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	time.Sleep(20 * time.Millisecond)
	require.Zero(t, svc.startedCount())

	// Two streams connecting at once start step 0 once.
	claims := &sessionClaims{Dbsid: 11, Cid: 22}
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.startFirstStep(context.Background(), cb, claims)
		}()
	}
	wg.Wait()
	require.Eventually(t, func() bool { return svc.startedCount() == 1 }, time.Second, 5*time.Millisecond)

	// A stream that reconnects finds the step started and leaves it be.
	c.startFirstStep(context.Background(), cb, claims)
	time.Sleep(20 * time.Millisecond)
	require.Equal(t, []int{0}, svc.started)
}
