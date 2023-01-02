package fields

import (
	"context"

	"github.com/graph-gophers/graphql-go/exec"
	"github.com/graph-gophers/graphql-go/types"
)

// SelectedFields retrieves the selected fields passed via the context during the request execution
func SelectedField(ctx context.Context) *types.SelectedField {
	return exec.SelectedFieldFromContext(ctx)
}
