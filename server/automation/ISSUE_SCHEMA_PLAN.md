# NG Automation — Structured Issue Schema: Implementation Plan

## Goal

Replace the loosely-typed `NgAutomationIssue` with a descriptive, machine-consumable
schema so the automation builder UI can render validation problems richly (severity
grouping, click-to-locate, per-problem rendering) and remain forward-compatible.

This supersedes the interim `{Kind, Severity, Description, Culprit}` struct shipped
earlier on this branch — that gets replaced. Nothing released yet → breaking JSON
change is safe.

## Locked decisions

| # | Decision | Choice |
|---|---|---|
| A | IDs in JSON | string (uint64 overflows JS number) |
| B | `field` value naming | plural — `arguments` / `results` (matches struct) |
| C | Location data | embedded in each semantic detail (no separate culprit object) |
| D | `severity` levels | `error` / `warning` / `info` |
| E | `code` naming | dot-namespaced (`scope.unknown`, `gateway.noElse`) |
| F | `details` recursion | field present in struct, **always nil for now** (flat); expand later |
| G | detail params in Go | one typed struct per variant (no `map[string]any`); custom marshal packs to `{type, parameters}` |
| H | API typing | raw JSON must include `type` per detail. No OpenAPI / lib/js work in scope |

## Wire format (the contract)

```json
{
  "code": "scope.unknown",
  "severity": "error",
  "message": "Step \"step3\" argument 0 references scope \"step2\" which no longer exists.",
  "details": [
    { "type": "missingReference",
      "parameters": { "refKind": "scope", "ref": "step2", "stepID": "3", "field": "arguments", "fieldIndex": 0 } }
  ]
}
```

- Top level = identity (`code`), `severity`, system-generated `message`, `details[]`.
- Each detail = `{ type, parameters }`; `type` is the discriminator, `parameters` is the
  variant-specific payload. Optional `details[]` on a detail is reserved (always nil now).
- **Contract rule:** details always marshal through the `NgAutomationIssueDetail` wrapper,
  never a bare variant struct — the wrapper injects `type`. Constructors enforce this.

---

## Phase 1 — Issue type (`automation/types/ng_automation.go`)

Replace `NgAutomationIssue`; delete `NgAutomationIssueCulprit`.

```go
NgAutomationIssue struct {
    Code     string                     `json:"code"`
    Severity string                     `json:"severity"`
    Message  string                     `json:"message"`
    Details  []*NgAutomationIssueDetail `json:"details,omitempty"`
}
```

Keep severity consts; add code + detail-type const blocks (exported = FE contract + Go
compile safety):

```go
const (
    NgAutomationSeverityError   = "error"
    NgAutomationSeverityWarning = "warning"
    NgAutomationSeverityInfo    = "info"
)

const (
    IssueCodeScopeUnknown        = "scope.unknown"
    IssueCodeTypeMismatch        = "type.mismatch"
    IssueCodeStepDuplicateID     = "step.duplicateID"
    IssueCodeStepEmptyID         = "step.emptyID"
    IssueCodeStepIDCollision     = "step.idCollision"
    IssueCodeTriggerDuplicateID  = "trigger.duplicateID"
    IssueCodeTriggerEmptyID      = "trigger.emptyID"
    IssueCodeTriggerMultiPaths   = "trigger.multiplePaths"
    IssueCodeFunctionUnknown     = "function.unknown"
    IssueCodePathEmpty           = "path.empty"
    IssueCodePathSelfLoop        = "path.selfLoop"
    IssueCodePathUnknownChild    = "path.unknownChild"
    IssueCodePathUnknownParent   = "path.unknownParent"
    IssueCodeGatewayTooFewPaths  = "gateway.tooFewPaths"
    IssueCodeGatewayNoElse       = "gateway.noElse"
    IssueCodeGatewayMultiElse    = "gateway.multipleElse"
    IssueCodeGraphNoEntry        = "graph.noEntry"
    IssueCodeGraphCycle          = "graph.cycle"
    IssueCodeGraphAmbiguousEntry = "graph.ambiguousEntry"
    IssueCodeAutomationNil       = "automation.nil"
    IssueCodeInternal            = "internal"
    IssueCodeRunAsLoadFailed     = "runAs.loadFailed"
    IssueCodeRunAsInvalid        = "runAs.invalid"
)

const (
    IssueDetailMissingReference = "missingReference"
    IssueDetailInvalidType      = "invalidType"
    IssueDetailDuplicateID      = "duplicateID"
    IssueDetailCycle            = "cycle"
    IssueDetailResourceRef      = "resourceRef"
    IssueDetailGatewayPaths     = "gatewayPaths"
    IssueDetailEmptyField       = "emptyField"
)
```

`Error()` joins `Message`.

---

## Phase 2 — Detail variants + marshal (`automation/types/ng_automation_issue.go`, new)

### Wrapper + variant structs

