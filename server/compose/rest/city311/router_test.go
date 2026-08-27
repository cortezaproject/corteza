package city311

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	city311Service "github.com/cortezaproject/corteza/server/compose/service/city311"
	contract "github.com/cortezaproject/corteza/server/compose/types/city311"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	id.Init(ctx)
	code := m.Run()
	cancel()
	os.Exit(code)
}

func testRouter(t *testing.T) (http.Handler, store.Storer, *city311Service.Service) {
	t.Helper()
	ctx := context.Background()
	dsn := fmt.Sprintf("sqlite3://file:%s?mode=memory&cache=shared", t.Name())
	st, err := sqlite.Connect(ctx, dsn)
	require.NoError(t, err)
	require.NoError(t, store.Upgrade(ctx, zap.NewNop(), st))
	svc := city311Service.New(st)
	require.NoError(t, svc.Seed(ctx, time.Date(2026, 2, 3, 15, 4, 5, 0, time.UTC)))
	router := chi.NewRouter()
	router.Route("/api/v1", MountRoutesWithService(svc))
	return router, st, svc
}

func executeJSON(t *testing.T, router http.Handler, method, path string, body any, headers map[string]string, userID uint64) *httptest.ResponseRecorder {
	t.Helper()
	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	request := httptest.NewRequest(method, path, bytes.NewReader(encoded))
	request.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	if userID != 0 {
		request = request.WithContext(auth.SetIdentityToContext(request.Context(), auth.Authenticated(userID)))
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func validPortalBody() map[string]any {
	return map[string]any{
		"summary": "Pothole on Example Street", "description": "A deep pothole blocks the eastbound traffic lane.",
		"service_type": "POTHOLE", "requester": map[string]any{"display_name": "Alex Resident", "email": "alex@example.invalid"},
		"location": map[string]any{"address": "100 Example Street, Buffalo, NY 14201", "latitude": 42.88645, "longitude": -78.87837},
	}
}

func TestPortalSubmissionUsesSinglePublishedSuccessStatus(t *testing.T) {
	router, _, _ := testRouter(t)
	headers := map[string]string{contract.IdempotencyHeader: "portal-replay-1"}
	first := executeJSON(t, router, http.MethodPost, "/api/v1/portal/service-requests", validPortalBody(), headers, 0)
	require.Equal(t, http.StatusCreated, first.Code, first.Body.String())
	replay := executeJSON(t, router, http.MethodPost, "/api/v1/portal/service-requests", validPortalBody(), headers, 0)
	require.Equal(t, http.StatusCreated, replay.Code, replay.Body.String())
	require.JSONEq(t, first.Body.String(), replay.Body.String())

	conflictBody := validPortalBody()
	conflictBody["summary"] = "A different pothole report"
	conflict := executeJSON(t, router, http.MethodPost, "/api/v1/portal/service-requests", conflictBody, headers, 0)
	require.Equal(t, http.StatusConflict, conflict.Code)
	require.Contains(t, conflict.Body.String(), string(contract.ErrorIdempotencyConflict))
}

func TestIntegrationSubmissionRequiresScopeAndReturnsReplayStatus(t *testing.T) {
	router, _, _ := testRouter(t)
	path := "/api/v1/service-requests"
	headers := map[string]string{contract.IdempotencyHeader: "integration-replay-1"}
	unauthenticated := executeJSON(t, router, http.MethodPost, path, validPortalBody(), headers, 0)
	require.Equal(t, http.StatusUnauthorized, unauthenticated.Code)

	executeWithScope := func(scope string) *httptest.ResponseRecorder {
		encoded, err := json.Marshal(validPortalBody())
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(encoded))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set(contract.IdempotencyHeader, headers[contract.IdempotencyHeader])
		token := jwt.New()
		require.NoError(t, token.Set("scope", scope))
		ctx := jwtauth.NewContext(request.Context(), token, nil)
		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(77))
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request.WithContext(ctx))
		return response
	}

	forbidden := executeWithScope("profile.read")
	require.Equal(t, http.StatusForbidden, forbidden.Code)
	created := executeWithScope(contract.ScopeRequestWrite)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	replayed := executeWithScope(contract.ScopeRequestWrite)
	require.Equal(t, http.StatusOK, replayed.Code, replayed.Body.String())
	require.JSONEq(t, created.Body.String(), replayed.Body.String())
}

