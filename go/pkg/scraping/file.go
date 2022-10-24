package scraping

import (
	"encoding/json"
	"github.com/jlewi/eyecare/go/api"
	"github.com/jlewi/eyecare/go/pkg/helpers"
	"github.com/pkg/errors"
	"io"
	"os"
)

// ReadFile reads a JSONL file containing a list of sites
func ReadFile(input string) ([]api.Site, error) {
	sites := make([]api.Site, 0, 20)
	f, err := os.Open(input)
	defer helpers.DeferIgnoreError(f.Close)
	if err != nil {
		return sites, errors.Wrapf(err, "Failed to read file %v", input)
	}
	d := json.NewDecoder(f)
	for {
		t := &api.Site{}
		if err := d.Decode(t); err != nil {
			if err == io.EOF {
				break
			}
			return sites, errors.Wrapf(err, "Failed to decode site")
		}
		sites = append(sites, *t)
	}
	return sites, nil
}

// WriteFile writes the sites to the file as JSONL
func WriteFile(outFile string, sites []api.Site) error {
	// Emit it as JSONL
	f, err := os.Create(outFile)
	defer helpers.DeferIgnoreError(f.Close)
	if err != nil {
		return errors.Wrapf(err, "Failed to create file: %v", outFile)
	}

	e := json.NewEncoder(f)
	for _, t := range sites {
		if err := e.Encode(t); err != nil {
			return errors.Wrapf(err, "Failed to encode the site")
		}
	}
	return nil
}
