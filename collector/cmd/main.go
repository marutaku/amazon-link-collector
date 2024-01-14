package main

import (
	"log"
	"os"

	"github.com/marutaku/amazon-link-collector/collector/crawler"
	"github.com/marutaku/amazon-link-collector/collector/rss"
	"github.com/urfave/cli/v2"
)

func crawl(bookmarkURL string) error {
	feedparser := rss.NewFeedParser(bookmarkURL)
	cache := crawler.NewLocalCache("./.cache")
	downloader := crawler.NewDownloader(cache)
	bookmarks, err := feedparser.Parse()
	if err != nil {
		return err
	}
	contents, err := downloader.BulkDownload(bookmarks)
	if err != nil {
		return err
	}
	for _, content := range contents {
		amazonLinks, err := crawler.ExtractAmazonLink(content)
		if err != nil {
			return err
		}
		for _, amazonLink := range amazonLinks {
			log.Println(amazonLink)
		}
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:  "Amazon Link Collector",
		Usage: "Collect Amazon links from RSS feed",
		Action: func(cCtx *cli.Context) error {
			rssFeedURL := cCtx.Args().Get(0)
			if rssFeedURL == "" {
				log.Fatal("RSS feed URL is required")
			}
			err := crawl(rssFeedURL)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
