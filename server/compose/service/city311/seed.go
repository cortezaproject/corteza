package city311

import (
	"context"
	"fmt"
	"strconv"
	"time"

	composeTypes "github.com/cortezaproject/corteza/server/compose/types"
	contract "github.com/cortezaproject/corteza/server/compose/types/city311"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/store"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
)

var seedStatuses = []contract.ServiceRequestStatus{
	contract.ServiceRequestStatusDraft,
	contract.ServiceRequestStatusSubmitted,
	contract.ServiceRequestStatusTriaged,
	contract.ServiceRequestStatusAssigned,
	contract.ServiceRequestStatusInProgress,
	contract.ServiceRequestStatusResolved,
	contract.ServiceRequestStatusClosed,
	contract.ServiceRequestStatusReopened,
}

type seedActor struct {
	Handle     string
	Name       string
	Email      string
	Role       contract.ApplicationRole
	Department contract.DepartmentCode
	Districts  []contract.DistrictCode
}

var seedActors = []seedActor{
	{Handle: "city311-constituent", Name: "City 311 Constituent", Email: "constituent1@city311.example.invalid", Role: contract.ApplicationRoleConstituent},
	{Handle: "city311-constituent-two", Name: "City 311 Constituent Two", Email: "constituent2@city311.example.invalid", Role: contract.ApplicationRoleConstituent},
	{Handle: "city311-service-agent", Name: "City 311 Service Agent", Email: "service-agent@city311.example.invalid", Role: contract.ApplicationRoleServiceAgent, Department: contract.DepartmentStreets, Districts: []contract.DistrictCode{contract.DistrictNorth}},
	{Handle: "city311-supervisor", Name: "City 311 Supervisor", Email: "supervisor@city311.example.invalid", Role: contract.ApplicationRoleSupervisor, Department: contract.DepartmentStreets, Districts: []contract.DistrictCode{contract.DistrictNorth, contract.DistrictCentral}},
	{Handle: "city311-department-manager", Name: "City 311 Department Manager", Email: "department-manager@city311.example.invalid", Role: contract.ApplicationRoleDepartmentManager, Department: contract.DepartmentStreets, Districts: []contract.DistrictCode{contract.DistrictNorth, contract.DistrictCentral, contract.DistrictSouth}},
	{Handle: "city311-platform-administrator", Name: "City 311 Platform Administrator", Email: "platform-admin@city311.example.invalid", Role: contract.ApplicationRolePlatformAdministrator},
	{Handle: "city311-workflow-designer", Name: "City 311 Workflow Designer", Email: "workflow-designer@city311.example.invalid", Role: contract.ApplicationRoleWorkflowDesigner, Department: contract.DepartmentGeneralServices, Districts: []contract.DistrictCode{contract.DistrictCentral}},
	{Handle: "city311-integration-client", Name: "City 311 Integration Client", Email: "integration-client@city311.example.invalid", Role: contract.ApplicationRoleIntegrationClient},
}

// Seed installs the deterministic public dataset without modifying any row that
// already exists. It is safe to run on clean, migrated, and previously seeded databases.
func (svc *Service) Seed(ctx context.Context, benchmarkNow time.Time) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, tx store.Storer) error {
		if err := seedSequence(ctx, tx); err != nil {
			return fmt.Errorf("seed request-number sequence: %w", err)
		}
		if err := svc.seedActors(ctx, tx, benchmarkNow); err != nil {
			return fmt.Errorf("seed application actors: %w", err)
		}
		if err := svc.seedRequests(ctx, tx, benchmarkNow); err != nil {
			return fmt.Errorf("seed service requests: %w", err)
		}
		return nil
	})
}

func seedSequence(ctx context.Context, tx store.Storer) error {
	_, err := store.LookupCity311RequestSequenceByID(ctx, tx, 2026)
	if err == nil {
		return nil
	}
	if !errors.IsNotFound(err) {
		return err
	}
	return store.CreateCity311RequestSequence(ctx, tx, &composeTypes.City311RequestSequence{ID: 2026, NextNumber: 41})
}

func (svc *Service) seedActors(ctx context.Context, tx store.Storer, createdAt time.Time) error {
	for _, item := range seedActors {
		if err := svc.seedActorRecord(ctx, tx, item, createdAt); err != nil {
			return err
		}
	}
	return nil
}

