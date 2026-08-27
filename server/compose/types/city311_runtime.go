package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	city311Types "github.com/cortezaproject/corteza/server/compose/types/city311"
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	ServiceType          = city311Types.ServiceType
	DepartmentCode       = city311Types.DepartmentCode
	DistrictCode         = city311Types.DistrictCode
	SourceChannel        = city311Types.SourceChannel
	OriginClass          = city311Types.OriginClass
	ServiceRequestStatus = city311Types.ServiceRequestStatus
	AuditActorType       = city311Types.AuditActorType

	City311Uint64Set          []uint64
	City311ApplicationRoleSet []city311Types.ApplicationRole
	City311DistrictCodeSet    []city311Types.DistrictCode
	City311JSON               map[string]any

	City311ServiceRequest struct {
		ID                uint64                            `json:"requestID,string"`
		RequestNumber     string                            `json:"requestNumber"`
		Summary           string                            `json:"summary"`
		Description       string                            `json:"description"`
		ServiceType       city311Types.ServiceType          `json:"serviceType"`
		OwningDepartment  city311Types.DepartmentCode       `json:"owningDepartment"`
		CouncilDistrict   city311Types.DistrictCode         `json:"councilDistrict,omitempty"`
		SourceChannel     city311Types.SourceChannel        `json:"sourceChannel"`
		OriginClass       city311Types.OriginClass          `json:"originClass"`
		Status            city311Types.ServiceRequestStatus `json:"status"`
		PrimaryRequester  City311JSON                       `json:"primaryRequester"`
		Location          City311JSON                       `json:"location"`
		CustomFields      City311JSON                       `json:"customFields"`
		PrimaryAssigneeID uint64                            `json:"primaryAssigneeID,string,omitempty"`
		CollaboratorIDs   City311Uint64Set                  `json:"collaboratorIDs"`
		DuplicateGroupID  string                            `json:"duplicateGroupID,omitempty"`
		Version           int                               `json:"version"`
		CreatedAt         time.Time                         `json:"createdAt"`
		UpdatedAt         time.Time                         `json:"updatedAt"`
	}

	City311ServiceRequestFilter struct {
		RequestNumber     string
		ServiceType       string
		OwningDepartment  string
		CouncilDistrict   string
		SourceChannel     string
		OriginClass       string
		Status            string
		PrimaryAssigneeID uint64
		Check             func(*City311ServiceRequest) (bool, error)
		filter.Sorting
		filter.Paging
	}
	City311ServiceRequestSet []*City311ServiceRequest

	City311RequestSequence struct {
		ID         uint64 `json:"year"`
		NextNumber uint64 `json:"nextNumber"`
	}
	typeCity311RequestSequenceFilter struct {
		Check func(*City311RequestSequence) (bool, error)
		filter.Sorting
		filter.Paging
	}
	City311RequestSequenceFilter = typeCity311RequestSequenceFilter
	City311RequestSequenceSet    []*City311RequestSequence

	City311IdempotencyRecord struct {
		ID             uint64      `json:"id,string"`
		Operation      string      `json:"operation"`
		KeyHash        string      `json:"keyHash"`
		RequestHash    string      `json:"requestHash"`
		ResponseStatus int         `json:"responseStatus"`
		ResponseBody   City311JSON `json:"responseBody"`
		RequestID      uint64      `json:"requestID,string"`
		CreatedAt      time.Time   `json:"createdAt"`
		ExpiresAt      time.Time   `json:"expiresAt"`
	}
	City311IdempotencyRecordFilter struct {
		Operation string
		KeyHash   string
		Check     func(*City311IdempotencyRecord) (bool, error)
		filter.Sorting
		filter.Paging
	}
	City311IdempotencyRecordSet []*City311IdempotencyRecord

	City311AuditEvent struct {
		ID            uint64                      `json:"id,string"`
		RequestID     uint64                      `json:"requestID,string"`
		EventType     string                      `json:"eventType"`
		ActorType     city311Types.AuditActorType `json:"actorType"`
		ActorID       uint64                      `json:"actorID,string"`
		SourceChannel city311Types.SourceChannel  `json:"sourceChannel"`
		Before        City311JSON                 `json:"before"`
		After         City311JSON                 `json:"after"`
		CreatedAt     time.Time                   `json:"createdAt"`
	}
	City311AuditEventFilter struct {
		RequestID uint64
		EventType string
		Check     func(*City311AuditEvent) (bool, error)
		filter.Sorting
		filter.Paging
	}
	City311AuditEventSet []*City311AuditEvent

	City311RequestAttachment struct {
		ID        uint64    `json:"id,string"`
		RequestID uint64    `json:"requestID,string"`
		Filename  string    `json:"filename"`
		MediaType string    `json:"mediaType"`
		Size      uint64    `json:"size"`
		Content   []byte    `json:"-"`
		CreatedAt time.Time `json:"createdAt"`
	}
	City311RequestAttachmentFilter struct {
		RequestID uint64
		Check     func(*City311RequestAttachment) (bool, error)
		filter.Sorting
		filter.Paging
	}
	City311RequestAttachmentSet []*City311RequestAttachment

	City311PublicHistoryItem struct {
		ID                    uint64                      `json:"id,string"`
		RequestID             uint64                      `json:"requestID,string"`
		Action                string                      `json:"action"`
		ResponsibleDepartment city311Types.DepartmentCode `json:"responsibleDepartment"`
		OccurredAt            time.Time                   `json:"occurredAt"`
	}
	City311PublicHistoryItemFilter struct {
		RequestID uint64
		Check     func(*City311PublicHistoryItem) (bool, error)
		filter.Sorting
		filter.Paging
	}
	City311PublicHistoryItemSet []*City311PublicHistoryItem

	City311ActorProfile struct {
		ID               uint64                      `json:"userID,string"`
		ApplicationRoles City311ApplicationRoleSet   `json:"applicationRoles"`
		Department       city311Types.DepartmentCode `json:"department,omitempty"`
		Districts        City311DistrictCodeSet      `json:"districts"`
		CreatedAt        time.Time                   `json:"createdAt"`
		UpdatedAt        time.Time                   `json:"updatedAt"`
	}
	City311ActorProfileFilter struct {
		Department string
		Check      func(*City311ActorProfile) (bool, error)
		filter.Sorting
		filter.Paging
	}
	City311ActorProfileSet []*City311ActorProfile
)

func scanCity311JSON(src, dst any) error {
	if src == nil {
		return nil
	}
	switch value := src.(type) {
	case []byte:
		return json.Unmarshal(value, dst)
	case string:
		return json.Unmarshal([]byte(value), dst)
	default:
		encoded, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return json.Unmarshal(encoded, dst)
	}
}

func (set *City311Uint64Set) Scan(src any) error                   { return scanCity311JSON(src, set) }
func (set City311Uint64Set) Value() (driver.Value, error)          { return json.Marshal(set) }
func (set *City311ApplicationRoleSet) Scan(src any) error          { return scanCity311JSON(src, set) }
func (set City311ApplicationRoleSet) Value() (driver.Value, error) { return json.Marshal(set) }
func (set *City311DistrictCodeSet) Scan(src any) error             { return scanCity311JSON(src, set) }
func (set City311DistrictCodeSet) Value() (driver.Value, error)    { return json.Marshal(set) }
func (value *City311JSON) Scan(src any) error                      { return scanCity311JSON(src, value) }
func (value City311JSON) Value() (driver.Value, error)             { return json.Marshal(value) }
