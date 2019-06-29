package patchwerk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var errBadJSONDoc = fmt.Errorf("Invalid JSON Document")

type JSONPatchOperation struct {
	Operation string      `json:"op"`
	Path      string      `json:"path"`
	Value     interface{} `json:"value,omitempty"`
}

func (j *JSONPatchOperation) JSON() string {
	b, _ := json.Marshal(j)
	return string(b)
}

func (j *JSONPatchOperation) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("{")
	b.WriteString(fmt.Sprintf(`"op":"%s"`, j.Operation))
	b.WriteString(fmt.Sprintf(`,"path":"%s"`, j.Path))
	// Consider omitting Value for non-nullable operations.
	if j.Value != nil || j.Operation == "replace" || j.Operation == "add" {
		v, err := json.Marshal(j.Value)
		if err != nil {
			return nil, err
		}
		b.WriteString(`,"value":`)
		b.Write(v)
	}
	b.WriteString("}")
	return b.Bytes(), nil
}

type ByPath []*JSONPatchOperation

func (a ByPath) Len() int           { return len(a) }
func (a ByPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPath) Less(i, j int) bool { return a[i].Path < a[j].Path }

func NewPatch(operation, path string, value interface{}) *JSONPatchOperation {
	return &JSONPatchOperation{Operation: operation, Path: path, Value: value}
}

// Diff creates a patch as specified in http://jsonpatch.com/
//
// 'a' is original, 'b' is the modified document. Both are to be given as json encoded content.
// The function will return an array of JSONPatchOperations
//
// An error will be returned if any of the two documents are invalid.
func Diff(a, b []byte) ([]byte, error) {
	var aI interface{}
	var bI interface{}

	err := json.Unmarshal(a, &aI)
	if err != nil {
		return nil, errBadJSONDoc
	}
	err = json.Unmarshal(b, &bI)
	if err != nil {
		return nil, errBadJSONDoc
	}

	ops, err := diff(aI, bI, "")
	if err != nil {
		return nil, err
	}
	return json.Marshal(ops)
}

// From http://tools.ietf.org/html/rfc6901#section-4 :
//
// Evaluation of each reference token begins by decoding any escaped
// character sequence.  This is performed by first transforming any
// occurrence of the sequence '~1' to '/', and then transforming any
// occurrence of the sequence '~0' to '~'.
//   TODO decode support:
//   var rfc6901Decoder = strings.NewReplacer("~1", "/", "~0", "~")

var rfc6901Encoder = strings.NewReplacer("~", "~0", "/", "~1")

func makePath(path string, newPart interface{}) string {
	key := rfc6901Encoder.Replace(fmt.Sprintf("%v", newPart))
	if path == "" {
		return "/" + key
	}
	if strings.HasSuffix(path, "/") {
		return path + key
	}
	return path + "/" + key
}

func diff(a, b interface{}, p string) ([]*JSONPatchOperation, error) {
	fullReplace := []*JSONPatchOperation{NewPatch("replace", p, b)}
	patch := []*JSONPatchOperation{}

	// If values are not of the same type simply replace
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return fullReplace, nil
	}

	var err error
	var tempPatch []*JSONPatchOperation
	switch at := a.(type) {
	case map[string]interface{}:
		bt := b.(map[string]interface{})
		tempPatch, err = diffObjects(at, bt, p)
		if err != nil {
			return nil, err
		}
		patch = append(patch, tempPatch...)
	case string, float64, bool:
		if !reflect.DeepEqual(a, b) {
			patch = append(patch, NewPatch("replace", p, b))
		}
	case []interface{}:
		bt, ok := b.([]interface{})
		if !ok {
			// array replaced by non-array
			patch = append(patch, NewPatch("replace", p, b))
		} else {
			// arrays are not the same length
			tempPatch, err = diffArrays(at, bt, p)
			if err != nil {
				return nil, err
			}
			patch = append(patch, tempPatch...)
		}
	case nil:
		switch b.(type) {
		case nil:
			// Both nil, fine.
		default:
			patch = append(patch, NewPatch("add", p, b))
		}
	default:
		panic(fmt.Sprintf("Unknown type:%T ", a))
	}
	return getSmallestPatch(fullReplace, patch), nil
}

func getSmallestPatch(patches ...[]*JSONPatchOperation) []*JSONPatchOperation {
	smallestPatch := patches[0]
	b, _ := json.Marshal(patches[0])
	smallestSize := len(b)
	for i := 1; i < len(patches); i++ {
		p := patches[i]
		b, _ := json.Marshal(p)
		size := len(b)
		if size < smallestSize {
			smallestPatch = p
			smallestSize = size
		}
	}
	return smallestPatch
}
