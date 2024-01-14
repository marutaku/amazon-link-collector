package crawler

import (
	"github.com/marutaku/amazon-link-collector/collector/domain"
)

var MAX_CONCURRENT_DOWNLOAD_NUM = 1

type Downloader struct {
	originsConnections map[string]chan string
	cache              Cache
}

func NewDownloader(cache Cache) *Downloader {
	return &Downloader{
		originsConnections: map[string]chan string{},
		cache:              cache,
	}
}

func download(ch chan string, cache Cache) {
	for {
		url := <-ch
		if content, err := cache.GetBookmarkCache(url); err != nil {
			ch <- content
		} else {
			ch <- ""
		}
	}
}

func (d *Downloader) BulkDownload(bookmarks []domain.Bookmark) ([]string, error) {
	for _, bookmarks := range bookmarks {
		hostname, err := bookmarks.Hostname()
		if err != nil {
			return nil, err
		}
		if _, ok := d.originsConnections[hostname]; !ok {
			d.originsConnections[hostname] = make(chan string, MAX_CONCURRENT_DOWNLOAD_NUM)
		}
	}
	contentsArray := make([]string, len(bookmarks))
	go func() {
		for index, bookmark := range bookmarks {
			hostname, err := bookmark.Hostname()
			if err != nil {
				return
			}
			d.originsConnections[hostname] <- bookmark.URL
			contents := <-d.originsConnections[hostname]
			contentsArray[index] = contents
		}
	}()
	return contentsArray, nil
}
