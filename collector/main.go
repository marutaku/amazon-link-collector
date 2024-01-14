package main

import (
	"log"
	"os"

	"github.com/marutaku/amazon-link-collector/collector/rss"
	"github.com/urfave/cli/v2"
)

func crawl(bookmarkURL string) error {
	feedparser := rss.NewFeedParser(bookmarkURL)
	bookmarks, err := feedparser.Parse()
	if err != nil {
		return err
	}
	urls := make([]string, len(bookmarks))
	for index, bookmark := range bookmarks {
		urls[index] = bookmark.URL
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
