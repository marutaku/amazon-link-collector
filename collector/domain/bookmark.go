package domain

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type Bookmark struct {
	Title       string
	PublishedAt time.Time
	URL         string
}

type BookmarkJSON struct {
	Title       string `json:"title"`
	PublishedAt string `json:"published_at"`
	URL         string `json:"url"`
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

func (b *Bookmark) MarshalJSON() ([]byte, error) {
	v, err := json.Marshal(&BookmarkJSON{
		Title:       b.Title,
		PublishedAt: b.PublishedAt.Format(time.RFC3339),
		URL:         b.URL,
	})
	return v, err
}

func (b *Bookmark) UnmarshalJSON(byte []byte) error {
	var BookmarkJSON *BookmarkJSON
	err := json.Unmarshal(byte, &BookmarkJSON)
	if err != nil {
		fmt.Println(err)
	}
	b.Title = BookmarkJSON.Title
	b.PublishedAt, err = time.Parse(time.RFC3339, BookmarkJSON.PublishedAt)
	if err != nil {
		return err
	}
	b.URL = BookmarkJSON.URL
	return err
}

func (b *Bookmark) Hostname() (string, error) {
	parsedURL, err := url.Parse(b.URL)
	if err != nil {
		return "", err
	}
	return parsedURL.Hostname(), nil
}
