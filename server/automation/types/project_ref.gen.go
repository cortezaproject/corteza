package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// ProjectRef returns the ID of the project Workflow belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r Workflow) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project Session belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r Session) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project Trigger belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r Trigger) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project NgAutomation belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r NgAutomation) ProjectRef() uint64 {
	return r.ProjectID
}
