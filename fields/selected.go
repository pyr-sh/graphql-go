package fields

import (
	"context"

	"github.com/graph-gophers/graphql-go/exec"
	"github.com/graph-gophers/graphql-go/types"
)

// SelectedFields retrieves the selected fields passed via the context during the request execution
func SelectedFields(ctx context.Context) types.SelectedFields {
	return exec.SelectedFieldsFromContext(ctx)
}
