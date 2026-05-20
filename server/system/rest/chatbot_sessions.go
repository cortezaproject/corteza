package rest

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/crusttech/human/server/pkg/api"
	pkgAuth "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/agentic/observability"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	// ChatbotSession is the unified admin/operator REST surface that
	// aggregates widget-persisted and preview in-memory chatbot sessions
	// behind a single set of routes. Per-row source is exposed via the
	// "source" field on the wire payload.
	ChatbotSession struct {
		chatbotSvc chatbotLookupService
		sessionSvc chatbotSessionAdminService
		previewSvc previewSessionsService
		bus        *observability.Bus
		ac         chatbotSessionsAccessController
	}

	// chatbotUnifiedSession is the wire shape returned by /chatbots/sessions
	// for both widget-persisted and preview in-memory sessions.
	//
	// The chatbot name is intentionally not embedded — clients resolve it
	// from chatbotID via the chatbot list so renames flow through without
	// stale snapshots. Preview rows for unsaved drafts carry chatbotID "0"
	// and the client falls back to a generic label.
	chatbotUnifiedSession struct {
		ID             string                 `json:"id"`
		Source         string                 `json:"source"`
		ChatbotID      uint64                 `json:"chatbotID,string"`
		Status         string                 `json:"status"`
		CurrentStep    int                    `json:"currentStep"`
		ConversationID uint64                 `json:"conversationID,string,omitempty"`
		Handoff        *chatbotUnifiedHandoff `json:"handoff,omitempty"`
		CreatedAt      time.Time              `json:"createdAt,omitempty"`
		UpdatedAt      *time.Time             `json:"updatedAt,omitempty"`

		CanView          bool `json:"canView"`
		CanManage        bool `json:"canManage"`
		CanManageHandoff bool `json:"canManageHandoff"`
	}

	chatbotUnifiedHandoff struct {
		ID        string     `json:"id"`
		Status    string     `json:"status"`
		StartedAt time.Time  `json:"startedAt,omitempty"`
		ClosedAt  *time.Time `json:"closedAt,omitempty"`
	}

	chatbotUnifiedReadPayload struct {
		Session *chatbotUnifiedSession `json:"session"`
		Steps   []chatbotUnifiedStep   `json:"steps"`
	}

	chatbotUnifiedStep struct {
		ID             string    `json:"id,omitempty"`
		ScenarioIndex  int       `json:"scenarioIndex"`
		ScenarioID     string    `json:"scenarioID,omitempty"`
		ConversationID uint64    `json:"conversationID,string,omitempty"`
		Status         string    `json:"status"`
		HookLog        []string  `json:"hookLog,omitempty"`
		CreatedAt      time.Time `json:"createdAt,omitempty"`
	}

	chatbotUnifiedSetPayload struct {
		Filter chatbotUnifiedFilter     `json:"filter"`
		Set    []*chatbotUnifiedSession `json:"set"`
	}

	chatbotUnifiedFilter struct {
		ChatbotID uint64   `json:"chatbotID,omitempty,string"`
		Status    []string `json:"status,omitempty"`
		Source    string   `json:"source,omitempty"`
		filter.Sorting
		filter.Paging
	}

	chatbotLookupService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Chatbot, error)
	}

	chatbotSessionAdminService interface {
		Search(ctx context.Context, f types.ChatbotSessionFilter) (types.ChatbotSessionSet, types.ChatbotSessionFilter, error)
		FindByID(ctx context.Context, ID uint64) (*types.ChatbotSession, error)
		FindStepsBySession(ctx context.Context, sessionID uint64) (types.ChatbotSessionStepSet, error)
		FindStepBySession(ctx context.Context, sessionID uint64) (*types.ChatbotSessionStep, error)
		FindHandoffBySession(ctx context.Context, sessionID uint64) (*types.ChatbotSessionHandoff, error)

		AdvanceStep(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64) error
		CloseSession(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64) error
		AcceptHandoffByID(ctx context.Context, cb *types.Chatbot, sessionID, convID, handoffID uint64, operatorID uint64, operatorName string) error
		SendOperatorMessage(ctx context.Context, cb *types.Chatbot, sessionID, convID, handoffID uint64, msg, operator string) error
		CloseHandoffPublic(ctx context.Context, convID, handoffID uint64) error
	}

	previewSessionsService interface {
		List(chatbotID uint64, status []string) []*service.ChatbotPreviewSession
		Find(id string) *service.ChatbotPreviewSession
		AdvanceStep(ps *service.ChatbotPreviewSession) error
		CloseSession(ps *service.ChatbotPreviewSession) error
		AcceptHandoff(ps *service.ChatbotPreviewSession, operator string) error
		SendOperatorMessage(ctx context.Context, ps *service.ChatbotPreviewSession, msg, operator string) error
		CloseHandoff(ps *service.ChatbotPreviewSession)
	}

	chatbotSessionsAccessController interface {
		CanViewSessionsOnChatbot(context.Context, *types.Chatbot) bool
		CanManageSessionsOnChatbot(context.Context, *types.Chatbot) bool
		CanManageSessionsHandoffOnChatbot(context.Context, *types.Chatbot) bool
	}

	// sessionRef holds either a widget-source DB session or a preview in-memory
	// session along with the resolved chatbot. Exactly one of widgetSession /
	// previewPS is non-nil.
	sessionRef struct {
		source        string
		widgetSession *types.ChatbotSession
		previewPS     *service.ChatbotPreviewSession
		chatbot       *types.Chatbot
	}
)

