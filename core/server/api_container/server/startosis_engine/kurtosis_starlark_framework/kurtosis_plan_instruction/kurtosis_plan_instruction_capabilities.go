package kurtosis_plan_instruction

import (
	"context"
	"github.com/avenbreaks/xarchon/core/server/api_container/server/startosis_engine/kurtosis_starlark_framework/builtin_argument"
	"github.com/avenbreaks/xarchon/core/server/api_container/server/startosis_engine/startosis_errors"
	"github.com/avenbreaks/xarchon/core/server/api_container/server/startosis_engine/startosis_validator"
	"go.starlark.net/starlark"
)

type KurtosisPlanInstructionCapabilities interface {
	Interpret(arguments *builtin_argument.ArgumentValuesSet) (starlark.Value, *startosis_errors.InterpretationError)

	Validate(arguments *builtin_argument.ArgumentValuesSet, validatorEnvironment *startosis_validator.ValidatorEnvironment) *startosis_errors.ValidationError

	Execute(ctx context.Context, arguments *builtin_argument.ArgumentValuesSet) (string, error)
}
