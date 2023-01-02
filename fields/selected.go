package fields

import (
	"context"

	"github.com/graph-gophers/graphql-go/exec"
	"github.com/graph-gophers/graphql-go/types"
)

// SelectedFieldsFromContext retrieves the selected fields passed via the context during the request
// execution
func SelectedFieldsFromContext(ctx context.Context) []*types.SelectedField {
	return exec.SelectedFieldsFromContext(ctx)
}
