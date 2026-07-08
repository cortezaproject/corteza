package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

type (
	WorkflowSet            []*Workflow
	SessionSet             []*Session
	TriggerSet             []*Trigger
	NgAutomationSet        []*NgAutomation
	TriggerConstraintSet   []*TriggerConstraint
	WorkflowPathSet        []*WorkflowPath
	WorkflowIssueSet       []*WorkflowIssue
	WorkflowStepSet        []*WorkflowStep
	StateSet               []*State
	NgAutomationIssueSet   []*NgAutomationIssue
	NgAutomationTriggerSet []*NgAutomationTrigger
	NgAutomationStepSet    []*NgAutomationStep
	NgAutomationPathSet    []*NgAutomationPath
)

func (set WorkflowSet) Walk(w func(*Workflow) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set WorkflowSet) FindByID(ID uint64) *Workflow {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set WorkflowSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set SessionSet) Walk(w func(*Session) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set SessionSet) FindByID(ID uint64) *Session {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set SessionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set TriggerSet) Walk(w func(*Trigger) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set TriggerSet) FindByID(ID uint64) *Trigger {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set TriggerSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set NgAutomationSet) Walk(w func(*NgAutomation) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set NgAutomationSet) FindByID(ID uint64) *NgAutomation {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set NgAutomationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set TriggerConstraintSet) Walk(w func(*TriggerConstraint) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set WorkflowPathSet) Walk(w func(*WorkflowPath) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set WorkflowIssueSet) Walk(w func(*WorkflowIssue) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set WorkflowStepSet) Walk(w func(*WorkflowStep) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set WorkflowStepSet) FindByID(ID uint64) *WorkflowStep {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set WorkflowStepSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set StateSet) Walk(w func(*State) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set StateSet) FindByID(ID uint64) *State {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set StateSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set NgAutomationIssueSet) Walk(w func(*NgAutomationIssue) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set NgAutomationIssueSet) Filter(f func(*NgAutomationIssue) (bool, error)) (out NgAutomationIssueSet, err error) {
	var ok bool
	out = NgAutomationIssueSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set NgAutomationTriggerSet) Walk(w func(*NgAutomationTrigger) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set NgAutomationTriggerSet) FindByID(ID uint64) *NgAutomationTrigger {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set NgAutomationTriggerSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set NgAutomationStepSet) Walk(w func(*NgAutomationStep) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set NgAutomationStepSet) FindByID(ID uint64) *NgAutomationStep {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set NgAutomationStepSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set NgAutomationPathSet) Walk(w func(*NgAutomationPath) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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
