package city311

import "time"

var fixtureCreatedAt = time.Date(2026, time.August, 24, 12, 0, 0, 0, time.UTC)

func MockCreatedServiceRequest() ServiceRequestResponse {
	return ServiceRequestResponse{
		RequestID:     "case-7c58d2",
		RequestNumber: "SR-2026-00041",
		Status:        ServiceRequestStatusSubmitted,
		Version:       1,
		CreatedAt:     fixtureCreatedAt,
		Links: ResourceLinks{
			Self: "/api/v1/service-requests/case-7c58d2",
		},
	}
}

func MockValidationError() APIError {
	return APIError{
		Error:     ErrorValidation,
		Message:   "The request contains invalid fields.",
		Retryable: false,
		Errors: []FieldError{
			{Field: "/summary", Code: ValidationRequired},
		},
	}
}

func MockIdempotencyConflict() APIError {
	return APIError{
		Error:     ErrorIdempotencyConflict,
		Message:   "The idempotency key has already been used with different content.",
		Retryable: false,
	}
}

func MockUnauthenticated() APIError {
	return APIError{
		Error:     ErrorUnauthenticated,
		Message:   "Authentication is required.",
		Retryable: false,
	}
}

func MockForbidden() APIError {
	return APIError{
		Error:     ErrorForbidden,
		Message:   "The authenticated actor is not permitted to perform this operation.",
		Retryable: false,
	}
}

func MockNotFound() APIError {
	return APIError{
		Error:     ErrorNotFound,
		Message:   "The requested resource was not found.",
		Retryable: false,
	}
}

func MockRateLimited() APIError {
	return APIError{
		Error:     ErrorRateLimited,
		Message:   "The client request limit has been exceeded.",
		Retryable: true,
	}
}

func MockInvalidResetToken() APIError {
	return APIError{
		Error:     ErrorInvalidResetToken,
		Message:   "The reset token is invalid.",
		Retryable: false,
	}
}

func MockExpiredResetToken() APIError {
	return APIError{
		Error:     ErrorExpiredResetToken,
		Message:   "The reset token has expired.",
		Retryable: false,
	}
}

func MockWorkflowFailure(code ErrorCode, retryable bool) APIError {
	return APIError{
		Error:     code,
		Message:   "The authenticated workflow action failed.",
		Retryable: retryable,
	}
}

func MockVersionConflict(currentVersion uint64) APIError {
	return APIError{
		Error:          ErrorVersionConflict,
		Message:        "The resource was updated by another operation.",
		Retryable:      false,
		CurrentVersion: &currentVersion,
	}
}

func MockGeocodeSuccess() GeocodeResponse {
	return GeocodeResponse{
		Address:         "100 Example Street, Buffalo, NY 14201",
		Latitude:        42.9001,
		Longitude:       -78.8801,
		PrecisionDigits: 4,
		Provider:        "BENCHMARK_MAP",
	}
}

func MockGeocodeNotFound() APIError {
	return APIError{
		Error:     ErrorAddressNotFound,
		Message:   "The address was not found.",
		Retryable: false,
	}
}

func MockGeocodeUnavailable() APIError {
	return APIError{
		Error:     ErrorMapTemporarilyUnavailable,
		Message:   "The mapping service is temporarily unavailable.",
		Retryable: true,
	}
}

func MockCivicWorksCreated() CivicWorksWorkOrder {
	createdAt := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	return CivicWorksWorkOrder{
		WorkOrderID:          "WO-000001",
		SourceCaseID:         "case-7c58d2",
		ServiceRequestNumber: "SR-2026-00041",
		Status:               CivicWorksStatusAssigned,
		ExternalStatusURL:    "http://civicworks:8080/ui/work-orders/WO-000001",
		Version:              1,
		CreatedAt:            createdAt,
		UpdatedAt:            createdAt,
	}
}

func MockCivicWorksCompletedEvent() CivicWorksEvent {
	return CivicWorksEvent{
		EventID:        "EVT-000031",
		EventType:      "work_order.status_changed",
		WorkOrderID:    "WO-000001",
		SourceCaseID:   "case-7c58d2",
		PreviousStatus: CivicWorksStatusInProgress,
		Status:         CivicWorksStatusCompleted,
		Version:        3,
		OccurredAt:     time.Date(2026, time.August, 20, 10, 30, 0, 0, time.UTC),
	}
}

func MockAnonymousStatusFound() AnonymousStatusLookupResponse {
	return AnonymousStatusLookupResponse{
		RequestDetail: &PublicServiceRequestDetail{
			RequestNumber:    "SR-2026-00041",
			Summary:          "Pothole blocking the eastbound lane",
			ServiceType:      ServiceTypePothole,
			Status:           ServiceRequestStatusInProgress,
			OwningDepartment: DepartmentStreets,
			UpdatedAt:        fixtureCreatedAt.Add(2 * time.Hour),
			History: []PublicHistoryItem{
				{
					Action:                "Request submitted",
					OccurredAt:            fixtureCreatedAt,
					ResponsibleDepartment: string(DepartmentStreets),
				},
			},
		},
	}
}

func MockAnonymousStatusNotFound() AnonymousStatusLookupResponse {
	return AnonymousStatusLookupResponse{RequestDetail: nil}
}

func MockPasswordResetRequested() PasswordResetResponse {
	return PasswordResetResponse{Message: "If the account is eligible, a reset link has been sent."}
}
