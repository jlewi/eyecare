package api

// Site represents a DRs office
// TODO(jeremy): This struct is missing JSON tags so it ends up using UpperCamelCase rather than lower camelCase
type Site struct {
	Domain                   string
	IsOptometrist            bool
	CornealTopography        bool
	CornealTopographyMention string
	Results                  []Result `json:"results"`
}

// Result is a found match
type Result struct {
	Path string `json:"path"`
	Term string `json:"term"`
}
