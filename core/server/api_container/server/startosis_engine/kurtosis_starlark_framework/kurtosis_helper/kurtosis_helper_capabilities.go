package kurtosis_helper

import (
	"github.com/avenbreaks/xarchon/core/server/api_container/server/startosis_engine/kurtosis_starlark_framework/builtin_argument"
	"github.com/avenbreaks/xarchon/core/server/api_container/server/startosis_engine/startosis_errors"
	"go.starlark.net/starlark"
)

type KurtosisHelperCapabilities interface {
	Interpret(arguments *builtin_argument.ArgumentValuesSet) (starlark.Value, *startosis_errors.InterpretationError)
}
