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

func (s ResolverContext) FindRoot(path ...string) (Selection, error) {
	return s.find(s.Operation.Selections, path...)
}

func (s ResolverContext) find(set SelectionSet, path ...string) (Selection, error) {
	var found Selection
	for _, sel := range set {
		switch sel := sel.(type) {
		case *Field:
			if sel.Name.Name == path[0] {
				found = sel
			}
		case *InlineFragment:
			foundInFrag, err := s.find(sel.Fragment.Selections, path...)
			if err != nil {
				return nil, err
			}
			found = foundInFrag
		case *FragmentSpread:
			fragSelections, err := s.getSelections(sel)
			if err != nil {
				return nil, err
			}
			found, err = s.find(fragSelections, path...)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.Errorf("unhandled type %T", sel)
		}
		if found != nil {
			break
		}
	}

	path = path[1:]
	if len(path) == 0 || found == nil {
		return found, nil
	}
	sel, err := s.getSelections(found)
	if err != nil {
		return nil, err
	}
	return s.find(sel, path...)
}

func (s ResolverContext) getSelections(sel Selection) (SelectionSet, error) {
	switch sel := sel.(type) {
	case *Field:
		return sel.SelectionSet, nil
	case *InlineFragment:
		return sel.Fragment.Selections, nil
	case *FragmentSpread:
		frag, ok := s.fragmentByName[sel.Name.Name]
		if !ok {
			return nil, errors.Errorf("failed to find fragment %s in the mapping", sel.Name.Name)
		}
		return frag.Fragment.Selections, nil
	default:
		return nil, errors.Errorf("unhandled type %T", sel)
	}
}
