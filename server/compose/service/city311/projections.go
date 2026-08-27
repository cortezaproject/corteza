package city311

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	composeTypes "github.com/cortezaproject/corteza/server/compose/types"
	contract "github.com/cortezaproject/corteza/server/compose/types/city311"
)

func cloneMap(input map[string]any) map[string]any {
	if input == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(input))
	for key, value := range input {
		out[key] = value
	}
	return out
}

func mapFrom(value any) (map[string]any, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	if err = json.Unmarshal(encoded, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func requesterMap(requestID uint64, requester contract.RequesterInput) map[string]any {
	phones := []contract.PhoneNumber{}
	if strings.TrimSpace(requester.Phone) != "" {
		phones = append(phones, contract.PhoneNumber{Label: contract.PhoneLabelMobile, Value: strings.TrimSpace(requester.Phone)})
	}
	value := contract.Constituent{
		ConstituentID: "C-" + strconv.FormatUint(requestID, 10), DisplayName: strings.TrimSpace(requester.DisplayName),
		Emails: []string{strings.ToLower(strings.TrimSpace(requester.Email))}, PhoneNumbers: phones,
		Addresses: []contract.Address{}, PrimaryCategory: contract.ContactCategoryResident,
		PreferredLanguage: contract.LanguageEN, EmailOptOut: false,
	}
	out, _ := mapFrom(value)
	return out
}

func requesterInput(value map[string]any) contract.RequesterInput {
	constituent := contract.Constituent{}
	if encoded, err := json.Marshal(value); err == nil {
		_ = json.Unmarshal(encoded, &constituent)
	}
	requester := contract.RequesterInput{DisplayName: constituent.DisplayName}
	if len(constituent.Emails) > 0 {
		requester.Email = constituent.Emails[0]
	}
	if len(constituent.PhoneNumbers) > 0 {
		requester.Phone = constituent.PhoneNumbers[0].Value
	}
	return requester
}

func locationMap(input *contract.LocationInput) map[string]any {
	if input == nil {
		return map[string]any{}
	}
	address := strings.TrimSpace(input.Address)
	postalCode := "14201"
	parts := strings.Fields(address)
	if len(parts) > 0 && len(parts[len(parts)-1]) >= 5 {
		postalCode = parts[len(parts)-1]
	}
	latitude, longitude := input.Latitude, input.Longitude
	if latitude != nil {
		rounded := math.Round(*latitude*10_000) / 10_000
		latitude = &rounded
	}
	if longitude != nil {
		rounded := math.Round(*longitude*10_000) / 10_000
		longitude = &rounded
	}
	value := contract.ServiceRequestLocation{
		Address:  contract.Address{Line1: address, City: "Buffalo", Region: "NY", PostalCode: postalCode, Country: "US", Primary: true},
		Latitude: latitude, Longitude: longitude,
	}
	out, _ := mapFrom(value)
	return out
}

func requestSnapshot(request *composeTypes.City311ServiceRequest) map[string]any {
	out, _ := mapFrom(toContract(request))
	return out
}

func responseFor(request *composeTypes.City311ServiceRequest) *contract.ServiceRequestResponse {
	return &contract.ServiceRequestResponse{
		RequestID: strconv.FormatUint(request.ID, 10), RequestNumber: request.RequestNumber,
		Status: request.Status, Version: uint64(request.Version), CreatedAt: request.CreatedAt,
		Links: contract.ResourceLinks{Self: "/api/v1/staff/service-requests/" + strconv.FormatUint(request.ID, 10)},
	}
}

func toContract(request *composeTypes.City311ServiceRequest) contract.ServiceRequest {
	requester := contract.Constituent{}
	if encoded, err := json.Marshal(request.PrimaryRequester); err == nil {
		_ = json.Unmarshal(encoded, &requester)
	}
	var location *contract.ServiceRequestLocation
	if len(request.Location) > 0 {
		value := contract.ServiceRequestLocation{}
		if encoded, err := json.Marshal(request.Location); err == nil && json.Unmarshal(encoded, &value) == nil {
			location = &value
		}
	}
	var district *contract.DistrictCode
	if request.CouncilDistrict != "" {
		value := request.CouncilDistrict
		district = &value
	}
	return contract.ServiceRequest{
		RequestID: strconv.FormatUint(request.ID, 10), RequestNumber: request.RequestNumber,
		Summary: request.Summary, Description: request.Description, ServiceType: request.ServiceType,
		OwningDepartment: request.OwningDepartment, CouncilDistrict: district, SourceChannel: request.SourceChannel,
		OriginClass: request.OriginClass, Status: request.Status, PrimaryRequester: requester, Location: location,
		CustomFields: cloneMap(request.CustomFields), DuplicateGroupID: request.DuplicateGroupID,
		Version: uint64(request.Version), CreatedAt: request.CreatedAt, UpdatedAt: request.UpdatedAt,
	}
}

func optionalID(value uint64) *string {
	if value == 0 {
		return nil
	}
	formatted := strconv.FormatUint(value, 10)
	return &formatted
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func stringifyIDs(values []uint64) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, strconv.FormatUint(value, 10))
	}
	return out
}

func queueItem(actor contract.Actor, request *composeTypes.City311ServiceRequest) contract.RequestQueueItem {
	return contract.RequestQueueItem{
		RequestID: strconv.FormatUint(request.ID, 10), RequestNumber: request.RequestNumber, Summary: request.Summary,
		ServiceType: request.ServiceType, Status: request.Status, OwningDepartment: request.OwningDepartment,
		CouncilDistrict: request.CouncilDistrict, OriginClass: request.OriginClass, SourceChannel: request.SourceChannel,
		PrimaryAssigneeID: optionalID(request.PrimaryAssigneeID), DuplicateGroupID: optionalString(request.DuplicateGroupID),
		Version: uint64(request.Version), UpdatedAt: request.UpdatedAt, AvailableActions: availableActions(actor, request),
	}
}

func availableActions(actor contract.Actor, request *composeTypes.City311ServiceRequest) []string {
	if !canRead(actor, request) || !canOperateRequest(actor) {
		return []string{}
	}
	switch request.Status {
	case contract.ServiceRequestStatusSubmitted:
		return []string{"TRIAGE"}
	case contract.ServiceRequestStatusTriaged:
		return []string{"ASSIGN"}
	case contract.ServiceRequestStatusAssigned, contract.ServiceRequestStatusReopened:
		return []string{"START_PROGRESS"}
	case contract.ServiceRequestStatusInProgress:
		return []string{"RESOLVE"}
	case contract.ServiceRequestStatusResolved:
		return []string{"CLOSE", "REQUEST_REOPEN"}
	case contract.ServiceRequestStatusClosed:
		return []string{"REQUEST_REOPEN"}
	default:
		return []string{}
	}
}
