package types

import (
	"context"

	"github.com/pkg/errors"
)

// ResolverContext is a custom context structure w/ additional data meant to be used in resolver methods
type ResolverContext struct {
	context.Context

	op             *OperationDefinition
	def            *ExecutableDefinition
	fragmentByName map[string]*FragmentDefinition
}

func NewResolverContext(ctx context.Context, o *OperationDefinition, d *ExecutableDefinition) ResolverContext {
	rc := ResolverContext{ctx, o, d, map[string]*FragmentDefinition{}}
	for _, f := range d.Fragments { // construct the fragments mapping
		rc.fragmentByName[f.Name.Name] = f
	}
	return rc
}

func (rc ResolverContext) GetFieldNames(path ...string) (res []string) {
	found := rc.findField(rc.op.Selections, path...)
	if found == nil {
		return
	}
	for _, sel := range found.SelectionSet {
		res = append(res, rc.fieldNames(sel)...)
	}
	return
}

func (rc ResolverContext) fieldNames(sel Selection) (res []string) {
	switch sel := sel.(type) {
	case *Field:
		res = append(res, sel.Name.Name)
	case *InlineFragment, *FragmentSpread:
		for _, s := range rc.getSelections(sel) {
			res = append(res, rc.fieldNames(s)...)
		}
	default:
		panic(errors.Errorf("unhandled type %T", sel))
	}
	return
}

func (rc ResolverContext) findField(set SelectionSet, path ...string) *Field {
	var found *Field
	for _, sel := range set {
		switch sel := sel.(type) {
		case *Field:
			if sel.Name.Name == path[0] {
				found = sel
			}
		case *InlineFragment, *FragmentSpread:
			found = rc.findField(rc.getSelections(sel), path...)
		default:
			panic(errors.Errorf("unhandled type %T", sel))
		}
		if found != nil {
			break
		}
	}

	path = path[1:]
	if len(path) == 0 || found == nil {
		return found
	}
	return rc.findField(rc.getSelections(found), path...)
}

func (rc ResolverContext) getSelections(sel Selection) SelectionSet {
	switch sel := sel.(type) {
	case *Field:
		return sel.SelectionSet
	case *InlineFragment:
		return sel.Fragment.Selections
	case *FragmentSpread:
		frag, ok := rc.fragmentByName[sel.Name.Name]
		if !ok {
			return nil
		}
		return frag.Fragment.Selections
	default:
		panic(errors.Errorf("unhandled type %T", sel))
	}
}
