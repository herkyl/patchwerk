package patchwerk

import (
	"reflect"
)

type tmpEl struct {
	val     interface{}
	isFixed bool
}

func diffArrays(a, b []interface{}, p string) ([]*JSONPatchOperation, error) {
	patch := []*JSONPatchOperation{}

	// Find elements that are fixed in both arrays
	tmp := make([]tmpEl, len(a))
	for i, ae := range a {
		newEl := tmpEl{val: ae}
		for j := i; j < len(b); j++ {
			if len(b) <= j { //b is out of bounds
				break
			}
			be := b[j]
			if reflect.DeepEqual(ae, be) {
				newEl.isFixed = true // this element should remain in place
			}
		}
		tmp[i] = newEl
	}

	// Create a new array using adds and removes
	aIndex := 0
	bIndex := 0
	addedDelta := 0
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	for aIndex+addedDelta < maxLen {
		tmpIndex := aIndex + addedDelta
		newPath := makePath(p, tmpIndex)
		if aIndex >= len(a) && bIndex >= len(b) {
			break
		}
		if aIndex >= len(a) { // a is out of bounds, all new items in b must be adds
			if tmpIndex < len(b) {
			  patch = append(patch, NewPatch("add", newPath, b[tmpIndex]))
			}
			addedDelta++
			continue
		}
		if bIndex >= len(b) { // b is out of bounds, all new items in a must be removed
			patch = append(patch, NewPatch("remove", newPath, a[tmpIndex]))
			addedDelta--
			aIndex++
			continue
		}
		// can compare elements, so let's compare them
		te := tmp[aIndex]
		for j := bIndex; j < maxLen; j++ {
			be := b[j]
			if reflect.DeepEqual(te.val, be) {
				// element is already in b, move on
				bIndex++
				aIndex++
				break
			} else {
				if te.isFixed {
					patch = append(patch, NewPatch("add", newPath, be))
					addedDelta++
					bIndex++
					break
				} else {
					patch = append(patch, NewPatch("remove", newPath, te.val)) //save value for remove so we can use it later
					addedDelta--
					aIndex++
					break
				}

			}
		}
	}

	// See if remove+add pairs can be combined into a diff (can also be a replace)
	replacedPatch := []*JSONPatchOperation{}
	for i := 0; i < len(patch); i++ {
		a := patch[i]
		var b *JSONPatchOperation
		if i+1 < len(patch) {
			b = patch[i+1]
		}
		if b != nil && a.Path == b.Path && a.Operation == "remove" && b.Operation == "add" {
			diffPatch, err := diff(a.Value, b.Value, a.Path)
			if err != nil {
				return nil, err
			}
			replacedPatch = append(replacedPatch, diffPatch...)
			i++
		} else {
			if a.Operation == "remove" {
				a.Value = nil
			}
			replacedPatch = append(replacedPatch, a)
		}
	}

	return replacedPatch, nil
}