```go
type NgAutomationIssueDetail struct {
    MissingReference *DetailMissingReference
    InvalidType      *DetailInvalidType
    DuplicateID      *DetailDuplicateID
    Cycle            *DetailCycle
    ResourceRef      *DetailResourceRef
    GatewayPaths     *DetailGatewayPaths
    EmptyField       *DetailEmptyField

    Details []*NgAutomationIssueDetail // reserved; nil for now
}

type DetailMissingReference struct {
    RefKind    string `json:"refKind"`           // scope|step|trigger|function
    Ref        string `json:"ref"`
    StepID     uint64 `json:"stepID,string,omitempty"`
    Field      string `json:"field,omitempty"`   // arguments|results
    FieldIndex int    `json:"fieldIndex,omitempty"`
}

type DetailInvalidType struct {
    StepID     uint64 `json:"stepID,string"`
    Field      string `json:"field,omitempty"`
    FieldIndex int    `json:"fieldIndex,omitempty"`
    Target     string `json:"target,omitempty"`
    Expected   string `json:"expected"`
    Actual     string `json:"actual"`
}

type DetailDuplicateID struct {
    Resource string `json:"resource"` // step|trigger
    ID       uint64 `json:"id,string"`
    Indices  []int  `json:"indices,omitempty"`
}

type DetailCycle struct {
    StepIDs IssueIDs `json:"stepIDs"` // serialized as ["3","5"]
}

type DetailResourceRef struct {
    Resource   string `json:"resource"` // step|trigger|path
    ID         uint64 `json:"id,string,omitempty"`
    Index      int    `json:"index,omitempty"`
    Field      string `json:"field,omitempty"`
    FieldIndex int    `json:"fieldIndex,omitempty"`
}

type DetailGatewayPaths struct {
    StepID    uint64 `json:"stepID,string"`
    Violation string `json:"violation"` // tooFew|noElse|multipleElse
    Got       int    `json:"got"`
    Want      int    `json:"want,omitempty"`
}

type DetailEmptyField struct {
    Resource string `json:"resource"` // step|trigger|path
    Index    int    `json:"index"`
    Field    string `json:"field,omitempty"`
}
```

### Discriminator — each variant defines its type

```go
type issueDetailPayload interface{ issueDetailType() string }

func (*DetailMissingReference) issueDetailType() string { return IssueDetailMissingReference }
func (*DetailInvalidType) issueDetailType() string      { return IssueDetailInvalidType }
func (*DetailDuplicateID) issueDetailType() string      { return IssueDetailDuplicateID }
func (*DetailCycle) issueDetailType() string            { return IssueDetailCycle }
func (*DetailResourceRef) issueDetailType() string      { return IssueDetailResourceRef }
func (*DetailGatewayPaths) issueDetailType() string     { return IssueDetailGatewayPaths }
func (*DetailEmptyField) issueDetailType() string       { return IssueDetailEmptyField }
```

### Marshal / Unmarshal — typed field ↔ `{type, parameters}`

```go
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
```

### `IssueIDs` — string-array IDs for slices

```go
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
```

### Constructors (return wrapper, set exactly one field)

```go
func NewDetailMissingReference(refKind, ref string, stepID uint64, field string, fieldIndex int) *NgAutomationIssueDetail
func NewDetailInvalidType(stepID uint64, field string, fieldIndex int, target, expected, actual string) *NgAutomationIssueDetail
func NewDetailDuplicateID(resource string, id uint64, indices ...int) *NgAutomationIssueDetail
func NewDetailCycle(stepIDs ...uint64) *NgAutomationIssueDetail
func NewDetailResourceRef(resource string, id uint64, index int) *NgAutomationIssueDetail
func NewDetailGatewayPaths(stepID uint64, violation string, got, want int) *NgAutomationIssueDetail
func NewDetailEmptyField(resource string, index int, field string) *NgAutomationIssueDetail
```

---

## Phase 3 — Converter rewrite (`automation/service/ng_automation_converter.go`)

New factory:

```go
func issue(code, severity, message string, details ...*automationTypes.NgAutomationIssueDetail) *automationTypes.NgAutomationIssue
```

Delete the old `atStep/atTrigger/atPath/atStepID` helpers (replaced by `types.NewDetail*`).

**Message generation lives in `types`, not the converter.** `types.NewIssue(code, severity, details...)`
sets `Message` from `issueMessage(code, details)` — a single per-issue-code generator that
reads the data it needs from the typed details. Call sites pass only machine parts (code,
severity, details); they never write message strings. The converter `issue()` helper is a thin
wrapper over `types.NewIssue`.

### Call-site mapping

