package importer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

type (
	// AutomationScript provides very basic and immature automation script importer
	AutomationScript struct {
		imp       *Importer
		namespace *types.Namespace
		set       automation.ScriptSet
		dirty     map[uint64]bool

		triggers map[string]automation.TriggerSet

		modRefs []automationTriggerModuleRef
	}

	automationTriggerModuleRef struct {
		// automation handle, report index, module handle
		as string
		ri int
		mh string
	}

	automationFinder interface {
		// automation finder does not have find by name...
		FindScripts(ctx context.Context, f automation.ScriptFilter) (automation.ScriptSet, automation.ScriptFilter, error)
	}
)

func NewAutomationImporter(imp *Importer, ns *types.Namespace) *AutomationScript {
	out := &AutomationScript{
		imp:       imp,
		namespace: ns,
		set:       automation.ScriptSet{},
		triggers:  make(map[string]automation.TriggerSet),
		modRefs:   make([]automationTriggerModuleRef, 0),
		dirty:     make(map[uint64]bool),
	}

	if imp.automationFinder != nil && ns.ID > 0 {
		out.set, _, _ = imp.automationFinder.FindScripts(
			context.Background(),
			automation.ScriptFilter{
				NamespaceID: ns.ID,
			},
		)
	}

	return out
}

func (asImp *AutomationScript) getModule(handle string) (*types.Module, error) {
	if g, ok := asImp.imp.namespaces.modules[asImp.namespace.Slug]; !ok {
		return nil, errors.Errorf("could not get modules %q from non existing namespace %q", handle, asImp.namespace.Slug)
	} else {
		return g.Get(handle)
	}
}

func (asImp *AutomationScript) CastSet(set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Automation scripts defined as collection
			deinterfacer.KVsetString(&handle, "name", def)
		}

		return asImp.Cast(handle, def)
	})
}

func (asImp *AutomationScript) Cast(handle string, def interface{}) (err error) {
	if !deinterfacer.IsMap(def) {
		return errors.New("expecting map of values for automation")
	}

	var script *automation.Script

	if !importer.IsValidHandle(handle) {
		return errors.New("invalid automation handle")
	}

	handle = importer.NormalizeHandle(handle)
	if script, err = asImp.Get(handle); err != nil {
		return err
	} else if script == nil {
		script = &automation.Script{
			NamespaceID: asImp.namespace.ID,
			Name:        handle,
		}

		asImp.set = append(asImp.set, script)
	} else if script.ID == 0 {
		return errors.Errorf("automation handle %q already defined in this import session", script.Name)
	}

	asImp.dirty[script.ID] = true

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "name":
			script.Name = deinterfacer.ToString(val)

		case "source":
			script.Source = deinterfacer.ToString(val)

		case "async":
			script.Async = deinterfacer.ToBool(val, false)
		case "runInUA":
			script.RunInUA = deinterfacer.ToBool(val, false)
		case "critical":
			script.Critical = deinterfacer.ToBool(val, false)
		case "enabled":
			script.Enabled = deinterfacer.ToBool(val, true)
		case "timeout":
			script.Timeout = uint(deinterfacer.ToInt(val))

		case "triggers":
			asImp.triggers[handle], err = asImp.castTriggers(handle, script, val)
			if err != nil {
				return err
			}

		case "allow", "deny":
			return asImp.imp.permissions.CastSet(types.AutomationScriptPermissionResource.String()+handle, key, val)

		default:
			return fmt.Errorf("unexpected key %q for automation %q", key, handle)
		}

		return
	})
}

func (asImp *AutomationScript) castTriggers(handle string, script *automation.Script, def interface{}) (automation.TriggerSet, error) {
	var (
		t  *automation.Trigger
		tt = automation.TriggerSet{}
	)

	return tt, deinterfacer.Each(def, func(n int, _ string, def interface{}) (err error) {
		t = &automation.Trigger{
			Enabled:   true,
			Condition: "0",
		}

		err = deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
			switch key {
			case "enabled":
				t.Enabled = deinterfacer.ToBool(val, true)
			case "event":
				t.Event = deinterfacer.ToString(val)
			case "resource":
				t.Resource = deinterfacer.ToString(val)
			case "condition":
				t.Condition = deinterfacer.ToString(val)
			case "module":
				module := deinterfacer.ToString(val)
				asImp.modRefs = append(asImp.modRefs, automationTriggerModuleRef{handle, len(tt), module})
			default:
				return fmt.Errorf("unexpected key %q for automation script's %q trigger", key, handle)

			}
			return
		})

		if err != nil {
			return
		}

		tt = append(tt, t)
		return
	})
}

// Get existing automation scripts
func (asImp *AutomationScript) Get(handle string) (*automation.Script, error) {
	handle = importer.NormalizeHandle(handle)
	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid automation script handle")
	}

	return asImp.set.FindByName(handle, asImp.namespace.ID), nil
}

func (asImp *AutomationScript) Store(ctx context.Context, k automationScriptKeeper) (err error) {
	if err = asImp.resolveRefs(); err != nil {
		return
	}

	for _, s := range asImp.set {
		var name = s.Name

		if tt, has := asImp.triggers[name]; has {
			s.AddTrigger(automation.STMS_UPDATE, tt...)
		}

		if s.ID == 0 {
			s.NamespaceID = asImp.namespace.ID
			err = k.CreateScript(ctx, s)
		} else if asImp.dirty[s.ID] {
			err = k.UpdateScript(ctx, s)
		}

		if err != nil {
			return errors.Wrapf(err, "could not save script %s (%d)", s.Name, s.ID)
		}

		asImp.dirty[s.ID] = false
		asImp.imp.permissions.UpdateResources(types.AutomationScriptPermissionResource.String(), name, s.ID)

	}

	return
}

// Resolve refs for all scripts
func (asImp *AutomationScript) resolveRefs() error {
	for _, ref := range asImp.modRefs {
		s := asImp.set.FindByName(ref.as, asImp.namespace.ID)
		if s == nil {
			// try to find it in no-mans land
			s = asImp.set.FindByName(ref.as, 0)
		}

		if s == nil {
			return errors.Errorf("invalid reference, unknown automation script (%v)", ref)
		}
		if _, has := asImp.triggers[ref.as]; !has {
			return errors.Errorf("invalid reference, triggers not initialized (%v)", ref)
		}
		if ref.ri > len(asImp.triggers[ref.as]) {
			return errors.Errorf("invalid reference, trigger index out of range (%v)", ref)
		}

		if module, err := asImp.getModule(ref.mh); err != nil {
			return errors.Errorf("invalid reference, module loading error: %v", err)
		} else if module == nil {
			return errors.Errorf("invalid reference, unknown module (%v)", ref)
		} else {
			asImp.triggers[ref.as][ref.ri].Condition = strconv.FormatUint(module.ID, 10)
		}
	}

	return nil
}