func (svc *Service) seedActorRecord(ctx context.Context, tx store.Storer, item seedActor, createdAt time.Time) error {
	role, err := svc.findOrCreateSeedRole(ctx, tx, item, createdAt)
	if err != nil {
		return err
	}
	user, err := svc.findOrCreateSeedUser(ctx, tx, item, createdAt)
	if err != nil {
		return err
	}
	if err = ensureSeedRoleMembership(ctx, tx, item.Handle, role.ID, user.ID); err != nil {
		return err
	}
	return svc.ensureSeedActorProfile(ctx, tx, item, user.ID, createdAt)
}

func (svc *Service) findOrCreateSeedRole(ctx context.Context, tx store.Storer, item seedActor, createdAt time.Time) (*systemTypes.Role, error) {
	roleHandle := "city311-" + string(item.Role)
	role, err := store.LookupRoleByHandle(ctx, tx, roleHandle)
	if err == nil {
		return role, nil
	}
	if !errors.IsNotFound(err) {
		return nil, fmt.Errorf("lookup role for %s: %w", item.Handle, err)
	}
	role = &systemTypes.Role{ID: svc.nextID(), Handle: roleHandle, Name: item.Name, Meta: &systemTypes.RoleMeta{Description: "City 311 application role " + string(item.Role)}, CreatedAt: createdAt}
	if err = store.CreateRole(ctx, tx, role); err != nil {
		return nil, fmt.Errorf("create role for %s: %w", item.Handle, err)
	}
	return role, nil
}

func (svc *Service) findOrCreateSeedUser(ctx context.Context, tx store.Storer, item seedActor, createdAt time.Time) (*systemTypes.User, error) {
	user, err := store.LookupUserByEmail(ctx, tx, item.Email)
	if err == nil {
		return user, nil
	}
	if !errors.IsNotFound(err) {
		return nil, fmt.Errorf("lookup user for %s: %w", item.Handle, err)
	}
	user = &systemTypes.User{ID: svc.nextID(), Handle: item.Handle, Username: item.Handle, Email: item.Email, Name: item.Name, EmailConfirmed: true, Meta: &systemTypes.UserMeta{PreferredLanguage: "en"}, CreatedAt: createdAt}
	if err = store.CreateUser(ctx, tx, user); err != nil {
		return nil, fmt.Errorf("create user for %s: %w", item.Handle, err)
	}
	return user, nil
}

func ensureSeedRoleMembership(ctx context.Context, tx store.Storer, handle string, roleID, userID uint64) error {
	membership := &systemTypes.RoleMember{RoleID: roleID, Resource: fmt.Sprintf("corteza::system:user/%d", userID)}
	memberships, _, err := store.SearchRoleMembers(ctx, tx, systemTypes.RoleMemberFilter{RoleID: roleID, Resource: membership.Resource, Limit: 1})
	if err != nil {
		return fmt.Errorf("search role membership for %s: %w", handle, err)
	}
	if len(memberships) > 0 {
		return nil
	}
	if err = store.CreateRoleMember(ctx, tx, membership); err != nil {
		return fmt.Errorf("create role membership for %s: %w", handle, err)
	}
	return nil
}

func (svc *Service) ensureSeedActorProfile(ctx context.Context, tx store.Storer, item seedActor, userID uint64, createdAt time.Time) error {
	_, err := store.LookupCity311ActorProfileByID(ctx, tx, userID)
	if err == nil {
		return nil
	}
	if !errors.IsNotFound(err) {
		return fmt.Errorf("lookup actor profile for %s: %w", item.Handle, err)
	}
	err = store.CreateCity311ActorProfile(ctx, tx, &composeTypes.City311ActorProfile{
		ID: userID, ApplicationRoles: composeTypes.City311ApplicationRoleSet{item.Role}, Department: item.Department,
		Districts: composeTypes.City311DistrictCodeSet(item.Districts), CreatedAt: createdAt, UpdatedAt: createdAt,
	})
	if err != nil {
		return fmt.Errorf("create actor profile for %s: %w", item.Handle, err)
	}
	return nil
}

func (svc *Service) seedRequests(ctx context.Context, tx store.Storer, benchmarkNow time.Time) error {
	latitude, longitude := 42.8865, -78.8784
	for index, status := range seedStatuses {
		if err := svc.seedRequestRecord(ctx, tx, benchmarkNow, index, status, latitude, longitude); err != nil {
			return err
		}
	}
	return nil
}

func (svc *Service) seedRequestRecord(ctx context.Context, tx store.Storer, benchmarkNow time.Time, index int, status contract.ServiceRequestStatus, latitude, longitude float64) error {
	request, err := svc.findOrCreateSeedRequest(ctx, tx, benchmarkNow, index, status, latitude, longitude)
	if err != nil {
		return err
	}
	if err = svc.seedConstituent(ctx, tx, request); err != nil {
		return err
	}
	if err = svc.ensureSeedAudit(ctx, tx, request); err != nil {
		return err
	}
	return svc.ensureSeedHistory(ctx, tx, request)
}

