package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type (
	// NgAutomationIssueDetail wraps exactly one variant payload and marshals it as
	// {type, parameters}. The wrapper is the only thing that injects the discriminator,
	// so details never serialize as a bare variant struct.
	NgAutomationIssueDetail struct {
		MissingReference *DetailMissingReference
		InvalidType      *DetailInvalidType
		DuplicateID      *DetailDuplicateID
		Cycle            *DetailCycle
		ResourceRef      *DetailResourceRef
		GatewayPaths     *DetailGatewayPaths
		EmptyField       *DetailEmptyField

		Details []*NgAutomationIssueDetail // reserved; nil for now
	}

	DetailMissingReference struct {
		RefKind    string `json:"refKind"` // scope|step|trigger|function
		Ref        string `json:"ref"`
		StepID     uint64 `json:"stepID,string,omitempty"`
		Field      string `json:"field,omitempty"` // arguments|results
		FieldIndex int    `json:"fieldIndex,omitempty"`
	}

	DetailInvalidType struct {
		StepID     uint64 `json:"stepID,string"`
		Field      string `json:"field,omitempty"`
		FieldIndex int    `json:"fieldIndex,omitempty"`
		Target     string `json:"target,omitempty"`
		Expected   string `json:"expected"`
		Actual     string `json:"actual"`
	}

	DetailDuplicateID struct {
		Resource string `json:"resource"` // step|trigger
		ID       uint64 `json:"id,string"`
		Indices  []int  `json:"indices,omitempty"`
	}

	DetailCycle struct {
		StepIDs IssueIDs `json:"stepIDs"` // serialized as ["3","5"]
	}

	DetailResourceRef struct {
		Resource   string `json:"resource"` // step|trigger|path
		ID         uint64 `json:"id,string,omitempty"`
		Index      int    `json:"index,omitempty"`
		Field      string `json:"field,omitempty"`
		FieldIndex int    `json:"fieldIndex,omitempty"`
	}

	DetailGatewayPaths struct {
		StepID    uint64 `json:"stepID,string"`
		Violation string `json:"violation"` // tooFew|noElse|multipleElse
		Got       int    `json:"got"`
		Want      int    `json:"want,omitempty"`
	}

	DetailEmptyField struct {
		Resource string `json:"resource"` // step|trigger|path
		Index    int    `json:"index"`
		Field    string `json:"field,omitempty"`
	}

	issueDetailPayload interface{ issueDetailType() string }
)

func (*DetailMissingReference) issueDetailType() string { return IssueDetailMissingReference }
func (*DetailInvalidType) issueDetailType() string      { return IssueDetailInvalidType }
func (*DetailDuplicateID) issueDetailType() string      { return IssueDetailDuplicateID }
func (*DetailCycle) issueDetailType() string            { return IssueDetailCycle }
func (*DetailResourceRef) issueDetailType() string      { return IssueDetailResourceRef }
func (*DetailGatewayPaths) issueDetailType() string     { return IssueDetailGatewayPaths }
func (*DetailEmptyField) issueDetailType() string       { return IssueDetailEmptyField }

func (d NgAutomationIssueDetail) payload() (issueDetailPayload, error) {
	var (
		found issueDetailPayload
		n     int
	)
	set := func(p issueDetailPayload, ok bool) {
		if ok {
			found, n = p, n+1
		}
	}
	set(d.MissingReference, d.MissingReference != nil)
	set(d.InvalidType, d.InvalidType != nil)
	set(d.DuplicateID, d.DuplicateID != nil)
	set(d.Cycle, d.Cycle != nil)
	set(d.ResourceRef, d.ResourceRef != nil)
	set(d.GatewayPaths, d.GatewayPaths != nil)
	set(d.EmptyField, d.EmptyField != nil)

	if n != 1 {
		return nil, fmt.Errorf("issue detail must set exactly one variant, got %d", n)
	}
	return found, nil
}

func (d NgAutomationIssueDetail) MarshalJSON() ([]byte, error) {
	p, err := d.payload()
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Type       string                     `json:"type"`
		Parameters issueDetailPayload         `json:"parameters,omitempty"`
		Details    []*NgAutomationIssueDetail `json:"details,omitempty"`
	}{p.issueDetailType(), p, d.Details})
}

