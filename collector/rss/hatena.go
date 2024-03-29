package rss

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/marutaku/amazon-link-collector/collector/domain"
	"github.com/mmcdole/gofeed"
)

type FeedParser struct {
	feedBaseURL string
	logger      *log.Logger
}

func NewFeedParser(feedBaseURL string, logger *log.Logger) *FeedParser {
	return &FeedParser{
		feedBaseURL: feedBaseURL,
		logger:      logger,
	}
}

func (f *FeedParser) Parse() ([]*domain.Bookmark, error) {
	var bookmarks []*domain.Bookmark
	pageNum := 1
	for {
		u, err := url.Parse(f.feedBaseURL)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		q.Set("page", fmt.Sprint(pageNum))
		u.RawQuery = q.Encode()
		url := u.String()
		f.logger.Println(url)
		feed, err := gofeed.NewParser().ParseURL(url)
		if err != nil {
			return nil, err
		}
		if len(feed.Items) == 0 {
			break
		}
		for _, item := range feed.Items {
			if item == nil {
				break
			}
			bookmarks = append(bookmarks, domain.NewBookmark(item.Title, *item.PublishedParsed, item.Link))
		}
		time.Sleep(1 * time.Second)
		pageNum++
	}

	return bookmarks, nil
}
