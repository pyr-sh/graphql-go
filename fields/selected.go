package fields

import (
	"context"

	"github.com/graph-gophers/graphql-go/exec"
	"github.com/graph-gophers/graphql-go/types"
)

// Selected retrieves the selected field passed via the context during the request execution
func Selected(ctx context.Context) *types.SelectedField {
	return exec.SelectedFieldFromContext(ctx)
}

// Root retrieves the root parent field of the selected field
func Root(ctx context.Context) *types.SelectedField {
	return exec.RootFieldFromContext(ctx)
}
