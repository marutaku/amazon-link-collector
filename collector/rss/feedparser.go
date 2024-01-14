package rss

import (
	"github.com/marutaku/amazon-link-collector/collector/domain"
	"github.com/mmcdole/gofeed"
)

type FeedParser struct {
	feedBaseURL string
}

func NewFeedParser(feedBaseURL string) *FeedParser {
	return &FeedParser{
		feedBaseURL: feedBaseURL,
	}
}

func (f *FeedParser) Parse() ([]*domain.Bookmark, error) {
	feed, err := gofeed.NewParser().ParseURL(f.feedBaseURL)
	if err != nil {
		return nil, err
	}
	var bookmarks []*domain.Bookmark
	for _, item := range feed.Items {
		if item == nil {
			break
		}
		bookmarks = append(bookmarks, domain.NewBookmark(item.Title, *item.PublishedParsed, item.Link))
	}
	return bookmarks, nil
}