func (d *NgAutomationIssueDetail) UnmarshalJSON(b []byte) error {
	var wire struct {
		Type       string                     `json:"type"`
		Parameters json.RawMessage            `json:"parameters"`
		Details    []*NgAutomationIssueDetail `json:"details"`
	}
	if err := json.Unmarshal(b, &wire); err != nil {
		return err
	}
	d.Details = wire.Details

	var target issueDetailPayload
	switch wire.Type {
	case IssueDetailMissingReference:
		d.MissingReference = &DetailMissingReference{}
		target = d.MissingReference
	case IssueDetailInvalidType:
		d.InvalidType = &DetailInvalidType{}
		target = d.InvalidType
	case IssueDetailDuplicateID:
		d.DuplicateID = &DetailDuplicateID{}
		target = d.DuplicateID
	case IssueDetailCycle:
		d.Cycle = &DetailCycle{}
		target = d.Cycle
	case IssueDetailResourceRef:
		d.ResourceRef = &DetailResourceRef{}
		target = d.ResourceRef
	case IssueDetailGatewayPaths:
		d.GatewayPaths = &DetailGatewayPaths{}
		target = d.GatewayPaths
	case IssueDetailEmptyField:
		d.EmptyField = &DetailEmptyField{}
		target = d.EmptyField
	default:
		return fmt.Errorf("unknown issue detail type %q", wire.Type)
	}

	if len(wire.Parameters) == 0 {
		return nil
	}
	return json.Unmarshal(wire.Parameters, target)
}

// IssueIDs serializes a slice of uint64 IDs as JSON strings to survive JS number limits.
type IssueIDs []uint64

func (l IssueIDs) MarshalJSON() ([]byte, error) {
	out := make([]string, len(l))
	for i, v := range l {
		out[i] = strconv.FormatUint(v, 10)
	}
	return json.Marshal(out)
}

func (l *IssueIDs) UnmarshalJSON(b []byte) error {
	var ss []string
	if err := json.Unmarshal(b, &ss); err != nil {
		return err
	}
	*l = make(IssueIDs, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		(*l)[i] = v
	}
	return nil
}

func NewDetailMissingReference(refKind, ref string, stepID uint64, field string, fieldIndex int) *NgAutomationIssueDetail {
	return &NgAutomationIssueDetail{MissingReference: &DetailMissingReference{
		RefKind:    refKind,
		Ref:        ref,
		StepID:     stepID,
		Field:      field,
		FieldIndex: fieldIndex,
	}}
}

func NewDetailInvalidType(stepID uint64, field string, fieldIndex int, target, expected, actual string) *NgAutomationIssueDetail {
	return &NgAutomationIssueDetail{InvalidType: &DetailInvalidType{
		StepID:     stepID,
		Field:      field,
		FieldIndex: fieldIndex,
		Target:     target,
		Expected:   expected,
		Actual:     actual,
	}}
}

func NewDetailDuplicateID(resource string, id uint64, indices ...int) *NgAutomationIssueDetail {
	return &NgAutomationIssueDetail{DuplicateID: &DetailDuplicateID{
		Resource: resource,
		ID:       id,
		Indices:  indices,
	}}
}

func NewDetailCycle(stepIDs ...uint64) *NgAutomationIssueDetail {
	return &NgAutomationIssueDetail{Cycle: &DetailCycle{StepIDs: IssueIDs(stepIDs)}}
}

func NewDetailResourceRef(resource string, id uint64, index int) *NgAutomationIssueDetail {
	return &NgAutomationIssueDetail{ResourceRef: &DetailResourceRef{
		Resource: resource,
		ID:       id,
		Index:    index,
	}}
}

func NewDetailGatewayPaths(stepID uint64, violation string, got, want int) *NgAutomationIssueDetail {
	return &NgAutomationIssueDetail{GatewayPaths: &DetailGatewayPaths{
		StepID:    stepID,
		Violation: violation,
		Got:       got,
		Want:      want,
	}}
}

func NewDetailEmptyField(resource string, index int, field string) *NgAutomationIssueDetail {
	return &NgAutomationIssueDetail{EmptyField: &DetailEmptyField{
		Resource: resource,
		Index:    index,
		Field:    field,
	}}
}

// NewIssue builds an issue whose Message is derived from its code and details by
// issueMessage. Callers supply the machine-readable parts (code, severity, typed
// details); the human line is generated here, never at the call site.
func NewIssue(code, severity string, details ...*NgAutomationIssueDetail) *NgAutomationIssue {
	i := &NgAutomationIssue{Code: code, Severity: severity, Details: details}
	i.Message = issueMessage(code, details)
	return i
}

