package crawler

type Cache interface {
	IsCached(url string) bool
	StoreCache(url string) error
	GetCache(url string) (string, error)
}

type LocalCache struct {
	cacheDir string
}

func NewLocalCache(cacheDir string) *LocalCache {
	return &LocalCache{
		cacheDir: cacheDir,
	}
}
