package nimbleway

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func Test_Parsing(t *testing.T) {
	type testCase struct {
		fileName string
		expected SERPResults
	}

	cases := []testCase{
		{
			fileName: "results.json",
			expected: SERPResults{
				Status:      "success",
				HtmlContent: "somehtml",
				Parsing: ParsedResults{
					Status: "success",
					Entities: Entities{
						OrganicResult: []OrganicResult{
							{

								DisplayedUrl: "https://www.yelp.com \u203a ... \u203a Optometrists",
								EntityType:   "OrganicResult",
								Position:     1,
								Snippet:      "Best Optometrists in San Mateo, CA - Eyeworks of San Mateo, Hilary Chiem, OD, Clear Optometry, Parkside Optometry, Site for Sore Eyes - San Mateo,\u00a0...",
								Title:        "The Best 10 Optometrists in San Mateo, California - Yelp",
								Url:          "https://www.yelp.com/search?cflt=optometrists&find_loc=San+Mateo%2C+CA",
							},
							{
								DisplayedUrl: "https://www.parksideoptometry.com",
								EntityType:   "OrganicResult",
								Position:     3,
								SiteLinks: []SiteLink{
									{
										Title: "About Us",
										Url:   "https://www.parksideoptometry.com/about-us.html",
									},
									{
										Title: "Eyecare Services",
										Url:   "https://www.parksideoptometry.com/services.html",
									},
								},
								Snippet: "Welcome to Parkside Optometry Your Optometrists in San Mateo, California 94403. Call us at 650-830-5675 today. ... Parkside Optometry is a full service eye and\u00a0...",
								Title:   "San Mateo Optometrists, Parkside Optometry",
								Url:     "https://www.parksideoptometry.com/",
							},
						},
						Pagination: []Pagination{
							{
								CurrentPage: 1,
								EntityType:  "Pagination",
								NextPageUrl: "/search?q=optometrist+in+san+mateo&hl=en&ei=hPtRY_XkCMDV5NoP88OPsA4&start=10&sa=N&ved=2ahUKEwj1ttSgm_D6AhXAKlkFHfPhA-YQ8NMDegQIARAW",
								OtherPageUrls: map[string]string{
									"2": "/search?q=optometrist+in+san+mateo&hl=en&ei=hPtRY_XkCMDV5NoP88OPsA4&start=10&sa=N&ved=2ahUKEwj1ttSgm_D6AhXAKlkFHfPhA-YQ8tMDegQIARAE",
									"3": "/search?q=optometrist+in+san+mateo&hl=en&ei=hPtRY_XkCMDV5NoP88OPsA4&start=20&sa=N&ved=2ahUKEwj1ttSgm_D6AhXAKlkFHfPhA-YQ8tMDegQIARAG",
								},
							},
						},
						RelatedSearch: []RelatedSearch{
							{
								EntityType: "RelatedSearch",
								Query:      "optometrist near me",
								Url:        "/search?hl=en&q=Optometrist+near+me&sa=X&ved=2ahUKEwj1ttSgm_D6AhXAKlkFHfPhA-YQ1QJ6BAg2EAE",
							},
							{
								EntityType: "RelatedSearch",
								Query:      "ophthalmologist",
								Url:        "/search?hl=en&q=Ophthalmologist&sa=X&ved=2ahUKEwj1ttSgm_D6AhXAKlkFHfPhA-YQ1QJ6BAg9EAE",
							},
						},
						SearchInformation: []SearchInformation{
							{
								EntityType:         "SearchInformation",
								QueryDisplayed:     "optometrist in san mateo",
								TimeTakenDisplayed: "0.47 seconds",
								TotalResults:       "About 774,000 results ",
							},
						},
					},
					TotalEntitiesCount: 19,
					EntitiesCount: EntitiesCount{
						OrganicResult:     9,
						Pagination:        1,
						RelatedSearch:     8,
						SearchInformation: 1,
					},
					Metrics: map[string]any{},
				},

				URL: "https://www.google.com/search?q=optometrist+in+san+mateo&hl=en",
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
				t.Fatalf("Failed to parse results; %v", err)
			}
			if d := cmp.Diff(c.expected, *r); d != "" {
				t.Errorf("Unexpected diff:\n%v", d)
			}
		})
	}
}
