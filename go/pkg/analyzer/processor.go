package analyzer

import (
	"github.com/go-logr/zapr"
	"github.com/jlewi/eyecare/go/api"
	"github.com/jlewi/eyecare/go/pkg/scraping"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"regexp"
)

// RegexAnalyzer looks for matches to a bunch of regexes.
type RegexAnalyzer struct {
	Matches []*regexp.Regexp
}

// Process processes all the files in the specified directory
func (a *RegexAnalyzer) Process(dir string) ([]api.Result, error) {
	log := zapr.NewLogger(zap.L())
	log.Info("Walking directory", "directory", dir)

	results := make([]api.Result, 0, 30)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "There was a problem walking directory; %v", path)
		}

		if info.IsDir() {
			return nil
		}

		log.Info("Handling file", "file", path)
		f, err := os.Open(path)

		if err != nil {
			return errors.Wrapf(err, "Could not open file: %v", f)
		}

		text := scraping.TextFromHtml(f)

		for _, m := range a.Matches {
			if m.MatchString(text) {
				results = append(results, api.Result{
					Path: path,
					Term: m.String(),
				})
			}
		}

		return nil
	})

	return results, err
}
