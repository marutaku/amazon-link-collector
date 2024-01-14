package crawler

import "github.com/marutaku/amazon-link-collector/collector/domain"

var MAX_CONCURRENT_DOWNLOAD_NUM = 1

type Downloader struct {
	originsConnections map[string]chan string
}

func NewDownloader() *Downloader {
	return &Downloader{
		originsConnections: map[string]chan string{},
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
	return []string{}, nil
}
