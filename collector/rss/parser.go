package rss

import (
	"github.com/marutaku/amazon-link-collector/collector/domain"
)

type RSSFeedParser interface {
	Parse() ([]*domain.Bookmark, error)
}
