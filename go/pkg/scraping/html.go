package scraping

import (
	"fmt"
	zerowidth "github.com/trubitsyn/go-zero-width"
	"golang.org/x/net/html"
	"io"
	"strings"
)

const (
	SoftHyphen = "\u00ad"
)

// TextFromHtml returns all the text from html.
// Based on: https://stackoverflow.com/questions/44441665/how-to-extract-only-text-from-html-in-golang
func TextFromHtml(reader io.Reader) string {
	domDocTest := html.NewTokenizer(reader)
	previousStartTokenTest := domDocTest.Token()

	var b strings.Builder
	for {
		tt := domDocTest.Next()
		switch {
		case tt == html.ErrorToken:
			// End of the document,  done
			return b.String()
		case tt == html.StartTagToken:
			previousStartTokenTest = domDocTest.Token()
		case tt == html.TextToken:
			if previousStartTokenTest.Data == "script" || previousStartTokenTest.Data == "style" {
				continue
			}

			TxtContent := strings.TrimSpace(html.UnescapeString(string(domDocTest.Text())))
			// Remove some unprintable characters. These could potentially interfere
			// with string matching. 
			TxtContent = zerowidth.RemoveZeroWidthCharacters(TxtContent)
			TxtContent = strings.Replace(TxtContent, SoftHyphen, "", -1)
			if len(TxtContent) > 0 {
				fmt.Fprintf(&b, "%s\n", TxtContent)
			}
		}
	}
}
