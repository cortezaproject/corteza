package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// ProjectRef returns the ID of the project ExposedModule belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r ExposedModule) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project SharedModule belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r SharedModule) ProjectRef() uint64 {
	return r.ProjectID
}

// ProjectRef returns the ID of the project ModuleMapping belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r ModuleMapping) ProjectRef() uint64 {
	return r.ProjectID
}
