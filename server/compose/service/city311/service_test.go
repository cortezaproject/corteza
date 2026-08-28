package city311

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	composeTypes "github.com/cortezaproject/corteza/server/compose/types"
	contract "github.com/cortezaproject/corteza/server/compose/types/city311"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func testService(t *testing.T) (*Service, store.Storer) {
	t.Helper()
	ctx := context.Background()
	dsn := fmt.Sprintf("sqlite3://file:%s?mode=memory&cache=shared", t.Name())
	st, err := sqlite.Connect(ctx, dsn)
	require.NoError(t, err)
	require.NoError(t, store.Upgrade(ctx, zap.NewNop(), st))

	svc := New(st)
	fixedNow := time.Date(2026, 2, 3, 15, 4, 5, 0, time.UTC)
	svc.now = func() time.Time { return fixedNow }
	var next uint64 = 900_000_000_000_000_000
	svc.nextID = func() uint64 { next++; return next }
	return svc, st
}

func validSubmission() contract.ServiceRequestCreate {
	latitude, longitude := 42.88645, -78.87837
	return contract.ServiceRequestCreate{
		Summary: "Pothole on Example Street", Description: "A deep pothole is blocking the eastbound traffic lane.",
		ServiceType:  contract.ServiceTypePothole,
		Requester:    contract.RequesterInput{DisplayName: "Alex Resident", Email: "alex@example.invalid", Phone: "+17165550101"},
		Location:     &contract.LocationInput{Address: "100 Example Street, Buffalo, NY 14201", Latitude: &latitude, Longitude: &longitude},
		CustomFields: map[string]any{"reported_damage": false},
	}
}

func TestSeedIsRepeatableAndPreservesSeededRows(t *testing.T) {
	svc, st := testService(t)
	ctx := context.Background()
	require.NoError(t, svc.Seed(ctx, svc.now()))

	seeded, err := store.LookupCity311ServiceRequestByRequestNumber(ctx, st, "SR-2026-00034")
	require.NoError(t, err)
	seeded.Summary = "Preserved operator edit"
	require.NoError(t, store.UpdateCity311ServiceRequest(ctx, st, seeded))
	constituent, err := store.LookupCity311ConstituentByConstituentID(ctx, st, "C-1")
	require.NoError(t, err)
	constituent.Profile["display_name"] = "Preserved constituent edit"
	require.NoError(t, store.UpdateCity311Constituent(ctx, st, constituent))
	require.NoError(t, svc.Seed(ctx, svc.now()))

	set, _, err := store.SearchCity311ServiceRequests(ctx, st, composeTypes.City311ServiceRequestFilter{Paging: filter.Paging{Limit: 100}})
	require.NoError(t, err)
	require.Len(t, set, 8)
	seeded, err = store.LookupCity311ServiceRequestByRequestNumber(ctx, st, "SR-2026-00034")
	require.NoError(t, err)
	require.Equal(t, "Preserved operator edit", seeded.Summary)
	constituent, err = store.LookupCity311ConstituentByConstituentID(ctx, st, "C-1")
	require.NoError(t, err)
	require.Equal(t, "Preserved constituent edit", constituent.Profile["display_name"])
	constituents, _, err := store.SearchCity311Constituents(ctx, st, composeTypes.City311ConstituentFilter{Paging: filter.Paging{Limit: 100}})
	require.NoError(t, err)
	require.Len(t, constituents, 8)
	history, _, err := store.SearchCity311PublicHistoryItems(ctx, st, composeTypes.City311PublicHistoryItemFilter{RequestID: seeded.ID})
	require.NoError(t, err)
	require.Len(t, history, 1)
	sequence, err := store.LookupCity311RequestSequenceByID(ctx, st, 2026)
	require.NoError(t, err)
	require.Equal(t, uint64(41), sequence.NextNumber)
}

func TestUpgradeAndSeedPreserveBaselineUser(t *testing.T) {
	svc, st := testService(t)
	ctx := context.Background()
	createdAt := svc.now().Add(-24 * time.Hour)
	baseline := &systemTypes.User{ID: 800_000_000_000_000_001, Handle: "baseline-operator", Username: "baseline-operator", Email: "baseline-operator@example.invalid", Name: "Baseline Operator", EmailConfirmed: true, CreatedAt: createdAt}
	require.NoError(t, store.CreateUser(ctx, st, baseline))
	require.NoError(t, store.Upgrade(ctx, zap.NewNop(), st))
	require.NoError(t, svc.Seed(ctx, svc.now()))

	fetched, err := store.LookupUserByID(ctx, st, baseline.ID)
	require.NoError(t, err)
	require.Equal(t, baseline.Handle, fetched.Handle)
	require.Equal(t, baseline.Email, fetched.Email)
}

