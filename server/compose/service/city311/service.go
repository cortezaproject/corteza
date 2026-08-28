package city311

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	composeTypes "github.com/cortezaproject/corteza/server/compose/types"
	contract "github.com/cortezaproject/corteza/server/compose/types/city311"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
)

const idempotencyLifetime = 24 * time.Hour

const maximumAttachmentSize = 10 << 20

var e164Pattern = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)

var allowedAttachmentMediaTypes = map[string]bool{
	"image/jpeg": true, "image/png": true, "application/pdf": true, "text/plain": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

type (
	Service struct {
		store  store.Storer
		now    func() time.Time
		nextID func() uint64
		mu     sync.Mutex
	}

	SubmissionOptions struct {
		Operation             string
		SourceChannel         contract.SourceChannel
		ActorType             contract.AuditActorType
		ActorID               uint64
		StaffActor            *contract.Actor
		ExistingConstituentID string
		RequireIdempotency    bool
	}

	RequestFilter struct {
		Statuses           []contract.ServiceRequestStatus
		ServiceTypes       []contract.ServiceType
		OwningDepartments  []contract.DepartmentCode
		CouncilDistricts   []contract.DistrictCode
		OriginClasses      []contract.OriginClass
		SourceChannels     []contract.SourceChannel
		PrimaryAssigneeIDs []uint64
		CollaboratorIDs    []uint64
		Categories         []contract.ContactCategory
		CreatedFrom        *time.Time
		CreatedTo          *time.Time
		DuplicateGroups    []string
		PageSize           uint
		PageToken          string
		Sort               string
	}

	ServiceError struct {
		Status  int
		Payload contract.APIError
	}

	validatedAttachment struct {
		Filename  string
		MediaType string
		Content   []byte
	}
)

var Default *Service

func (e *ServiceError) Error() string { return e.Payload.Message }

func New(s store.Storer) *Service {
	return &Service{store: s, now: func() time.Time { return time.Now().UTC().Round(time.Second) }, nextID: id.Next}
}

func Initialize(ctx context.Context, s store.Storer) error {
	Default = New(s)
	benchmarkNow := Default.now()
	if configured := strings.TrimSpace(os.Getenv("BENCHMARK_NOW")); configured != "" {
		parsed, err := time.Parse(time.RFC3339, configured)
		if err != nil {
			return fmt.Errorf("invalid BENCHMARK_NOW: %w", err)
		}
		benchmarkNow = parsed.UTC()
		Default.now = func() time.Time { return benchmarkNow }
	}
	return Default.Seed(ctx, benchmarkNow)
}

func validationError(fields ...contract.FieldError) *ServiceError {
	return &ServiceError{Status: 422, Payload: contract.APIError{
		Error:     contract.ErrorValidation,
		Message:   "The request contains invalid fields.",
		Retryable: false,
		Errors:    fields,
	}}
}

func apiError(status int, code contract.ErrorCode, message string) *ServiceError {
	return &ServiceError{Status: status, Payload: contract.APIError{Error: code, Message: message, Retryable: false}}
}

func validateWrite(in contract.ServiceRequestCreate) *ServiceError {
	var fields []contract.FieldError
	trimmedSummary := strings.TrimSpace(in.Summary)
	summaryLength := utf8.RuneCountInString(trimmedSummary)
	if summaryLength < 5 {
		fields = append(fields, contract.FieldError{Field: "/summary", Code: contract.ValidationTooShort})
	} else if summaryLength > 160 {
		fields = append(fields, contract.FieldError{Field: "/summary", Code: contract.ValidationTooLong})
	}
	trimmedDescription := strings.TrimSpace(in.Description)
	descriptionLength := utf8.RuneCountInString(trimmedDescription)
	if descriptionLength < 10 {
		fields = append(fields, contract.FieldError{Field: "/description", Code: contract.ValidationTooShort})
	} else if descriptionLength > 5000 {
		fields = append(fields, contract.FieldError{Field: "/description", Code: contract.ValidationTooLong})
	}
	if _, ok := departmentForServiceType(in.ServiceType); !ok {
		fields = append(fields, contract.FieldError{Field: "/service_type", Code: contract.ValidationInvalidValue})
	}
	displayName := strings.TrimSpace(in.Requester.DisplayName)
	if displayName == "" {
		fields = append(fields, contract.FieldError{Field: "/requester/display_name", Code: contract.ValidationRequired})
	} else if utf8.RuneCountInString(displayName) > 120 {
		fields = append(fields, contract.FieldError{Field: "/requester/display_name", Code: contract.ValidationTooLong})
	}
	email := strings.TrimSpace(in.Requester.Email)
	parsedEmail, emailErr := mail.ParseAddress(email)
	if emailErr != nil || parsedEmail.Address != email || parsedEmail.Name != "" {
		fields = append(fields, contract.FieldError{Field: "/requester/email", Code: contract.ValidationInvalidFormat})
	}
	if phone := strings.TrimSpace(in.Requester.Phone); phone != "" && !e164Pattern.MatchString(phone) {
		fields = append(fields, contract.FieldError{Field: "/requester/phone", Code: contract.ValidationInvalidFormat})
	}
	if serviceTypeRequiresLocation(in.ServiceType) && (in.Location == nil || strings.TrimSpace(in.Location.Address) == "") {
		fields = append(fields, contract.FieldError{Field: "/location/address", Code: contract.ValidationLocationRequired})
	} else if serviceTypeRequiresLocation(in.ServiceType) && (in.Location.Latitude == nil || in.Location.Longitude == nil) {
		fields = append(fields, contract.FieldError{Field: "/location", Code: contract.ValidationCoordinatesRequired})
	}
	if in.Location != nil {
		if in.Location.Latitude != nil && (*in.Location.Latitude < -90 || *in.Location.Latitude > 90) {
			fields = append(fields, contract.FieldError{Field: "/location/latitude", Code: contract.ValidationOutOfRange})
		}
		if in.Location.Longitude != nil && (*in.Location.Longitude < -180 || *in.Location.Longitude > 180) {
			fields = append(fields, contract.FieldError{Field: "/location/longitude", Code: contract.ValidationOutOfRange})
		}
	}
	if len(in.Attachments) > 5 {
		fields = append(fields, contract.FieldError{Field: "/attachments", Code: contract.ValidationTooManyItems})
	}
	if len(fields) > 0 {
		return validationError(fields...)
	}
	return nil
}

func validateInlineAttachments(input []contract.AttachmentInput) ([]validatedAttachment, *ServiceError) {
	out := make([]validatedAttachment, 0, len(input))
	var fields []contract.FieldError
	for index, item := range input {
		prefix := fmt.Sprintf("/attachments/%d", index)
		filename := path.Base(strings.ReplaceAll(strings.TrimSpace(item.Filename), `\`, "/"))
		filenameLength := utf8.RuneCountInString(filename)
		if filename == "." || filename == "/" || filenameLength == 0 {
			fields = append(fields, contract.FieldError{Field: prefix + "/filename", Code: contract.ValidationRequired})
		} else if filenameLength > 120 {
			fields = append(fields, contract.FieldError{Field: prefix + "/filename", Code: contract.ValidationTooLong})
		}
		mediaType := strings.TrimSpace(item.MediaType)
		if !allowedAttachmentMediaTypes[mediaType] {
			fields = append(fields, contract.FieldError{Field: prefix + "/media_type", Code: contract.ValidationInvalidValue})
		}
		content, err := base64.StdEncoding.DecodeString(item.ContentBase64)
		if err != nil {
			fields = append(fields, contract.FieldError{Field: prefix + "/content_base64", Code: contract.ValidationInvalidFormat})
		} else if len(content) > maximumAttachmentSize {
			fields = append(fields, contract.FieldError{Field: prefix + "/content_base64", Code: contract.ValidationOutOfRange})
		}
		out = append(out, validatedAttachment{Filename: filename, MediaType: mediaType, Content: content})
	}
	if len(fields) > 0 {
		return nil, validationError(fields...)
	}
	return out, nil
}

func departmentForServiceType(serviceType contract.ServiceType) (contract.DepartmentCode, bool) {
	switch serviceType {
	case contract.ServiceTypeTreeMaintenance:
		return contract.DepartmentPublicWorks, true
	case contract.ServiceTypePothole:
		return contract.DepartmentStreets, true
	case contract.ServiceTypeMissedTrash:
		return contract.DepartmentSanitation, true
	case contract.ServiceTypeGeneralInquiry:
		return contract.DepartmentGeneralServices, true
	default:
		return "", false
	}
}

func serviceTypeRequiresLocation(serviceType contract.ServiceType) bool {
	return serviceType != contract.ServiceTypeGeneralInquiry
}

func hashJSON(value any) (string, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}

func hashKey(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func (svc *Service) Submit(ctx context.Context, in contract.ServiceRequestCreate, idempotencyKey string, options SubmissionOptions) (*contract.ServiceRequestResponse, int, error) {
	var referencedConstituent *composeTypes.City311Constituent
	if options.ExistingConstituentID != "" {
		options.ExistingConstituentID = strings.TrimSpace(options.ExistingConstituentID)
		if options.ExistingConstituentID == "" {
			return nil, 0, validationError(contract.FieldError{Field: "/constituent/constituent_id", Code: contract.ValidationRequired})
		}
		if options.StaffActor == nil {
			return nil, 0, apiError(403, contract.ErrorForbidden, "A staff actor is required to reference an existing constituent.")
		}
		var err error
		referencedConstituent, err = svc.ResolveConstituent(ctx, *options.StaffActor, options.ExistingConstituentID)
		if err != nil {
			return nil, 0, err
		}
		in.Requester = requesterInput(referencedConstituent.Profile)
	}
	if err := validateWrite(in); err != nil {
		return nil, 0, err
	}
	if options.StaffActor != nil {
		department, _ := departmentForServiceType(in.ServiceType)
		candidate := &composeTypes.City311ServiceRequest{OwningDepartment: department}
		if !canOperateRequest(*options.StaffActor) || !canRead(*options.StaffActor, candidate) {
			return nil, 0, apiError(403, contract.ErrorForbidden, "The service request is outside the actor's record scope.")
		}
	}
	attachments, attachmentErr := validateInlineAttachments(in.Attachments)
	if attachmentErr != nil {
		return nil, 0, attachmentErr
	}
	if options.RequireIdempotency && idempotencyKey == "" {
		return nil, 0, validationError(contract.FieldError{Field: "/headers/Idempotency-Key", Code: contract.ValidationRequired})
	}
	if len(idempotencyKey) > 255 {
		return nil, 0, validationError(contract.FieldError{Field: "/headers/Idempotency-Key", Code: contract.ValidationTooLong})
	}
	if options.Operation == "" {
		options.Operation = "service_request_create"
	}
	requestHash, err := hashJSON(struct {
		Request               contract.ServiceRequestCreate `json:"request"`
		ExistingConstituentID string                        `json:"existing_constituent_id,omitempty"`
	}{Request: in, ExistingConstituentID: options.ExistingConstituentID})
	if err != nil {
		return nil, 0, err
	}

	// The benchmark runtime has one application service. The lock complements the
	// database transaction so sequence allocation and same-key submissions cannot race.
	svc.mu.Lock()
	defer svc.mu.Unlock()

	var response *contract.ServiceRequestResponse
	status := 201
	err = store.Tx(ctx, svc.store, func(ctx context.Context, tx store.Storer) error {
		now := svc.now()
		keyHash := ""
		if idempotencyKey != "" {
			keyHash = hashKey(idempotencyKey)
			existing, lookupErr := store.LookupCity311IdempotencyRecordByOperationKeyHash(ctx, tx, options.Operation, keyHash)
			if lookupErr == nil {
				if !existing.ExpiresAt.After(now) {
					if deleteErr := store.DeleteCity311IdempotencyRecord(ctx, tx, existing); deleteErr != nil {
						return deleteErr
					}
				} else {
					if existing.RequestHash != requestHash {
						return apiError(409, contract.ErrorIdempotencyConflict, "The idempotency key was already used with a different request.")
					}
					encoded, marshalErr := json.Marshal(existing.ResponseBody)
					if marshalErr != nil {
						return marshalErr
					}
					response = &contract.ServiceRequestResponse{}
					if unmarshalErr := json.Unmarshal(encoded, response); unmarshalErr != nil {
						return unmarshalErr
					}
					status = 200
					return nil
				}
			}
			if lookupErr != nil && !errors.IsNotFound(lookupErr) {
				return lookupErr
			}
		}

		if options.ExistingConstituentID != "" {
			current, lookupErr := store.LookupCity311ConstituentByConstituentID(ctx, tx, options.ExistingConstituentID)
			if errors.IsNotFound(lookupErr) {
				return apiError(404, contract.ErrorNotFound, "The constituent was not found.")
			}
			if lookupErr != nil {
				return lookupErr
			}
			if !canReadConstituent(*options.StaffActor, current) {
				return apiError(403, contract.ErrorForbidden, "The constituent is outside the actor's record scope.")
			}
			referencedConstituent = current
		}

		year := uint64(now.Year())
		sequence, sequenceErr := store.LookupCity311RequestSequenceByID(ctx, tx, year)
		if errors.IsNotFound(sequenceErr) {
			start := uint64(1)
			if year == 2026 {
				start = 41
			}
			sequence = &composeTypes.City311RequestSequence{ID: year, NextNumber: start}
			if sequenceErr = store.CreateCity311RequestSequence(ctx, tx, sequence); sequenceErr != nil {
				return sequenceErr
			}
		} else if sequenceErr != nil {
			return sequenceErr
		}
		number := sequence.NextNumber
		sequence.NextNumber++
		if sequenceErr = store.UpdateCity311RequestSequence(ctx, tx, sequence); sequenceErr != nil {
			return sequenceErr
		}

		requestID := svc.nextID()
		department, _ := departmentForServiceType(in.ServiceType)
		originClass := contract.OriginClassExternal
		if options.SourceChannel == contract.SourceChannelStaffInPerson {
			originClass = contract.OriginClassInternal
		}
		constituentProfile := map[string]any(nil)
		if referencedConstituent != nil {
			constituentProfile = cloneMap(referencedConstituent.Profile)
			constituentProfile["constituent_id"] = referencedConstituent.ConstituentID
		} else {
			constituentProfile = requesterMap(requestID, in.Requester)
			constituent := &composeTypes.City311Constituent{
				ID: svc.nextID(), ConstituentID: fmt.Sprint(constituentProfile["constituent_id"]), Profile: cloneMap(constituentProfile),
				OwningDepartment: department, CreatedAt: now, UpdatedAt: now,
			}
			if err := store.CreateCity311Constituent(ctx, tx, constituent); err != nil {
				return err
			}
		}
		stored := &composeTypes.City311ServiceRequest{
			ID:               requestID,
			RequestNumber:    fmt.Sprintf("SR-%04d-%05d", year, number),
			Summary:          strings.TrimSpace(in.Summary),
			Description:      strings.TrimSpace(in.Description),
			ServiceType:      in.ServiceType,
			OwningDepartment: department,
			SourceChannel:    options.SourceChannel,
			OriginClass:      originClass,
			Status:           contract.ServiceRequestStatusSubmitted,
			PrimaryRequester: constituentProfile,
			Location:         locationMap(in.Location),
			CustomFields:     cloneMap(in.CustomFields),
			CollaboratorIDs:  composeTypes.City311Uint64Set{},
			Version:          1,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if err := store.CreateCity311ServiceRequest(ctx, tx, stored); err != nil {
			return err
		}
		for _, item := range attachments {
			if err := store.CreateCity311RequestAttachment(ctx, tx, &composeTypes.City311RequestAttachment{
				ID: svc.nextID(), RequestID: requestID, Filename: item.Filename, MediaType: item.MediaType,
				Size: uint64(len(item.Content)), Content: item.Content, CreatedAt: now,
			}); err != nil {
				return err
			}
		}
		audit := &composeTypes.City311AuditEvent{
			ID: svc.nextID(), RequestID: requestID, EventType: "SERVICE_REQUEST_SUBMITTED",
			ActorType: options.ActorType, ActorID: options.ActorID, SourceChannel: options.SourceChannel,
			Before: map[string]any{}, After: requestSnapshot(stored), CreatedAt: now,
		}
		if err := store.CreateCity311AuditEvent(ctx, tx, audit); err != nil {
			return err
		}
		if err := store.CreateCity311PublicHistoryItem(ctx, tx, &composeTypes.City311PublicHistoryItem{
			ID: svc.nextID(), RequestID: requestID, Action: string(stored.Status),
			ResponsibleDepartment: stored.OwningDepartment, OccurredAt: now,
		}); err != nil {
			return err
		}

		response = responseFor(stored)
		body, marshalErr := mapFrom(response)
		if marshalErr != nil {
			return marshalErr
		}
		if idempotencyKey == "" {
			return nil
		}
		return store.CreateCity311IdempotencyRecord(ctx, tx, &composeTypes.City311IdempotencyRecord{
			ID: svc.nextID(), Operation: options.Operation, KeyHash: keyHash, RequestHash: requestHash,
			ResponseStatus: 201, ResponseBody: body, RequestID: requestID,
			CreatedAt: now, ExpiresAt: now.Add(idempotencyLifetime),
		})
	})
	return response, status, err
}

func (svc *Service) FindActor(ctx context.Context, userID uint64) (contract.Actor, error) {
	profile, err := store.LookupCity311ActorProfileByID(ctx, svc.store, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return contract.Actor{}, apiError(403, contract.ErrorForbidden, "The authenticated actor has no City 311 role.")
		}
		return contract.Actor{}, err
	}
	return contract.Actor{ID: userID, Roles: []contract.ApplicationRole(profile.ApplicationRoles), Department: profile.Department, Districts: []contract.DistrictCode(profile.Districts)}, nil
}

func (svc *Service) Find(ctx context.Context, actor contract.Actor, requestID uint64) (*contract.StaffServiceRequestDetail, error) {
	stored, err := store.LookupCity311ServiceRequestByID(ctx, svc.store, requestID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, apiError(404, contract.ErrorNotFound, "The service request was not found.")
		}
		return nil, err
	}
	if !canRead(actor, stored) {
		return nil, apiError(403, contract.ErrorForbidden, "The service request is outside the actor's record scope.")
	}
	return svc.detail(ctx, actor, stored)
}

func (svc *Service) Transition(ctx context.Context, actor contract.Actor, requestID uint64, expectedVersion uint64, input contract.RequestTransition) (*contract.StaffServiceRequestDetail, error) {
	if expectedVersion == 0 {
		return nil, apiError(428, contract.ErrorExpectedVersionRequired, "If-Match must identify the expected record version.")
	}
	if strings.TrimSpace(string(input.ToStatus)) == "" {
		return nil, validationError(contract.FieldError{Field: "/to_status", Code: contract.ValidationRequired})
	}

	svc.mu.Lock()
	defer svc.mu.Unlock()
	err := store.Tx(ctx, svc.store, func(ctx context.Context, tx store.Storer) error {
		request, err := store.LookupCity311ServiceRequestByID(ctx, tx, requestID)
		if errors.IsNotFound(err) {
			return apiError(404, contract.ErrorNotFound, "The service request was not found.")
		}
		if err != nil {
			return err
		}
		if !canRead(actor, request) {
			return apiError(403, contract.ErrorForbidden, "The service request is outside the actor's record scope.")
		}
		if !canOperateRequest(actor) {
			return apiError(403, contract.ErrorForbidden, "The actor cannot transition service requests.")
		}
		if uint64(request.Version) != expectedVersion {
			current := uint64(request.Version)
			return &ServiceError{Status: 409, Payload: contract.APIError{
				Error: contract.ErrorVersionConflict, Message: "The service request was updated by another operation.",
				Retryable: false, CurrentVersion: &current,
			}}
		}
		if !transitionAllowed(request.Status, input.ToStatus) {
			return validationError(contract.FieldError{Field: "/to_status", Code: contract.ValidationInvalidValue})
		}
		before := requestSnapshot(request)
		request.Status = input.ToStatus
		request.Version++
		request.UpdatedAt = svc.now()
		if err = store.UpdateCity311ServiceRequest(ctx, tx, request); err != nil {
			return err
		}
		if err = store.CreateCity311AuditEvent(ctx, tx, &composeTypes.City311AuditEvent{
			ID: svc.nextID(), RequestID: request.ID, EventType: "SERVICE_REQUEST_TRANSITIONED",
			ActorType: contract.AuditActorStaff, ActorID: actor.ID, SourceChannel: contract.SourceChannelStaffInPerson,
			Before: before, After: requestSnapshot(request), CreatedAt: request.UpdatedAt,
		}); err != nil {
			return err
		}
		return store.CreateCity311PublicHistoryItem(ctx, tx, &composeTypes.City311PublicHistoryItem{
			ID: svc.nextID(), RequestID: request.ID, Action: string(request.Status),
			ResponsibleDepartment: request.OwningDepartment, OccurredAt: request.UpdatedAt,
		})
	})
	if err != nil {
		return nil, err
	}
	return svc.Find(ctx, actor, requestID)
}

func transitionAllowed(from, to contract.ServiceRequestStatus) bool {
	for _, candidate := range contract.ServiceRequestTransitions[from] {
		if candidate == to {
			return true
		}
	}
	return false
}

func (svc *Service) List(ctx context.Context, actor contract.Actor, requested RequestFilter) (*contract.ListResponse, error) {
	if !isStaff(actor) {
		return nil, apiError(403, contract.ErrorForbidden, "A staff role is required.")
	}
	if fieldErrors := validateRequestFilter(requested); len(fieldErrors) > 0 {
		return nil, validationError(fieldErrors...)
	}
	if requested.PageSize == 0 {
		requested.PageSize = 50
	}
	if requested.PageSize > 100 {
		return nil, validationError(contract.FieldError{Field: "/query/page_size", Code: contract.ValidationOutOfRange})
	}
	publishedSort, err := normalizeSort(requested.Sort)
	if err != nil {
		return nil, validationError(contract.FieldError{Field: "/query/sort", Code: contract.ValidationInvalidFormat})
	}
	offset, err := decodePageToken(requested.PageToken, publishedSort)
	if err != nil {
		return nil, validationError(contract.FieldError{Field: "/query/page_token", Code: contract.ValidationInvalidFormat})
	}
	f := composeTypes.City311ServiceRequestFilter{
		Check: func(item *composeTypes.City311ServiceRequest) (bool, error) {
			return canRead(actor, item) && matchesRequestFilter(item, requested), nil
		},
	}
	matching, _, err := store.SearchCity311ServiceRequests(ctx, svc.store, f)
	if err != nil {
		return nil, err
	}
	sortServiceRequests(matching, publishedSort)
	if offset > len(matching) {
		return nil, validationError(contract.FieldError{Field: "/query/page_token", Code: contract.ValidationInvalidFormat})
	}
	end := offset + int(requested.PageSize)
	if end > len(matching) {
		end = len(matching)
	}
	set := matching[offset:end]
	response := &contract.ListResponse{
		Items: make([]contract.RequestQueueItem, 0, len(set)), TotalCount: len(matching),
		AppliedFilters: appliedFilters(requested), Sort: publishedSort,
	}
	for _, item := range set {
		response.Items = append(response.Items, queueItem(actor, item))
	}
	if end < len(matching) {
		token, err := encodePageToken(end, publishedSort)
		if err != nil {
			return nil, err
		}
		response.NextPageToken = &token
	}
	return response, nil
}

func validateRequestFilter(requested RequestFilter) []contract.FieldError {
	var out []contract.FieldError
	if !containsEnums(requested.Statuses, contract.ServiceRequestStatuses) {
		out = append(out, contract.FieldError{Field: "/query/filters/status", Code: contract.ValidationInvalidValue})
	}
	if !containsEnums(requested.ServiceTypes, contract.ServiceTypes) {
		out = append(out, contract.FieldError{Field: "/query/filters/service_type", Code: contract.ValidationInvalidValue})
	}
	if !containsEnums(requested.OwningDepartments, contract.DepartmentCodes) {
		out = append(out, contract.FieldError{Field: "/query/filters/department", Code: contract.ValidationInvalidValue})
	}
	if !containsEnums(requested.CouncilDistricts, contract.DistrictCodes) {
		out = append(out, contract.FieldError{Field: "/query/filters/district", Code: contract.ValidationInvalidValue})
	}
	if !containsEnums(requested.OriginClasses, contract.OriginClasses) {
		out = append(out, contract.FieldError{Field: "/query/filters/origin_class", Code: contract.ValidationInvalidValue})
	}
	if !containsEnums(requested.SourceChannels, contract.SourceChannels) {
		out = append(out, contract.FieldError{Field: "/query/filters/source_channel", Code: contract.ValidationInvalidValue})
	}
	if !containsEnums(requested.Categories, contract.ContactCategories) {
		out = append(out, contract.FieldError{Field: "/query/filters/category", Code: contract.ValidationInvalidValue})
	}
	if requested.CreatedFrom != nil && requested.CreatedTo != nil && requested.CreatedFrom.After(*requested.CreatedTo) {
		out = append(out, contract.FieldError{Field: "/query/filters/created_from", Code: contract.ValidationOutOfRange})
	}
	return out
}

func containsEnums[T ~string](values []T, set []T) bool {
	for _, value := range values {
		found := false
		for _, candidate := range set {
			if value == candidate {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func matchesRequestFilter(item *composeTypes.City311ServiceRequest, requested RequestFilter) bool {
	if !matchesString(string(item.Status), requested.Statuses) ||
		!matchesString(string(item.ServiceType), requested.ServiceTypes) ||
		!matchesString(string(item.OwningDepartment), requested.OwningDepartments) ||
		!matchesString(string(item.CouncilDistrict), requested.CouncilDistricts) ||
		!matchesString(string(item.OriginClass), requested.OriginClasses) ||
		!matchesString(string(item.SourceChannel), requested.SourceChannels) ||
		!matchesUint64(item.PrimaryAssigneeID, requested.PrimaryAssigneeIDs) ||
		!matchesAnyUint64(item.CollaboratorIDs, requested.CollaboratorIDs) ||
		!matchesString(fmt.Sprint(item.PrimaryRequester["primary_category"]), requested.Categories) ||
		!matchesString(item.DuplicateGroupID, requested.DuplicateGroups) {
		return false
	}
	if requested.CreatedFrom != nil && item.CreatedAt.Before(*requested.CreatedFrom) {
		return false
	}
	return requested.CreatedTo == nil || !item.CreatedAt.After(*requested.CreatedTo)
}

func matchesString[T ~string](value string, expected []T) bool {
	if len(expected) == 0 {
		return true
	}
	for _, candidate := range expected {
		if value == string(candidate) {
			return true
		}
	}
	return false
}

func matchesUint64(value uint64, expected []uint64) bool {
	if len(expected) == 0 {
		return true
	}
	for _, candidate := range expected {
		if value == candidate {
			return true
		}
	}
	return false
}

func matchesAnyUint64(values, expected []uint64) bool {
	if len(expected) == 0 {
		return true
	}
	for _, value := range values {
		if matchesUint64(value, expected) {
			return true
		}
	}
	return false
}

func normalizeSort(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		raw = "-updated_at"
	}
	parts := strings.Split(raw, ",")
	if len(parts) > 3 {
		return nil, fmt.Errorf("too many sort fields")
	}
	columns := map[string]bool{"created_at": true, "updated_at": true, "request_number": true, "status": true}
	published := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		descending := strings.HasPrefix(part, "-")
		part = strings.TrimPrefix(strings.TrimPrefix(part, "-"), "+")
		if !columns[part] {
			return nil, fmt.Errorf("unsupported sort field")
		}
		publishedValue := part
		if descending {
			publishedValue = "-" + part
		}
		published = append(published, publishedValue)
	}
	return published, nil
}

type listPageToken struct {
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
}

func encodePageToken(offset int, publishedSort []string) (string, error) {
	encoded, err := json.Marshal(listPageToken{Offset: offset, Sort: strings.Join(publishedSort, ",")})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encoded), nil
}

func decodePageToken(raw string, publishedSort []string) (int, error) {
	if raw == "" {
		return 0, nil
	}
	encoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return 0, err
	}
	token := listPageToken{}
	if err = json.Unmarshal(encoded, &token); err != nil {
		return 0, err
	}
	if token.Offset < 0 || token.Sort != strings.Join(publishedSort, ",") {
		return 0, fmt.Errorf("page token does not match requested sort")
	}
	return token.Offset, nil
}

func sortServiceRequests(set composeTypes.City311ServiceRequestSet, publishedSort []string) {
	sort.SliceStable(set, func(i, j int) bool {
		left, right := set[i], set[j]
		for _, expression := range publishedSort {
			descending := strings.HasPrefix(expression, "-")
			field := strings.TrimPrefix(expression, "-")
			comparison := compareServiceRequestField(left, right, field)
			if comparison == 0 {
				continue
			}
			if descending {
				return comparison > 0
			}
			return comparison < 0
		}
		return left.ID < right.ID
	})
}

func compareServiceRequestField(left, right *composeTypes.City311ServiceRequest, field string) int {
	switch field {
	case "created_at":
		return left.CreatedAt.Compare(right.CreatedAt)
	case "updated_at":
		return left.UpdatedAt.Compare(right.UpdatedAt)
	case "request_number":
		return strings.Compare(left.RequestNumber, right.RequestNumber)
	case "status":
		return strings.Compare(string(left.Status), string(right.Status))
	default:
		return 0
	}
}

// ResolveConstituent performs an indexed lookup of an addressable constituent
// and enforces the caller's department and district scope before exposing PII.
func (svc *Service) ResolveConstituent(ctx context.Context, actor contract.Actor, constituentID string) (*composeTypes.City311Constituent, error) {
	constituentID = strings.TrimSpace(constituentID)
	if constituentID == "" {
		return nil, validationError(contract.FieldError{Field: "/constituent/constituent_id", Code: contract.ValidationRequired})
	}
	constituent, err := store.LookupCity311ConstituentByConstituentID(ctx, svc.store, constituentID)
	if errors.IsNotFound(err) {
		return nil, apiError(404, contract.ErrorNotFound, "The constituent was not found.")
	}
	if err != nil {
		return nil, err
	}
	if !canReadConstituent(actor, constituent) {
		return nil, apiError(403, contract.ErrorForbidden, "The constituent is outside the actor's record scope.")
	}
	return constituent, nil
}

func appliedFilters(requested RequestFilter) map[string]any {
	out := map[string]any{}
	if len(requested.Statuses) > 0 {
		out["status"] = requested.Statuses
	}
	if len(requested.ServiceTypes) > 0 {
		out["service_type"] = requested.ServiceTypes
	}
	if len(requested.OwningDepartments) > 0 {
		out["department"] = requested.OwningDepartments
	}
	if len(requested.CouncilDistricts) > 0 {
		out["district"] = requested.CouncilDistricts
	}
	if len(requested.OriginClasses) > 0 {
		out["origin_class"] = requested.OriginClasses
	}
	if len(requested.SourceChannels) > 0 {
		out["source_channel"] = requested.SourceChannels
	}
	if len(requested.PrimaryAssigneeIDs) > 0 {
		out["assignee"] = stringifyIDs(requested.PrimaryAssigneeIDs)
	}
	if len(requested.CollaboratorIDs) > 0 {
		out["collaborator"] = stringifyIDs(requested.CollaboratorIDs)
	}
	if len(requested.Categories) > 0 {
		out["category"] = requested.Categories
	}
	if requested.CreatedFrom != nil {
		out["created_from"] = requested.CreatedFrom.Format(time.RFC3339)
	}
	if requested.CreatedTo != nil {
		out["created_to"] = requested.CreatedTo.Format(time.RFC3339)
	}
	if len(requested.DuplicateGroups) > 0 {
		out["duplicate_group"] = requested.DuplicateGroups
	}
	return out
}

func (svc *Service) detail(ctx context.Context, actor contract.Actor, stored *composeTypes.City311ServiceRequest) (*contract.StaffServiceRequestDetail, error) {
	audits, _, err := store.SearchCity311AuditEvents(ctx, svc.store, composeTypes.City311AuditEventFilter{RequestID: stored.ID})
	if err != nil {
		return nil, err
	}
	history, _, err := store.SearchCity311PublicHistoryItems(ctx, svc.store, composeTypes.City311PublicHistoryItemFilter{RequestID: stored.ID})
	if err != nil {
		return nil, err
	}
	primaryAssignee := optionalID(stored.PrimaryAssigneeID)
	result := &contract.StaffServiceRequestDetail{
		Request: toContract(stored), AvailableActions: availableActions(actor, stored), PrimaryAssigneeID: primaryAssignee,
		CollaboratorIDs: stringifyIDs(stored.CollaboratorIDs), Reminders: []any{}, History: make([]contract.PublicHistoryItem, 0, len(history)), Audit: make([]contract.AuditEvent, 0, len(audits)), ExternalWorkOrder: nil,
	}
	for _, item := range history {
		result.History = append(result.History, contract.PublicHistoryItem{Action: item.Action, OccurredAt: item.OccurredAt, ResponsibleDepartment: string(item.ResponsibleDepartment)})
	}
	for _, audit := range audits {
		result.Audit = append(result.Audit, contract.AuditEvent{
			EntityType: "service_request", EntityID: strconv.FormatUint(audit.RequestID, 10), EventType: audit.EventType,
			ActorType: audit.ActorType, ActorID: strconv.FormatUint(audit.ActorID, 10), OccurredAt: audit.CreatedAt,
			SourceChannel: audit.SourceChannel, Before: audit.Before, After: audit.After,
		})
	}
	return result, nil
}

func canRead(actor contract.Actor, request *composeTypes.City311ServiceRequest) bool {
	if hasRole(actor, contract.ApplicationRolePlatformAdministrator) {
		return true
	}
	if !isStaff(actor) || actor.Department != request.OwningDepartment {
		return false
	}
	if hasRole(actor, contract.ApplicationRoleDepartmentManager) || request.CouncilDistrict == "" {
		return true
	}
	for _, district := range actor.Districts {
		if district == request.CouncilDistrict {
			return true
		}
	}
	return false
}

func canReadConstituent(actor contract.Actor, constituent *composeTypes.City311Constituent) bool {
	if hasRole(actor, contract.ApplicationRolePlatformAdministrator) {
		return true
	}
	if !isStaff(actor) || actor.Department != constituent.OwningDepartment {
		return false
	}
	if hasRole(actor, contract.ApplicationRoleDepartmentManager) || constituent.CouncilDistrict == "" {
		return true
	}
	for _, district := range actor.Districts {
		if district == constituent.CouncilDistrict {
			return true
		}
	}
	return false
}

func isStaff(actor contract.Actor) bool {
	return hasRole(actor, contract.ApplicationRoleServiceAgent) || hasRole(actor, contract.ApplicationRoleSupervisor) || hasRole(actor, contract.ApplicationRoleDepartmentManager) || hasRole(actor, contract.ApplicationRolePlatformAdministrator) || hasRole(actor, contract.ApplicationRoleWorkflowDesigner)
}

func canOperateRequest(actor contract.Actor) bool {
	return hasRole(actor, contract.ApplicationRoleServiceAgent) || hasRole(actor, contract.ApplicationRoleSupervisor) || hasRole(actor, contract.ApplicationRoleDepartmentManager) || hasRole(actor, contract.ApplicationRolePlatformAdministrator)
}

func hasRole(actor contract.Actor, role contract.ApplicationRole) bool {
	for _, candidate := range actor.Roles {
		if candidate == role {
			return true
		}
	}
	return false
}
