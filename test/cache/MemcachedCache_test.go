package test_cache

import (
	"os"
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	memcache "github.com/pip-services3-go/pip-services3-memcached-go/cache"
	memfixture "github.com/pip-services3-go/pip-services3-memcached-go/test/fixture"
)

func TestMemcachedCache(t *testing.T) {
	var cache *memcache.MemcachedCache
	var fixture *memfixture.CacheFixture

	host := os.Getenv("MEMCACHED_SERVICE_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("MEMCACHED_SERVICE_PORT")
	if port == "" {
		port = "11211"
	}

	cache = memcache.NewMemcachedCache()
	config := cconf.NewConfigParamsFromTuples(
		"connection.host", host,
		"connection.port", port,
	)
	cache.Configure(config)
	fixture = memfixture.NewCacheFixture(cache)
	cache.Open("")
	defer cache.Close("")

	t.Run("TestMemcachedCache:Store and Retrieve", fixture.TestStoreAndRetrieve)
	t.Run("TestMemcachedCache:Retrieve Expired", fixture.TestRetrieveExpired)
	t.Run("TestMemcachedCache:Remove", fixture.TestRemove)
}
