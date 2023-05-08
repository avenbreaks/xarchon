/*
 * Copyright (c) 2022 - present Kurtosis Technologies Inc.
 * All Rights Reserved.
 */

package api_container_launcher

import (
	"github.com/avenbreaks/xarchon/core/launcher/args"
	"github.com/avenbreaks/xarchon/core/launcher/args/kurtosis_backend_config"
)

type KubernetesBackendConfigSupplier struct {
}

func NewKubernetesKurtosisBackendConfigSupplier() KubernetesBackendConfigSupplier {
	return KubernetesBackendConfigSupplier{
		// More fields here when needed
	}
}

func (backendConfigSupplier KubernetesBackendConfigSupplier) getKurtosisBackendConfig() (args.KurtosisBackendType, interface{}) {
	return args.KurtosisBackendType_Kubernetes, kurtosis_backend_config.KubernetesBackendConfig{}
}