func TestPortalSubmissionValidatesAttachmentTokenLimit(t *testing.T) {
	router, _, _ := testRouter(t)
	body := validPortalBody()
	body["attachment_tokens"] = []string{"1", "2", "3", "4", "5", "6"}
	response := executeJSON(t, router, http.MethodPost, "/api/v1/portal/service-requests", body, map[string]string{contract.IdempotencyHeader: "too-many-attachments"}, 0)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code)
	require.Contains(t, response.Body.String(), "/attachment_tokens")
}

func TestStaffRoutesEnforceIdentityScopeAndVersion(t *testing.T) {
	router, st, _ := testRouter(t)
	user, err := store.LookupUserByEmail(context.Background(), st, "service-agent@city311.example.invalid")
	require.NoError(t, err)
	request, err := store.LookupCity311ServiceRequestByRequestNumber(context.Background(), st, "SR-2026-00034")
	require.NoError(t, err)

	unauthenticated := executeJSON(t, router, http.MethodGet, "/api/v1/staff/service-requests", nil, nil, 0)
	require.Equal(t, http.StatusUnauthorized, unauthenticated.Code)

	queue := executeJSON(t, router, http.MethodGet, "/api/v1/staff/service-requests?status=SUBMITTED&page_size=10", nil, nil, user.ID)
	require.Equal(t, http.StatusOK, queue.Code, queue.Body.String())
	require.Contains(t, queue.Body.String(), "SR-2026-00034")
	require.NotContains(t, queue.Body.String(), "SR-2026-00035")

	path := fmt.Sprintf("/api/v1/staff/service-requests/%d/transitions", request.ID)
	transition := executeJSON(t, router, http.MethodPost, path, map[string]any{"to_status": "TRIAGED", "reason": "Validated and routed"}, map[string]string{contract.IfMatchHeader: `"1"`}, user.ID)
	require.Equal(t, http.StatusOK, transition.Code, transition.Body.String())
	require.Contains(t, transition.Body.String(), `"version":2`)

	stale := executeJSON(t, router, http.MethodPost, path, map[string]any{"to_status": "ASSIGNED"}, map[string]string{contract.IfMatchHeader: `"1"`}, user.ID)
	require.Equal(t, http.StatusConflict, stale.Code)
	require.Contains(t, stale.Body.String(), `"current_version":2`)

	missingVersion := executeJSON(t, router, http.MethodPost, path, map[string]any{"to_status": "ASSIGNED"}, nil, user.ID)
	require.Equal(t, http.StatusPreconditionRequired, missingVersion.Code)
}

func TestStaffCreateSupportsConstituentReferenceAndRejectsMixedUnion(t *testing.T) {
	router, st, _ := testRouter(t)
	user, err := store.LookupUserByEmail(context.Background(), st, "service-agent@city311.example.invalid")
	require.NoError(t, err)
	body := map[string]any{"request": validPortalBody(), "constituent": map[string]any{"constituent_id": "C-1"}}
	created := executeJSON(t, router, http.MethodPost, "/api/v1/staff/service-requests", body, nil, user.ID)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	require.Contains(t, created.Body.String(), "constituent1@city311.example.invalid")
	require.Contains(t, created.Body.String(), `"origin_class":"INTERNAL"`)

	body["constituent"] = map[string]any{"constituent_id": "C-1", "display_name": "Mixed", "email": "mixed@example.invalid"}
	mixed := executeJSON(t, router, http.MethodPost, "/api/v1/staff/service-requests", body, nil, user.ID)
	require.Equal(t, http.StatusUnprocessableEntity, mixed.Code)
}

func TestParseIfMatchAcceptsStrongAndWeakEntityTags(t *testing.T) {
	for _, value := range []string{`"12"`, `W/"12"`, "12"} {
		parsed, err := parseIfMatch(value)
		require.NoError(t, err)
		require.Equal(t, uint64(12), parsed)
	}
	_, err := parseIfMatch("")
	require.Error(t, err)
}
