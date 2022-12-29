package types

import "context"

type ResolverContext struct {
	context.Context
	Definition *ExecutableDefinition
}
