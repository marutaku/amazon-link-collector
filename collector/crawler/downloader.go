package crawler

import (
	"io"
	"net/http"
	"sync"
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
	hostname, err := ExtractHostname(url)
	if err != nil {
		errorsArray[i] = err
		return
	}
	if d.cache.IsCached(url) {
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
	d.cache.StoreBookmarkCache(url, stringBody)
	contentsArray[i] = stringBody
}

func (d *Downloader) BulkDownload(urls []string) ([]string, error) {
	for _, url := range urls {
		hostname, err := ExtractHostname(url)
		if err != nil {
			return nil, err
		}
		if _, ok := d.originsConnections[hostname]; !ok {
			d.originsConnections[hostname] = make(chan struct{}, MAX_CONCURRENT_DOWNLOAD_IN_SAME_ORIGIN)
		}
	}

	var wg sync.WaitGroup
	contentsArray := make([]string, len(urls))
	errorsArray := make([]error, len(urls))

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			hostname, err := ExtractHostname(url)
			if err != nil {
				errorsArray[i] = err
				return
			}
			if d.cache.IsCached(url) {
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
			d.cache.StoreBookmarkCache(url, stringBody)
			contentsArray[i] = stringBody
		}(i, url)
	}

	wg.Wait()

	for _, err := range errorsArray {
		if err != nil {
			return nil, err
		}
	}

	return contentsArray, nil
}
