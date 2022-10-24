package scraping

import "testing"

func Test_URLToDomain(t *testing.T) {
	type testCase struct {
		input    string
		expected string
	}

	cases := []testCase{
		{
			input:    "https://vs.com/",
			expected: "vs.com",
		},
		{
			input:    "vs.com",
			expected: "vs.com",
		},
		{
			input:    "www.vs.com",
			expected: "vs.com",
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			actual := URLToDomain(c.input)

			if actual != c.expected {
				t.Errorf("Got %v; want %v", actual, c.expected)
			}
		})
	}
}
