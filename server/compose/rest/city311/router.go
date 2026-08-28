package city311

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	city311Service "github.com/cortezaproject/corteza/server/compose/service/city311"
	contract "github.com/cortezaproject/corteza/server/compose/types/city311"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

const (
	maximumJSONBody      = 72 << 20
	serviceRequestsRoute = "/service-requests"
	invalidFieldsMessage = "The request contains invalid fields."
)

var strongVersionPattern = regexp.MustCompile(`^"[1-9][0-9]*"$`)

type handler struct{ service *city311Service.Service }

type stringList []string

func (values *stringList) UnmarshalJSON(data []byte) error {
	var list []string
	if err := json.Unmarshal(data, &list); err == nil {
		*values = list
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*values = []string{value}
	return nil
}

type staffQueueFilters struct {
	Status         stringList `json:"status"`
	ServiceType    stringList `json:"service_type"`
	Department     stringList `json:"department"`
	District       stringList `json:"district"`
	OriginClass    stringList `json:"origin_class"`
	SourceChannel  stringList `json:"source_channel"`
	Assignee       stringList `json:"assignee"`
	Collaborator   stringList `json:"collaborator"`
	Category       stringList `json:"category"`
	CreatedFrom    string     `json:"created_from"`
	CreatedTo      string     `json:"created_to"`
	DuplicateGroup stringList `json:"duplicate_group"`
}

func MountRoutes() func(chi.Router) {
	return MountRoutesWithService(city311Service.Default)
}

func MountRoutesWithService(service *city311Service.Service) func(chi.Router) {
	return func(r chi.Router) {
		h := &handler{service: service}
		r.With(requireScope(contract.ScopeRequestWrite)).Post(serviceRequestsRoute, h.integrationSubmit)
		r.Post("/portal/service-requests", h.portalSubmit)
		r.Route("/staff", func(r chi.Router) {
			r.Use(requireIdentity)
			r.Post(serviceRequestsRoute, h.staffSubmit)
			r.Get(serviceRequestsRoute, h.staffList)
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
	if len(input.AttachmentTokens) > 0 {
		writeValidation(w, "/attachment_tokens", contract.ValidationInvalidValue)
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
	if len(input.Request.AttachmentTokens) > 0 {
		writeValidation(w, "/request/attachment_tokens", contract.ValidationInvalidValue)
		return
	}
	requester, existingConstituentID, err := staffRequester(input.Constituent)
	if err != nil {
		writeResult(w, 0, nil, err)
		return
	}
	response, _, err := h.service.Submit(r.Context(), contract.ServiceRequestCreate{
		Summary: input.Request.Summary, Description: input.Request.Description, ServiceType: input.Request.ServiceType,
		Requester: requester, Location: input.Request.Location, CustomFields: input.Request.CustomFields,
	}, "", city311Service.SubmissionOptions{
		Operation: "staff_service_request_create", SourceChannel: contract.SourceChannelStaffInPerson,
		ActorType: contract.AuditActorStaff, ActorID: actor.ID, StaffActor: &actor, ExistingConstituentID: existingConstituentID,
	})
	if err != nil {
		writeResult(w, 0, nil, err)
		return
	}
	requestID, _ := strconv.ParseUint(response.RequestID, 10, 64)
	detail, err := h.service.Find(r.Context(), actor, requestID)
	writeResult(w, http.StatusCreated, detail, err)
}

func staffRequester(value contract.StaffConstituentInput) (contract.RequesterInput, string, error) {
	hasReference := value.ConstituentID != nil
	hasDisplayName := value.DisplayName != nil
	hasEmail := value.Email != nil
	if hasReference && (hasDisplayName || hasEmail) {
		return contract.RequesterInput{}, "", &city311Service.ServiceError{Status: 422, Payload: contract.APIError{
			Error: contract.ErrorValidation, Message: invalidFieldsMessage, Retryable: false,
			Errors: []contract.FieldError{{Field: "/constituent", Code: contract.ValidationInvalidValue}},
		}}
	}
	if hasReference {
		constituentID := strings.TrimSpace(*value.ConstituentID)
		if constituentID == "" {
			return contract.RequesterInput{}, "", &city311Service.ServiceError{Status: 422, Payload: contract.APIError{
				Error: contract.ErrorValidation, Message: invalidFieldsMessage, Retryable: false,
				Errors: []contract.FieldError{{Field: "/constituent/constituent_id", Code: contract.ValidationRequired}},
			}}
		}
		return contract.RequesterInput{}, constituentID, nil
	}
	if !hasDisplayName || !hasEmail {
		return contract.RequesterInput{}, "", &city311Service.ServiceError{Status: 422, Payload: contract.APIError{
			Error: contract.ErrorValidation, Message: invalidFieldsMessage, Retryable: false,
			Errors: []contract.FieldError{{Field: "/constituent", Code: contract.ValidationRequired}},
		}}
	}
	return contract.RequesterInput{DisplayName: *value.DisplayName, Email: *value.Email}, "", nil
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
	requestFilter, err := parseStaffQueueFilters(r)
	if err != nil {
		writeResult(w, 0, nil, err)
		return
	}
	requestFilter.PageSize = uint(pageSize)
	requestFilter.PageToken = r.URL.Query().Get("page_token")
	requestFilter.Sort = r.URL.Query().Get("sort")
	result, err := h.service.List(r.Context(), actor, requestFilter)
	writeResult(w, http.StatusOK, result, err)
}

func parseStaffQueueFilters(r *http.Request) (city311Service.RequestFilter, error) {
	input := staffQueueFilters{}
	raw := r.URL.Query().Get("filters")
	if raw != "" {
		decoder := json.NewDecoder(strings.NewReader(raw))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&input); err != nil {
			return city311Service.RequestFilter{}, queueFilterError("filters", contract.ValidationInvalidFormat)
		}
		if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
			return city311Service.RequestFilter{}, queueFilterError("filters", contract.ValidationInvalidFormat)
		}
	} else {
		// Accept OpenAPI form-style exploded objects in addition to the JSON object
		// representation used by the frozen examples and frontend mocks.
		query := r.URL.Query()
		input.Status = queryValues(query, "status")
		input.ServiceType = queryValues(query, "service_type")
		input.Department = queryValues(query, "department")
		input.District = queryValues(query, "district")
		input.OriginClass = queryValues(query, "origin_class")
		input.SourceChannel = queryValues(query, "source_channel")
		input.Assignee = queryValues(query, "assignee")
		input.Collaborator = queryValues(query, "collaborator")
		input.Category = queryValues(query, "category")
		input.CreatedFrom = firstQueryValue(queryValues(query, "created_from"))
		input.CreatedTo = firstQueryValue(queryValues(query, "created_to"))
		input.DuplicateGroup = queryValues(query, "duplicate_group")
	}

	assignees, err := parseFilterIDs(input.Assignee, "assignee")
	if err != nil {
		return city311Service.RequestFilter{}, err
	}
	collaborators, err := parseFilterIDs(input.Collaborator, "collaborator")
	if err != nil {
		return city311Service.RequestFilter{}, err
	}
	createdFrom, err := parseFilterTime(input.CreatedFrom, "created_from")
	if err != nil {
		return city311Service.RequestFilter{}, err
	}
	createdTo, err := parseFilterTime(input.CreatedTo, "created_to")
	if err != nil {
		return city311Service.RequestFilter{}, err
	}

	return city311Service.RequestFilter{
		Statuses:           convertStrings(input.Status, contract.ServiceRequestStatus("")),
		ServiceTypes:       convertStrings(input.ServiceType, contract.ServiceType("")),
		OwningDepartments:  convertStrings(input.Department, contract.DepartmentCode("")),
		CouncilDistricts:   convertStrings(input.District, contract.DistrictCode("")),
		OriginClasses:      convertStrings(input.OriginClass, contract.OriginClass("")),
		SourceChannels:     convertStrings(input.SourceChannel, contract.SourceChannel("")),
		PrimaryAssigneeIDs: assignees,
		CollaboratorIDs:    collaborators,
		Categories:         convertStrings(input.Category, contract.ContactCategory("")),
		CreatedFrom:        createdFrom,
		CreatedTo:          createdTo,
		DuplicateGroups:    []string(input.DuplicateGroup),
	}, nil
}

func queryValues(query map[string][]string, field string) stringList {
	values := query[field]
	if len(values) == 0 {
		values = query["filters["+field+"]"]
	}
	out := make(stringList, 0, len(values))
	for _, value := range values {
		for _, item := range strings.Split(value, ",") {
			if item = strings.TrimSpace(item); item != "" {
				out = append(out, item)
			}
		}
	}
	return out
}

func firstQueryValue(values stringList) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func parseFilterIDs(values stringList, field string) ([]uint64, error) {
	out := make([]uint64, 0, len(values))
	for _, value := range values {
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil || parsed == 0 {
			return nil, queueFilterError(field, contract.ValidationInvalidFormat)
		}
		out = append(out, parsed)
	}
	return out, nil
}

func parseFilterTime(value, field string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, queueFilterError(field, contract.ValidationInvalidFormat)
	}
	parsed = parsed.UTC()
	return &parsed, nil
}

func queueFilterError(field string, code contract.ValidationCode) *city311Service.ServiceError {
	return &city311Service.ServiceError{Status: 422, Payload: contract.APIError{
		Error: contract.ErrorValidation, Message: invalidFieldsMessage, Retryable: false,
		Errors: []contract.FieldError{{Field: "/query/filters/" + field, Code: code}},
	}}
}

func convertStrings[T ~string](values stringList, _ T) []T {
	out := make([]T, 0, len(values))
	for _, value := range values {
		out = append(out, T(value))
	}
	return out
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
	if !strongVersionPattern.MatchString(value) {
		return 0, errors.New("If-Match must be one quoted positive decimal version")
	}
	return strconv.ParseUint(value[1:len(value)-1], 10, 64)
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
	writeJSON(w, 422, contract.APIError{Error: contract.ErrorValidation, Message: invalidFieldsMessage, Retryable: false, Errors: []contract.FieldError{{Field: field, Code: code}}})
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
