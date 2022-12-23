# eyecare
Code to do some analysis for eye care

# Build the code

```
make build-go
```
# SERP

You can use Google Search to find sites to crawl.
This uses Nimbleway's SERP API so you will need a Nimbleway account.

```
build/eyecare --outfile=results/results.jsonl --query="optometrist in san mateo"
```
# Manually Adding A Domain

To manually add a domain to the list

```
build/eyecare domains add https://pacificaeyedoctor.com/
```

# Crawling & Scraping

To scrape all the sites listed in the JSONL file.

```
build/eyecare scrape --input=results/results.jsonl --output=~/scraped_sites
```

This will scrape the sites and dump the output in the directory specified 
as the output argument.

## URLs

Dashes in URL names are replaced with underscores in the file.

e.g. the file `clearvision1.com/about_me.unknown` 
corresponds to the URL [clearvision1.com/about-me](clearvision1.com/about-me)

