package dump

import (
	"context"
	"fmt"
	"github.com/avenbreaks/xarchon/api/golang/engine/kurtosis_engine_rpc_api_bindings"
	"github.com/avenbreaks/xarchon/cli/cli/command_framework/highlevel/engine_consuming_kurtosis_command"
	"github.com/avenbreaks/xarchon/cli/cli/command_framework/lowlevel/args"
	"github.com/avenbreaks/xarchon/cli/cli/command_framework/lowlevel/flags"
	"github.com/avenbreaks/xarchon/cli/cli/command_str_consts"
	"github.com/avenbreaks/xarchon/container-engine-lib/lib/backend_interface"
	metrics_client "github.com/kurtosis-tech/metrics-library/golang/lib/client"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	outputDirpathArg = "output-dirpath"

	kurtosisBackendCtxKey = "kurtosis-backend"
	engineClientCtxKey    = "engine-client"

	defaultKurtosisDumpDir = "kurtosis-dump"
	outputDirIsOptional    = true
	dumpDirTimeDelimiter   = "--"
)

var KurtosisDump = &engine_consuming_kurtosis_command.EngineConsumingKurtosisCommand{
	CommandStr:                command_str_consts.KurtosisDumpCmdStr,
	ShortDescription:          "Dumps entire Kurtosis State",
	LongDescription:           "Dumps entire Kurtosis State to the given directory",
	KurtosisBackendContextKey: kurtosisBackendCtxKey,
	EngineClientContextKey:    engineClientCtxKey,
	Flags:                     nil,
	// TODO perhaps add an --enclave flag here and deprecate enclave dump but that clashes with engine dumping
	Args: []*args.ArgConfig{
		{
			Key:          outputDirpathArg,
			DefaultValue: defaultKurtosisDumpDir,
			IsOptional:   outputDirIsOptional,
		},
	},
	RunFunc: run,
}

func run(
	ctx context.Context,
	kurtosisBackend backend_interface.KurtosisBackend,
	_ kurtosis_engine_rpc_api_bindings.EngineServiceClient,
	_ metrics_client.MetricsClient,
	_ *flags.ParsedFlags,
	args *args.ParsedArgs,
) error {
	outputDirPath, err := args.GetNonGreedyArg(outputDirpathArg)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred getting output dirpath using arg key '%v'", outputDirpathArg)
	}

	if outputDirPath == defaultKurtosisDumpDir {
		outputDirPath = fmt.Sprintf("%s%s%d", outputDirPath, dumpDirTimeDelimiter, time.Now().Unix())
	}

	if err := kurtosisBackend.DumpKurtosis(ctx, outputDirPath); err != nil {
		return stacktrace.Propagate(err, "An error occurred dumping all of Kurtosis")
	}

	logrus.Infof("Dumped all of Kurtosis to '%v'", outputDirPath)
	return nil
}
