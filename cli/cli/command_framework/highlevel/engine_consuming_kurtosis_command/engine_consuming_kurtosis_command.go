package engine_consuming_kurtosis_command

import (
	"context"
	"github.com/avenbreaks/xarchon/api/golang/engine/kurtosis_engine_rpc_api_bindings"
	"github.com/avenbreaks/xarchon/cli/cli/command_framework/lowlevel"
	"github.com/avenbreaks/xarchon/cli/cli/command_framework/lowlevel/args"
	"github.com/avenbreaks/xarchon/cli/cli/command_framework/lowlevel/flags"
	"github.com/avenbreaks/xarchon/cli/cli/defaults"
	"github.com/avenbreaks/xarchon/cli/cli/helpers/engine_manager"
	"github.com/avenbreaks/xarchon/cli/cli/helpers/metrics_client_factory"
	"github.com/avenbreaks/xarchon/cli/cli/helpers/portal_manager"
	"github.com/avenbreaks/xarchon/container-engine-lib/lib/backend_interface"
	"github.com/avenbreaks/xarchon/contexts-config-store/store"
	metrics_client "github.com/kurtosis-tech/metrics-library/golang/lib/client"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	engineClientCloseFuncCtxKey     = "engine-client-close-func"
	metricsClientKey                = "metrics-client-key"
	metricsClientClosingFunctionKey = "metrics-client-closing-func"
)

// This is a convenience KurtosisCommand for commands that interact with the engine
type EngineConsumingKurtosisCommand struct {
	// The string for the command (e.g. "inspect" or "ls")
	CommandStr string

	// Will be used when displaying the command for tab completion
	ShortDescription string

	LongDescription string

	// The name of the key that will be set during PreValidationAndRun where the KurtosisBackend can be found
	KurtosisBackendContextKey string

	// TODO Replace with KurtosisContext!!! This will:
	//  1) be easier to work with and
	//  2) force us to use the same SDK we give to users, so there's no "secret" or "private" API, which will force
	//     us to improve the SDK
	// The name of the key that will be set during PreValidationAndRun where the engine client will be made available
	EngineClientContextKey string

	// Order isn't important here
	Flags []*flags.FlagConfig

	Args []*args.ArgConfig

	RunFunc func(
		ctx context.Context,
		// TODO This is a hack that's only here temporarily because we have commands that use KurtosisBackend directly (they
		//  should not), and EngineConsumingKurtosisCommand therefore needs to provide them with a KurtosisBackend. Once all our
		//  commands only access the Kurtosis APIs, we can remove this.
		kurtosisBackend backend_interface.KurtosisBackend,
		engineClient kurtosis_engine_rpc_api_bindings.EngineServiceClient,
		client metrics_client.MetricsClient,
		flags *flags.ParsedFlags,
		args *args.ParsedArgs,
	) error
}