// issueMessage renders the human-readable line for an issue from its code and the
// data carried in its details. It is the single source of issue wording.
func issueMessage(code string, details []*NgAutomationIssueDetail) string {
	var d *NgAutomationIssueDetail
	if len(details) > 0 {
		d = details[0]
	}

	switch code {
	case IssueCodeAutomationNil:
		return "nil automation"
	case IssueCodeGraphNoEntry:
		return "no entry steps (every step has at least one parent)"
	case IssueCodeGraphCycle:
		if d != nil && d.Cycle != nil && len(d.Cycle.StepIDs) > 0 {
			return fmt.Sprintf("cycle detected: %s", joinIDs(d.Cycle.StepIDs))
		}
		return "cycle detected"
	case IssueCodeGraphAmbiguousEntry:
		return "cannot infer trigger-to-step entry paths: multiple unconnected triggers and entry steps"

	case IssueCodeTriggerEmptyID:
		return "trigger has empty ID"
	case IssueCodeTriggerDuplicateID:
		if d != nil && d.DuplicateID != nil {
			return fmt.Sprintf("duplicate trigger ID %d", d.DuplicateID.ID)
		}
		return "duplicate trigger ID"
	case IssueCodeTriggerMultiPaths:
		return "trigger has multiple outbound paths (only one allowed)"

	case IssueCodeStepEmptyID:
		return "step has empty ID"
	case IssueCodeStepDuplicateID:
		if d != nil && d.DuplicateID != nil {
			return fmt.Sprintf("duplicate step ID %d", d.DuplicateID.ID)
		}
		return "duplicate step ID"
	case IssueCodeStepIDCollision:
		if d != nil && d.ResourceRef != nil {
			return fmt.Sprintf("step ID %d collides with trigger ID", d.ResourceRef.ID)
		}
		return "step ID collides with trigger ID"

	case IssueCodeFunctionUnknown:
		if ref := missingRef(d); ref != "" {
			return fmt.Sprintf("unknown function %q", ref)
		}
		return "unknown function"

	case IssueCodePathEmpty:
		return "path has empty parent or child"
	case IssueCodePathSelfLoop:
		return "path is a self-loop"
	case IssueCodePathUnknownChild:
		if ref := missingRef(d); ref != "" {
			return fmt.Sprintf("unknown child step %s", ref)
		}
		return "unknown child step"
	case IssueCodePathUnknownParent:
		if ref := missingRef(d); ref != "" {
			return fmt.Sprintf("unknown parent step %s", ref)
		}
		return "unknown parent step"
	case IssueCodeInternal:
		return "internal step index missing"

	case IssueCodeScopeUnknown:
		if d != nil && d.MissingReference != nil {
			m := d.MissingReference
			return fmt.Sprintf("step %d %s reference unknown scope %q (referenced step/trigger may have been deleted)",
				m.StepID, m.Field, m.Ref)
		}
		return "reference to unknown scope"

	case IssueCodeGatewayTooFewPaths:
		if d != nil && d.GatewayPaths != nil {
			g := d.GatewayPaths
			return fmt.Sprintf("gateway step %d must have at least %d outbound paths, got %d", g.StepID, g.Want, g.Got)
		}
		return "gateway has too few outbound paths"
	case IssueCodeGatewayMultiElse:
		if d != nil && d.GatewayPaths != nil {
			g := d.GatewayPaths
			return fmt.Sprintf("gateway step %d has %d else paths (exactly one allowed)", g.StepID, g.Got)
		}
		return "gateway has multiple else paths"
	case IssueCodeGatewayNoElse:
		if d != nil && d.GatewayPaths != nil {
			return fmt.Sprintf("gateway step %d has no else path", d.GatewayPaths.StepID)
		}
		return "gateway has no else path"
	}

	return code
}

func missingRef(d *NgAutomationIssueDetail) string {
	if d == nil || d.MissingReference == nil {
		return ""
	}
	return d.MissingReference.Ref
}

func joinIDs(ids IssueIDs) string {
	out := make([]string, len(ids))
	for i, v := range ids {
		out[i] = strconv.FormatUint(v, 10)
	}
	return strings.Join(out, " → ")
}
