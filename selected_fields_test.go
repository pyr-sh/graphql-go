package graphql_test

import (
	"context"
	"testing"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/fields"
	"github.com/graph-gophers/graphql-go/gqltesting"
)

// Nested fields with aliases, arguments and directives

type selectedFieldsResolver1 struct {
	T *testing.T
}

func (r *selectedFieldsResolver1) Test(ctx context.Context) selectedFieldsTestResolver1 {
	sel := fields.Selected(ctx)
	if sel == nil {
		failWithError(r.T, "selected field is nil")
	}

	if !(sel.Name == "test" && sel.TypeName == "Test!" && len(sel.Fields) == 4) {
		failWithError(r.T, "test field comparison failed")
	}

	a, b, c, d := sel.Fields[0], sel.Fields[1], sel.Fields[2], sel.Fields[3]

	// a
	if !(a.Name == "a" && a.TypeName == "String!" && a.Alias == "a" && len(a.Directives) == 1) {
		failWithError(r.T, "a field comparison failed")
	}
	if !(b.Name == "b" && b.TypeName == "String" && b.Alias == "b" && len(b.Directives) == 1) {
		failWithError(r.T, "b field comparison failed")
	}
	aDir := a.Directives.Get("testFieldDirective")
	if aDir == nil {
		failWithError(r.T, "testFieldDirective doesn't exist on a")
	}
	aArgValue, ok := a.Args["value"]
	if !ok {
		failWithError(r.T, "value argument on a doesn't exist")
	}
	if aArgValue != "aValue" {
		failWithError(r.T, "unexpected %s value of a's value argument", aArgValue)
	}

	// b
	if !(b.Name == "b" && b.TypeName == "String" && b.Alias == "b" && len(b.Directives) == 1) {
		failWithError(r.T, "b field comparison failed")
	}
	bDir := b.Directives.Get("testFieldDefinitionDirective")
	if bDir == nil {
		failWithError(r.T, "testFieldDefinitionDirective doesn't exist on b")
	}
	bDirArg, ok := bDir.Arguments.Get("description")
	if !ok {
		failWithError(r.T, "description arg doesn't exist on testFieldDefinitionDirective")
	}
	bDirArgRaw := bDirArg.Deserialize(nil)
	bDirArgV, ok := bDirArgRaw.(string)
	if !ok {
		failWithError(r.T, "expected string description arg, got %T", bDirArgRaw)
	}
	if bDirArgV != "Test description" {
		failWithError(r.T, "unexpected directive arg description %s", bDirArgV)
	}

	// c
	if !(c.Name == "c" && c.TypeName == "Int!" && c.Alias == "c" && len(c.Directives) == 1) {
		failWithError(r.T, "c field comparison failed")
	}
	cDir := c.Directives.Get("testFieldDefinitionDirective")
	if cDir == nil {
		failWithError(r.T, "testFieldDefinitionDirective doesn't exist on c")
	}
	cDirArg, ok := cDir.Arguments.Get("description")
	if !ok {
		failWithError(r.T, "description arg doesn't exist on testFieldDefinitionDirective")
	}
	cDirArgRaw := cDirArg.Deserialize(nil)
	cDirArgV, ok := cDirArgRaw.(string)
	if !ok {
		failWithError(r.T, "expected string description arg, got %T", cDirArgRaw)
	}
	if cDirArgV != "Default description" {
		failWithError(r.T, "unexpected directive arg description %s", cDirArgV)
	}

	// d
	if !(d.Name == "d" && d.TypeName == "D!" && d.Alias == "aliasedD" && len(d.Fields) == 1 &&
		d.Fields[0].Name == "value" && d.Fields[0].TypeName == "String!") {
		failWithError(r.T, "c field comparison failed")
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

func TestSelectedFieldsNestedAliasesArgsDirectives(t *testing.T) {
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

// Union query result

type selectedFieldsResolver2 struct {
	T *testing.T
}

func (r *selectedFieldsResolver2) Test(ctx context.Context) []selectedFieldsTestUnion {
	f := fields.Selected(ctx)
	if f == nil {
		failWithError(r.T, "selected field is nil")
	}
	if !(f.Name == "test" && f.TypeName == "[Test!]!") {
		failWithError(r.T, "test field validation failed")
	}
	if len(f.Fields) != 4 {
		failWithError(r.T, "expected 4 fields, got %d", len(f.Fields))
	}

	aID, bID, aValue, bIndex := f.Fields[0], f.Fields[1], f.Fields[2], f.Fields[3]
	if !(aID.Name == "id" && aID.TypeName == "String!" && aID.AssertedTypeName == "A") {
		r.T.FailNow()
	}
	if !(bID.Name == "id" && bID.TypeName == "String!" && bID.AssertedTypeName == "B") {
		r.T.FailNow()
	}
	if !(aValue.Name == "value" && aValue.TypeName == "Int!" && aValue.AssertedTypeName == "A") {
		r.T.FailNow()
	}
	if !(bIndex.Name == "index" && bIndex.TypeName == "Int!" && bIndex.AssertedTypeName == "B") {
		r.T.FailNow()
	}

	return []selectedFieldsTestUnion{}
}

type aOrB interface {
	ID() string
}
type selectedFieldsTestUnion struct{ aOrB }

type selectedFieldsTestAResolver2 struct{}

func (r selectedFieldsTestAResolver2) ID() string   { return "" }
func (r selectedFieldsTestAResolver2) Value() int32 { return 0 }

type selectedFieldsTestBResolver2 struct{}

func (r selectedFieldsTestBResolver2) ID() string   { return "" }
func (r selectedFieldsTestBResolver2) Index() int32 { return 0 }

func (r selectedFieldsTestUnion) ToA() (*selectedFieldsTestAResolver2, bool) {
	v, ok := r.aOrB.(selectedFieldsTestAResolver2)
	if !ok {
		return nil, false
	}
	return &v, true
}

func (r selectedFieldsTestUnion) ToB() (*selectedFieldsTestBResolver2, bool) {
	v, ok := r.aOrB.(selectedFieldsTestBResolver2)
	if !ok {
		return nil, false
	}
	return &v, true
}

func TestSelectedFieldsUnionQueryResult(t *testing.T) {
	schemaString := `
	schema {
		query: Query
	}

	type Query {
		test: [Test!]!
	}

	union Test = A | B

	type A implements Identifiable {
		id: String!
		value: Int!
	}

	type B implements Identifiable {
		id: String!
		index: Int!
	}

	interface Identifiable {
		id: String!
	}`

	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Schema: graphql.MustParseSchema(schemaString, &selectedFieldsResolver2{t}),
			Query: `
				{
					test {
						... on Identifiable {
							id
						}
						... on A {
							value
						}
						... on B {
							index
						}
					}
				}
			`,
			IgnoreResult: true,
		},
	})
}

