package patchwerk

import (
	"reflect"
)

// diff returns the (recursive) difference between a and b as an array of JsonPatchOperations.
func diffObjects(a, b map[string]interface{}, path string) ([]*JSONPatchOperation, error) {
	patch := []*JSONPatchOperation{}
	for key, bv := range b {
		p := makePath(path, key)
		av, ok := a[key]
		// Key doesn't exist in original document, value was added
		if !ok {
			patch = append(patch, NewPatch("add", p, bv))
			continue
		}
		// If types have changed, replace completely
		if reflect.TypeOf(av) != reflect.TypeOf(bv) {
			patch = append(patch, NewPatch("replace", p, bv))
			continue
		}
		// Types are the same, compare values
		tempPatch, err := diff(av, bv, p)
		if err != nil {
			return nil, err
		}
		patch = append(patch, tempPatch...)
	}
	// Now add all deleted values as nil
	for key := range a {
		_, ok := b[key]
		if !ok {
			p := makePath(path, key)
			patch = append(patch, NewPatch("remove", p, nil))
		}
	}
	return patch, nil
}
