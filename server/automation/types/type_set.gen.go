package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/types/types.yaml

type (

	// NgAutomationSet slice of NgAutomation
	//
	// This type is auto-generated.
	NgAutomationSet []*NgAutomation

	// NgAutomationPathSet slice of NgAutomationPath
	//
	// This type is auto-generated.
	NgAutomationPathSet []*NgAutomationPath

	// NgAutomationStepSet slice of NgAutomationStep
	//
	// This type is auto-generated.
	NgAutomationStepSet []*NgAutomationStep

	// NgAutomationTriggerSet slice of NgAutomationTrigger
	//
	// This type is auto-generated.
	NgAutomationTriggerSet []*NgAutomationTrigger

	// SessionSet slice of Session
	//
	// This type is auto-generated.
	SessionSet []*Session

	// StateSet slice of State
	//
	// This type is auto-generated.
	StateSet []*State

	// TriggerSet slice of Trigger
	//
	// This type is auto-generated.
	TriggerSet []*Trigger

	// TriggerConstraintSet slice of TriggerConstraint
	//
	// This type is auto-generated.
	TriggerConstraintSet []*TriggerConstraint

	// WorkflowSet slice of Workflow
	//
	// This type is auto-generated.
	WorkflowSet []*Workflow

	// WorkflowIssueSet slice of WorkflowIssue
	//
	// This type is auto-generated.
	WorkflowIssueSet []*WorkflowIssue

	// WorkflowPathSet slice of WorkflowPath
	//
	// This type is auto-generated.
	WorkflowPathSet []*WorkflowPath

	// WorkflowStepSet slice of WorkflowStep
	//
	// This type is auto-generated.
	WorkflowStepSet []*WorkflowStep
)

// Walk iterates through every slice item and calls w(NgAutomation) err
//
// This function is auto-generated.
func (set NgAutomationSet) Walk(w func(*NgAutomation) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(NgAutomation) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set NgAutomationSet) Filter(f func(*NgAutomation) (bool, error)) (out NgAutomationSet, err error) {
	var ok bool
	out = NgAutomationSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set NgAutomationSet) FindByID(ID uint64) *NgAutomation {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set NgAutomationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(NgAutomationPath) err
//
// This function is auto-generated.
func (set NgAutomationPathSet) Walk(w func(*NgAutomationPath) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(NgAutomationPath) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set NgAutomationPathSet) Filter(f func(*NgAutomationPath) (bool, error)) (out NgAutomationPathSet, err error) {
	var ok bool
	out = NgAutomationPathSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(NgAutomationStep) err
//
// This function is auto-generated.
func (set NgAutomationStepSet) Walk(w func(*NgAutomationStep) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(NgAutomationStep) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set NgAutomationStepSet) Filter(f func(*NgAutomationStep) (bool, error)) (out NgAutomationStepSet, err error) {
	var ok bool
	out = NgAutomationStepSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set NgAutomationStepSet) FindByID(ID uint64) *NgAutomationStep {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set NgAutomationStepSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(NgAutomationTrigger) err
//
// This function is auto-generated.
func (set NgAutomationTriggerSet) Walk(w func(*NgAutomationTrigger) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(NgAutomationTrigger) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set NgAutomationTriggerSet) Filter(f func(*NgAutomationTrigger) (bool, error)) (out NgAutomationTriggerSet, err error) {
	var ok bool
	out = NgAutomationTriggerSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set NgAutomationTriggerSet) FindByID(ID uint64) *NgAutomationTrigger {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set NgAutomationTriggerSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Session) err
//
// This function is auto-generated.
func (set SessionSet) Walk(w func(*Session) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Session) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set SessionSet) Filter(f func(*Session) (bool, error)) (out SessionSet, err error) {
	var ok bool
	out = SessionSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set SessionSet) FindByID(ID uint64) *Session {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set SessionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(State) err
//
// This function is auto-generated.
func (set StateSet) Walk(w func(*State) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(State) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set StateSet) Filter(f func(*State) (bool, error)) (out StateSet, err error) {
	var ok bool
	out = StateSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set StateSet) FindByID(ID uint64) *State {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set StateSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Trigger) err
//
// This function is auto-generated.
func (set TriggerSet) Walk(w func(*Trigger) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Trigger) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set TriggerSet) Filter(f func(*Trigger) (bool, error)) (out TriggerSet, err error) {
	var ok bool
	out = TriggerSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set TriggerSet) FindByID(ID uint64) *Trigger {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set TriggerSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(TriggerConstraint) err
//
// This function is auto-generated.
func (set TriggerConstraintSet) Walk(w func(*TriggerConstraint) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(TriggerConstraint) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set TriggerConstraintSet) Filter(f func(*TriggerConstraint) (bool, error)) (out TriggerConstraintSet, err error) {
	var ok bool
	out = TriggerConstraintSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(Workflow) err
//
// This function is auto-generated.
func (set WorkflowSet) Walk(w func(*Workflow) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Workflow) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set WorkflowSet) Filter(f func(*Workflow) (bool, error)) (out WorkflowSet, err error) {
	var ok bool
	out = WorkflowSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set WorkflowSet) FindByID(ID uint64) *Workflow {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set WorkflowSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(WorkflowIssue) err
//
// This function is auto-generated.
func (set WorkflowIssueSet) Walk(w func(*WorkflowIssue) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(WorkflowIssue) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set WorkflowIssueSet) Filter(f func(*WorkflowIssue) (bool, error)) (out WorkflowIssueSet, err error) {
	var ok bool
	out = WorkflowIssueSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(WorkflowPath) err
//
// This function is auto-generated.
func (set WorkflowPathSet) Walk(w func(*WorkflowPath) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(WorkflowPath) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set WorkflowPathSet) Filter(f func(*WorkflowPath) (bool, error)) (out WorkflowPathSet, err error) {
	var ok bool
	out = WorkflowPathSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(WorkflowStep) err
//
// This function is auto-generated.
func (set WorkflowStepSet) Walk(w func(*WorkflowStep) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(WorkflowStep) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set WorkflowStepSet) Filter(f func(*WorkflowStep) (bool, error)) (out WorkflowStepSet, err error) {
	var ok bool
	out = WorkflowStepSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set WorkflowStepSet) FindByID(ID uint64) *WorkflowStep {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set WorkflowStepSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
