package types

import (
	"github.com/graph-gophers/graphql-go/errors"
)

// FieldDefinition is a representation of a GraphQL FieldDefinition.
//
// http://spec.graphql.org/draft/#FieldDefinition
type FieldDefinition struct {
	Name       string
	Arguments  ArgumentsDefinition
	Type       Type
	Directives DirectiveList
	Desc       string
	Loc        errors.Location
}

// FieldsDefinition is a list of an ObjectTypeDefinition's Fields.
//
// https://spec.graphql.org/draft/#FieldsDefinition
type FieldsDefinition []*FieldDefinition

// Get returns a FieldDefinition in a FieldsDefinition by name or nil if not found.
func (l FieldsDefinition) Get(name string) *FieldDefinition {
	for _, f := range l {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// Names returns a slice of FieldDefinition names.
func (l FieldsDefinition) Names() []string {
	names := make([]string, len(l))
	for i, f := range l {
		names[i] = f.Name
	}
	return names
}

// SelectedField is the public representation of a field selection in a GraphQL query
type SelectedField struct {
	Name             string
	TypeName         string
	Alias            string // equal to Name if alias is not provided
	Directives       DirectiveList
	Args             interface{}
	ArgsMap          map[string]interface{}
	Fields           []*SelectedField
	AssertedTypeName string // non-empty for field selections on union fields
}

type SelectedFieldIdentifier struct {
	Name  string
	Alias string
}

// Lookup searches the field tree for the field defined by the input path
func (f *SelectedField) Lookup(path ...SelectedFieldIdentifier) *SelectedField {
	if len(path) == 0 || f == nil {
		return f
	}
	curr := path[0]
	if curr.Name == "" {
		panic("path component's name cannot be empty")
	}
	for _, subf := range f.Fields {
		if subf.Name == curr.Name && curr.Alias == subf.Alias {
			return subf.Lookup(path[1:]...)
		}
	}
	return nil
}
