package centralized_logs

import (
	"context"
	"github.com/avenbreaks/xarchon/container-engine-lib/lib/backend_interface/objects/enclave"
	"github.com/avenbreaks/xarchon/container-engine-lib/lib/backend_interface/objects/service"
	"github.com/avenbreaks/xarchon/engine/server/engine/centralized_logs/logline"
)

type LogsDatabaseClient interface {
	StreamUserServiceLogs(
		ctx context.Context,
		enclaveUuid enclave.EnclaveUUID,
		userServiceUuids map[service.ServiceUUID]bool,
		conjunctiveLogLineFilters logline.ConjunctiveLogLineFilters,
		shouldFollowLogs bool,
	) (
		userServiceLogsByServiceUuidChan chan map[service.ServiceUUID][]logline.LogLine,
		errChan chan error,
		cancelFunc context.CancelFunc,
		err error,
	)

	FilterExistingServiceUuids(
		ctx context.Context,
		enclaveUuid enclave.EnclaveUUID,
		userServiceUuids map[service.ServiceUUID]bool,
	) (
		map[service.ServiceUUID]bool,
		error,
	)
}
