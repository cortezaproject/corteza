package city311

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	city311Service "github.com/cortezaproject/corteza/server/compose/service/city311"
	contract "github.com/cortezaproject/corteza/server/compose/types/city311"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

const maximumJSONBody = 72 << 20

type handler struct{ service *city311Service.Service }

func MountRoutes() func(chi.Router) {
	return MountRoutesWithService(city311Service.Default)
}

func MountRoutesWithService(service *city311Service.Service) func(chi.Router) {
	return func(r chi.Router) {
		h := &handler{service: service}
		r.With(requireScope(contract.ScopeRequestWrite)).Post("/service-requests", h.integrationSubmit)
		r.Post("/portal/service-requests", h.portalSubmit)
		r.Route("/staff", func(r chi.Router) {
			r.Use(requireIdentity)
			r.Post("/service-requests", h.staffSubmit)
			r.Get("/service-requests", h.staffList)
			r.Get("/service-requests/{request_id}", h.staffDetail)
			r.Post("/service-requests/{request_id}/transitions", h.staffTransition)
		})
	}
}

func (h *handler) integrationSubmit(w http.ResponseWriter, r *http.Request) {
	input := contract.ServiceRequestCreate{}
	if !decodeJSON(w, r, &input) {
		return
	}
	response, status, err := h.service.Submit(r.Context(), input, r.Header.Get(contract.IdempotencyHeader), city311Service.SubmissionOptions{
		Operation: "service_request_create", SourceChannel: contract.SourceChannelAPI,
		ActorType: contract.AuditActorIntegrationClient, ActorID: auth.GetIdentityFromContext(r.Context()).Identity(), RequireIdempotency: true,
	})
	writeResult(w, status, response, err)
}

func (h *handler) portalSubmit(w http.ResponseWriter, r *http.Request) {
	input := contract.PortalServiceRequestSubmit{}
	if !decodeJSON(w, r, &input) {
		return
	}
	if len(input.AttachmentTokens) > 5 {
		writeValidation(w, "/attachment_tokens", contract.ValidationTooManyItems)
		return
	}
	identity := auth.GetIdentityFromContext(r.Context())
	source := contract.SourceChannelPortalAnonymous
	actorID := uint64(0)
	if identity.Valid() {
		source = contract.SourceChannelPortalAuthenticated
		actorID = identity.Identity()
	}
	response, status, err := h.service.Submit(r.Context(), contract.ServiceRequestCreate{
		Summary: input.Summary, Description: input.Description, ServiceType: input.ServiceType,
		Requester: input.Requester, Location: input.Location, CustomFields: input.CustomFields,
	}, r.Header.Get(contract.IdempotencyHeader), city311Service.SubmissionOptions{
		Operation: "portal_service_request_submit", SourceChannel: source,
		ActorType: contract.AuditActorConstituent, ActorID: actorID, RequireIdempotency: true,
	})
	// This endpoint publishes one success status. Replays return the original
	// representation with 201 while the integration endpoint distinguishes 200.
	if status == http.StatusOK {
		status = http.StatusCreated
	}
	writeResult(w, status, response, err)
}

func (h *handler) staffSubmit(w http.ResponseWriter, r *http.Request) {
	input := contract.StaffServiceRequestCreate{}
	if !decodeJSON(w, r, &input) {
		return
	}
	identity := auth.GetIdentityFromContext(r.Context())
	actor, err := h.service.FindActor(r.Context(), identity.Identity())
	if err != nil {
		writeResult(w, 0, nil, err)
		return
	}
	if !actorCanCreate(actor) {
		writeResult(w, 0, nil, &city311Service.ServiceError{Status: 403, Payload: contract.APIError{Error: contract.ErrorForbidden, Message: "The actor cannot create staff service requests.", Retryable: false}})
		return
	}
	if len(input.Request.AttachmentTokens) > 5 {
		writeValidation(w, "/request/attachment_tokens", contract.ValidationTooManyItems)
		return
	}
	requester, err := h.staffRequester(r, input.Constituent)
	if err != nil {
		writeResult(w, 0, nil, err)
		return
	}
	response, _, err := h.service.Submit(r.Context(), contract.ServiceRequestCreate{
		Summary: input.Request.Summary, Description: input.Request.Description, ServiceType: input.Request.ServiceType,
		Requester: requester, Location: input.Request.Location, CustomFields: input.Request.CustomFields,
	}, "", city311Service.SubmissionOptions{
		Operation: "staff_service_request_create", SourceChannel: contract.SourceChannelStaffInPerson,
		ActorType: contract.AuditActorStaff, ActorID: actor.ID,
	})
	if err != nil {
		writeResult(w, 0, nil, err)
		return
	}
	requestID, _ := strconv.ParseUint(response.RequestID, 10, 64)
	detail, err := h.service.Find(r.Context(), actor, requestID)
	writeResult(w, http.StatusCreated, detail, err)
}

