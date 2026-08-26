package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// ProjectRef returns the ID of the project Attachment belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r Attachment) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project Chart belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r Chart) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project Module belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r Module) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project ModuleField belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r ModuleField) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project Namespace belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r Namespace) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project Page belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r Page) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project PageLayout belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r PageLayout) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project Record belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r Record) ProjectRef() uint64 {
	return r.ProjectID
}