func (ChatbotSession) New() *ChatbotSession {
	return &ChatbotSession{
		chatbotSvc: service.DefaultChatbot,
		sessionSvc: service.DefaultChatbotSession,
		previewSvc: service.DefaultChatbotPreview,
		bus:        service.DefaultObsBus,
		ac:         service.DefaultAccessControl,
	}
}

// ---- handlers ---------------------------------------------------------------

func (ctrl *ChatbotSession) List(ctx context.Context, r *request.ChatbotSessionList) (interface{}, error) {
	source := r.Source

	out := make([]*chatbotUnifiedSession, 0)

	if source == "" || source == "widget" {
		f := types.ChatbotSessionFilter{
			ChatbotID: r.ChatbotID,
			Status:    r.Status,
		}
		f.IncTotal = r.IncTotal
		set, _, err := ctrl.sessionSvc.Search(ctx, f)
		if err != nil {
			return nil, err
		}
		for _, s := range set {
			cb, cbErr := ctrl.chatbotSvc.FindByID(ctx, s.ChatbotID)
			if cbErr != nil || cb == nil {
				continue
			}
			if !ctrl.ac.CanViewSessionsOnChatbot(ctx, cb) {
				continue
			}
			h, _ := ctrl.sessionSvc.FindHandoffBySession(ctx, s.ID)
			out = append(out, ctrl.toUnifiedFromDB(ctx, s, cb, h))
		}
	}

	if source == "" || source == "preview" {
		for _, ps := range ctrl.previewSvc.List(r.ChatbotID, r.Status) {
			// Preview is admin-gated at the route level and operates on
			// unsaved drafts (cb.ID may be 0), so RBAC on the chatbot
			// resource is skipped here.
			out = append(out, ctrl.toUnifiedFromPreview(ctx, ps))
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})

	return &chatbotUnifiedSetPayload{
		Filter: chatbotUnifiedFilter{
			ChatbotID: r.ChatbotID,
			Status:    r.Status,
			Source:    source,
		},
		Set: out,
	}, nil
}