func (cmd *EngineConsumingKurtosisCommand) MustGetCobraCommand() *cobra.Command {
	// Validation
	if cmd.KurtosisBackendContextKey == engineClientCloseFuncCtxKey {
		panic(stacktrace.NewError(
			"Kurtosis backend context key '%v' on command '%v' is equal to engine client close function context key '%v'; this is a bug in Kurtosis!",
			cmd.KurtosisBackendContextKey,
			cmd.CommandStr,
			engineClientCloseFuncCtxKey,
		))
	}

	if cmd.EngineClientContextKey == engineClientCloseFuncCtxKey {
		panic(stacktrace.NewError(
			"Engine client context key '%v' on command '%v' is equal to engine client close function context key '%v'; this is a bug in Kurtosis!",
			cmd.EngineClientContextKey,
			cmd.CommandStr,
			engineClientCloseFuncCtxKey,
		))
	}
	if cmd.KurtosisBackendContextKey == cmd.EngineClientContextKey {
		panic(stacktrace.NewError(
			"Kurtosis backend context key '%v' on command '%v' is equal to engine client close function context key '%v'; this is a bug in Kurtosis!",
			cmd.KurtosisBackendContextKey,
			cmd.CommandStr,
			cmd.EngineClientContextKey,
		))
	}

	if cmd.KurtosisBackendContextKey == metricsClientClosingFunctionKey {
		panic(stacktrace.NewError(
			"Kurtosis backend context key '%v' on command '%v' is equal to metrics client close function context key '%v'; this is a bug in Kurtosis!",
			cmd.KurtosisBackendContextKey,
			cmd.CommandStr,
			metricsClientClosingFunctionKey,
		))
	}

	if cmd.EngineClientContextKey == metricsClientClosingFunctionKey {
		panic(stacktrace.NewError(
			"Engine client context key '%v' on command '%v' is equal to metrics client close function context key '%v'; this is a bug in Kurtosis!",
			cmd.EngineClientContextKey,
			cmd.CommandStr,
			metricsClientClosingFunctionKey,
		))
	}

	if cmd.KurtosisBackendContextKey == metricsClientKey {
		panic(stacktrace.NewError(
			"Kurtosis backend context key '%v' on command '%v' is equal to metrics client context key '%v'; this is a bug in Kurtosis!",
			cmd.KurtosisBackendContextKey,
			cmd.CommandStr,
			metricsClientKey,
		))
	}

	if cmd.EngineClientContextKey == metricsClientKey {
		panic(stacktrace.NewError(
			"Engine client context key '%v' on command '%v' is equal to metrics client context key '%v'; this is a bug in Kurtosis!",
			cmd.EngineClientContextKey,
			cmd.CommandStr,
			metricsClientKey,
		))
	}

	lowlevelCmd := &lowlevel.LowlevelKurtosisCommand{
		CommandStr:               cmd.CommandStr,
		ShortDescription:         cmd.ShortDescription,
		LongDescription:          cmd.LongDescription,
		Flags:                    cmd.Flags,
		Args:                     cmd.Args,
		PreValidationAndRunFunc:  cmd.getSetupFunc(),
		RunFunc:                  cmd.getRunFunc(),
		PostValidationAndRunFunc: cmd.getTeardownFunc(),
	}

	return lowlevelCmd.MustGetCobraCommand()
}

func (cmd *EngineConsumingKurtosisCommand) getSetupFunc() func(context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		result := ctx

		metricsClient, metricsClientCloser, err := metrics_client_factory.GetMetricsClient()
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred while creating metrics client")
		}

		engineManager, err := engine_manager.NewEngineManager(ctx)
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred getting an engine manager.")
		}

		// TODO This is a hack that's only here temporarily because we have commands that use KurtosisBackend directly (they
		//  should not), and EngineConsumingKurtosisCommand therefore needs to provide them with a KurtosisBackend. Once all our
		//  commands only access the Kurtosis APIs, we can remove this.
		kurtosisBackend := engineManager.GetKurtosisBackend()

		engineClient, closeClientFunc, err := engineManager.StartEngineIdempotentlyWithDefaultVersion(ctx, defaults.DefaultEngineLogLevel)
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred creating a new Kurtosis engine client")
		}
		result = context.WithValue(result, cmd.EngineClientContextKey, engineClient)
		result = context.WithValue(result, engineClientCloseFuncCtxKey, closeClientFunc)
		result = context.WithValue(result, cmd.KurtosisBackendContextKey, kurtosisBackend)
		result = context.WithValue(result, metricsClientKey, metricsClient)
		result = context.WithValue(result, metricsClientClosingFunctionKey, metricsClientCloser)

		return result, nil
	}
}