func (h *handler) staffRequester(r *http.Request, value contract.StaffConstituentInput) (contract.RequesterInput, error) {
	hasReference := value.ConstituentID != nil
	hasDisplayName := value.DisplayName != nil
	hasEmail := value.Email != nil
	if hasReference && (hasDisplayName || hasEmail) {
		return contract.RequesterInput{}, &city311Service.ServiceError{Status: 422, Payload: contract.APIError{
			Error: contract.ErrorValidation, Message: "The request contains invalid fields.", Retryable: false,
			Errors: []contract.FieldError{{Field: "/constituent", Code: contract.ValidationInvalidValue}},
		}}
	}
	if hasReference {
		return h.service.ResolveConstituent(r.Context(), *value.ConstituentID)
	}
	if !hasDisplayName || !hasEmail {
		return contract.RequesterInput{}, &city311Service.ServiceError{Status: 422, Payload: contract.APIError{
			Error: contract.ErrorValidation, Message: "The request contains invalid fields.", Retryable: false,
			Errors: []contract.FieldError{{Field: "/constituent", Code: contract.ValidationRequired}},
		}}
	}
	return contract.RequesterInput{DisplayName: *value.DisplayName, Email: *value.Email}, nil
}

func (h *handler) staffList(w http.ResponseWriter, r *http.Request) {
	actor, err := h.service.FindActor(r.Context(), auth.GetIdentityFromContext(r.Context()).Identity())
	if err != nil {
		writeResult(w, 0, nil, err)
		return
	}
	pageSize := uint64(50)
	if raw := r.URL.Query().Get("page_size"); raw != "" {
		pageSize, err = strconv.ParseUint(raw, 10, 16)
		if err != nil || pageSize == 0 {
			writeValidation(w, "/query/page_size", contract.ValidationInvalidFormat)
			return
		}
	}
	assigneeID := uint64(0)
	if raw := r.URL.Query().Get("assignee"); raw != "" {
		assigneeID, err = strconv.ParseUint(raw, 10, 64)
		if err != nil || assigneeID == 0 {
			writeValidation(w, "/query/assignee", contract.ValidationInvalidFormat)
			return
		}
	}
	result, err := h.service.List(r.Context(), actor, city311Service.RequestFilter{
		Status:            contract.ServiceRequestStatus(r.URL.Query().Get("status")),
		ServiceType:       contract.ServiceType(r.URL.Query().Get("service_type")),
		OwningDepartment:  contract.DepartmentCode(r.URL.Query().Get("department")),
		CouncilDistrict:   contract.DistrictCode(r.URL.Query().Get("district")),
		OriginClass:       contract.OriginClass(r.URL.Query().Get("origin_class")),
		SourceChannel:     contract.SourceChannel(r.URL.Query().Get("source_channel")),
		PrimaryAssigneeID: assigneeID,
		PageSize:          uint(pageSize), PageToken: r.URL.Query().Get("page_token"), Sort: r.URL.Query().Get("sort"),
	})
	writeResult(w, http.StatusOK, result, err)
}

