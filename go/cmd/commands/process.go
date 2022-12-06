package commands

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/jlewi/eyecare/go/pkg/analyzer"
	"github.com/jlewi/eyecare/go/pkg/scraping"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"regexp"
)

func NewProcessCmd() *cobra.Command {
	var domainsFile string
	var htmlDir string
	cmd := &cobra.Command{
		Short: "process",
		Use:   "process",
		Run: func(cmd *cobra.Command, args []string) {
			log := zapr.NewLogger(zap.L())
			err := func() error {
				sites, err := scraping.ReadFile(domainsFile)
				if err != nil {
					return err
				}
				log.Info("Read sites", "count", len(sites), "file", domainsFile)

				p := analyzer.RegexAnalyzer{
					Matches: make([]*regexp.Regexp, 0, 30),
				}

				patterns := []string{"(?i)keratoconus", "(?i)topography", "(?i)topographers"}
				//patterns := []string{"(?i)topography", "(?i)topographers"}

				for _, pat := range patterns {
					e, err := regexp.Compile(pat)
					if err != nil {
						return errors.Wrapf(err, "Failed to compile regex: %v", pat)
					}
					p.Matches = append(p.Matches, e)
				}

				for i, s := range sites {
					if s.Results != nil {
						log.Info("Skipping already processed domain", "domain", s.Domain)
					}

					domainHtml := filepath.Join(htmlDir, s.Domain)

					if _, err := os.Stat(domainHtml); os.IsNotExist(err) {
						log.Info("Skipping domain; no scraping output found", "directory", domainHtml)
						continue
					}

					results, err := p.Process(domainHtml)
					if err != nil {
						return errors.Wrapf(err, "Failed to process directory: %v", domainHtml)
					}

					sites[i].Results = results
					log.Info("Writing results to file", "file", domainsFile)
					if err := scraping.WriteFile(domainsFile, sites); err != nil {
						return err
					}
				}

				return nil
			}()
			if err != nil {
				fmt.Fprintf(os.Stdout, "Failed with error:\n%v", err)
			}
		},
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get homedirectory; only default output directory is affected")
		homeDir = "/"
	}
	defaultDir := filepath.Join(homeDir, "scraped_sites")
	cmd.Flags().StringVarP(&domainsFile, "domains", "i", "results/results.jsonl", "The JSONL file containing a list of sites to scrape. This will also be where the results will be written")
	cmd.Flags().StringVarP(&htmlDir, "html-dir", "d", defaultDir, "The directory containing the scraped sites for all the files.")
	return cmd
}