func (ctrl *ChatbotSession) Read(ctx context.Context, r *request.ChatbotSessionRead) (interface{}, error) {
	ref, err := ctrl.resolve(ctx, r.SessionID)
	if err != nil {
		return nil, err
	}
	if !ctrl.ac.CanViewSessionsOnChatbot(ctx, ref.chatbot) {
		return nil, errChatbotSessionForbidden
	}

	if ref.source == "widget" {
		h, _ := ctrl.sessionSvc.FindHandoffBySession(ctx, ref.widgetSession.ID)
		session := ctrl.toUnifiedFromDB(ctx, ref.widgetSession, ref.chatbot, h)
		steps, err := ctrl.sessionSvc.FindStepsBySession(ctx, ref.widgetSession.ID)
		if err != nil {
			return nil, err
		}
		out := make([]chatbotUnifiedStep, 0, len(steps))
		for _, s := range steps {
			out = append(out, chatbotUnifiedStep{
				ID:             strconv.FormatUint(s.ID, 10),
				ScenarioIndex:  s.ScenarioIndex,
				ConversationID: s.ConversationID,
				Status:         s.Status,
				CreatedAt:      s.CreatedAt,
			})
		}
		return &chatbotUnifiedReadPayload{Session: session, Steps: out}, nil
	}

	ps := ref.previewPS
	session := ctrl.toUnifiedFromPreview(ctx, ps)
	steps := make([]chatbotUnifiedStep, 0, len(ps.Steps))
	for _, st := range ps.Steps {
		steps = append(steps, chatbotUnifiedStep{
			ScenarioIndex: st.ScenarioIndex,
			ScenarioID:    st.ScenarioID,
			Status:        st.Status,
			HookLog:       st.HookLog,
		})
	}
	return &chatbotUnifiedReadPayload{Session: session, Steps: steps}, nil
}

func (ctrl *ChatbotSession) AdvanceStep(ctx context.Context, r *request.ChatbotSessionAdvanceStep) (interface{}, error) {
	ref, err := ctrl.resolve(ctx, r.SessionID)
	if err != nil {
		return nil, err
	}
	if !ctrl.ac.CanManageSessionsOnChatbot(ctx, ref.chatbot) {
		return nil, errChatbotSessionForbidden
	}

	if ref.source == "widget" {
		convID, err := ctrl.sessionConvID(ctx, ref.widgetSession.ID)
		if err != nil {
			return nil, err
		}
		if err := ctrl.sessionSvc.AdvanceStep(ctx, ref.chatbot, ref.widgetSession.ID, convID); err != nil {
			return nil, err
		}
		return api.OK(), nil
	}

	if err := ctrl.previewSvc.AdvanceStep(ref.previewPS); err != nil {
		return nil, err
	}
	return api.OK(), nil
}

func (ctrl *ChatbotSession) Close(ctx context.Context, r *request.ChatbotSessionClose) (interface{}, error) {
	ref, err := ctrl.resolve(ctx, r.SessionID)
	if err != nil {
		return nil, err
	}
	if !ctrl.ac.CanManageSessionsOnChatbot(ctx, ref.chatbot) {
		return nil, errChatbotSessionForbidden
	}

	if ref.source == "widget" {
		convID, err := ctrl.sessionConvID(ctx, ref.widgetSession.ID)
		if err != nil {
			return nil, err
		}
		if err := ctrl.sessionSvc.CloseSession(ctx, ref.chatbot, ref.widgetSession.ID, convID); err != nil {
			return nil, err
		}
		return api.OK(), nil
	}

	if err := ctrl.previewSvc.CloseSession(ref.previewPS); err != nil {
		return nil, err
	}
	return api.OK(), nil
}

func (ctrl *ChatbotSession) HandoffAccept(ctx context.Context, r *request.ChatbotSessionHandoffAccept) (interface{}, error) {
	ref, err := ctrl.resolve(ctx, r.SessionID)
	if err != nil {
		return nil, err
	}
	if !ctrl.ac.CanManageSessionsHandoffOnChatbot(ctx, ref.chatbot) {
		return nil, errChatbotSessionForbidden
	}

	if ref.source == "widget" {
		h, err := ctrl.sessionSvc.FindHandoffBySession(ctx, ref.widgetSession.ID)
		if err != nil || h == nil {
			return nil, service.ChatbotSessionErrHandoffNotFound()
		}
		convID, err := ctrl.sessionConvID(ctx, ref.widgetSession.ID)
		if err != nil {
			return nil, err
		}
		var operatorID uint64
		if id := pkgAuth.GetIdentityFromContext(ctx); id != nil {
			operatorID = id.Identity()
		}
		if err := ctrl.sessionSvc.AcceptHandoffByID(ctx, ref.chatbot, ref.widgetSession.ID, convID, h.ID, operatorID, r.Operator); err != nil {
			return nil, err
		}
		return api.OK(), nil
	}

	if err := ctrl.previewSvc.AcceptHandoff(ref.previewPS, r.Operator); err != nil {
		return nil, err
	}
	return api.OK(), nil
}

