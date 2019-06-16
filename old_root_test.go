package patchwerk

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONPatchCreate(t *testing.T) {
	cases := map[string]struct {
		a    string
		b    string
		diff string
	}{
		"object": {
			`{"asdf":"qwerty"}`,
			`{"asdf":"zzz"}`,
			`[{"op":"replace","path":"/asdf","value":"zzz"}]`,
		},
		"object with array": {
			`{"items":[{"asdf":"qwerty"}]}`,
			`{"items":[{"asdf":"bla"},{"asdf":"zzz"}]}`,
			`[{"op":"replace","path":"/items","value":[{"asdf":"bla"},{"asdf":"zzz"}]}]`,
		},
		"from empty array": {
			`[]`,
			`[{"asdf":"bla"}]`,
			`[{"op":"add","path":"/0","value":{"asdf":"bla"}}]`,
		},
		"to empty array": {
			`[{"asdf":"bla"}]`,
			`[]`,
			`[{"op":"remove","path":"/0"}]`,
		},
		"from object to array": {
			`{"foo":"bar"}`,
			`[{"foo":"bar"}]`,
			`[{"op":"replace","path":"","value":[{"foo":"bar"}]}]`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Logf(`Running test: "%s"`, name)
			patch, err := Diff([]byte(tc.a), []byte(tc.b))
			assert.NoError(t, err)

			patchBytes, err := json.Marshal(patch)
			assert.NoError(t, err)

			assert.Equal(t, tc.diff, string(patchBytes))
		})
	}
}
