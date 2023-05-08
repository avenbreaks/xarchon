package kurtosis_helper

import (
	"github.com/avenbreaks/xarchon/core/server/api_container/server/startosis_engine/kurtosis_starlark_framework"
	"github.com/avenbreaks/xarchon/core/server/api_container/server/startosis_engine/startosis_errors"
	"go.starlark.net/starlark"
)

type kurtosisHelperInternal struct {
	*kurtosis_starlark_framework.KurtosisBaseBuiltinInternal

	capabilities KurtosisHelperCapabilities
}

func newKurtosisHelperInternal(wrappedBuiltin *kurtosis_starlark_framework.KurtosisBaseBuiltinInternal, capabilities KurtosisHelperCapabilities) *kurtosisHelperInternal {
	return &kurtosisHelperInternal{
		KurtosisBaseBuiltinInternal: wrappedBuiltin,

		capabilities: capabilities,
	}
}

func (builtin *kurtosisHelperInternal) interpret() (starlark.Value, *startosis_errors.InterpretationError) {
	return builtin.capabilities.Interpret(builtin.GetArguments())
}
