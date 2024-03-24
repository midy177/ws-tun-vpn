package cachex

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// The global cachex
var _cache = cache.New(30*time.Minute, 10*time.Minute)

// GetCache returns the cachex
func GetCache() *cache.Cache {
	return _cache
}
