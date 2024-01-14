package crawler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/marutaku/amazon-link-collector/collector/utils"
)

type Cache interface {
	IsCached(urlString string) (bool, error)
	StoreBookmarkCache(urlString string, content string) error
	GetBookmarkCache(urlString string) (string, error)
}

type LocalCache struct {
	cacheDir string
}

func hashURL(url string) string {
	r := sha256.Sum256([]byte(url))
	return hex.EncodeToString(r[:])
}

func (c *LocalCache) buildCacheFilepath(urlString string) (string, error) {
	hostname, err := utils.ExtractHostname(urlString)
	if err != nil {
		return "", err
	}
	hashedURL := hashURL(urlString)
	return path.Join(c.cacheDir, hostname, fmt.Sprintf("%s.txt", hashedURL)), nil
}

func NewLocalCache(cacheDir string) *LocalCache {
	return &LocalCache{
		cacheDir: cacheDir,
	}
}

func (c *LocalCache) IsCached(urlString string) (bool, error) {
	cacheFilepath, err := c.buildCacheFilepath(urlString)
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(cacheFilepath); !os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

func (c *LocalCache) StoreBookmarkCache(urlString string, content string) error {
	cacheFilepath, err := c.buildCacheFilepath(urlString)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(path.Dir(cacheFilepath), 0755); err != nil {
		return err
	}
	f, err := os.Create(cacheFilepath)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}

func (c *LocalCache) GetBookmarkCache(urlString string) (string, error) {
	cacheFilepath, err := c.buildCacheFilepath(urlString)
	if err != nil {
		return "", err
	}
	contents, err := os.ReadFile(cacheFilepath)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