func (h *handler) staffDetail(w http.ResponseWriter, r *http.Request) {
	requestID, err := strconv.ParseUint(chi.URLParam(r, "request_id"), 10, 64)
	if err != nil || requestID == 0 {
		writeValidation(w, "/path/request_id", contract.ValidationInvalidFormat)
		return
	}
	actor, err := h.service.FindActor(r.Context(), auth.GetIdentityFromContext(r.Context()).Identity())
	if err != nil {
		writeResult(w, 0, nil, err)
		return
	}
	result, err := h.service.Find(r.Context(), actor, requestID)
	writeResult(w, http.StatusOK, result, err)
}

func (h *handler) staffTransition(w http.ResponseWriter, r *http.Request) {
	requestID, err := strconv.ParseUint(chi.URLParam(r, "request_id"), 10, 64)
	if err != nil || requestID == 0 {
		writeValidation(w, "/path/request_id", contract.ValidationInvalidFormat)
		return
	}
	expectedVersion, err := parseIfMatch(r.Header.Get(contract.IfMatchHeader))
	if err != nil {
		writeResult(w, 0, nil, &city311Service.ServiceError{Status: 428, Payload: contract.APIError{Error: contract.ErrorExpectedVersionRequired, Message: "If-Match must identify the expected record version.", Retryable: false}})
		return
	}
	input := contract.RequestTransition{}
	if !decodeJSON(w, r, &input) {
		return
	}
	actor, err := h.service.FindActor(r.Context(), auth.GetIdentityFromContext(r.Context()).Identity())
	if err != nil {
		writeResult(w, 0, nil, err)
		return
	}
	result, err := h.service.Transition(r.Context(), actor, requestID, expectedVersion, input)
	writeResult(w, http.StatusOK, result, err)
}

func parseIfMatch(value string) (uint64, error) {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "W/")
	value = strings.Trim(value, "\"")
	if value == "" {
		return 0, errors.New("missing If-Match")
	}
	return strconv.ParseUint(value, 10, 64)
}

func requireIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.GetIdentityFromContext(r.Context()).Valid() {
			writeJSON(w, http.StatusUnauthorized, contract.APIError{Error: contract.ErrorUnauthenticated, Message: "Authentication is required.", Retryable: false})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func requireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())
			if err != nil || token == nil || !auth.GetIdentityFromContext(r.Context()).Valid() {
				writeJSON(w, http.StatusUnauthorized, contract.APIError{Error: contract.ErrorUnauthenticated, Message: "A valid access token is required.", Retryable: false})
				return
			}
			if !auth.CheckJwtScope(token, scope) {
				writeJSON(w, http.StatusForbidden, contract.APIError{Error: contract.ErrorForbidden, Message: "The access token does not grant the required scope.", Retryable: false})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maximumJSONBody))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeValidation(w, "/", contract.ValidationInvalidFormat)
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeValidation(w, "/", contract.ValidationInvalidFormat)
		return false
	}
	return true
}

func writeValidation(w http.ResponseWriter, field string, code contract.ValidationCode) {
	writeJSON(w, 422, contract.APIError{Error: contract.ErrorValidation, Message: "The request contains invalid fields.", Retryable: false, Errors: []contract.FieldError{{Field: field, Code: code}}})
}

func writeResult(w http.ResponseWriter, status int, value any, err error) {
	if err != nil {
		var serviceErr *city311Service.ServiceError
		if errors.As(err, &serviceErr) {
			writeJSON(w, serviceErr.Status, serviceErr.Payload)
			return
		}
		writeJSON(w, http.StatusInternalServerError, contract.APIError{Error: contract.ErrorOperationFailed, Message: "The operation could not be completed.", Retryable: false})
		return
	}
	writeJSON(w, status, value)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func actorCanCreate(actor contract.Actor) bool {
	for _, role := range actor.Roles {
		switch role {
		case contract.ApplicationRoleServiceAgent, contract.ApplicationRoleSupervisor, contract.ApplicationRoleDepartmentManager, contract.ApplicationRolePlatformAdministrator:
			return true
		}
	}
	return false
}
