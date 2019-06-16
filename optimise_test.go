package patchwerk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const lorem = "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."

// Test should replace the entire array instead of doing separate array operations
func TestReplaceInsteadOfArrayOps(t *testing.T) {
	patch, e := Diff([]byte(`{"a":[1, 2, 3]}`), []byte(`{"a":[1, 0, 0]}`))
	assert.NoError(t, e)
	assert.Equal(t, 1, len(patch))
	p := patch[0]
	assert.Equal(t, "replace", p.Operation)
	assert.Equal(t, "/a", p.Path)
	assert.Equal(t, []interface{}{float64(1), float64(0), float64(0)}, p.Value)
}

// Test should do individual array operations because one of the constant values is too big for an efficient replace
func TestReplaceObjectInArray(t *testing.T) {
	a := fmt.Sprintf(`[1, 2, {"a": "%s", "b": "2"}]`, lorem)
	b := fmt.Sprintf(`[1, 2, {"a": "%s", "b": "1", "c": "3"}]`, lorem)
	patch, e := Diff([]byte(a), []byte(b))
	assert.NoError(t, e)
	t.Log("PATCH:", patch)
	assert.Equal(t, 2, len(patch))
	p1 := patch[0]
	assert.Equal(t, "replace", p1.Operation)
	assert.Equal(t, "/2/b", p1.Path)
	assert.Equal(t, "1", p1.Value)

	p2 := patch[1]
	assert.Equal(t, "add", p2.Operation)
	assert.Equal(t, "/2/c", p2.Path)
	assert.Equal(t, "3", p2.Value)
}

func TestInnerObjectAddition(t *testing.T) {
	patch, e := Diff([]byte(`[1, 2, ["a", {"k1": "v1"}]]`), []byte(`[1, 2, ["a", {"k2": "v2", "k1": "v1"}]]`))
	assert.NoError(t, e)
	t.Log("PATCH:", patch)
	assert.Equal(t, 1, len(patch))
	p1 := patch[0]
	assert.Equal(t, "add", p1.Operation)
	assert.Equal(t, "/2/1/k2", p1.Path)
	assert.Equal(t, "v2", p1.Value)
}

func TestInnerArrayAddition(t *testing.T) {
	patch, e := Diff([]byte(`[1, ["a", ["x", true], "b"], 2]`), []byte(`[1, ["a", ["x", false, true], "b"], 2]`))
	assert.NoError(t, e)
	t.Log("PATCH:", patch)
	assert.Equal(t, 1, len(patch))
	p1 := patch[0]
	assert.Equal(t, "add", p1.Operation)
	assert.Equal(t, "/1/1/1", p1.Path)
	assert.Equal(t, false, p1.Value)
}
