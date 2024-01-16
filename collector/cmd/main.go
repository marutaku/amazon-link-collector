package main

import (
	"log"
	"os"

	"encoding/json"

	"github.com/marutaku/amazon-link-collector/collector/crawler"
	"github.com/marutaku/amazon-link-collector/collector/rss"
	"github.com/urfave/cli/v2"
)

func crawl(bookmarkURL string, exportJsonLinePath string) error {
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
	outputFile, err := os.Create(exportJsonLinePath)
	if err != nil {
		return err
	}
	for index, content := range contents {
		amazonLinks, err := crawler.ExtractAmazonLink(content)
		if err != nil {
			return err
		}
		bookmarks[index].AmazonLinks = amazonLinks
		jsonB, err := json.Marshal(bookmarks[index])
		if err != nil {
			return err
		}
		outputFile.WriteString(string(jsonB) + "\n")
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
			exportJsonLinePath := cCtx.Args().Get(1)
			if exportJsonLinePath == "" {
				log.Fatal("Export JSON line path is required")
			}
			err := crawl(rssFeedURL, exportJsonLinePath)
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