func (ctrl *ChatbotSession) OperatorMessage(ctx context.Context, r *request.ChatbotSessionOperatorMessage) (interface{}, error) {
	ref, err := ctrl.resolve(ctx, r.SessionID)
	if err != nil {
		return nil, err
	}
	if !ctrl.ac.CanManageSessionsHandoffOnChatbot(ctx, ref.chatbot) {
		return nil, errChatbotSessionForbidden
	}

	if ref.source == "widget" {
		h, err := ctrl.sessionSvc.FindHandoffBySession(ctx, ref.widgetSession.ID)
		if err != nil || h == nil {
			return nil, service.ChatbotSessionErrHandoffNotFound()
		}
		convID, err := ctrl.sessionConvID(ctx, ref.widgetSession.ID)
		if err != nil {
			return nil, err
		}
		if err := ctrl.sessionSvc.SendOperatorMessage(ctx, ref.chatbot, ref.widgetSession.ID, convID, h.ID, r.Message, r.Operator); err != nil {
			return nil, err
		}
		return api.OK(), nil
	}

	if err := ctrl.previewSvc.SendOperatorMessage(ctx, ref.previewPS, r.Message, r.Operator); err != nil {
		return nil, err
	}
	return api.OK(), nil
}

func (ctrl *ChatbotSession) HandoffComplete(ctx context.Context, r *request.ChatbotSessionHandoffComplete) (interface{}, error) {
	ref, err := ctrl.resolve(ctx, r.SessionID)
	if err != nil {
		return nil, err
	}
	if !ctrl.ac.CanManageSessionsHandoffOnChatbot(ctx, ref.chatbot) {
		return nil, errChatbotSessionForbidden
	}

	if ref.source == "widget" {
		h, err := ctrl.sessionSvc.FindHandoffBySession(ctx, ref.widgetSession.ID)
		if err != nil || h == nil {
			return nil, service.ChatbotSessionErrHandoffNotFound()
		}
		convID, err := ctrl.sessionConvID(ctx, ref.widgetSession.ID)
		if err != nil {
			return nil, err
		}
		if err := ctrl.sessionSvc.CloseHandoffPublic(ctx, convID, h.ID); err != nil {
			return nil, err
		}
		return api.OK(), nil
	}

	ctrl.previewSvc.CloseHandoff(ref.previewPS)
	return api.OK(), nil
}

