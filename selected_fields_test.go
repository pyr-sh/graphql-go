package graphql_test

import (
	"context"
	"testing"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/fields"
	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/graph-gophers/graphql-go/types"
)

// Test 1: Nested fields with aliases, arguments and directives

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
			Schema: graphql.MustParseSchema(schemaString, &nestedAliasesArgsDirectivesQueryResolver{t}),
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

type nestedAliasesArgsDirectivesQueryResolver struct {
	T *testing.T
}

func (r *nestedAliasesArgsDirectivesQueryResolver) Test(ctx context.Context) nestedAliasesArgsDirectivesTestResolver {
	sel := fields.Selected(ctx)
	if sel == nil {
		failWithError(r.T, "selected field is nil")
	}

	if !(sel.Name == "test" && sel.TypeName == "Test!" && len(sel.Fields) == 4) {
		failWithError(r.T, "test field assertion failed")
	}

	a, b, c, d := sel.Fields[0], sel.Fields[1], sel.Fields[2], sel.Fields[3]

	// a
	if !(a.Name == "a" && a.TypeName == "String!" && a.Alias == "a" && len(a.Directives) == 1) {
		failWithError(r.T, "a field assertion failed")
	}
	if !(b.Name == "b" && b.TypeName == "String" && b.Alias == "b" && len(b.Directives) == 1) {
		failWithError(r.T, "b field assertion failed")
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
		failWithError(r.T, "b field assertion failed")
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
		failWithError(r.T, "c field assertion failed")
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
		failWithError(r.T, "c field assertion failed")
	}

	return nestedAliasesArgsDirectivesTestResolver{}
}

type nestedAliasesArgsDirectivesTestResolver struct{}

func (r nestedAliasesArgsDirectivesTestResolver) A(struct{ Value *string }) string { return "a" }
func (r nestedAliasesArgsDirectivesTestResolver) B() *string                       { return nil }
func (r nestedAliasesArgsDirectivesTestResolver) C() int32                         { return 1 }
func (r nestedAliasesArgsDirectivesTestResolver) D() nestedAliasesArgsDirectivesDResolver {
	return nestedAliasesArgsDirectivesDResolver{}
}

type nestedAliasesArgsDirectivesDResolver struct{}

func (r nestedAliasesArgsDirectivesDResolver) Value() string { return "value" }

// Test 2: Union query result

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
			Schema: graphql.MustParseSchema(schemaString, &unionQueryResultQueryResolver{t}),
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

type unionQueryResultQueryResolver struct {
	T *testing.T
}

