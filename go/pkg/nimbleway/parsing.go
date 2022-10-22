package nimbleway

import (
	"github.com/go-logr/zapr"
	"github.com/jlewi/eyecare/go/api"
	"github.com/jlewi/eyecare/go/pkg/scraping"
	"go.uber.org/zap"
	"net/url"
)

const ()

// ToSites parses the list of SERPResults to a list of sites
func ToSites(results *SERPResults) ([]api.Site, error) {
	sites := make([]api.Site, 0, len(results.Parsing.Entities.OrganicResult))
	log := zapr.NewLogger(zap.L())

	excludedDomains := map[string]bool{"yelp.com": true, "zocdoc.com": true}

	for _, r := range results.Parsing.Entities.OrganicResult {
		log.Info("Parsing OrganicResult", "url", r.Url)
		u, err := url.Parse(r.Url)
		if err != nil {
			log.Error(err, "Failed to parse url", "url", r.Url)
			continue
		}
		domain := scraping.HostToDomain(u.Host)

		if _, ok := excludedDomains[domain]; ok {
			log.Info("Skipping result; excluded domain", "domain", domain, "url", u.String())
			continue
		}
		t := api.Site{
			Domain: domain,
		}

		sites = append(sites, t)
	}

	return sites, nil
}
