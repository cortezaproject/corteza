package tests

import (
	"context"
	"testing"
	"time"

	composeTypes "github.com/cortezaproject/corteza/server/compose/types"
	contract "github.com/cortezaproject/corteza/server/compose/types/city311"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/stretchr/testify/require"
)

func testCity311ActorProfiles(t *testing.T, s store.City311ActorProfiles) {
	ctx := context.Background()
	require.NoError(t, s.TruncateCity311ActorProfiles(ctx))
	profile := &composeTypes.City311ActorProfile{
		ID: 101, ApplicationRoles: composeTypes.City311ApplicationRoleSet{contract.ApplicationRoleServiceAgent},
		Department: contract.DepartmentStreets, Districts: composeTypes.City311DistrictCodeSet{contract.DistrictNorth},
		CreatedAt: *now(), UpdatedAt: *now(),
	}
	require.NoError(t, s.CreateCity311ActorProfile(ctx, profile))
	fetched, err := s.LookupCity311ActorProfileByID(ctx, profile.ID)
	require.NoError(t, err)
	require.Equal(t, profile.ApplicationRoles, fetched.ApplicationRoles)
	set, _, err := s.SearchCity311ActorProfiles(ctx, composeTypes.City311ActorProfileFilter{Department: string(contract.DepartmentStreets)})
	require.NoError(t, err)
	require.Len(t, set, 1)
	fetched.Districts = append(fetched.Districts, contract.DistrictCentral)
	require.NoError(t, s.UpdateCity311ActorProfile(ctx, fetched))
	require.NoError(t, s.DeleteCity311ActorProfileByID(ctx, profile.ID))
}

func testCity311AuditEvents(t *testing.T, s store.City311AuditEvents) {
	ctx := context.Background()
	require.NoError(t, s.TruncateCity311AuditEvents(ctx))
	event := &composeTypes.City311AuditEvent{
		ID: 201, RequestID: 301, EventType: "SERVICE_REQUEST_SUBMITTED", ActorType: contract.AuditActorSystem,
		SourceChannel: contract.SourceChannelAPI, Before: composeTypes.City311JSON{}, After: composeTypes.City311JSON{"status": "SUBMITTED"}, CreatedAt: *now(),
	}
	require.NoError(t, s.CreateCity311AuditEvent(ctx, event))
	fetched, err := s.LookupCity311AuditEventByID(ctx, event.ID)
	require.NoError(t, err)
	require.Equal(t, "SUBMITTED", fetched.After["status"])
	set, _, err := s.SearchCity311AuditEvents(ctx, composeTypes.City311AuditEventFilter{RequestID: event.RequestID})
	require.NoError(t, err)
	require.Len(t, set, 1)
	require.NoError(t, s.DeleteCity311AuditEvent(ctx, event))
}

func testCity311RequestAttachments(t *testing.T, s store.City311RequestAttachments) {
	ctx := context.Background()
	require.NoError(t, s.TruncateCity311RequestAttachments(ctx))
	attachment := &composeTypes.City311RequestAttachment{
		ID: 251, RequestID: 301, Filename: "photo.png", MediaType: "image/png", Size: 4, Content: []byte("data"), CreatedAt: *now(),
	}
	require.NoError(t, s.CreateCity311RequestAttachment(ctx, attachment))
	fetched, err := s.LookupCity311RequestAttachmentByID(ctx, attachment.ID)
	require.NoError(t, err)
	require.Equal(t, attachment.Content, fetched.Content)
	set, _, err := s.SearchCity311RequestAttachments(ctx, composeTypes.City311RequestAttachmentFilter{RequestID: attachment.RequestID})
	require.NoError(t, err)
	require.Len(t, set, 1)
	require.NoError(t, s.DeleteCity311RequestAttachmentByID(ctx, attachment.ID))
}