func TestSubmitIsTransactionalAndIdempotent(t *testing.T) {
	svc, st := testService(t)
	ctx := context.Background()
	require.NoError(t, svc.Seed(ctx, svc.now()))
	options := SubmissionOptions{
		Operation: "service_request_create", SourceChannel: contract.SourceChannelAPI,
		ActorType: contract.AuditActorIntegrationClient, ActorID: 77, RequireIdempotency: true,
	}

	first, status, err := svc.Submit(ctx, validSubmission(), "request-key-1", options)
	require.NoError(t, err)
	require.Equal(t, 201, status)
	require.Equal(t, "SR-2026-00041", first.RequestNumber)
	createdConstituent, err := store.LookupCity311ConstituentByConstituentID(ctx, st, "C-"+first.RequestID)
	require.NoError(t, err)
	require.Equal(t, "Alex Resident", createdConstituent.Profile["display_name"])

	replay, status, err := svc.Submit(ctx, validSubmission(), "request-key-1", options)
	require.NoError(t, err)
	require.Equal(t, 200, status)
	require.Equal(t, first, replay)

	conflicting := validSubmission()
	conflicting.Summary = "A different report"
	_, _, err = svc.Submit(ctx, conflicting, "request-key-1", options)
	var serviceErr *ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, 409, serviceErr.Status)
	require.Equal(t, contract.ErrorIdempotencyConflict, serviceErr.Payload.Error)

	set, _, err := store.SearchCity311ServiceRequests(ctx, st, composeTypes.City311ServiceRequestFilter{Paging: filter.Paging{Limit: 100}})
	require.NoError(t, err)
	require.Len(t, set, 9)
	constituents, _, err := store.SearchCity311Constituents(ctx, st, composeTypes.City311ConstituentFilter{Paging: filter.Paging{Limit: 100}})
	require.NoError(t, err)
	require.Len(t, constituents, 9)
	sequence, err := store.LookupCity311RequestSequenceByID(ctx, st, 2026)
	require.NoError(t, err)
	require.Equal(t, uint64(42), sequence.NextNumber)
}

func TestExpiredIdempotencyKeyCanBeReused(t *testing.T) {
	svc, st := testService(t)
	ctx := context.Background()
	require.NoError(t, svc.Seed(ctx, svc.now()))
	options := SubmissionOptions{Operation: "service_request_create", SourceChannel: contract.SourceChannelAPI, ActorType: contract.AuditActorIntegrationClient, RequireIdempotency: true}
	first, _, err := svc.Submit(ctx, validSubmission(), "expiring-key", options)
	require.NoError(t, err)
	record, err := store.LookupCity311IdempotencyRecordByOperationKeyHash(ctx, st, options.Operation, hashKey("expiring-key"))
	require.NoError(t, err)
	record.ExpiresAt = svc.now()
	require.NoError(t, store.UpdateCity311IdempotencyRecord(ctx, st, record))
	nextDay := svc.now().Add(25 * time.Hour)
	svc.now = func() time.Time { return nextDay }

	reused := validSubmission()
	reused.Summary = "Pothole reported after the replay window"
	second, status, err := svc.Submit(ctx, reused, "expiring-key", options)
	require.NoError(t, err)
	require.Equal(t, 201, status)
	require.NotEqual(t, first.RequestID, second.RequestID)
}

