package nimbleway

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/jlewi/eyecare/go/api"
	"os"
	"path/filepath"
	"testing"
)

func Test_ToSites(t *testing.T) {
	type testCase struct {
		fileName string
		expected []api.Site
	}

	cases := []testCase{
		{
			fileName: "results.json",
			expected: []api.Site{
				{
					Domain: "parksideoptometry.com",
				},
			},
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory; %v", err)
	}
	testDir := filepath.Join(cwd, "test_data")
	for _, c := range cases {
		t.Run(c.fileName, func(t *testing.T) {
			path := filepath.Join(testDir, c.fileName)
			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to readFile %v; error %v", path, err)
			}
			r := &SERPResults{}
			d := json.NewDecoder(bytes.NewReader(b))
			d.DisallowUnknownFields()
			if err := d.Decode(r); err != nil {
				t.Fatalf("Failed to parse serp results; %v", err)
			}

			actual, err := ToSites(r)
			if err != nil {
				t.Fatalf("Failed to convert SerpResults; error %v", err)
			}

			if d := cmp.Diff(c.expected, actual); d != "" {
				t.Errorf("Unexpected diff:\n%v", d)
			}
		})
	}
}
