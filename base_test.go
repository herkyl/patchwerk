package patchwerk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoChange(t *testing.T) {
	patch, e := Diff([]byte(`{"a":100, "b":200, "c":"hello"}`), []byte(`{"a":100, "b":200, "c":"hello"}`))
	assert.NoError(t, e)
	assert.Equal(t, 0, len(patch))
}

func TestAddingToArray(t *testing.T) {
	patch, e := Diff([]byte(`{"a":[1, 2, 3]}`), []byte(`{"a":[1, 2, 3, 4]}`))
	assert.NoError(t, e)
	assert.Equal(t, 1, len(patch))
	p := patch[0]
	assert.Equal(t, "add", p.Operation)
	assert.Equal(t, "/a/3", p.Path)
	assert.Equal(t, float64(4), p.Value)
}

func TestReplaceScalars(t *testing.T) {
	patch, e := Diff([]byte(`1`), []byte(`"s"`))
	assert.NoError(t, e)
	assert.Equal(t, 1, len(patch))
	p := patch[0]
	assert.Equal(t, "replace", p.Operation)
	assert.Equal(t, "", p.Path)
	assert.Equal(t, "s", p.Value)
}
