package test_lock

import (
	"os"
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	memlock "github.com/pip-services3-go/pip-services3-memcached-go/lock"
	memfixture "github.com/pip-services3-go/pip-services3-memcached-go/test/fixture"
)

func TestMemcachedLock(t *testing.T) {
	var lock *memlock.MemcachedLock
	var fixture *memfixture.LockFixture

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
	lock.Configure(config)
	fixture = memfixture.NewLockFixture(lock)

	lock.Open("")
	defer lock.Close("")

	t.Run("Try Acquire Lock", fixture.TestTryAcquireLock)
	t.Run("Acquire Lock", fixture.TestAcquireLock)
	t.Run("Release Lock", fixture.TestReleaseLock)
}
