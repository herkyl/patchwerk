package patchwerk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	evanphx "github.com/evanphx/json-patch"
	"github.com/stretchr/testify/assert"
)

type jsonTest struct {
	Comment  string        `json:"comment"`
	Original []interface{} `json:"original"`
	Target   []interface{} `json:"target"`
	Patch    []interface{} `json:"patch"`
	Disabled bool          `json:"disabled"`
}

func TestArrays(t *testing.T) {
	file, err := ioutil.ReadFile("arrays.json")
	assert.NoError(t, err)

	var jsonTests []jsonTest
	err = json.Unmarshal([]byte(file), &jsonTests)
	assert.NoError(t, err)

	for i, tc := range jsonTests {
		if tc.Disabled {
			continue
		}
		testName := fmt.Sprintf(`Test #%d %s`, i+1, tc.Comment)
		// tc.Original = RootWrap{tc.Original}
		// tc.Target = RootWrap{tc.Target}
		t.Run(testName, func(t *testing.T) {
			b1, err := json.Marshal(tc.Original)
			assert.NoError(t, err)
			b2, err := json.Marshal(tc.Target)
			assert.NoError(t, err)

			patch, err := diffArrays(tc.Original, tc.Target, "")
			assert.NoError(t, err)

			pb1, err := json.Marshal(tc.Patch)
			assert.NoError(t, err)

			pb2, err := json.Marshal(patch)
			assert.NoError(t, err)

			assert.JSONEq(t, string(pb1), string(pb2))

			ep, err := evanphx.DecodePatch(pb2)
			assert.NoError(t, err)

			modified, err := ep.Apply(b1)
			assert.NoError(t, err)

			assert.Equal(t, b2, modified)
		})
	}
}
