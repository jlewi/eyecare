package analyzer

import (
	"github.com/google/go-cmp/cmp"
	"github.com/jlewi/eyecare/go/api"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func Test_Process(t *testing.T) {
	type testCase struct {
		name     string
		testDir  string
		matches  []string
		expected []api.Result
	}

	cases := []testCase{
		{
			name:    "basic",
			testDir: "test_data",
			expected: []api.Result{
				{Path: "/Users/jlewi/git_eyecare/go/pkg/analyzer/test_data/file1.html",
					Term: "(?i)keratoconus",
				},
				{Path: "/Users/jlewi/git_eyecare/go/pkg/analyzer/test_data/file2.html",
					Term: "(?i)keratoconus",
				},
				{Path: "/Users/jlewi/git_eyecare/go/pkg/analyzer/test_data/file2.html",
					Term: "(?i)topography",
				},
			},
			matches: []string{"(?i)keratoconus", "(?i)topography"},
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory; %v", err)
	}
	testDir := filepath.Join(cwd, "test_data")
	for _, c := range cases {
		t.Run(c.testDir, func(t *testing.T) {
			p := RegexAnalyzer{Matches: make([]*regexp.Regexp, 0, len(c.matches))}
			for _, m := range c.matches {
				re, err := regexp.Compile(m)
				if err != nil {
					t.Fatalf("Failed to compile regex: %v; error %v", m, err)
				}
				p.Matches = append(p.Matches, re)
			}

			actual, err := p.Process(testDir)
			if err != nil {
				t.Fatalf("Failed to run Process: %v", err)
			}
			if d := cmp.Diff(c.expected, actual); d != "" {
				t.Errorf("Unexpected diff:\n%v", d)
			}
		})
	}
}
