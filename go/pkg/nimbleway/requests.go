package nimbleway

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

const (
	SERPEndpoint    = "https://api.webit.live/api/v1/realtime/serp"
	LoginEndpoint   = "https://api.nimbleway.com/api/v1/account/login"
	JSONContentType = "application/json"
	ContentHeader   = "Content-Type"
)

type SecretsFile struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SERPRequest struct {
	SearchEngine string `json:"search_engine"`
	Country      string `json:"country"`
	Locale       string `json:"locale"`
	Query        string `json:"query"`
	Parse        bool   `json:"parse"`
}

type SERPResults struct {
	HtmlContent string        `json:"html_content"`
	Status      string        `json:"status"`
	Parsing     ParsedResults `json:"parsing"`
	URL         string        `json:"url"`
}

type ParsedResults struct {
	Status             string                 `json:"status"`
	Entities           Entities               `json:"entities"`
	TotalEntitiesCount int                    `json:"total_entities_count"`
	EntitiesCount      EntitiesCount          `json:"entities_count"`
	Metrics            map[string]interface{} `json:"metrics"`
}

type EntitiesCount struct {
	OrganicResult     int `json:"OrganicResult"`
	Pagination        int `json:"Pagination"`
	RelatedSearch     int `json:"RelatedSearch"`
	SearchInformation int `json:"SearchInformation"`
}
type Entities struct {
	OrganicResult     []OrganicResult     `json:"OrganicResult"`
	Pagination        []Pagination        `json:"Pagination"`
	RelatedSearch     []RelatedSearch     `json:"RelatedSearch"`
	SearchInformation []SearchInformation `json:"SearchInformation"`
}

type SearchInformation struct {
	EntityType         string `json:"entityType"`
	QueryDisplayed     string `json:"query_displayed"`
	TimeTakenDisplayed string `json:"time_taken_displayed"`
	TotalResults       string `json:"total_results"`
}

type Pagination struct {
	CurrentPage   int               `json:"current_page"`
	EntityType    string            `json:"entityType"`
	NextPageUrl   string            `json:"next_page_url"`
	OtherPageUrls map[string]string `json:"other_page_urls"`
}

type OrganicResult struct {
	DisplayedUrl string     `json:"displayed_url"`
	EntityType   string     `json:"entityType"`
	Position     int        `json:"position"`
	Snippet      string     `json:"snippet"`
	Title        string     `json:"title"`
	Url          string     `json:"url"`
	SiteLinks    []SiteLink `json:"sitelinks"`
}

type RelatedSearch struct {
	EntityType string `json:"entityType"`
	Query      string `json:"query"`
	Url        string `json:"url"`
}
type SiteLink struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type TokenResponse struct {
	Token string `json: "token"`
}

// NewRequest instantiates a default request with the provided string
func NewRequest(query string) *SERPRequest {
	return &SERPRequest{
		SearchEngine: "google_search",
		Country:      "US",
		Locale:       "en",
		Query:        query,
		Parse:        true,
	}
}

type Client struct {
	Token   string
	Secrets SecretsFile
}

func (c *Client) GetToken() error {
	b, err := json.Marshal(c.Secrets)
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal secrets to get token")
	}
	body := bytes.NewReader(b)
	resp, err := http.Post(LoginEndpoint, JSONContentType, body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "Failed to login to nimbleway to get token; status code %v", resp.StatusCode)
	}

	tok, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Failed to read response body")
	}
	tokResp := &TokenResponse{}
	if err := json.Unmarshal(tok, tokResp); err != nil {
		return errors.Wrapf(err, "Failed to unmarshal the token response")
	}
	c.Token = tokResp.Token
	return nil
}

// Serp executes a SERP request
func (c *Client) Serp(ctx context.Context, req *SERPRequest) (*SERPResults, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to marshal request to JSON")
	}
	body := bytes.NewReader(payload)
	hReq, err := http.NewRequest(http.MethodPost, SERPEndpoint, body)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create new request")
	}
	hReq.Header.Set("Authorization", "Bearer "+c.Token)
	hReq.Header.Set(ContentHeader, JSONContentType)
	resp, err := http.DefaultClient.Do(hReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "SERP request failed with code %v", resp.StatusCode)
	}

	d := json.NewDecoder(resp.Body)
	results := &SERPResults{}
	if err := d.Decode(results); err != nil {
		return nil, errors.Wrapf(err, "Failed to read response body")
	}
	return results, nil
}
