package scraping

import (
	"github.com/go-logr/logr"
	"github.com/gocolly/colly/v2"
	"github.com/jlewi/eyecare/go/api"
	"github.com/jlewi/eyecare/go/pkg/helpers"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

// Scraper scrapes a bunch of sites
type Scraper struct {
	OutputDir string
	Log       logr.Logger
	Force     bool
}

func (s *Scraper) Sites(sites []api.Site) error {
	// TODO(jeremy): Should
	for _, t := range sites {
		// TODO(jeremy): Should we try to keep going on error.
		if err := s.Scrape(t); err != nil {
			return err
		}
	}
	return nil
}

// Scrape a single site
func (s *Scraper) Scrape(site api.Site) error {
	log := s.Log.WithValues("domain", site.Domain)
	output := filepath.Join(s.OutputDir, site.Domain)

	log.Info("Using output directory", "directory", output)

	if _, err := os.Stat(output); os.IsNotExist(err) {
		log.Info("Creating output directory", "directory", output)
		os.MkdirAll(output, helpers.UserGroupAllPerm)
	} else {
		if s.Force {
			log.Info("Warning directory already exists site will be rescraped and contents overwriten")
		} else {
			log.Info("Directory exists; assuming site has already been scraped; not rescraping", "directory", output)
			return nil
		}
	}

	_, err := os.Stat(output)
	if err != nil {
		return errors.Wrapf(err, "Couldn't check directory %v", output)
	}

	// TODO(jeremy): Should we set MaxDepth to limit the scan.
	// TODO(jeremy): I think the collector is inherently parallel. So if we pass the collector in and don't
	// call wait we could easily run multiple collector jobs at once.
	c := colly.NewCollector()
	c.AllowedDomains = []string{site.Domain, "www." + site.Domain}
	log.Info("Setting allowed domain", "domain", site.Domain)
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// We don't care about PDFs
		if strings.HasSuffix(link, "pdf") {
			log.Info("Skipping PDF", "pdf", link)
			return
		}
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		log.Info("Visiting", "url", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		name := r.FileName()

		ext := filepath.Ext(name)
		if ext == ".unknown" {
			name = filepath.Base(name) + ".html"
		}
		path := filepath.Join(output, r.FileName())
		if err := r.Save(path); err != nil {
			log.Error(err, "Failed to save file.", "path", path, "url", r.Request.URL.String())
		}
	})
	url := "http://" + site.Domain
	log.Info("Starting collector", "url", url)

	// TODO(jeremy): Should we look for a sitemap.
	if err := c.Visit(url); err != nil {
		return errors.Wrapf(err, "Failed to visit domain: %v", site.Domain)
	}
	log.Info("Waiting for collector to finish")
	c.Wait()
	return nil
}
