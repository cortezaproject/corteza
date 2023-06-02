package bundler

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cortezaproject/corteza/server/automation/plugin"
	"github.com/cortezaproject/corteza/server/automation/plugin/grpc"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/plugin/discovery"
	hcp "github.com/hashicorp/go-plugin"
)

type (
	pluginService interface {
		RegisterPlugin(context.Context, string, hcp.Plugin, func() any) (any, error)
		DeregisterPlugin(context.Context, string) error
		SetRawPlugin(context.Context, string, any) error
	}

	automationRegistry interface {
		AddFunctions(ff ...*types.Function)
	}

	af struct {
		ps pluginService
		ar automationRegistry
	}
)

func Function(p pluginService, r automationRegistry) *af {
	return &af{
		ps: p,
		ar: r,
	}
}

// Validate checks if any of the found subfolders of a bundle contains an automation workflow yaml file
func (w *af) Validate(ctx context.Context, b []string) (byte, []string) {
	ss := []string{}

	for _, bb := range b {
		if strings.Contains(strings.ToLower(bb), "automation_function") {
			ga, _ := filepath.Glob(bb + "/*")

			for _, gga := range ga {
				if is, _ := isBinaryFile(gga); is {
					ss = append(ss, gga)
				}
			}
		}
	}

	if len(ss) == 0 {
		return 0, ss
	}

	return discovery.AUTOMATION_FUNCTION, ss
}

func (w *af) Type() byte {
	return discovery.AUTOMATION_FUNCTION
}

func (w *af) Deregister(ctx context.Context, paths []string) error {
	for _, p := range paths {
		w.ps.DeregisterPlugin(ctx, p)
	}
	return nil
}

func (w *af) Register(ctx context.Context, paths []string) error {
	// call plugin service
	for _, p := range paths {
		af := plugin.MakeAutomationFunction()
		rawPlugin, err := w.ps.RegisterPlugin(ctx, p, &grpc.AutomationFunctionPlugin{}, func() any { return af })

		if err != nil {
			return err
		}

		if _, is := rawPlugin.(grpc.AutomationFunction); !is {
			continue
		}

		af.SetPlugin(rawPlugin.(grpc.AutomationFunction))
		w.ar.AddFunctions(af.Generate())
	}

	return nil
}

func isBinaryFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	const maxBytes = 512 // Read the first 512 bytes to determine if it's binary or text
	buffer := make([]byte, maxBytes)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}

	// Check for NUL bytes, a common indicator of binary files
	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return true, nil // File is binary
		}
	}

	// Use the `isText` helper function to further analyze if it's a text file
	return false, nil
}
