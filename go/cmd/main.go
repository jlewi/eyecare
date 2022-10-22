// scraping is a project to collect statistics about keratoconus and the prevalence of
// corneal topography by scraping some websites. It uses nimbleway.
package main

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/jlewi/eyecare/go/cmd/commands"
	"github.com/jlewi/eyecare/go/pkg/helpers"
	"github.com/jlewi/eyecare/go/pkg/scraping"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	url2 "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/jlewi/p22h/backend/pkg/logging"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	var level string
	var jsonLog bool
	rootCmd := &cobra.Command{
		Short: "scraping CLI",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			_, err := logging.InitLogger(level, !jsonLog)
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&level, "level", "", "info", "The logging level.")
	rootCmd.PersistentFlags().BoolVarP(&jsonLog, "json-logs", "", false, "Enable json logging.")
	return rootCmd
}

func newScrapeCmd() *cobra.Command {
	var url string
	var output string
	cmd := &cobra.Command{
		Use:   "scrape",
		Short: "scraping CLI",
		Run: func(cmd *cobra.Command, args []string) {
			log := zapr.NewLogger(zap.L())
			err := func() error {
				u, err := url2.Parse(url)
				if err != nil {
					return errors.Wrapf(err, "Failed to parse url %+v", u)
				}
				domain := scraping.HostToDomain(u.Host)

				dName := strings.Replace(domain, ".", "_", -1)

				output = filepath.Join(output, dName)

				log.Info("Using output directory", "directory", output)

				if _, err := os.Stat(output); os.IsNotExist(err) {
					log.Info("Creating output directory", "directory", output)
					os.MkdirAll(output, helpers.UserGroupAllPerm)
				}

				s, err := os.Stat(output)
				if err != nil {
					return errors.Wrapf(err, "Couldn't check directory %v", output)
				}

				if !s.IsDir() {
					return errors.Errorf("path %v exists and is not a directory", output)
				}

				c := colly.NewCollector()
				c.AllowedDomains = []string{domain, u.Host}
				log.Info("Setting allowed domain", "domain", colly.AllowedDomains)
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
				log.Info("Starting collector", "url", url)

				// TODO(jeremy): Should we look for a sitemap.
				if err := c.Visit(url); err != nil {
					return errors.Wrapf(err, "Failed to visit url: %v", url)
				}
				log.Info("Waiting for collector to finish")
				c.Wait()
				return nil
			}()
			if err != nil {
				fmt.Fprintf(os.Stdout, "Failed with error:\n%+v", err)
			}
		},
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get homedirectory; only default output directory is affected")
		homeDir = "/"
	}
	defaultDir := filepath.Join(homeDir, "scraped_sites")
	cmd.Flags().StringVarP(&url, "url", "u", "https://www.peneye.com/", "The url to scrape.")
	cmd.Flags().StringVarP(&output, "output", "o", defaultDir, "The directory to save the scrapings to.")
	return cmd
}

func main() {
	rootCmd := newRootCmd()
	rootCmd.AddCommand(newScrapeCmd())
	rootCmd.AddCommand(commands.NewSearchCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Command failed with error: %+v", err)
		os.Exit(1)
	}
}
