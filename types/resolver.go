package types

import (
	"context"

	"github.com/pkg/errors"
)

type ResolverContext struct {
	context.Context
	Operation  *OperationDefinition
	Definition *ExecutableDefinition

	// Internal
	fragmentByName map[string]*FragmentDefinition
}

func NewResolverContext(ctx context.Context, o *OperationDefinition, d *ExecutableDefinition) ResolverContext {
	rc := ResolverContext{ctx, o, d, map[string]*FragmentDefinition{}}
	for _, f := range d.Fragments {
		rc.fragmentByName[f.Name.Name] = f
	}
	return rc
}

func (s ResolverContext) FindRoot(path ...string) Selection {
	return s.find(s.Operation.Selections, path...)
}

func (s ResolverContext) find(set SelectionSet, path ...string) Selection {
	var found Selection
	for _, sel := range set {
		switch sel := sel.(type) {
		case *Field:
			if sel.Name.Name == path[0] {
				found = sel
			}
		case *InlineFragment:
			found = s.find(sel.Fragment.Selections, path...)
		case *FragmentSpread:
			found = s.find(s.getSelections(sel), path...)
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
	return s.find(s.getSelections(found), path...)
}

func (s ResolverContext) getSelections(sel Selection) SelectionSet {
	switch sel := sel.(type) {
	case *Field:
		return sel.SelectionSet
	case *InlineFragment:
		return sel.Fragment.Selections
	case *FragmentSpread:
		frag, ok := s.fragmentByName[sel.Name.Name]
		if !ok {
			return nil
		}
		return frag.Fragment.Selections
	default:
		panic(errors.Errorf("unhandled type %T", sel))
	}
}
