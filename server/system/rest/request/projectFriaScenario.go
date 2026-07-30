package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/crusttech/human/server/pkg/payload"
	"github.com/go-chi/chi/v5"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	ProjectFriaScenarioList struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// Query GET parameter
		//
		// Search query (matches the title)
		Query string

		// ProjectFriaScenarioID GET parameter
		//
		// Filter by FRIA risk scenario IDs
		ProjectFriaScenarioID []string

		// AiSystemID GET parameter
		//
		// Filter by the assessed AI system
		AiSystemID uint64 `json:",string"`

		// Severity GET parameter
		//
		// Severity filter (low|medium|high|critical)
		Severity string

		// Deleted GET parameter
		//
		// Deleted filter (0=exclude 1=only 2=include)
		Deleted uint

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// IncTotal GET parameter
		//
		// Include total count
		IncTotal bool

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort
		Sort string
	}

	ProjectFriaScenarioCreate struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// AiSystemID POST parameter
		//
		// AI system this scenario assesses
		AiSystemID uint64 `json:",string"`

		// Title POST parameter
		//
		// Scenario title
		Title string

		// Severity POST parameter
		//
		// Severity of the assessed harm (low|medium|high|critical)
		Severity string

		// Description POST parameter
		//
		// Harm description
		Description string

		// TriggerTypes POST parameter
		//
		// Trigger condition taxonomy keys
		TriggerTypes []string

		// TriggerDescription POST parameter
		//
		// Trigger description
		TriggerDescription string

		// ImpactedParties POST parameter
		//
		// Impacted party taxonomy keys
		ImpactedParties []string

		// VulnerableGroups POST parameter
		//
		// Vulnerable group taxonomy keys
		VulnerableGroups []string

		// VulnerableGroupsNotes POST parameter
		//
		// Vulnerable group notes
		VulnerableGroupsNotes string

		// Rights POST parameter
		//
		// Fundamental rights taxonomy keys
		Rights []string

		// HarmVectors POST parameter
		//
		// AI harm vector taxonomy keys
		HarmVectors []string

		// HarmVectorsDescription POST parameter
		//
		// AI harm vector description
		HarmVectorsDescription string
	}

	ProjectFriaScenarioRead struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectFriaScenarioID PATH parameter
		//
		// FRIA risk scenario ID
		ProjectFriaScenarioID uint64 `json:",string"`
	}

	ProjectFriaScenarioUpdate struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectFriaScenarioID PATH parameter
		//
		// FRIA risk scenario ID
		ProjectFriaScenarioID uint64 `json:",string"`

		// AiSystemID POST parameter
		//
		// AI system this scenario assesses
		AiSystemID uint64 `json:",string"`

		// Title POST parameter
		//
		// Scenario title
		Title string

		// Severity POST parameter
		//
		// Severity of the assessed harm (low|medium|high|critical)
		Severity string

		// Description POST parameter
		//
		// Harm description
		Description string

		// TriggerTypes POST parameter
		//
		// Trigger condition taxonomy keys
		TriggerTypes []string

		// TriggerDescription POST parameter
		//
		// Trigger description
		TriggerDescription string

		// ImpactedParties POST parameter
		//
		// Impacted party taxonomy keys
		ImpactedParties []string

		// VulnerableGroups POST parameter
		//
		// Vulnerable group taxonomy keys
		VulnerableGroups []string

		// VulnerableGroupsNotes POST parameter
		//
		// Vulnerable group notes
		VulnerableGroupsNotes string

		// Rights POST parameter
		//
		// Fundamental rights taxonomy keys
		Rights []string

		// HarmVectors POST parameter
		//
		// AI harm vector taxonomy keys
		HarmVectors []string

		// HarmVectorsDescription POST parameter
		//
		// AI harm vector description
		HarmVectorsDescription string

		// UpdatedAt POST parameter
		//
		// Last update timestamp
		UpdatedAt *time.Time
	}

	ProjectFriaScenarioDelete struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectFriaScenarioID PATH parameter
		//
		// FRIA risk scenario ID
		ProjectFriaScenarioID uint64 `json:",string"`
	}
)

