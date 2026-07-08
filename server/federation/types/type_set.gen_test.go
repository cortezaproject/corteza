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

func TestNodeSetWalk(t *testing.T) {
	var (
		value = make(NodeSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Node) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Node) error { return fmt.Errorf("walk error") }))
}

func TestNodeSetFilter(t *testing.T) {
	var (
		value = make(NodeSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Node) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Node) (bool, error) {
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
		_, err := value.Filter(func(*Node) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestNodeSetIDs(t *testing.T) {
	var (
		value = make(NodeSet, 3)
		req   = require.New(t)
	)

	value[0] = new(Node)
	value[1] = new(Node)
	value[2] = new(Node)
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

func TestNodeSyncSetWalk(t *testing.T) {
	var (
		value = make(NodeSyncSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*NodeSync) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*NodeSync) error { return fmt.Errorf("walk error") }))
}

func TestNodeSyncSetFilter(t *testing.T) {
	var (
		value = make(NodeSyncSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*NodeSync) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*NodeSync) (bool, error) {
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
		_, err := value.Filter(func(*NodeSync) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestExposedModuleSetWalk(t *testing.T) {
	var (
		value = make(ExposedModuleSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*ExposedModule) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*ExposedModule) error { return fmt.Errorf("walk error") }))
}

func TestExposedModuleSetFilter(t *testing.T) {
	var (
		value = make(ExposedModuleSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*ExposedModule) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*ExposedModule) (bool, error) {
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
		_, err := value.Filter(func(*ExposedModule) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestExposedModuleSetIDs(t *testing.T) {
	var (
		value = make(ExposedModuleSet, 3)
		req   = require.New(t)
	)

	value[0] = new(ExposedModule)
	value[1] = new(ExposedModule)
	value[2] = new(ExposedModule)
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

func TestSharedModuleSetWalk(t *testing.T) {
	var (
		value = make(SharedModuleSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*SharedModule) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*SharedModule) error { return fmt.Errorf("walk error") }))
}

func TestSharedModuleSetFilter(t *testing.T) {
	var (
		value = make(SharedModuleSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*SharedModule) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*SharedModule) (bool, error) {
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
		_, err := value.Filter(func(*SharedModule) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestSharedModuleSetIDs(t *testing.T) {
	var (
		value = make(SharedModuleSet, 3)
		req   = require.New(t)
	)

	value[0] = new(SharedModule)
	value[1] = new(SharedModule)
	value[2] = new(SharedModule)
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

func TestModuleMappingSetWalk(t *testing.T) {
	var (
		value = make(ModuleMappingSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*ModuleMapping) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*ModuleMapping) error { return fmt.Errorf("walk error") }))
}

func TestModuleMappingSetFilter(t *testing.T) {
	var (
		value = make(ModuleMappingSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*ModuleMapping) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*ModuleMapping) (bool, error) {
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
		_, err := value.Filter(func(*ModuleMapping) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}
