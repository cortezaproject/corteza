package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	storeEncoder struct {
		s   store.Storer
		cfg *EncoderConfig

		// Each resource should define its own state that is used when encoding the resource.
		// Such approach removes the need for a janky generic global state.
		// This also simplifies any slight deviations between resource complexities.
		state map[resource.Interface]resourceState
	}

	// EncoderConfig allows us to configure the resource encoding process
	EncoderConfig struct {
		// OnExisting defines what to do if the resource exists
		OnExisting resource.MergeAlg
		// Skip if defines a pkg/expr expression when to skip the resource
		SkipIf string
		// Defer is called after the resource is encoded, regardles of the result
		Defer func()
		// DeferOk is called after the resource is encoded, only when successful
		DeferOk func()
		// DeferNok is called after the resource is encoded, only when failed
		// If you return an error, the encoding will terminate.
		// If you return nil (ignore the error), the encoding will continue.
		DeferNok func(error) error

		// IgnoreStore prevents encoders from accessing the store for initial resources
		IgnoreStore bool
	}

	accessControlRBACServicer interface {
		Can([]uint64, string, rbac.Resource) bool
	}

	composeAccessController interface {
		composeRecordValueAccessController
		composeRecordAccessController
	}
	composeRecordValueAccessController interface {
		CanReadRecordValue(context.Context, *types.ModuleField) bool
		CanUpdateRecordValue(context.Context, *types.ModuleField) bool
	}

	composeRecordAccessController interface {
		CanCreateRecordOnModule(context.Context, *types.Module) bool
		CanUpdateRecord(context.Context, *types.Record) bool
		CanDeleteRecord(context.Context, *types.Record) bool
	}

	payload struct {
		s     store.Storer
		state *envoy.ResourceState

		invokerID uint64
	}

	// resourceState allows each conforming struct to be initialized and encoded
	// by the store encoder
	resourceState interface {
		Prepare(context.Context, *payload) (err error)
		Encode(context.Context, *payload) (err error)
	}
)

var (
	ErrUnknownResource        = errors.New("unknown resource")
	ErrResourceStateUndefined = errors.New("undefined resource state")
)

// NewStoreEncoder initializes a fresh store encoder
//
// If no config is provided, it uses Skip as the default merge alg.
func NewStoreEncoder(s store.Storer, cfg *EncoderConfig) envoy.PrepareEncoder {
	if cfg == nil {
		cfg = &EncoderConfig{
			OnExisting: resource.Skip,
		}
	}

	return &storeEncoder{
		s:   s,
		cfg: cfg,

		state: make(map[resource.Interface]resourceState),
	}
}

// Prepare prepares the encoder for the given set of resources
//
// It initializes and prepares the resource state for each provided resource
func (se *storeEncoder) Prepare(ctx context.Context, ee ...*envoy.ResourceState) (err error) {
	f := func(rs resourceState, ers *envoy.ResourceState) error {
		err = rs.Prepare(ctx, se.makePayload(ctx, se.s, ers))
		if err != nil {
			return err
		}

		se.state[ers.Res] = rs
		return nil
	}

	for _, ers := range ee {
		switch res := ers.Res.(type) {
		// Compose resources
		case *resource.ComposeNamespace:
			err = f(newComposeNamespaceFromResource(res, se.cfg), ers)
		case *resource.ComposeModule:
			err = f(NewComposeModuleFromResource(res, se.cfg), ers)
		case *resource.ComposeRecord:
			err = f(NewComposeRecordFromResource(res, se.cfg), ers)
		case *resource.ComposeChart:
			err = f(newComposeChartFromResource(res, se.cfg), ers)
		case *resource.ComposePage:
			err = f(newComposePageFromResource(res, se.cfg), ers)

		// System resources
		case *resource.User:
			err = f(NewUserFromResource(res, se.cfg), ers)
		case *resource.Template:
			err = f(NewTemplateFromResource(res, se.cfg), ers)
		case *resource.Role:
			err = f(NewRoleFromResource(res, se.cfg), ers)
		case *resource.Application:
			err = f(NewApplicationFromResource(res, se.cfg), ers)
		case *resource.Setting:
			err = f(NewSettingFromResource(res, se.cfg), ers)
		case *resource.RbacRule:
			err = f(newRbacRuleFromResource(res, se.cfg), ers)
		case *resource.ResourceTranslation:
		//	err = f(newResourceTranslationFromResource(res, se.cfg), ers)

		// Automation resources
		case *resource.AutomationWorkflow:
			err = f(newAutomationWorkflowFromResource(res, se.cfg), ers)

		default:
			err = ErrUnknownResource
		}

		if err != nil {
			return se.WrapError("prepare", ers.Res, err)
		}
	}
	return nil
}

// Encode encodes available resource states using the given store encoder
func (se *storeEncoder) Encode(ctx context.Context, p envoy.Provider) error {
	var ers *envoy.ResourceState
	return store.Tx(ctx, se.s, func(ctx context.Context, s store.Storer) (err error) {
		for {
			ers, err = p.NextInverted(ctx)
			if err != nil {
				return err
			}
			if ers == nil {
				return nil
			}

			state := se.state[ers.Res]
			if state == nil {
				err = ErrResourceStateUndefined
			} else {
				err = state.Encode(ctx, se.makePayload(ctx, s, ers))
			}

			if err != nil {
				//return se.WrapError("encode", ers.Res, err)
			}
		}
	})
}

func (se *storeEncoder) makePayload(ctx context.Context, s store.Storer, ers *envoy.ResourceState) *payload {
	return &payload{
		s:         s,
		state:     ers,
		invokerID: auth.GetIdentityFromContext(ctx).Identity(),
	}
}

func (se *storeEncoder) WrapError(act string, res resource.Interface, err error) error {
	return fmt.Errorf("store encoder %s %s %v: %s", act, res.ResourceType(), res.Identifiers().StringSlice(), err)
}

func resourceErrIdentifierNotUnique(i string) error {
	return fmt.Errorf("bad resource identifier %v: not unique", i)
}
