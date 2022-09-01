package test_lock

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	memlock "github.com/pip-services3-gox/pip-services3-memcached-gox/lock"
	memfixture "github.com/pip-services3-gox/pip-services3-memcached-gox/test/fixture"
)

func TestMemcachedLock(t *testing.T) {
	var lock *memlock.MemcachedLock
	var fixture *memfixture.LockFixture

	ctx := context.Background()

	host := os.Getenv("MEMCACHED_SERVICE_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("MEMCACHED_SERVICE_PORT")
	if port == "" {
		port = "11211"
	}

	lock = memlock.NewMemcachedLock()

	config := cconf.NewConfigParamsFromTuples(
		"connection.host", host,
		"connection.port", port,
	)
	lock.Configure(ctx, config)
	fixture = memfixture.NewLockFixture(lock)

	lock.Open(ctx, "")
	defer lock.Close(ctx, "")

	t.Run("Try Acquire Lock", fixture.TestTryAcquireLock)
	t.Run("Acquire Lock", fixture.TestAcquireLock)
	t.Run("Release Lock", fixture.TestReleaseLock)
}