func (cmd *EngineConsumingKurtosisCommand) getRunFunc() func(context.Context, *flags.ParsedFlags, *args.ParsedArgs) error {
	// Do the gruntwork necessary to give a Kurtosis dev the Docker manager & engine client without them
	// needing to think about how they should get it
	return func(ctx context.Context, flags *flags.ParsedFlags, args *args.ParsedArgs) error {
		uncastedEngineClient := ctx.Value(cmd.EngineClientContextKey)
		if uncastedEngineClient == nil {
			return stacktrace.NewError("Expected an engine client to have been stored in the context under key '%v', but none was found; this is a bug in Kurtosis!", cmd.EngineClientContextKey)
		}
		engineClient, ok := uncastedEngineClient.(kurtosis_engine_rpc_api_bindings.EngineServiceClient)
		if !ok {
			return stacktrace.NewError("Found an object that should be the engine client stored in the context under key '%v', but this object wasn't of the correct type", cmd.EngineClientContextKey)
		}

		uncastedKurtosisBackend := ctx.Value(cmd.KurtosisBackendContextKey)
		if uncastedKurtosisBackend == nil {
			return stacktrace.NewError("Expected Kurtosis backend to have been stored in the context under key '%v', but none was found; this is a bug in Kurtosis!", cmd.KurtosisBackendContextKey)
		}
		kurtosisBackend, ok := uncastedKurtosisBackend.(backend_interface.KurtosisBackend)
		if !ok {
			return stacktrace.NewError("Found an object that should be the Kurtosis backend stored in the context under key '%v', but this object wasn't of the correct type", cmd.KurtosisBackendContextKey)
		}

		uncastedMetricsClient := ctx.Value(metricsClientKey)
		if uncastedMetricsClient == nil {
			return stacktrace.NewError("Expected Metrics Client to have been stored in the context under key '%v', but none was found; this is a bug in Kurtosis", metricsClientKey)
		}
		metricsClient, ok := uncastedMetricsClient.(metrics_client.MetricsClient)
		if !ok {
			return stacktrace.NewError("Found an object that should be the metrics client stored in the context under key '%v', but this object wasn't of the correct type", metricsClientKey)
		}

		currentContext, err := store.GetContextsConfigStore().GetCurrentContext()
		if err != nil {
			return stacktrace.Propagate(err, "Error fetching current context information")
		}

		portalManager := portal_manager.NewPortalManager()
		if store.IsRemote(currentContext) && !portalManager.IsReachable() {
			// TODO: add command to start it when it's implemented
			return stacktrace.NewError("Selected context is a remote context but Kurtosis Portal daemon is " +
				"not reachable. Make sure it is started and re-run the command")
		}

		if err = cmd.RunFunc(ctx, kurtosisBackend, engineClient, metricsClient, flags, args); err != nil {
			return stacktrace.Propagate(
				err,
				"An error occurred calling the run function for command '%v'",
				cmd.CommandStr,
			)
		}

		return nil
	}
}

func (cmd *EngineConsumingKurtosisCommand) getTeardownFunc() func(ctx context.Context) {
	return func(ctx context.Context) {
		uncastedEngineClientCloseFunc := ctx.Value(engineClientCloseFuncCtxKey)
		if uncastedEngineClientCloseFunc != nil {
			engineClientCloseFunc, ok := uncastedEngineClientCloseFunc.(func() error)
			if ok {
				if err := engineClientCloseFunc(); err != nil {
					logrus.Warnf("We tried to close the engine client after we're done using it, but doing so threw an error:\n%v", err)
				}
			} else {
				logrus.Errorf("Expected the object at context key '%v' to be an engine client close function, but it wasn't; this is a bug in Kurtosis!", engineClientCloseFuncCtxKey)
			}
		} else {
			logrus.Errorf(
				"Expected to find an engine client close function during teardown at context key '%v', but none was found; this is a bug in Kurtosis!",
				engineClientCloseFuncCtxKey,
			)
		}

		uncastedMetricsClientCloser := ctx.Value(metricsClientClosingFunctionKey)
		if uncastedMetricsClientCloser != nil {
			metricsClientCloser, ok := uncastedMetricsClientCloser.(func() error)
			if ok {
				if err := metricsClientCloser(); err != nil {
					logrus.Warnf("An error occurred while closing the metrics client\n%s", err)
				}
			} else {
				logrus.Errorf("Expected the object at context key '%v' to be a metrics client close function, but it wasn't; this is a bug in Kurtosis!", metricsClientClosingFunctionKey)
			}
		} else {
			logrus.Errorf(
				"Expected to metrics client close function during teardown at context key '%v', but none was found; this is a bug in Kurtosis!",
				metricsClientClosingFunctionKey,
			)
		}
	}
}