// NewProjectFriaScenarioList request
func NewProjectFriaScenarioList() *ProjectFriaScenarioList {
	return &ProjectFriaScenarioList{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":             r.ProjectID,
		"query":                 r.Query,
		"projectFriaScenarioID": r.ProjectFriaScenarioID,
		"aiSystemID":            r.AiSystemID,
		"severity":              r.Severity,
		"deleted":               r.Deleted,
		"limit":                 r.Limit,
		"incTotal":              r.IncTotal,
		"pageCursor":            r.PageCursor,
		"sort":                  r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioList) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioList) GetProjectFriaScenarioID() []string {
	return r.ProjectFriaScenarioID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioList) GetAiSystemID() uint64 {
	return r.AiSystemID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioList) GetSeverity() string {
	return r.Severity
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ProjectFriaScenarioList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["projectFriaScenarioID[]"]; ok {
			r.ProjectFriaScenarioID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["projectFriaScenarioID"]; ok {
			r.ProjectFriaScenarioID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["aiSystemID"]; ok && len(val) > 0 {
			r.AiSystemID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["severity"]; ok && len(val) > 0 {
			r.Severity, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["deleted"]; ok && len(val) > 0 {
			r.Deleted, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["incTotal"]; ok && len(val) > 0 {
			r.IncTotal, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["pageCursor"]; ok && len(val) > 0 {
			r.PageCursor, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["sort"]; ok && len(val) > 0 {
			r.Sort, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectFriaScenarioCreate request
func NewProjectFriaScenarioCreate() *ProjectFriaScenarioCreate {
	return &ProjectFriaScenarioCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":              r.ProjectID,
		"aiSystemID":             r.AiSystemID,
		"title":                  r.Title,
		"severity":               r.Severity,
		"description":            r.Description,
		"triggerTypes":           r.TriggerTypes,
		"triggerDescription":     r.TriggerDescription,
		"impactedParties":        r.ImpactedParties,
		"vulnerableGroups":       r.VulnerableGroups,
		"vulnerableGroupsNotes":  r.VulnerableGroupsNotes,
		"rights":                 r.Rights,
		"harmVectors":            r.HarmVectors,
		"harmVectorsDescription": r.HarmVectorsDescription,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetAiSystemID() uint64 {
	return r.AiSystemID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetSeverity() string {
	return r.Severity
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetTriggerTypes() []string {
	return r.TriggerTypes
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetTriggerDescription() string {
	return r.TriggerDescription
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetImpactedParties() []string {
	return r.ImpactedParties
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetVulnerableGroups() []string {
	return r.VulnerableGroups
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetVulnerableGroupsNotes() string {
	return r.VulnerableGroupsNotes
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetRights() []string {
	return r.Rights
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetHarmVectors() []string {
	return r.HarmVectors
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioCreate) GetHarmVectorsDescription() string {
	return r.HarmVectorsDescription
}

// Fill processes request and fills internal variables
func (r *ProjectFriaScenarioCreate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["aiSystemID"]; ok && len(val) > 0 {
				r.AiSystemID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["title"]; ok && len(val) > 0 {
				r.Title, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["severity"]; ok && len(val) > 0 {
				r.Severity, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["description"]; ok && len(val) > 0 {
				r.Description, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["triggerDescription"]; ok && len(val) > 0 {
				r.TriggerDescription, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["vulnerableGroupsNotes"]; ok && len(val) > 0 {
				r.VulnerableGroupsNotes, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["harmVectorsDescription"]; ok && len(val) > 0 {
				r.HarmVectorsDescription, err = val[0], nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["aiSystemID"]; ok && len(val) > 0 {
			r.AiSystemID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["title"]; ok && len(val) > 0 {
			r.Title, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["severity"]; ok && len(val) > 0 {
			r.Severity, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["description"]; ok && len(val) > 0 {
			r.Description, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["triggerTypes[]"]; ok && len(val) > 0  {
		//    r.TriggerTypes, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["triggerDescription"]; ok && len(val) > 0 {
			r.TriggerDescription, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["impactedParties[]"]; ok && len(val) > 0  {
		//    r.ImpactedParties, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		//if val, ok := req.Form["vulnerableGroups[]"]; ok && len(val) > 0  {
		//    r.VulnerableGroups, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["vulnerableGroupsNotes"]; ok && len(val) > 0 {
			r.VulnerableGroupsNotes, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["rights[]"]; ok && len(val) > 0  {
		//    r.Rights, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		//if val, ok := req.Form["harmVectors[]"]; ok && len(val) > 0  {
		//    r.HarmVectors, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["harmVectorsDescription"]; ok && len(val) > 0 {
			r.HarmVectorsDescription, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectFriaScenarioRead request
func NewProjectFriaScenarioRead() *ProjectFriaScenarioRead {
	return &ProjectFriaScenarioRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":             r.ProjectID,
		"projectFriaScenarioID": r.ProjectFriaScenarioID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioRead) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioRead) GetProjectFriaScenarioID() uint64 {
	return r.ProjectFriaScenarioID
}

// Fill processes request and fills internal variables
func (r *ProjectFriaScenarioRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectFriaScenarioID")
		r.ProjectFriaScenarioID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectFriaScenarioUpdate request
func NewProjectFriaScenarioUpdate() *ProjectFriaScenarioUpdate {
	return &ProjectFriaScenarioUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":              r.ProjectID,
		"projectFriaScenarioID":  r.ProjectFriaScenarioID,
		"aiSystemID":             r.AiSystemID,
		"title":                  r.Title,
		"severity":               r.Severity,
		"description":            r.Description,
		"triggerTypes":           r.TriggerTypes,
		"triggerDescription":     r.TriggerDescription,
		"impactedParties":        r.ImpactedParties,
		"vulnerableGroups":       r.VulnerableGroups,
		"vulnerableGroupsNotes":  r.VulnerableGroupsNotes,
		"rights":                 r.Rights,
		"harmVectors":            r.HarmVectors,
		"harmVectorsDescription": r.HarmVectorsDescription,
		"updatedAt":              r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetProjectFriaScenarioID() uint64 {
	return r.ProjectFriaScenarioID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetAiSystemID() uint64 {
	return r.AiSystemID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetSeverity() string {
	return r.Severity
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetTriggerTypes() []string {
	return r.TriggerTypes
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetTriggerDescription() string {
	return r.TriggerDescription
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetImpactedParties() []string {
	return r.ImpactedParties
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetVulnerableGroups() []string {
	return r.VulnerableGroups
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetVulnerableGroupsNotes() string {
	return r.VulnerableGroupsNotes
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetRights() []string {
	return r.Rights
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetHarmVectors() []string {
	return r.HarmVectors
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetHarmVectorsDescription() string {
	return r.HarmVectorsDescription
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *ProjectFriaScenarioUpdate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["aiSystemID"]; ok && len(val) > 0 {
				r.AiSystemID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["title"]; ok && len(val) > 0 {
				r.Title, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["severity"]; ok && len(val) > 0 {
				r.Severity, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["description"]; ok && len(val) > 0 {
				r.Description, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["triggerDescription"]; ok && len(val) > 0 {
				r.TriggerDescription, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["vulnerableGroupsNotes"]; ok && len(val) > 0 {
				r.VulnerableGroupsNotes, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["harmVectorsDescription"]; ok && len(val) > 0 {
				r.HarmVectorsDescription, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["updatedAt"]; ok && len(val) > 0 {
				r.UpdatedAt, err = payload.ParseISODatePtrWithErr(val[0])
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["aiSystemID"]; ok && len(val) > 0 {
			r.AiSystemID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["title"]; ok && len(val) > 0 {
			r.Title, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["severity"]; ok && len(val) > 0 {
			r.Severity, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["description"]; ok && len(val) > 0 {
			r.Description, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["triggerTypes[]"]; ok && len(val) > 0  {
		//    r.TriggerTypes, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["triggerDescription"]; ok && len(val) > 0 {
			r.TriggerDescription, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["impactedParties[]"]; ok && len(val) > 0  {
		//    r.ImpactedParties, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		//if val, ok := req.Form["vulnerableGroups[]"]; ok && len(val) > 0  {
		//    r.VulnerableGroups, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["vulnerableGroupsNotes"]; ok && len(val) > 0 {
			r.VulnerableGroupsNotes, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["rights[]"]; ok && len(val) > 0  {
		//    r.Rights, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		//if val, ok := req.Form["harmVectors[]"]; ok && len(val) > 0  {
		//    r.HarmVectors, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["harmVectorsDescription"]; ok && len(val) > 0 {
			r.HarmVectorsDescription, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["updatedAt"]; ok && len(val) > 0 {
			r.UpdatedAt, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectFriaScenarioID")
		r.ProjectFriaScenarioID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectFriaScenarioDelete request
func NewProjectFriaScenarioDelete() *ProjectFriaScenarioDelete {
	return &ProjectFriaScenarioDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":             r.ProjectID,
		"projectFriaScenarioID": r.ProjectFriaScenarioID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioDelete) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFriaScenarioDelete) GetProjectFriaScenarioID() uint64 {
	return r.ProjectFriaScenarioID
}

// Fill processes request and fills internal variables
func (r *ProjectFriaScenarioDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectFriaScenarioID")
		r.ProjectFriaScenarioID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
