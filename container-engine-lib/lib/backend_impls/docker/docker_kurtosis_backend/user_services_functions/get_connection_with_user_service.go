package user_service_functions

import (
	"context"
	"fmt"
	"github.com/avenbreaks/xarchon/container-engine-lib/lib/backend_impls/docker/docker_kurtosis_backend/shared_helpers"
	"github.com/avenbreaks/xarchon/container-engine-lib/lib/backend_impls/docker/docker_manager"
	"github.com/avenbreaks/xarchon/container-engine-lib/lib/backend_interface/objects/enclave"
	"github.com/avenbreaks/xarchon/container-engine-lib/lib/backend_interface/objects/service"
	"github.com/kurtosis-tech/stacktrace"
	"net"
)

const (
	defaultCommandToRunInsteadOfBash               = ""
	commandToRunIndex                              = 2
	lineToEchoWhileUserWaitsForTheirCommandToBeRun = `echo "Running '%v'" && %v`
)

// We'll try to use the nicer-to-use shells first before we drop down to the lower shells
// If the user passes a different commandToRunInsteadOfBash than defaultCommandToRunInsteadOfBash we try to run that instead
var commandToRunWhenCreatingUserServiceShell = []string{
	"sh",
	"-c",
	`if command -v 'bash' > /dev/null; then
		echo "Found bash on container; creating bash shell..."; bash; 
       else 
		echo "No bash found on container; dropping down to sh shell..."; sh; 
	fi`,
}

func GetConnectionWithUserService(ctx context.Context, enclaveId enclave.EnclaveUUID, serviceUuid service.ServiceUUID, dockerManager *docker_manager.DockerManager, commandToRunInsteadOfBash string) (net.Conn, error) {
	_, serviceDockerResources, err := shared_helpers.GetSingleUserServiceObjAndResourcesNoMutex(ctx, enclaveId, serviceUuid, dockerManager)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting service object and Docker resources for service '%v' in enclave '%v'", serviceUuid, enclaveId)
	}
	container := serviceDockerResources.ServiceContainer

	if commandToRunInsteadOfBash != defaultCommandToRunInsteadOfBash {
		commandToRunWhenCreatingUserServiceShell[commandToRunIndex] = fmt.Sprintf(lineToEchoWhileUserWaitsForTheirCommandToBeRun, commandToRunInsteadOfBash, commandToRunInsteadOfBash)
	}

	hijackedResponse, err := dockerManager.CreateContainerExec(ctx, container.GetId(), commandToRunWhenCreatingUserServiceShell)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting a shell on user service with UUID '%v' in enclave '%v'", serviceUuid, enclaveId)
	}

	newConnection := hijackedResponse.Conn

	return newConnection, nil
}