| Current issue | code | severity | detail |
|---|---|---|---|
| nil automation | `automation.nil` | error | — |
| no entry | `graph.noEntry` | error | — |
| cycle | `graph.cycle` | error | `Cycle{stepIDs}` |
| trigger empty ID | `trigger.emptyID` | error | `EmptyField{trigger, index}` |
| duplicate trigger ID | `trigger.duplicateID` | error | `DuplicateID{trigger, id, indices}` |
| step empty ID | `step.emptyID` | error | `EmptyField{step, index}` |
| duplicate step ID | `step.duplicateID` | error | `DuplicateID{step, id, indices}` |
| step/trigger collision | `step.idCollision` | error | `ResourceRef{step, id, index}` |
| unknown function | `function.unknown` | error | `MissingReference{function, ref, stepID}` |
| path empty | `path.empty` | error | `ResourceRef{path, index}` |
| self-loop | `path.selfLoop` | error | `ResourceRef{path, index}` |
| unknown child (trigger path) | `path.unknownChild` | error | `MissingReference{step, ref=childID}` |
| multiple trigger paths | `trigger.multiplePaths` | error | `ResourceRef{path, index}` |
| unknown parent | `path.unknownParent` | error | `MissingReference{step, ref=parentID}` |
| unknown child | `path.unknownChild` | error | `MissingReference{step, ref=childID}` |
| internal index missing | `internal` | error | `ResourceRef{path, index}` |
| ambiguous entry | `graph.ambiguousEntry` | error | — |
| unknown scope | `scope.unknown` | error | `MissingReference{scope, ref}` |
| gateway < 2 paths | `gateway.tooFewPaths` | error | `GatewayPaths{stepID, got, want=2}` |
| gateway multiple else | `gateway.multipleElse` | error | `GatewayPaths` |
| gateway no else | `gateway.noElse` | error | `GatewayPaths` |

`validateScopeRefs` →
`issue(IssueCodeScopeUnknown, NgAutomationSeverityError, msg, types.NewDetailMissingReference("scope", e.Scope, steps[i].ID, kind, exprI))`.

---

## Phase 4 — Data-plumbing prerequisites

Two details need data not currently tracked:

| Detail | Gap | Fix |
|---|---|---|
| `duplicateID.indices` | `indexSteps` / `indexTriggers` store `idx[id]=ptr`, not first index | track first index (map value `{ptr,index}` or side `id→firstIndex`); report `[first, i]` |
| `cycle.stepIDs` | `detectCycle` returns a string naming one node | extend to return the loop ring (`[]id.ID` via gray path). Degrade to single node now if deferred |

Both small. `cycle` full-ring is optional for v1 (emit the one known node).

---

## Phase 5 — Direct constructions (`automation/service/ng_automation.go` ~L745)

run-as errors → `issue(IssueCodeRunAsLoadFailed, NgAutomationSeverityError, msg)` /
`IssueCodeRunAsInvalid`. No details.

---

## Phase 6 — Tests (`automation/service/ng_automation_converter_test.go`)

Keep the **single** scope scenario (per scope of work). Typed assertions:

```go
req.Equal(types.IssueCodeScopeUnknown, issues[0].Code)
req.Equal(types.NgAutomationSeverityError, issues[0].Severity)
req.Len(issues[0].Details, 1)

d := issues[0].Details[0].MissingReference
req.NotNil(d)
req.Equal("scope", d.RefKind)
req.Equal("step2", d.Ref)
req.Equal(uint64(3), d.StepID)
req.Equal("arguments", d.Field)
```

Add a round-trip marshal test (in `automation/types`) to lock the bridge:

```go
// marshal issue → JSON → unmarshal → deep-equal; assert raw JSON contains "type":"missingReference"
```

---

## Phase 7 — Verify

1. `go build ./automation/...`
2. `go vet ./automation/...` (ignore pre-existing `sync.RWMutex` warnings)
3. `go test ./automation/service/ -run TestValidateScopeRefs`
4. `go test ./automation/types/ -run TestIssueDetail` (round-trip)
5. Store round-trip: `NgAutomationIssueSet.Value()` / `Scan()` (gen file) marshal via the
   custom code — verify build, no schema change.
6. Grep: no stale `NgAutomationIssueCulprit`, `.Kind`, `.Culprit`, `map[string]int` on the NG issue type.

---

## Scope boundaries

- **In:** Go types, marshal/unmarshal, converter call sites, run-as constructions, tests.
- **Out:** OpenAPI generation, lib/js typed client. Raw JSON already carries `type` per
  detail (Decision H); FE consumes untyped JSON as today. Hand FE the `code` + detail-`type`
  catalogs from the const blocks.
- **Untouched:** legacy `WorkflowIssue` type (separate, not NG).

## Files

| File | Action |
|---|---|
| `automation/types/ng_automation.go` | replace issue struct, add code/detail/severity consts, `Error()` |
| `automation/types/ng_automation_issue.go` | **new** — variants, discriminators, marshal/unmarshal, `IssueIDs`, constructors |
| `automation/service/ng_automation_converter.go` | new `issue()` factory, rewrite ~22 call sites, plumb duplicate/cycle data |
| `automation/service/ng_automation.go` | run-as constructions (2) |
| `automation/service/ng_automation_converter_test.go` | typed scope-ref assertions |

## Risks

- Marshal errors if 0 or >1 variant set — constructors prevent; direct struct literals must
  set exactly one.
- `Parameters` interface non-nil → `omitempty` won't drop it; params always emitted. Fine
  (every detail has a type + params object).
- Breaking JSON change vs interim `Kind/Culprit` — branch-only, unreleased → safe.