// Interface query result

type selectedFieldsResolver3 struct {
	T *testing.T
}

func (r *selectedFieldsResolver3) Test(ctx context.Context) []identifiable {
	f := fields.Selected(ctx)
	if f == nil {
		failWithError(r.T, "selected field is nil")
	}
	if len(f.Fields) != 1 {
		failWithError(r.T, "expected 1 field, got %d", len(f.Fields))
	}
	if !(f.Name == "test" && f.TypeName == "[Identifiable!]!" && f.AssertedTypeName == "") {
		failWithError(r.T, "test field validation field")
	}
	field := f.Fields[0]
	if !(field.Name == "id" && field.TypeName == "String!") {
		failWithError(r.T, "id field validation field")
	}

	return []identifiable{}
}

type identifiable interface {
	ID() string
	ToA() (*selectedFieldsTestAResolver2, bool)
	ToB() (*selectedFieldsTestBResolver2, bool)
}

func TestSelectedFieldsInterfaceQueryResult(t *testing.T) {
	schemaString := `
	schema {
		query: Query
	}

	type Query {
		test: [Identifiable!]!
	}

	interface Identifiable {
		id: String!
	}

	type A implements Identifiable {
		id: String!
		value: Int!
	}

	type B implements Identifiable {
		id: String!
		index: Int!
	}`

	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Schema: graphql.MustParseSchema(schemaString, &selectedFieldsResolver3{t}),
			Query: `
				{
					test {
						id
					}
				}
			`,
			IgnoreResult: true,
		},
	})
}

// Utils

func failWithError(t *testing.T, msg string, args ...interface{}) {
	t.Errorf(msg, args...)
	t.FailNow()
}
