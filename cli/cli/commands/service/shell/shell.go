/*
 * Copyright (c) 2021 - present Kurtosis Technologies Inc.
 * All Rights Reserved.
 */

package shell

import (
	"bufio"
	"context"
	"github.com/avenbreaks/xarchon/api/golang/engine/kurtosis_engine_rpc_api_bindings"
	"github.com/avenbreaks/xarchon/api/golang/engine/lib/kurtosis_context"
	"github.com/avenbreaks/xarchon/cli/cli/command_framework/highlevel/enclave_id_arg"
	"github.com/avenbreaks/xarchon/cli/cli/command_framework/highlevel/engine_consuming_kurtosis_command"
	"github.com/avenbreaks/xarchon/cli/cli/command_framework/highlevel/service_identifier_arg"
	"github.com/avenbreaks/xarchon/cli/cli/command_framework/lowlevel/args"
	"github.com/avenbreaks/xarchon/cli/cli/command_framework/lowlevel/flags"
	"github.com/avenbreaks/xarchon/cli/cli/command_str_consts"
	"github.com/avenbreaks/xarchon/container-engine-lib/lib/backend_interface"
	"github.com/avenbreaks/xarchon/container-engine-lib/lib/backend_interface/objects/enclave"
	"github.com/avenbreaks/xarchon/container-engine-lib/lib/backend_interface/objects/service"
	metrics_client "github.com/kurtosis-tech/metrics-library/golang/lib/client"
	"github.com/kurtosis-tech/stacktrace"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
)

const (
	enclaveIdentifierArgKey          = "enclave"
	isEnclaveIdArgOptional           = false
	isEnclaveIdArgGreedy             = false
	commandToRunFlagKey              = "exec"
	commandToRunInsteadOfBashFlagKey = ""

	serviceIdentifierArgKey  = "service"
	isServiceGuidArgOptional = false
	isServiceGuidArgGreedy   = false

	kurtosisBackendCtxKey = "kurtosis-backend"
	engineClientCtxKey    = "engine-client"
)

var ServiceShellCmd = &engine_consuming_kurtosis_command.EngineConsumingKurtosisCommand{
	CommandStr:                command_str_consts.ServiceShellCmdStr,
	ShortDescription:          "Gets a shell on a service",
	LongDescription:           "Starts a shell on the specified service",
	KurtosisBackendContextKey: kurtosisBackendCtxKey,
	EngineClientContextKey:    engineClientCtxKey,
	Flags: []*flags.FlagConfig{
		{
			Key:     commandToRunFlagKey,
			Usage:   "If this flag is used Kurtosis will not run bash/sh on the container by default instead; Kurtosis will run the passed in command. Note if the command being run is multiple words you should wrap it in quotes",
			Type:    flags.FlagType_String,
			Default: commandToRunInsteadOfBashFlagKey,
		},
	},
	Args: []*args.ArgConfig{
		enclave_id_arg.NewEnclaveIdentifierArg(
			enclaveIdentifierArgKey,
			engineClientCtxKey,
			isEnclaveIdArgOptional,
			isEnclaveIdArgGreedy,
		),
		service_identifier_arg.NewServiceIdentifierArg(
			serviceIdentifierArgKey,
			isServiceGuidArgOptional,
			isServiceGuidArgGreedy,
		),
	},
	RunFunc: run,
}

func run(
	ctx context.Context,
	kurtosisBackend backend_interface.KurtosisBackend,
	_ kurtosis_engine_rpc_api_bindings.EngineServiceClient,
	_ metrics_client.MetricsClient,
	flags *flags.ParsedFlags,
	args *args.ParsedArgs,
) error {
	enclaveIdentifier, err := args.GetNonGreedyArg(enclaveIdentifierArgKey)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred getting the enclave identifier using arg key '%v'", enclaveIdentifierArgKey)
	}

	serviceIdentifier, err := args.GetNonGreedyArg(serviceIdentifierArgKey)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred getting the service identifier using arg key '%v'", serviceIdentifierArgKey)
	}

	commandToRunInsteadOfBash, err := flags.GetString(commandToRunFlagKey)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred while getting the command to run using key '%v'", commandToRunFlagKey)
	}

	kurtosisCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred connecting to the local Kurtosis engine")
	}

	enclaveCtx, err := kurtosisCtx.GetEnclaveContext(ctx, enclaveIdentifier)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred while getting enclave context for enclave with identifier '%v' exists", enclaveIdentifier)
	}

	enclaveUuid := enclave.EnclaveUUID(enclaveCtx.GetEnclaveUuid())

	serviceCtx, err := enclaveCtx.GetServiceContext(serviceIdentifier)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred while getting service context for service with identifier '%v'", serviceIdentifier)
	}
	serviceUuid := service.ServiceUUID(serviceCtx.GetServiceUUID())

	conn, err := kurtosisBackend.GetConnectionWithUserService(ctx, enclaveUuid, serviceUuid, commandToRunInsteadOfBash)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred getting connection with user service with UUID '%v' in enclave '%v'", serviceUuid, enclaveIdentifier)
	}
	defer conn.Close()

	newReader := bufio.NewReader(conn)

	// From this point on down, I don't know why it works.... but it does
	// I just followed the solution here: https://stackoverflow.com/questions/58732588/accept-user-input-os-stdin-to-container-using-golang-docker-sdk-interactive-co
	// This channel is being used to know the user exited the ContainerExec
	finishChan := make(chan bool)
	go func() {
		io.Copy(os.Stdout, newReader)
		finishChan <- true
	}()
	go io.Copy(os.Stderr, newReader)
	go io.Copy(conn, os.Stdin)

	stdinFd := int(os.Stdin.Fd())
	var oldState *terminal.State
	if terminal.IsTerminal(stdinFd) {
		oldState, err = terminal.MakeRaw(stdinFd)
		if err != nil {
			// print error
			return stacktrace.Propagate(err, "An error occurred making STDIN stream raw")
		}
		defer terminal.Restore(stdinFd, oldState)
	}

	_ = <-finishChan

	terminal.Restore(stdinFd, oldState)

	return nil
}
