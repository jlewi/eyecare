package scraping

import "strings"

// HostToDomain gets the domain from the hostname.
func HostToDomain(host string) string {
	domainPieces := strings.Split(host, ".")
	if strings.ToLower(domainPieces[0]) == "www" {
		domainPieces = domainPieces[1:len(domainPieces)]
	}
	domain := strings.Join(domainPieces, ".")
	return domain
}
