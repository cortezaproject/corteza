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

func TestAttachmentSetWalk(t *testing.T) {
	var (
		value = make(AttachmentSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Attachment) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Attachment) error { return fmt.Errorf("walk error") }))
}

func TestAttachmentSetFilter(t *testing.T) {
	var (
		value = make(AttachmentSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Attachment) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Attachment) (bool, error) {
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
		_, err := value.Filter(func(*Attachment) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestAttachmentSetIDs(t *testing.T) {
	var (
		value = make(AttachmentSet, 3)
		req   = require.New(t)
	)

	value[0] = new(Attachment)
	value[1] = new(Attachment)
	value[2] = new(Attachment)
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

func TestChartSetWalk(t *testing.T) {
	var (
		value = make(ChartSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Chart) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Chart) error { return fmt.Errorf("walk error") }))
}

func TestChartSetFilter(t *testing.T) {
	var (
		value = make(ChartSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Chart) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Chart) (bool, error) {
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
		_, err := value.Filter(func(*Chart) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestChartSetIDs(t *testing.T) {
	var (
		value = make(ChartSet, 3)
		req   = require.New(t)
	)

	value[0] = new(Chart)
	value[1] = new(Chart)
	value[2] = new(Chart)
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

func TestModuleSetWalk(t *testing.T) {
	var (
		value = make(ModuleSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Module) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Module) error { return fmt.Errorf("walk error") }))
}

func TestModuleSetFilter(t *testing.T) {
	var (
		value = make(ModuleSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Module) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Module) (bool, error) {
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
		_, err := value.Filter(func(*Module) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestModuleSetIDs(t *testing.T) {
	var (
		value = make(ModuleSet, 3)
		req   = require.New(t)
	)

	value[0] = new(Module)
	value[1] = new(Module)
	value[2] = new(Module)
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

func TestModuleFieldSetWalk(t *testing.T) {
	var (
		value = make(ModuleFieldSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*ModuleField) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*ModuleField) error { return fmt.Errorf("walk error") }))
}

func TestModuleFieldSetFilter(t *testing.T) {
	var (
		value = make(ModuleFieldSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*ModuleField) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*ModuleField) (bool, error) {
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
		_, err := value.Filter(func(*ModuleField) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestModuleFieldSetIDs(t *testing.T) {
	var (
		value = make(ModuleFieldSet, 3)
		req   = require.New(t)
	)

	value[0] = new(ModuleField)
	value[1] = new(ModuleField)
	value[2] = new(ModuleField)
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

func TestNamespaceSetWalk(t *testing.T) {
	var (
		value = make(NamespaceSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Namespace) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Namespace) error { return fmt.Errorf("walk error") }))
}

func TestNamespaceSetFilter(t *testing.T) {
	var (
		value = make(NamespaceSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Namespace) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Namespace) (bool, error) {
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
		_, err := value.Filter(func(*Namespace) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestNamespaceSetIDs(t *testing.T) {
	var (
		value = make(NamespaceSet, 3)
		req   = require.New(t)
	)

	value[0] = new(Namespace)
	value[1] = new(Namespace)
	value[2] = new(Namespace)
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

func TestPageSetWalk(t *testing.T) {
	var (
		value = make(PageSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Page) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Page) error { return fmt.Errorf("walk error") }))
}

func TestPageSetFilter(t *testing.T) {
	var (
		value = make(PageSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Page) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Page) (bool, error) {
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
		_, err := value.Filter(func(*Page) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestPageSetIDs(t *testing.T) {
	var (
		value = make(PageSet, 3)
		req   = require.New(t)
	)

	value[0] = new(Page)
	value[1] = new(Page)
	value[2] = new(Page)
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

func TestPageLayoutSetWalk(t *testing.T) {
	var (
		value = make(PageLayoutSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*PageLayout) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*PageLayout) error { return fmt.Errorf("walk error") }))
}

func TestPageLayoutSetFilter(t *testing.T) {
	var (
		value = make(PageLayoutSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*PageLayout) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*PageLayout) (bool, error) {
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
		_, err := value.Filter(func(*PageLayout) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestPageLayoutSetIDs(t *testing.T) {
	var (
		value = make(PageLayoutSet, 3)
		req   = require.New(t)
	)

	value[0] = new(PageLayout)
	value[1] = new(PageLayout)
	value[2] = new(PageLayout)
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

func TestRecordSetWalk(t *testing.T) {
	var (
		value = make(RecordSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Record) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Record) error { return fmt.Errorf("walk error") }))
}

func TestRecordSetFilter(t *testing.T) {
	var (
		value = make(RecordSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Record) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Record) (bool, error) {
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
		_, err := value.Filter(func(*Record) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestRecordSetIDs(t *testing.T) {
	var (
		value = make(RecordSet, 3)
		req   = require.New(t)
	)

	value[0] = new(Record)
	value[1] = new(Record)
	value[2] = new(Record)
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

func TestIconSetWalk(t *testing.T) {
	var (
		value = make(IconSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Icon) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Icon) error { return fmt.Errorf("walk error") }))
}

func TestIconSetFilter(t *testing.T) {
	var (
		value = make(IconSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Icon) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Icon) (bool, error) {
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
		_, err := value.Filter(func(*Icon) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestRecordValueSetWalk(t *testing.T) {
	var (
		value = make(RecordValueSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*RecordValue) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*RecordValue) error { return fmt.Errorf("walk error") }))
}

func TestRecordValueSetFilter(t *testing.T) {
	var (
		value = make(RecordValueSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*RecordValue) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*RecordValue) (bool, error) {
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
		_, err := value.Filter(func(*RecordValue) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestPrivacyModuleSetWalk(t *testing.T) {
	var (
		value = make(PrivacyModuleSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*PrivacyModule) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*PrivacyModule) error { return fmt.Errorf("walk error") }))
}

func TestPrivacyModuleSetFilter(t *testing.T) {
	var (
		value = make(PrivacyModuleSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*PrivacyModule) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*PrivacyModule) (bool, error) {
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
		_, err := value.Filter(func(*PrivacyModule) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestDeDupRuleSetWalk(t *testing.T) {
	var (
		value = make(DeDupRuleSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*DeDupRule) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*DeDupRule) error { return fmt.Errorf("walk error") }))
}

func TestDeDupRuleSetFilter(t *testing.T) {
	var (
		value = make(DeDupRuleSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*DeDupRule) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*DeDupRule) (bool, error) {
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
		_, err := value.Filter(func(*DeDupRule) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}