func (r *unionQueryResultQueryResolver) Test(ctx context.Context) []unionQueryResultTestUnion {
	f := fields.Selected(ctx)
	if f == nil {
		failWithError(r.T, "selected field is nil")
	}
	if !(f.Name == "test" && f.TypeName == "[Test!]!") {
		failWithError(r.T, "test field assertion failed")
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

	return []unionQueryResultTestUnion{}
}

type aOrB interface {
	ID() string
}
type unionQueryResultTestUnion struct{ aOrB }

type unionQueryResultAResolver struct{}

func (r unionQueryResultAResolver) ID() string   { return "" }
func (r unionQueryResultAResolver) Value() int32 { return 0 }

type unionQueryResultBResolver struct{}

func (r unionQueryResultBResolver) ID() string   { return "" }
func (r unionQueryResultBResolver) Index() int32 { return 0 }

func (r unionQueryResultTestUnion) ToA() (*unionQueryResultAResolver, bool) {
	v, ok := r.aOrB.(unionQueryResultAResolver)
	if !ok {
		return nil, false
	}
	return &v, true
}

func (r unionQueryResultTestUnion) ToB() (*unionQueryResultBResolver, bool) {
	v, ok := r.aOrB.(unionQueryResultBResolver)
	if !ok {
		return nil, false
	}
	return &v, true
}

// Test 3: Interface query result

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
			Schema: graphql.MustParseSchema(schemaString, &interfaceQueryResultQueryResolver{t}),
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

type interfaceQueryResultQueryResolver struct {
	T *testing.T
}

func (r *interfaceQueryResultQueryResolver) Test(ctx context.Context) []interfaceQueryResultIdentifiable {
	f := fields.Selected(ctx)
	if f == nil {
		failWithError(r.T, "selected field is nil")
	}
	if len(f.Fields) != 1 {
		failWithError(r.T, "expected 1 field, got %d", len(f.Fields))
	}
	if !(f.Name == "test" && f.TypeName == "[Identifiable!]!" && f.AssertedTypeName == "") {
		failWithError(r.T, "test field assertion field")
	}
	field := f.Fields[0]
	if !(field.Name == "id" && field.TypeName == "String!") {
		failWithError(r.T, "id field assertion field")
	}

	return []interfaceQueryResultIdentifiable{}
}

type interfaceQueryResultIdentifiable interface {
	ID() string
	ToA() (*unionQueryResultAResolver, bool)
	ToB() (*unionQueryResultBResolver, bool)
}

// Test 4: Nested fields lookup

func TestSelectedFieldsNestedFieldsLookup(t *testing.T) {
	schemaString := `
	schema {
		query: Query
	}

	type Query {
		test: Test!
	}

	type Test {
		a: A!
	}

	type A {
		b: B!
	}

	type B {
		c: C!
	}

	type C {
		value: String!
	}`

	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Schema: graphql.MustParseSchema(schemaString, &nestedFieldsLookupQueryResolver{t}),
			Query: `
				{
					test {
						a {
							b {
								aliasedC: c {
									value
								}
							}
						}
					}
				}
			`,
			IgnoreResult: true,
		},
	})
}

type nestedFieldsLookupQueryResolver struct {
	T *testing.T
}

func (r *nestedFieldsLookupQueryResolver) Test(ctx context.Context) nestedFieldsLookupTestResolver {
	testField := fields.Selected(ctx)
	if !(testField != nil && testField.Name == "test" && testField.TypeName == "Test!" && testField.Alias == "test") {
		failWithError(r.T, "test field assertion failed")
	}
	aField := testField.Lookup(types.SelectedFieldIdentifier{Name: "a", Alias: "a"})
	if !(aField != nil && aField.Name == "a" && aField.TypeName == "A!" && aField.Alias == "a") {
		failWithError(r.T, "a field assertion failed")
	}
	bField := aField.Lookup(types.SelectedFieldIdentifier{Name: "b", Alias: "b"})
	if !(bField != nil && bField.Name == "b" && bField.TypeName == "B!" && bField.Alias == "b") {
		failWithError(r.T, "b field assertion failed")
	}
	cField := bField.Lookup(types.SelectedFieldIdentifier{Name: "c", Alias: "aliasedC"})
	if !(cField != nil && cField.Name == "c" && cField.TypeName == "C!" && cField.Alias == "aliasedC") {
		failWithError(r.T, "c field assertion failed")
	}
	valueField := cField.Lookup(types.SelectedFieldIdentifier{Name: "value", Alias: "value"})
	if !(valueField != nil && valueField.Name == "value" && valueField.TypeName == "String!" &&
		valueField.Alias == "value" && len(valueField.Fields) == 0) {
		failWithError(r.T, "value field assertion failed")
	}
	return nestedFieldsLookupTestResolver{}
}

type nestedFieldsLookupTestResolver struct{}

func (r nestedFieldsLookupTestResolver) A() nestedFieldsLookupAResolver {
	return nestedFieldsLookupAResolver{}
}

type nestedFieldsLookupAResolver struct{}

func (r nestedFieldsLookupAResolver) B() nestedFieldsLookupBResolver {
	return nestedFieldsLookupBResolver{}
}

type nestedFieldsLookupBResolver struct{}

func (r nestedFieldsLookupBResolver) C() nestedFieldsLookupCResolver {
	return nestedFieldsLookupCResolver{}
}

type nestedFieldsLookupCResolver struct{}

func (r nestedFieldsLookupCResolver) Value() string { return "value" }

// Utils

func failWithError(t *testing.T, msg string, args ...interface{}) {
	t.Errorf(msg, args...)
	t.FailNow()
}
