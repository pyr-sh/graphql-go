package graphql_test

import (
	"context"
	"testing"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/fields"
	"github.com/graph-gophers/graphql-go/gqltesting"
)

// Nested fields with aliases

type selectedFieldsResolver1 struct {
	T *testing.T
}

func (r *selectedFieldsResolver1) Test(ctx context.Context) selectedFieldsTestResolver1 {
	sel := fields.Selected(ctx)
	if sel == nil {
		r.T.Error("selected field is nil")
		return selectedFieldsTestResolver1{}
	}

	if !(sel.Name == "test" && sel.TypeName == "Test!" && len(sel.Fields) == 4) {
		r.T.Error("test field comparison failed")
		return selectedFieldsTestResolver1{}
	}

	a, b, c, d := sel.Fields[0], sel.Fields[1], sel.Fields[2], sel.Fields[3]

	// a
	if !(a.Name == "a" && a.TypeName == "String!" && a.Alias == "a" && len(a.Directives) == 1) {
		r.T.Error("a field comparison failed")
		return selectedFieldsTestResolver1{}
	}
	if !(b.Name == "b" && b.TypeName == "String" && b.Alias == "b" && len(b.Directives) == 1) {
		r.T.Error("b field comparison failed")
		return selectedFieldsTestResolver1{}
	}
	aDir := a.Directives.Get("testFieldDirective")
	if aDir == nil {
		r.T.Error("testFieldDirective doesn't exist on a")
		return selectedFieldsTestResolver1{}
	}
	aArgValue, ok := a.Args["value"]
	if !ok {
		r.T.Error("value argument on a doesn't exist")
		return selectedFieldsTestResolver1{}
	}
	if aArgValue != "aValue" {
		r.T.Errorf("unexpected %s value of a's value argument", aArgValue)
		return selectedFieldsTestResolver1{}
	}

	// b
	if !(b.Name == "b" && b.TypeName == "String" && b.Alias == "b" && len(b.Directives) == 1) {
		r.T.Error("b field comparison failed")
		return selectedFieldsTestResolver1{}
	}
	bDir := b.Directives.Get("testFieldDefinitionDirective")
	if bDir == nil {
		r.T.Error("testFieldDefinitionDirective doesn't exist on b")
		return selectedFieldsTestResolver1{}
	}
	bDirArg, ok := bDir.Arguments.Get("description")
	if !ok {
		r.T.Error("description arg doesn't exist on testFieldDefinitionDirective")
		return selectedFieldsTestResolver1{}
	}
	bDirArgRaw := bDirArg.Deserialize(nil)
	bDirArgV, ok := bDirArgRaw.(string)
	if !ok {
		r.T.Errorf("expected string description arg, got %T", bDirArgRaw)
		return selectedFieldsTestResolver1{}
	}
	if bDirArgV != "Test description" {
		r.T.Errorf("unexpected directive arg description %s", bDirArgV)
		return selectedFieldsTestResolver1{}
	}

	// c
	if !(c.Name == "c" && c.TypeName == "Int!" && c.Alias == "c" && len(c.Directives) == 1) {
		r.T.Error("c field comparison failed")
		return selectedFieldsTestResolver1{}
	}
	cDir := c.Directives.Get("testFieldDefinitionDirective")
	if cDir == nil {
		r.T.Error("testFieldDefinitionDirective doesn't exist on c")
		return selectedFieldsTestResolver1{}
	}
	cDirArg, ok := cDir.Arguments.Get("description")
	if !ok {
		r.T.Error("description arg doesn't exist on testFieldDefinitionDirective")
		return selectedFieldsTestResolver1{}
	}
	cDirArgRaw := cDirArg.Deserialize(nil)
	cDirArgV, ok := cDirArgRaw.(string)
	if !ok {
		r.T.Errorf("expected string description arg, got %T", cDirArgRaw)
		return selectedFieldsTestResolver1{}
	}
	if cDirArgV != "Default description" {
		r.T.Errorf("unexpected directive arg description %s", cDirArgV)
		return selectedFieldsTestResolver1{}
	}

	// d
	if !(d.Name == "d" && d.TypeName == "D!" && d.Alias == "aliasedD" && len(d.Fields) == 1 &&
		d.Fields[0].Name == "value" && d.Fields[0].TypeName == "String!") {
		r.T.Error("c field comparison failed")
		return selectedFieldsTestResolver1{}
	}

	return selectedFieldsTestResolver1{}
}

type selectedFieldsTestResolver1 struct{}

func (r selectedFieldsTestDResolver) Value() string { return "value" }

type aArgs struct {
	Value *string
}

func (r selectedFieldsTestResolver1) A(aArgs) string { return "a" }
func (r selectedFieldsTestResolver1) B() *string     { return nil }
func (r selectedFieldsTestResolver1) C() int32       { return 1 }
func (r selectedFieldsTestResolver1) D() selectedFieldsTestDResolver {
	return selectedFieldsTestDResolver{}
}

type selectedFieldsTestDResolver struct{}

func TestSelectedFieldNestedWithAliases(t *testing.T) {
	schemaString := `
	schema {
		query: Query
	}

	type Query {
		test: Test!
	}

	type D {
		value: String!
	}
	
	type Test {
		a(value: String): String!
		b: String @testFieldDefinitionDirective(description: "Test description")
		c: Int! @testFieldDefinitionDirective
		d: D!
	}
	
	directive @testFieldDefinitionDirective(description: String = "Default description") on FIELD_DEFINITION
	directive @testFieldDirective on FIELD`

	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Schema: graphql.MustParseSchema(schemaString, &selectedFieldsResolver1{t}),
			Query: `
				{
					test {
						a(value: "aValue") @testFieldDirective
						b
						c
						aliasedD: d {
							value 
						}
					}
				}
			`,
			IgnoreResult: true,
		},
	})
}