func (svc *Service) findOrCreateSeedRequest(ctx context.Context, tx store.Storer, benchmarkNow time.Time, index int, status contract.ServiceRequestStatus, latitude, longitude float64) (*composeTypes.City311ServiceRequest, error) {
	requestNumber := fmt.Sprintf("SR-2026-%05d", 33+index)
	request, err := store.LookupCity311ServiceRequestByRequestNumber(ctx, tx, requestNumber)
	if err == nil {
		return request, nil
	}
	if !errors.IsNotFound(err) {
		return nil, err
	}
	createdAt := benchmarkNow.Add(-time.Duration(len(seedStatuses)-index) * time.Hour)
	request = &composeTypes.City311ServiceRequest{
		ID: svc.nextID(), RequestNumber: requestNumber, Summary: "Seeded City 311 request " + strconv.Itoa(index+1),
		Description: "Deterministic seeded request retained for brownfield and lifecycle regression checks.",
		ServiceType: contract.ServiceTypePothole, OwningDepartment: contract.DepartmentStreets,
		CouncilDistrict: contract.DistrictNorth, SourceChannel: contract.SourceChannelStaffInPerson,
		OriginClass: contract.OriginClassInternal, Status: status,
		PrimaryRequester: requesterMap(uint64(index+1), contract.RequesterInput{DisplayName: "Seed Constituent", Email: "constituent1@city311.example.invalid"}),
		Location:         locationMap(&contract.LocationInput{Address: "100 Example Street, Buffalo, NY 14201", Latitude: &latitude, Longitude: &longitude}),
		CustomFields:     map[string]any{}, CollaboratorIDs: composeTypes.City311Uint64Set{}, Version: 1,
		CreatedAt: createdAt, UpdatedAt: createdAt,
	}
	return request, store.CreateCity311ServiceRequest(ctx, tx, request)
}

func (svc *Service) ensureSeedAudit(ctx context.Context, tx store.Storer, request *composeTypes.City311ServiceRequest) error {
	audits, _, err := store.SearchCity311AuditEvents(ctx, tx, composeTypes.City311AuditEventFilter{RequestID: request.ID, EventType: "SEED_CREATED"})
	if err != nil {
		return err
	}
	if len(audits) > 0 {
		return nil
	}
	return store.CreateCity311AuditEvent(ctx, tx, &composeTypes.City311AuditEvent{
		ID: svc.nextID(), RequestID: request.ID, EventType: "SEED_CREATED", ActorType: contract.AuditActorSystem,
		SourceChannel: contract.SourceChannelStaffInPerson, Before: map[string]any{}, After: requestSnapshot(request), CreatedAt: request.CreatedAt,
	})
}

func (svc *Service) ensureSeedHistory(ctx context.Context, tx store.Storer, request *composeTypes.City311ServiceRequest) error {
	history, _, err := store.SearchCity311PublicHistoryItems(ctx, tx, composeTypes.City311PublicHistoryItemFilter{RequestID: request.ID})
	if err != nil {
		return err
	}
	if len(history) > 0 {
		return nil
	}
	return store.CreateCity311PublicHistoryItem(ctx, tx, &composeTypes.City311PublicHistoryItem{
		ID: svc.nextID(), RequestID: request.ID, Action: string(request.Status),
		ResponsibleDepartment: request.OwningDepartment, OccurredAt: request.CreatedAt,
	})
}

func (svc *Service) seedConstituent(ctx context.Context, tx store.Storer, request *composeTypes.City311ServiceRequest) error {
	constituentID := fmt.Sprint(request.PrimaryRequester["constituent_id"])
	if constituentID == "" || constituentID == "<nil>" {
		return nil
	}
	_, err := store.LookupCity311ConstituentByConstituentID(ctx, tx, constituentID)
	if err == nil {
		return nil
	}
	if !errors.IsNotFound(err) {
		return err
	}
	return store.CreateCity311Constituent(ctx, tx, &composeTypes.City311Constituent{
		ID: svc.nextID(), ConstituentID: constituentID, Profile: cloneMap(request.PrimaryRequester),
		OwningDepartment: request.OwningDepartment, CouncilDistrict: request.CouncilDistrict,
		CreatedAt: request.CreatedAt, UpdatedAt: request.CreatedAt,
	})
}
