/*
 * Copyright (c) 2022 - present Kurtosis Technologies Inc.
 * All Rights Reserved.
 */

package engine_server_launcher

import (
	"github.com/avenbreaks/xarchon/engine/launcher/args"
)

type KurtosisBackendConfigSupplier interface {
	// Private because only the launcher needs to call this
	getKurtosisBackendConfig() (args.KurtosisBackendType, interface{})
}
