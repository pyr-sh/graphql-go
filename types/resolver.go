package types

import "context"

type ResolverContext struct {
	context.Context
	Operation *OperationDefinition
}
