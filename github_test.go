package patchwerk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	evanphx "github.com/evanphx/json-patch"
	"github.com/stretchr/testify/assert"
)

type jsonTest2 struct {
	Comment  string      `json:"comment"`
	Original interface{} `json:"doc"`
	Target   interface{} `json:"expected"`
	Disabled bool        `json:"disabled"`
	Error    *string     `json:"error"`
}

// evanphx/json-patch fails when replacing elements on root level, let's wrap everything into an object
type RootWrap struct {
	Root interface{} `json:"root"`
}

// Tests copied from github.com/json-patch/json-patch-tests
func TestBase(t *testing.T) {
	file, err := ioutil.ReadFile("github.json")
	assert.NoError(t, err)

	var jsonTests []jsonTest2
	err = json.Unmarshal([]byte(file), &jsonTests)
	assert.NoError(t, err)

	for i, tc := range jsonTests {
		if tc.Disabled || tc.Error != nil || tc.Original == nil || tc.Target == nil {
			continue
		}
		testName := fmt.Sprintf(`Test #%d %s`, i, tc.Comment)
		tc.Original = RootWrap{tc.Original}
		tc.Target = RootWrap{tc.Target}
		t.Run(testName, func(t *testing.T) {
			b1, err := json.Marshal(tc.Original)
			assert.NoError(t, err)
			b2, err := json.Marshal(tc.Target)
			assert.NoError(t, err)

			patch, err := Diff(b1, b2)
			assert.NoError(t, err)
			pb, err := json.Marshal(patch)
			assert.NoError(t, err)
			t.Log("original", string(b1))
			t.Log("target", string(b2))
			for i, p := range patch {
				b, _ := json.Marshal(p)
				t.Log("diff", i, string(b))
			}

			ep, err := evanphx.DecodePatch(pb)
			assert.NoError(t, err)

			modified, err := ep.Apply(b1)
			assert.NoError(t, err)

			assert.Equal(t, b2, modified)
		})
	}
}