func TestInlineAttachmentValidationAndPersistence(t *testing.T) {
	svc, st := testService(t)
	ctx := context.Background()
	require.NoError(t, svc.Seed(ctx, svc.now()))
	input := validSubmission()
	input.Attachments = []contract.AttachmentInput{{
		Filename: "../evidence.txt", MediaType: "text/plain", ContentBase64: base64.StdEncoding.EncodeToString([]byte("attachment evidence")),
	}}
	created, _, err := svc.Submit(ctx, input, "attachment-key", SubmissionOptions{Operation: "service_request_create", SourceChannel: contract.SourceChannelAPI, RequireIdempotency: true})
	require.NoError(t, err)
	requestID, err := strconv.ParseUint(created.RequestID, 10, 64)
	require.NoError(t, err)
	attachments, _, err := store.SearchCity311RequestAttachments(ctx, st, composeTypes.City311RequestAttachmentFilter{RequestID: requestID})
	require.NoError(t, err)
	require.Len(t, attachments, 1)
	require.Equal(t, "evidence.txt", attachments[0].Filename)
	require.Equal(t, []byte("attachment evidence"), attachments[0].Content)

	invalid := validSubmission()
	invalid.Attachments = []contract.AttachmentInput{{Filename: "payload.exe", MediaType: "application/octet-stream", ContentBase64: "not-base64"}}
	_, _, err = svc.Submit(ctx, invalid, "invalid-attachment", SubmissionOptions{Operation: "service_request_create", SourceChannel: contract.SourceChannelAPI, RequireIdempotency: true})
	var serviceErr *ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, 422, serviceErr.Status)
	require.Len(t, serviceErr.Payload.Errors, 2)
}

func TestConcurrentSubmissionsAllocateUniqueSequentialNumbers(t *testing.T) {
	svc, _ := testService(t)
	ctx := context.Background()
	require.NoError(t, svc.Seed(ctx, svc.now()))
	const count = 16
	numbers := make([]string, count)
	errors := make([]error, count)
	var wait sync.WaitGroup
	for index := 0; index < count; index++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			response, _, err := svc.Submit(ctx, validSubmission(), fmt.Sprintf("parallel-%02d", index), SubmissionOptions{
				Operation: "portal_service_request_submit", SourceChannel: contract.SourceChannelPortalAnonymous,
				ActorType: contract.AuditActorConstituent, RequireIdempotency: true,
			})
			errors[index] = err
			if response != nil {
				numbers[index] = response.RequestNumber
			}
		}(index)
	}
	wait.Wait()
	for _, err := range errors {
		require.NoError(t, err)
	}
	sort.Strings(numbers)
	for index, number := range numbers {
		require.Equal(t, fmt.Sprintf("SR-2026-%05d", 41+index), number)
	}
}

func TestRecordScopeAndOptimisticTransition(t *testing.T) {
	svc, st := testService(t)
	ctx := context.Background()
	require.NoError(t, svc.Seed(ctx, svc.now()))
	request, err := store.LookupCity311ServiceRequestByRequestNumber(ctx, st, "SR-2026-00034")
	require.NoError(t, err)
	actor := contract.Actor{ID: 100, Roles: []contract.ApplicationRole{contract.ApplicationRoleServiceAgent}, Department: contract.DepartmentStreets, Districts: []contract.DistrictCode{contract.DistrictNorth}}
	_, err = svc.Transition(ctx, actor, request.ID, 1, contract.RequestTransition{ToStatus: contract.ServiceRequestStatusClosed})
	var serviceErr *ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, 422, serviceErr.Status)
	unchanged, err := store.LookupCity311ServiceRequestByID(ctx, st, request.ID)
	require.NoError(t, err)
	require.Equal(t, contract.ServiceRequestStatusSubmitted, unchanged.Status)
	require.Equal(t, 1, unchanged.Version)

	detail, err := svc.Transition(ctx, actor, request.ID, 1, contract.RequestTransition{ToStatus: contract.ServiceRequestStatusTriaged, Reason: "Validated and routed"})
	require.NoError(t, err)
	require.Equal(t, uint64(2), detail.Request.Version)
	require.Equal(t, contract.ServiceRequestStatusTriaged, detail.Request.Status)
	require.Len(t, detail.Audit, 2)
	require.Len(t, detail.History, 2)
	require.Equal(t, "TRIAGED", detail.History[1].Action)

	_, err = svc.Transition(ctx, actor, request.ID, 1, contract.RequestTransition{ToStatus: contract.ServiceRequestStatusAssigned})
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, 409, serviceErr.Status)
	require.NotNil(t, serviceErr.Payload.CurrentVersion)
	require.Equal(t, uint64(2), *serviceErr.Payload.CurrentVersion)

	wrongDepartment := actor
	wrongDepartment.Department = contract.DepartmentSanitation
	_, err = svc.Find(ctx, wrongDepartment, request.ID)
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, 403, serviceErr.Status)

	workflowOnly := contract.Actor{ID: 200, Roles: []contract.ApplicationRole{contract.ApplicationRoleWorkflowDesigner}, Department: contract.DepartmentStreets, Districts: []contract.DistrictCode{contract.DistrictNorth}}
	_, err = svc.Transition(ctx, workflowOnly, request.ID, 2, contract.RequestTransition{ToStatus: contract.ServiceRequestStatusAssigned})
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, 403, serviceErr.Status)
}

