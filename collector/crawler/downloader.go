package crawler

import (
	"io"
	"net/http"
	"sync"

	"github.com/marutaku/amazon-link-collector/collector/domain"
	"github.com/marutaku/amazon-link-collector/collector/utils"
)

var MAX_CONCURRENT_DOWNLOAD_IN_SAME_ORIGIN = 1

type Downloader struct {
	originsConnections map[string]chan struct{}
	cache              Cache
}

func NewDownloader(cache Cache) *Downloader {
	return &Downloader{
		originsConnections: map[string]chan struct{}{},
		cache:              cache,
	}
}

func (d *Downloader) download(i int, url string, contentsArray []string, errorsArray []error, wg *sync.WaitGroup) {
	defer wg.Done()
	hostname, err := utils.ExtractHostname(url)
	if err != nil {
		errorsArray[i] = err
		return
	}
	exists, err := d.cache.IsCached(url)
	if err != nil {
		errorsArray[i] = err
		return
	}
	if exists {
		content, err := d.cache.GetBookmarkCache(url)
		if err != nil {
			errorsArray[i] = err
			return
		}
		contentsArray[i] = content
		return
	}
	d.originsConnections[hostname] <- struct{}{}
	defer func() { <-d.originsConnections[hostname] }()

	resp, err := http.Get(url)
	if err != nil {
		errorsArray[i] = err
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorsArray[i] = err
		return
	}
	stringBody := string(body)
	err = d.cache.StoreBookmarkCache(url, stringBody)
	if err != nil {
		errorsArray[i] = err
		return
	}
	contentsArray[i] = stringBody
}

func (d *Downloader) BulkDownload(bookmarks []*domain.Bookmark) ([]string, error) {
	for _, bookmark := range bookmarks {
		hostname, err := utils.ExtractHostname(bookmark.URL)
		if err != nil {
			return nil, err
		}
		if _, ok := d.originsConnections[hostname]; !ok {
			d.originsConnections[hostname] = make(chan struct{}, MAX_CONCURRENT_DOWNLOAD_IN_SAME_ORIGIN)
		}
	}

	var wg sync.WaitGroup
	contentsArray := make([]string, len(bookmarks))
	errorsArray := make([]error, len(bookmarks))

	for i, bookmark := range bookmarks {
		wg.Add(1)
		go d.download(i, bookmark.URL, contentsArray, errorsArray, &wg)
	}

	wg.Wait()

	for _, err := range errorsArray {
		if err != nil {
			return nil, err
		}
	}

	return contentsArray, nil
}
