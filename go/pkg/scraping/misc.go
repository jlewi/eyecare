package scraping

import (
	"net/url"
	"strings"
)

// URLToDomain normalizes the URL to a domain.
// Returns empty string on error.
func URLToDomain(u string) string {
	p, err := url.Parse(u)
	if err != nil {
		return ""
	}

	host := p.Host
	if p.Scheme == "" {
		host = u
	}

	return HostToDomain(host)
}

// HostToDomain gets the domain from the hostname.
func HostToDomain(host string) string {
	domainPieces := strings.Split(host, ".")
	if strings.ToLower(domainPieces[0]) == "www" {
		domainPieces = domainPieces[1:len(domainPieces)]
	}
	domain := strings.Join(domainPieces, ".")
	return domain
}