// Stream subscribes the caller to the conversation bus for the given session
// and pumps SSE events until the connection closes. Marked raw: true in
// rest.yaml so codegen skips the api.Send envelope.
func (ctrl *ChatbotSession) Stream(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "sessionID")
	ref, err := ctrl.resolve(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if !ctrl.ac.CanViewSessionsOnChatbot(r.Context(), ref.chatbot) {
		http.Error(w, "chatbot session: forbidden", http.StatusForbidden)
		return
	}
	if ctrl.bus == nil {
		http.Error(w, "chatbot session: no bus", http.StatusServiceUnavailable)
		return
	}

	var convID uint64
	if ref.source == "widget" {
		convID, err = ctrl.sessionConvID(r.Context(), ref.widgetSession.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		convID = ref.previewPS.ConversationID
	}

	d := observability.NewSSEDispatcher(convID)
	ctrl.bus.Register(d)
	defer d.Close()
	observability.PumpSSE(w, r, d)
}

// ---- helpers ----------------------------------------------------------------

// resolve dispatches a string sessionID. Numeric IDs target widget DB
// sessions; non-numeric IDs target the in-memory preview map. The chatbot
// pointer is always populated for RBAC checks.
func (ctrl *ChatbotSession) resolve(ctx context.Context, id string) (*sessionRef, error) {
	if dbID, err := strconv.ParseUint(id, 10, 64); err == nil {
		s, err := ctrl.sessionSvc.FindByID(ctx, dbID)
		if err != nil {
			return nil, err
		}
		if s == nil {
			return nil, service.ChatbotSessionErrSessionNotFound()
		}
		cb, err := ctrl.chatbotSvc.FindByID(ctx, s.ChatbotID)
		if err != nil || cb == nil {
			return nil, errChatbotSessionChatbotMissing
		}
		return &sessionRef{source: "widget", widgetSession: s, chatbot: cb}, nil
	}
	ps := ctrl.previewSvc.Find(id)
	if ps == nil {
		return nil, service.ChatbotSessionErrSessionNotFound()
	}
	return &sessionRef{source: "preview", previewPS: ps, chatbot: &ps.Chatbot}, nil
}

// sessionConvID returns the conversation ID associated with the most recent
// step of the given widget session.
func (ctrl *ChatbotSession) sessionConvID(ctx context.Context, sessionID uint64) (uint64, error) {
	step, err := ctrl.sessionSvc.FindStepBySession(ctx, sessionID)
	if err != nil || step == nil {
		return 0, service.ChatbotSessionErrNoActiveStep()
	}
	return step.ConversationID, nil
}

func (ctrl *ChatbotSession) toUnifiedFromDB(
	ctx context.Context,
	s *types.ChatbotSession,
	cb *types.Chatbot,
	h *types.ChatbotSessionHandoff,
) *chatbotUnifiedSession {
	out := &chatbotUnifiedSession{
		ID:          strconv.FormatUint(s.ID, 10),
		Source:      "widget",
		ChatbotID:   s.ChatbotID,
		Status:      s.Status,
		CurrentStep: s.CurrentStep,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,

		CanView:          ctrl.ac.CanViewSessionsOnChatbot(ctx, cb),
		CanManage:        ctrl.ac.CanManageSessionsOnChatbot(ctx, cb),
		CanManageHandoff: ctrl.ac.CanManageSessionsHandoffOnChatbot(ctx, cb),
	}
	if h != nil {
		out.Handoff = &chatbotUnifiedHandoff{
			ID:        strconv.FormatUint(h.ID, 10),
			Status:    h.Status,
			StartedAt: h.InitiatedAt,
			ClosedAt:  h.ClosedAt,
		}
	}
	return out
}

func (ctrl *ChatbotSession) toUnifiedFromPreview(
	_ context.Context,
	ps *service.ChatbotPreviewSession,
) *chatbotUnifiedSession {
	out := &chatbotUnifiedSession{
		ID:             ps.ID,
		Source:         "preview",
		ChatbotID:      ps.Chatbot.ID,
		Status:         ps.Status,
		CurrentStep:    ps.CurrentStep,
		ConversationID: ps.ConversationID,
		CreatedAt:      ps.CreatedAt,

		// Preview is admin-gated at the route level; the chatbot snapshot is
		// the unsaved editor draft, so chatbot-resource RBAC checks never
		// match. Surface the route-level gate as "yes" on the wire so the
		// inbox doesn't grey out actions on legitimate preview rows.
		CanView:          true,
		CanManage:        true,
		CanManageHandoff: true,
	}
	if ps.Handoff != nil {
		out.Handoff = &chatbotUnifiedHandoff{
			ID:        ps.Handoff.ID,
			Status:    ps.Handoff.Status,
			StartedAt: ps.Handoff.Started,
			ClosedAt:  ps.Handoff.Closed,
		}
	}
	return out
}

var (
	errChatbotSessionForbidden      = errors.New("chatbot session: forbidden")
	errChatbotSessionChatbotMissing = errors.New("chatbot session: parent chatbot missing")
)