func testCity311PublicHistoryItems(t *testing.T, s store.City311PublicHistoryItems) {
	ctx := context.Background()
	require.NoError(t, s.TruncateCity311PublicHistoryItems(ctx))
	item := &composeTypes.City311PublicHistoryItem{
		ID: 275, RequestID: 301, Action: "SUBMITTED", ResponsibleDepartment: contract.DepartmentStreets, OccurredAt: *now(),
	}
	require.NoError(t, s.CreateCity311PublicHistoryItem(ctx, item))
	fetched, err := s.LookupCity311PublicHistoryItemByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.Action, fetched.Action)
	set, _, err := s.SearchCity311PublicHistoryItems(ctx, composeTypes.City311PublicHistoryItemFilter{RequestID: item.RequestID})
	require.NoError(t, err)
	require.Len(t, set, 1)
	require.NoError(t, s.DeleteCity311PublicHistoryItemByID(ctx, item.ID))
}

func testCity311IdempotencyRecords(t *testing.T, s store.City311IdempotencyRecords) {
	ctx := context.Background()
	require.NoError(t, s.TruncateCity311IdempotencyRecords(ctx))
	record := &composeTypes.City311IdempotencyRecord{
		ID: 401, Operation: "service_request_create", KeyHash: "key-hash", RequestHash: "request-hash", ResponseStatus: 201,
		ResponseBody: composeTypes.City311JSON{"request_id": "301"}, RequestID: 301, CreatedAt: *now(), ExpiresAt: now().Add(24 * time.Hour),
	}
	require.NoError(t, s.CreateCity311IdempotencyRecord(ctx, record))
	fetched, err := s.LookupCity311IdempotencyRecordByOperationKeyHash(ctx, record.Operation, record.KeyHash)
	require.NoError(t, err)
	require.Equal(t, record.RequestHash, fetched.RequestHash)
	require.NoError(t, s.DeleteCity311IdempotencyRecordByID(ctx, record.ID))
}

func testCity311RequestSequences(t *testing.T, s store.City311RequestSequences) {
	ctx := context.Background()
	require.NoError(t, s.TruncateCity311RequestSequences(ctx))
	sequence := &composeTypes.City311RequestSequence{ID: 2026, NextNumber: 41}
	require.NoError(t, s.CreateCity311RequestSequence(ctx, sequence))
	sequence.NextNumber = 42
	require.NoError(t, s.UpdateCity311RequestSequence(ctx, sequence))
	fetched, err := s.LookupCity311RequestSequenceByID(ctx, 2026)
	require.NoError(t, err)
	require.Equal(t, uint64(42), fetched.NextNumber)
	require.NoError(t, s.DeleteCity311RequestSequence(ctx, sequence))
}

func testCity311ServiceRequests(t *testing.T, s store.City311ServiceRequests) {
	ctx := context.Background()
	require.NoError(t, s.TruncateCity311ServiceRequests(ctx))
	request := &composeTypes.City311ServiceRequest{
		ID: 301, RequestNumber: "SR-2026-00041", Summary: "Pothole on Example Street", Description: "A deep pothole blocks one traffic lane.",
		ServiceType: contract.ServiceTypePothole, OwningDepartment: contract.DepartmentStreets, CouncilDistrict: contract.DistrictNorth,
		SourceChannel: contract.SourceChannelAPI, OriginClass: contract.OriginClassExternal, Status: contract.ServiceRequestStatusSubmitted,
		PrimaryRequester: composeTypes.City311JSON{"constituent_id": "C-301"}, Location: composeTypes.City311JSON{"address": "100 Example Street"},
		CustomFields: composeTypes.City311JSON{}, CollaboratorIDs: composeTypes.City311Uint64Set{}, Version: 1, CreatedAt: *now(), UpdatedAt: *now(),
	}
	require.NoError(t, s.CreateCity311ServiceRequest(ctx, request))
	fetched, err := s.LookupCity311ServiceRequestByRequestNumber(ctx, request.RequestNumber)
	require.NoError(t, err)
	require.Equal(t, "C-301", fetched.PrimaryRequester["constituent_id"])
	set, _, err := s.SearchCity311ServiceRequests(ctx, composeTypes.City311ServiceRequestFilter{Status: string(contract.ServiceRequestStatusSubmitted)})
	require.NoError(t, err)
	require.Len(t, set, 1)
	fetched.Status = contract.ServiceRequestStatusTriaged
	fetched.Version++
	require.NoError(t, s.UpdateCity311ServiceRequest(ctx, fetched))
	require.NoError(t, s.DeleteCity311ServiceRequestByID(ctx, request.ID))
}