func TestListUsesOpaqueCursorTotalAndPublishedSort(t *testing.T) {
	svc, _ := testService(t)
	ctx := context.Background()
	require.NoError(t, svc.Seed(ctx, svc.now()))
	actor := contract.Actor{ID: 100, Roles: []contract.ApplicationRole{contract.ApplicationRoleServiceAgent}, Department: contract.DepartmentStreets, Districts: []contract.DistrictCode{contract.DistrictNorth}}

	first, err := svc.List(ctx, actor, RequestFilter{PageSize: 2, Sort: "-updated_at"})
	require.NoError(t, err)
	require.Len(t, first.Items, 2)
	require.Equal(t, 8, first.TotalCount)
	require.Equal(t, []string{"-updated_at"}, first.Sort)
	require.NotNil(t, first.NextPageToken)

	second, err := svc.List(ctx, actor, RequestFilter{PageSize: 2, PageToken: *first.NextPageToken, Sort: "-updated_at"})
	require.NoError(t, err)
	require.Len(t, second.Items, 2)
	require.Equal(t, 8, second.TotalCount)
	require.NotEqual(t, first.Items[0].RequestID, second.Items[0].RequestID)

	_, err = svc.List(ctx, actor, RequestFilter{Statuses: []contract.ServiceRequestStatus{"UNKNOWN"}})
	var serviceErr *ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, 422, serviceErr.Status)
}

func TestRequestQueueMatcherCoversEveryPublishedFilter(t *testing.T) {
	createdAt := time.Date(2026, 2, 3, 10, 0, 0, 0, time.UTC)
	request := &composeTypes.City311ServiceRequest{
		Status: contract.ServiceRequestStatusSubmitted, ServiceType: contract.ServiceTypePothole,
		OwningDepartment: contract.DepartmentStreets, CouncilDistrict: contract.DistrictNorth,
		OriginClass: contract.OriginClassInternal, SourceChannel: contract.SourceChannelStaffInPerson,
		PrimaryAssigneeID: 101, CollaboratorIDs: composeTypes.City311Uint64Set{202},
		PrimaryRequester: map[string]any{"primary_category": string(contract.ContactCategoryBusiness)},
		DuplicateGroupID: "duplicate-1", CreatedAt: createdAt,
	}
	from, to := createdAt.Add(-time.Minute), createdAt.Add(time.Minute)
	matching := RequestFilter{
		Statuses:           []contract.ServiceRequestStatus{contract.ServiceRequestStatusSubmitted},
		ServiceTypes:       []contract.ServiceType{contract.ServiceTypePothole},
		OwningDepartments:  []contract.DepartmentCode{contract.DepartmentStreets},
		CouncilDistricts:   []contract.DistrictCode{contract.DistrictNorth},
		OriginClasses:      []contract.OriginClass{contract.OriginClassInternal},
		SourceChannels:     []contract.SourceChannel{contract.SourceChannelStaffInPerson},
		PrimaryAssigneeIDs: []uint64{101}, CollaboratorIDs: []uint64{202},
		Categories:  []contract.ContactCategory{contract.ContactCategoryBusiness},
		CreatedFrom: &from, CreatedTo: &to, DuplicateGroups: []string{"duplicate-1"},
	}
	require.True(t, matchesRequestFilter(request, matching))

	nonMatching := []RequestFilter{
		{Statuses: []contract.ServiceRequestStatus{contract.ServiceRequestStatusClosed}},
		{ServiceTypes: []contract.ServiceType{contract.ServiceTypeGeneralInquiry}},
		{OwningDepartments: []contract.DepartmentCode{contract.DepartmentSanitation}},
		{CouncilDistricts: []contract.DistrictCode{contract.DistrictSouth}},
		{OriginClasses: []contract.OriginClass{contract.OriginClassExternal}},
		{SourceChannels: []contract.SourceChannel{contract.SourceChannelAPI}},
		{PrimaryAssigneeIDs: []uint64{999}},
		{CollaboratorIDs: []uint64{999}},
		{Categories: []contract.ContactCategory{contract.ContactCategoryResident}},
		{CreatedFrom: func() *time.Time { value := createdAt.Add(time.Minute); return &value }()},
		{CreatedTo: func() *time.Time { value := createdAt.Add(-time.Minute); return &value }()},
		{DuplicateGroups: []string{"duplicate-2"}},
	}
	for index, filter := range nonMatching {
		require.False(t, matchesRequestFilter(request, filter), "published filter index %d was ignored", index)
	}
}

