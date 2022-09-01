package test_cache

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	memcache "github.com/pip-services3-gox/pip-services3-memcached-gox/cache"
	memfixture "github.com/pip-services3-gox/pip-services3-memcached-gox/test/fixture"
)

func TestMemcachedCache(t *testing.T) {
	ctx := context.Background()

	var cache *memcache.MemcachedCache[any]
	var fixture *memfixture.CacheFixture

	host := os.Getenv("MEMCACHED_SERVICE_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("MEMCACHED_SERVICE_PORT")
	if port == "" {
		port = "11211"
	}

	cache = memcache.NewMemcachedCache[any]()
	config := cconf.NewConfigParamsFromTuples(
		"connection.host", host,
		"connection.port", port,
	)
	cache.Configure(ctx, config)
	fixture = memfixture.NewCacheFixture(cache)
	cache.Open(ctx, "")
	defer cache.Close(ctx, "")

	t.Run("TestMemcachedCache:Store and Retrieve", fixture.TestStoreAndRetrieve)
	t.Run("TestMemcachedCache:Retrieve Expired", fixture.TestRetrieveExpired)
	t.Run("TestMemcachedCache:Remove", fixture.TestRemove)
}
