package fields

import (
	"context"

	"github.com/graph-gophers/graphql-go/exec"
	"github.com/graph-gophers/graphql-go/types"
)

// SelectedField retrieves the selected field passed via the context during the request execution
func SelectedField(ctx context.Context) *types.SelectedField {
	return exec.SelectedFieldFromContext(ctx)
}

// SelectedField retrieves the root parent field of the selected field
func RootField(ctx context.Context) *types.SelectedField {
	return exec.RootFieldFromContext(ctx)
}
