package domain

import (
	"net/url"
	"time"
)

type Bookmark struct {
	Title       string
	PublishedAt time.Time
	URL         string
}

func NewBookmark(title string, publishedAt time.Time, url string) *Bookmark {
	return &Bookmark{
		Title:       title,
		PublishedAt: publishedAt,
		URL:         url,
	}
}

func (b *Bookmark) String() string {
	return b.Title + " " + b.PublishedAt.String() + " " + b.URL
}

func (b *Bookmark) Hostname() (string, error) {
	parsedURL, err := url.Parse(b.URL)
	if err != nil {
		return "", err
	}
	return parsedURL.Hostname(), nil
}
