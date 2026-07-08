package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWorkflowSetWalk(t *testing.T) {
	var (
		value = make(WorkflowSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Workflow) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Workflow) error { return fmt.Errorf("walk error") }))
}

func TestWorkflowSetFilter(t *testing.T) {
	var (
		value = make(WorkflowSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Workflow) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Workflow) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*Workflow) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestWorkflowSetIDs(t *testing.T) {
	var (
		value = make(WorkflowSet, 3)
		req   = require.New(t)
	)

	value[0] = new(Workflow)
	value[1] = new(Workflow)
	value[2] = new(Workflow)
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestSessionSetWalk(t *testing.T) {
	var (
		value = make(SessionSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Session) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Session) error { return fmt.Errorf("walk error") }))
}

func TestSessionSetFilter(t *testing.T) {
	var (
		value = make(SessionSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Session) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Session) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*Session) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestSessionSetIDs(t *testing.T) {
	var (
		value = make(SessionSet, 3)
		req   = require.New(t)
	)

	value[0] = new(Session)
	value[1] = new(Session)
	value[2] = new(Session)
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestTriggerSetWalk(t *testing.T) {
	var (
		value = make(TriggerSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Trigger) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Trigger) error { return fmt.Errorf("walk error") }))
}

func TestTriggerSetFilter(t *testing.T) {
	var (
		value = make(TriggerSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Trigger) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Trigger) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*Trigger) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestTriggerSetIDs(t *testing.T) {
	var (
		value = make(TriggerSet, 3)
		req   = require.New(t)
	)

	value[0] = new(Trigger)
	value[1] = new(Trigger)
	value[2] = new(Trigger)
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestNgAutomationSetWalk(t *testing.T) {
	var (
		value = make(NgAutomationSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*NgAutomation) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*NgAutomation) error { return fmt.Errorf("walk error") }))
}

func TestNgAutomationSetFilter(t *testing.T) {
	var (
		value = make(NgAutomationSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*NgAutomation) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*NgAutomation) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*NgAutomation) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestNgAutomationSetIDs(t *testing.T) {
	var (
		value = make(NgAutomationSet, 3)
		req   = require.New(t)
	)

	value[0] = new(NgAutomation)
	value[1] = new(NgAutomation)
	value[2] = new(NgAutomation)
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestTriggerConstraintSetWalk(t *testing.T) {
	var (
		value = make(TriggerConstraintSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*TriggerConstraint) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*TriggerConstraint) error { return fmt.Errorf("walk error") }))
}

func TestTriggerConstraintSetFilter(t *testing.T) {
	var (
		value = make(TriggerConstraintSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*TriggerConstraint) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*TriggerConstraint) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*TriggerConstraint) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestWorkflowPathSetWalk(t *testing.T) {
	var (
		value = make(WorkflowPathSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*WorkflowPath) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*WorkflowPath) error { return fmt.Errorf("walk error") }))
}

func TestWorkflowPathSetFilter(t *testing.T) {
	var (
		value = make(WorkflowPathSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*WorkflowPath) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*WorkflowPath) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*WorkflowPath) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestWorkflowIssueSetWalk(t *testing.T) {
	var (
		value = make(WorkflowIssueSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*WorkflowIssue) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*WorkflowIssue) error { return fmt.Errorf("walk error") }))
}

func TestWorkflowIssueSetFilter(t *testing.T) {
	var (
		value = make(WorkflowIssueSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*WorkflowIssue) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*WorkflowIssue) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*WorkflowIssue) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestWorkflowStepSetWalk(t *testing.T) {
	var (
		value = make(WorkflowStepSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*WorkflowStep) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*WorkflowStep) error { return fmt.Errorf("walk error") }))
}

func TestWorkflowStepSetFilter(t *testing.T) {
	var (
		value = make(WorkflowStepSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*WorkflowStep) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*WorkflowStep) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*WorkflowStep) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestWorkflowStepSetIDs(t *testing.T) {
	var (
		value = make(WorkflowStepSet, 3)
		req   = require.New(t)
	)

	value[0] = new(WorkflowStep)
	value[1] = new(WorkflowStep)
	value[2] = new(WorkflowStep)
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestStateSetWalk(t *testing.T) {
	var (
		value = make(StateSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*State) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*State) error { return fmt.Errorf("walk error") }))
}

func TestStateSetFilter(t *testing.T) {
	var (
		value = make(StateSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*State) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*State) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*State) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestStateSetIDs(t *testing.T) {
	var (
		value = make(StateSet, 3)
		req   = require.New(t)
	)

	value[0] = new(State)
	value[1] = new(State)
	value[2] = new(State)
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestNgAutomationIssueSetWalk(t *testing.T) {
	var (
		value = make(NgAutomationIssueSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*NgAutomationIssue) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*NgAutomationIssue) error { return fmt.Errorf("walk error") }))
}

func TestNgAutomationIssueSetFilter(t *testing.T) {
	var (
		value = make(NgAutomationIssueSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*NgAutomationIssue) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*NgAutomationIssue) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*NgAutomationIssue) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestNgAutomationTriggerSetWalk(t *testing.T) {
	var (
		value = make(NgAutomationTriggerSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*NgAutomationTrigger) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*NgAutomationTrigger) error { return fmt.Errorf("walk error") }))
}

func TestNgAutomationTriggerSetFilter(t *testing.T) {
	var (
		value = make(NgAutomationTriggerSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*NgAutomationTrigger) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*NgAutomationTrigger) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*NgAutomationTrigger) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestNgAutomationTriggerSetIDs(t *testing.T) {
	var (
		value = make(NgAutomationTriggerSet, 3)
		req   = require.New(t)
	)

	value[0] = new(NgAutomationTrigger)
	value[1] = new(NgAutomationTrigger)
	value[2] = new(NgAutomationTrigger)
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestNgAutomationStepSetWalk(t *testing.T) {
	var (
		value = make(NgAutomationStepSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*NgAutomationStep) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*NgAutomationStep) error { return fmt.Errorf("walk error") }))
}

func TestNgAutomationStepSetFilter(t *testing.T) {
	var (
		value = make(NgAutomationStepSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*NgAutomationStep) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*NgAutomationStep) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*NgAutomationStep) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestNgAutomationStepSetIDs(t *testing.T) {
	var (
		value = make(NgAutomationStepSet, 3)
		req   = require.New(t)
	)

	value[0] = new(NgAutomationStep)
	value[1] = new(NgAutomationStep)
	value[2] = new(NgAutomationStep)
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestNgAutomationPathSetWalk(t *testing.T) {
	var (
		value = make(NgAutomationPathSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*NgAutomationPath) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*NgAutomationPath) error { return fmt.Errorf("walk error") }))
}

func TestNgAutomationPathSetFilter(t *testing.T) {
	var (
		value = make(NgAutomationPathSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*NgAutomationPath) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*NgAutomationPath) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*NgAutomationPath) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}