func TestValidationPrecedesPersistence(t *testing.T) {
	svc, st := testService(t)
	ctx := context.Background()
	require.NoError(t, svc.Seed(ctx, svc.now()))
	invalid := validSubmission()
	invalid.Description = "short"
	invalid.Location = nil
	_, _, err := svc.Submit(ctx, invalid, "invalid", SubmissionOptions{RequireIdempotency: true})
	var serviceErr *ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, 422, serviceErr.Status)
	require.Len(t, serviceErr.Payload.Errors, 2)
	set, _, searchErr := store.SearchCity311ServiceRequests(ctx, st, composeTypes.City311ServiceRequestFilter{Paging: filter.Paging{Limit: 100}})
	require.NoError(t, searchErr)
	require.Len(t, set, 8)
}

func TestResolveConstituentUsesIndexedResourceBeyondOneHundredRequests(t *testing.T) {
	svc, _ := testService(t)
	ctx := context.Background()
	require.NoError(t, svc.Seed(ctx, svc.now()))
	actor := contract.Actor{
		ID: 100, Roles: []contract.ApplicationRole{contract.ApplicationRoleServiceAgent},
		Department: contract.DepartmentStreets, Districts: []contract.DistrictCode{contract.DistrictNorth},
	}
	lastConstituentID := ""
	for index := 0; index < 110; index++ {
		input := validSubmission()
		input.Requester.DisplayName = fmt.Sprintf("Constituent %d", index)
		input.Requester.Email = fmt.Sprintf("constituent-%d@example.invalid", index)
		response, status, err := svc.Submit(ctx, input, "", SubmissionOptions{
			SourceChannel: contract.SourceChannelAPI, ActorType: contract.AuditActorIntegrationClient,
		})
		require.NoError(t, err)
		require.Equal(t, 201, status)
		lastConstituentID = "C-" + response.RequestID
	}
	resolved, err := svc.ResolveConstituent(ctx, actor, lastConstituentID)
	require.NoError(t, err)
	require.Equal(t, lastConstituentID, resolved.ConstituentID)
	require.Equal(t, "Constituent 109", resolved.Profile["display_name"])
}

