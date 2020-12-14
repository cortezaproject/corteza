package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
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
	}

	// resourceState allows each conforming struct to be initialized and encoded
	// by the store encoder
	resourceState interface {
		Prepare(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error)
		Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error)
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
	f := func(rs resourceState, es *envoy.ResourceState) error {
		err = rs.Prepare(ctx, se.s, es)
		if err != nil {
			return err
		}

		se.state[es.Res] = rs
		return nil
	}

	for _, e := range ee {
		switch res := e.Res.(type) {
		// Compose resources
		case *resource.ComposeNamespace:
			err = f(NewComposeNamespaceState(res, se.cfg), e)
		case *resource.ComposeModule:
			err = f(NewComposeModuleState(res, se.cfg), e)
		case *resource.ComposeRecord:
			err = f(NewComposeRecordState(res, se.cfg), e)
		case *resource.ComposeChart:
			err = f(NewComposeChartState(res, se.cfg), e)
		case *resource.ComposePage:
			err = f(NewComposePageState(res, se.cfg), e)

		// System resources
		case *resource.User:
			err = f(NewUserState(res, se.cfg), e)
		case *resource.Role:
			err = f(NewRole(res, se.cfg), e)
		case *resource.Application:
			err = f(NewApplicationState(res, se.cfg), e)
		case *resource.Settings:
			err = f(NewSettingsState(res, se.cfg), e)
		case *resource.RbacRule:
			err = f(NewRbacRuleState(res, se.cfg), e)

			// Messaging resources
		case *resource.MessagingChannel:
			err = f(NewMessagingChannelState(res, se.cfg), e)

		default:
			err = ErrUnknownResource
		}

		if err != nil {
			return se.WrapError("prepare", e.Res, err)
		}
	}

	return nil
}

// Encode encodes available resource states using the given store encoder
func (se *storeEncoder) Encode(ctx context.Context, p envoy.Provider) error {
	var e *envoy.ResourceState
	return store.Tx(ctx, se.s, func(ctx context.Context, s store.Storer) (err error) {
		for {
			e, err = p.NextInverted(ctx)
			if err != nil {
				return err
			}
			if e == nil {
				return nil
			}

			state := se.state[e.Res]
			if state == nil {
				err = ErrResourceStateUndefined
			} else {
				err = state.Encode(ctx, se.s, e)
			}

			if err != nil {
				return se.WrapError("encode", e.Res, err)
			}
		}
	})
}

func (se *storeEncoder) WrapError(act string, res resource.Interface, err error) error {
	rt := strings.Join(strings.Split(strings.TrimSpace(strings.TrimRight(res.ResourceType(), ":")), ":"), " ")
	return fmt.Errorf("store encoder %s %s %v: %s", act, rt, res.Identifiers().StringSlice(), err)
}

func resourceErrIdentifierNotUnique(i string) error {
	return fmt.Errorf("bad resource identifier %v: not unique", i)
}
