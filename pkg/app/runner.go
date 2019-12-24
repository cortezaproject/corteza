package app

import (
	"context"
	"sync"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	grpcWrap "github.com/cortezaproject/corteza-server/pkg/grpc"
)

type (
	runner struct {
		execStep int

		log *zap.Logger
		opt *Options

		parts []Runnable

		// CLI Commands
		rootCmd *cobra.Command

		// Servers
		httpApiServer httpApiServer
		grpcServer    grpcServer
	}

	httpApiServer interface {
		MountRoutes(mm ...func(chi.Router))
		Serve(context.Context)
	}

	grpcServer interface {
		RegisterServices(func(*grpc.Server))
		Serve(ctx context.Context)
	}

	Runnable interface {
		// Setup
		//
		// Step 0
		Setup(log *zap.Logger, opts *Options) error

		// Initialize initializes all (passive) services
		//
		// Step 1
		Initialize(ctx context.Context) error

		// Upgrade runs database migrations
		// dep: Initialize
		//
		// Step 2
		Upgrade(ctx context.Context) error

		// Activate wakes-up all (active) services, watchers...
		// dep: Upgrade
		//
		// Step 3
		Activate(ctx context.Context) error

		// Provision fills database with data
		// dep: Activate
		//
		// Ran before each cli command (and before REST API)
		//
		// Step 4
		Provision(ctx context.Context) error
	}

	cliCommandRegistrator interface {
		// RegisterCliCommands registers all command line interface commands
		RegisterCliCommands(cmd *cobra.Command)
	}

	apiRouteMounter interface {
		// MountApiRoutes mounts all routes
		MountApiRoutes(router chi.Router)
	}

	grpcRegistrator interface {
		// RegisterGrpcServices registers all gRPC services that are needed
		RegisterGrpcServices(server *grpc.Server)
	}
)

const (
	execStepReady = iota
	execStepSetupDone
	execStepInitializeDone
	execStepUpgradeDone
	execStepActivateDone
	execStepProvisionDone
)

func New(parts ...Runnable) *runner {
	return &runner{
		execStep: execStepReady,
		parts:    parts,
	}
}

// Setup runs setup on all parts
func (r *runner) Setup(log *zap.Logger, opt *Options) (err error) {
	r.log = log
	r.opt = opt
	return r.setup()
}

func (r *runner) setup() (err error) {
	if r.execStep < execStepSetupDone {
		if err = RunSetup(r.log, r.opt, r.parts...); err != nil {
			return
		}

		r.execStep = execStepSetupDone

	}

	return nil
}

// Initialize setups and initializes all parts
func (r *runner) Initialize(ctx context.Context) (err error) {
	if r.execStep < execStepInitializeDone {
		if err = r.setup(); err != nil {
			return
		}

		if err = RunInitialize(ctx, r.parts...); err != nil {
			return
		}

		r.execStep = execStepInitializeDone
	}

	return nil

}

// Upgrade initializes & upgrades all parts
func (r *runner) Upgrade(ctx context.Context) (err error) {
	if r.execStep < execStepUpgradeDone {
		if err = r.Initialize(ctx); err != nil {
			return
		}

		if r.opt.Upgrade.Always {
			if err = RunUpgrade(ctx, r.parts...); err != nil {
				return
			}
		}

		r.execStep = execStepUpgradeDone
	}

	return nil

}

// Activate upgrades and activates all parts
func (r *runner) Activate(ctx context.Context) (err error) {
	if r.execStep < execStepActivateDone {
		if err = r.Upgrade(ctx); err != nil {
			return
		}

		if err = RunActivate(ctx, r.parts...); err != nil {
			return
		}

		r.execStep = execStepActivateDone
	}

	return nil
}

// Provision activates and provisions all parts
func (r *runner) Provision(ctx context.Context) (err error) {
	if r.execStep < execStepProvisionDone {
		if err = r.Activate(ctx); err != nil {
			return
		}

		if r.opt.Provision.Always {
			if err = RunProvision(ctx, r.parts...); err != nil {
				return
			}
		}

		r.execStep = execStepProvisionDone
	}

	return nil
}

func (r *runner) setupHttpApi() {
	// @todo refactor wait-for out of HTTP API server.
	r.httpApiServer = api.New(r.log, r.opt.HTTPServer, r.opt.WaitFor)

	// Mount all HTTP API endpoints
	for _, part := range r.parts {
		if reg, is := part.(apiRouteMounter); is {
			r.httpApiServer.MountRoutes(reg.MountApiRoutes)
		}
	}
}

func (r *runner) setupGRPCServices() {

	r.grpcServer = grpcWrap.New(r.log, r.opt.GRPCServer)

	var hasServices bool

	// Register GRPC services
	for _, part := range r.parts {
		if reg, is := part.(grpcRegistrator); is {
			r.grpcServer.RegisterServices(reg.RegisterGrpcServices)
			hasServices = true
		}
	}

	if !hasServices {
		r.grpcServer = nil
	}
}

// serve starts all servers (HTTP API, GRPC)
func (r *runner) serve(ctx context.Context) (err error) {
	if err = r.Provision(ctx); err != nil {
		return
	}
	r.setupHttpApi()
	r.setupGRPCServices()

	wg := &sync.WaitGroup{}

	{
		wg.Add(1)
		go func(ctx context.Context) {
			r.httpApiServer.Serve(ctx)
			wg.Done()
		}(ctx)
	}

	if r.grpcServer != nil {
		wg.Add(1)
		go func(ctx context.Context) {
			r.grpcServer.Serve(ctx)
			wg.Done()
		}(ctx)

	}

	// Wait for all servers to be done
	wg.Wait()

	return nil
}

// Run will orchestrate all hooks, register commands, server and execute root command
func (r *runner) Run(ctx context.Context) error {
	r.rootCmd = cli.RootCommand(func() error {
		// Run initialization (and all steps before that)
		return r.Initialize(ctx)
	})

	serveCmd := cli.ServeCommand(func() error {
		return r.serve(ctx)
	})

	upgradeCmd := cli.UpgradeCommand(func() (err error) {
		if err = r.Initialize(ctx); err != nil {
			return
		}

		if err = RunUpgrade(ctx, r.parts...); err != nil {
			return
		}

		return
	})

	provisionCmd := cli.ProvisionCommand(func() (err error) {
		if err = r.Activate(ctx); err != nil {
			return
		}

		if err = RunProvision(ctx, r.parts...); err != nil {
			return
		}

		return
	})

	r.rootCmd.AddCommand(
		serveCmd,
		upgradeCmd,
		provisionCmd,
	)

	// Register CLI commands from all parts (when compatible)
	for _, part := range r.parts {
		if reg, is := part.(cliCommandRegistrator); is {
			reg.RegisterCliCommands(r.rootCmd)
		}
	}

	return r.rootCmd.Execute()
}

// Run is an application runner
//
// It accepts set of Runnables and calls hooks on each
func Run(log *zap.Logger, opt *Options, parts ...Runnable) {
	r := New(parts...)
	r.Setup(log, opt)

	ctx := cli.Context()

	cli.HandleError(r.Run(ctx))
}