func TestValidationInitializationAndListErrorPaths(t *testing.T) {
	svc, st := testService(t)
	ctx := context.Background()
	previousDefault := Default
	t.Cleanup(func() { Default = previousDefault })

	t.Setenv("BENCHMARK_NOW", "not-a-time")
	require.Error(t, Initialize(ctx, st))
	require.NoError(t, svc.Seed(ctx, svc.now()))
	require.Equal(t, "validation failed", (&ServiceError{Payload: contract.APIError{Message: "validation failed"}}).Error())

	invalid := validSubmission()
	invalid.Summary = strings.Repeat("s", 161)
	invalid.Description = strings.Repeat("d", 5001)
	invalid.ServiceType = contract.ServiceType("UNKNOWN")
	invalid.Requester = contract.RequesterInput{DisplayName: strings.Repeat("n", 121), Email: "Named <named@example.invalid>", Phone: "555-0100"}
	invalid.Location = nil
	invalid.Attachments = make([]contract.AttachmentInput, 6)
	require.GreaterOrEqual(t, len(validateWrite(invalid).Payload.Errors), 7)

	latitude, longitude := 91.0, 181.0
	locationErrors := validateLocation(contract.ServiceTypePothole, &contract.LocationInput{Address: "Address", Latitude: &latitude, Longitude: &longitude})
	require.Len(t, locationErrors, 2)
	require.Len(t, validateLocation(contract.ServiceTypePothole, &contract.LocationInput{Address: "Address"}), 1)
	require.Empty(t, validateLocation(contract.ServiceTypeGeneralInquiry, nil))

	_, attachmentErr := validateInlineAttachments([]contract.AttachmentInput{
		{Filename: "", MediaType: "invalid/type", ContentBase64: "not-base64"},
		{Filename: strings.Repeat("a", 121), MediaType: "text/plain", ContentBase64: base64.StdEncoding.EncodeToString([]byte("ok"))},
	})
	require.NotNil(t, attachmentErr)
	require.GreaterOrEqual(t, len(attachmentErr.Payload.Errors), 4)

	for _, serviceType := range contract.ServiceTypes {
		_, ok := departmentForServiceType(serviceType)
		require.True(t, ok, serviceType)
	}
	_, ok := departmentForServiceType(contract.ServiceType("UNKNOWN"))
	require.False(t, ok)
	require.Error(t, validateIdempotencyKey("", true))
	require.Error(t, validateIdempotencyKey(strings.Repeat("k", 256), false))
	_, err := hashJSON(make(chan int))
	require.Error(t, err)

	require.Error(t, validateTransitionInput(0, contract.RequestTransition{ToStatus: contract.ServiceRequestStatusTriaged}))
	require.Error(t, validateTransitionInput(1, contract.RequestTransition{}))
	_, err = svc.FindActor(ctx, 999999999)
	require.Error(t, err)
	_, err = svc.Find(ctx, contract.Actor{Roles: []contract.ApplicationRole{contract.ApplicationRolePlatformAdministrator}}, 999999999)
	require.Error(t, err)
	_, err = svc.ResolveConstituent(ctx, contract.Actor{}, " ")
	require.Error(t, err)
	_, err = svc.ResolveConstituent(ctx, contract.Actor{Roles: []contract.ApplicationRole{contract.ApplicationRolePlatformAdministrator}}, "missing")
	require.Error(t, err)

	staff := contract.Actor{ID: 100, Roles: []contract.ApplicationRole{contract.ApplicationRoleServiceAgent}, Department: contract.DepartmentStreets, Districts: []contract.DistrictCode{contract.DistrictNorth}}
	_, err = svc.List(ctx, contract.Actor{}, RequestFilter{})
	require.Error(t, err)
	_, err = svc.List(ctx, staff, RequestFilter{PageSize: 101})
	require.Error(t, err)
	_, err = svc.List(ctx, staff, RequestFilter{Sort: "unknown"})
	require.Error(t, err)
	_, err = svc.List(ctx, staff, RequestFilter{PageToken: "not-base64"})
	require.Error(t, err)
	token, err := encodePageToken(100, []string{"-updated_at"})
	require.NoError(t, err)
	_, err = svc.List(ctx, staff, RequestFilter{PageSize: 1, PageToken: token})
	require.Error(t, err)

	from, to := svc.now(), svc.now().Add(-time.Hour)
	filterErrors := validateRequestFilter(RequestFilter{
		Statuses: []contract.ServiceRequestStatus{"UNKNOWN"}, ServiceTypes: []contract.ServiceType{"UNKNOWN"},
		OwningDepartments: []contract.DepartmentCode{"UNKNOWN"}, CouncilDistricts: []contract.DistrictCode{"UNKNOWN"},
		OriginClasses: []contract.OriginClass{"UNKNOWN"}, SourceChannels: []contract.SourceChannel{"UNKNOWN"},
		Categories: []contract.ContactCategory{"UNKNOWN"}, CreatedFrom: &from, CreatedTo: &to,
	})
	require.Len(t, filterErrors, 8)
	require.Empty(t, validateRequestFilter(RequestFilter{}))

	_, err = decodePageToken(base64.RawURLEncoding.EncodeToString([]byte("{")), []string{"-updated_at"})
	require.Error(t, err)
	mismatched, err := encodePageToken(1, []string{"created_at"})
	require.NoError(t, err)
	_, err = decodePageToken(mismatched, []string{"-updated_at"})
	require.Error(t, err)
	_, err = normalizeSort("created_at,updated_at,status,request_number")
	require.Error(t, err)

	left := &composeTypes.City311ServiceRequest{RequestNumber: "SR-1", Status: contract.ServiceRequestStatusSubmitted, CreatedAt: svc.now(), UpdatedAt: svc.now()}
	right := &composeTypes.City311ServiceRequest{RequestNumber: "SR-2", Status: contract.ServiceRequestStatusTriaged, CreatedAt: svc.now().Add(time.Minute), UpdatedAt: svc.now().Add(time.Minute)}
	for _, field := range []string{"created_at", "updated_at", "request_number", "status"} {
		require.Less(t, compareServiceRequestField(left, right, field), 0, field)
	}
	require.Zero(t, compareServiceRequestField(left, right, "unknown"))
}
