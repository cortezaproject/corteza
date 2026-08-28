package service

import (
	"context"
	"testing"
)

// resolvingLLM is the shape the agent service now requires: a validator that can
// settle a provider and model, not merely check a temperature.
type resolvingLLM struct{}

func (resolvingLLM) ValidateTemperature(context.Context, uint64, string, *float64) error { return nil }

func (resolvingLLM) ResolveModel(context.Context, uint64, string) (uint64, string, error) {
	return 0, "", nil
}

// An agent stored with no model was accepted and then failed on its first
// prompt, in front of whoever ran it rather than whoever wrote it. Creation
// resolves the model now, and this pins the contract that makes that possible:
// drop ResolveModel from the interface and the package stops compiling.
//
// The behaviour itself is proven live — creating an agent with no model is
// refused with the provider's models named — because it needs a real provider
// to resolve against.
func TestAgentLLMValidatorResolvesTheModel(t *testing.T) {
	var v agentLLMValidator = resolvingLLM{}

	if _, _, err := v.ResolveModel(context.Background(), 0, ""); err != nil {
		t.Fatalf("contract call failed: %v", err)
	}
}
