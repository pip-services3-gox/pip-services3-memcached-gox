package lock

import (
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	ccon "github.com/pip-services3-go/pip-services3-components-go/connect"
	clock "github.com/pip-services3-go/pip-services3-components-go/lock"
)

/*
MemcachedLock are distributed lock that implemented based on Memcaches caching service.

The current implementation does not support authentication.

Configuration parameters:

- connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from IDiscovery
  - host:                  host name or IP address
  - port:                  port number
  - uri:                   resource URI or connection string with all parameters in it
- options:
  - max_size:              maximum number of values stored in this cache (default: 1000)
  - max_key_size:          maximum key length (default: 250)
  - max_expiration:        maximum expiration duration in milliseconds (default: 2592000)
  - max_value:             maximum value length (default: 1048576)
  - pool_size:             pool size (default: 5)
  - reconnect:             reconnection timeout in milliseconds (default: 10 sec)
  - retries:               number of retries (default: 3)
  - timeout:               default caching timeout in milliseconds (default: 1 minute)
    - failures:              number of failures before stop retrying (default: 5)
    - retry:                 retry timeout in milliseconds (default: 30 sec)
    - idle:                  idle timeout before disconnect in milliseconds (default: 5 sec)

References:

- *:discovery:*:*:1.0        (optional) IDiscovery services to resolve connection

Example:

    lock := NewMemcachedLock();
    lock.Configure(cconf.NewConfigParamsFromTuples(
      "host", "localhost",
      "port", 11211,
    ));

    err := lock.Open("123")
    if err != nil {
      ...
    }

    result, err := lock.TryAcquireLock("123", "key1", 3000)
    if result {
    	// Processing...
    }
    err = lock.ReleaseLock("123", "key1")
    // Continue...
*/
type MemcachedLock struct {
	*clock.Lock
	connectionResolver *ccon.ConnectionResolver
	// maxKeySize         int
	// maxExpiration      int64
	// maxValue           int64
	// poolSize           int
	// reconnect          int
	timeout int
	// retries            int
	// failures           int
	// retry              int
	remove bool
	//idle   int
	client *memcache.Client
}

// NewMemcachedLock method are creates a new instance of this lock.
func NewMemcachedLock() *MemcachedLock {
	c := &MemcachedLock{
		connectionResolver: ccon.NewEmptyConnectionResolver(),
		// maxKeySize:         250,
		// maxExpiration:      2592000,
		// maxValue:           1048576,
		// poolSize:           5,
		// reconnect:          10000,
		timeout: 5000,
		// retries:            5,
		// failures:           5,
		// retry:              30000,
		remove: false,
		//idle:   5000,
		client: nil,
	}
	c.Lock = clock.InheritLock(c)
	return c
}

// Configure method are configures component by passing configuration parameters.
//   - config    configuration parameters to be set.
func (c *MemcachedLock) Configure(config *cconf.ConfigParams) {
	c.Lock.Configure(config)

	c.connectionResolver.Configure(config)

	// c.maxKeySize = config.GetAsIntegerWithDefault("options.max_key_size", c.maxKeySize)
	// c.maxExpiration = config.GetAsLongWithDefault("options.max_expiration", c.maxExpiration)
	// c.maxValue = config.GetAsLongWithDefault("options.max_value", c.maxValue)
	// c.poolSize = config.GetAsIntegerWithDefault("options.pool_size", c.poolSize)
	// c.reconnect = config.GetAsIntegerWithDefault("options.reconnect", c.reconnect)
	c.timeout = config.GetAsIntegerWithDefault("options.timeout", c.timeout)
	// c.retries = config.GetAsIntegerWithDefault("options.retries", c.retries)
	// c.failures = config.GetAsIntegerWithDefault("options.failures", c.failures)
	// c.retry = config.GetAsIntegerWithDefault("options.retry", c.retry)
	// c.remove = config.GetAsBooleanWithDefault("options.remove", c.remove)
	//c.idle = config.GetAsIntegerWithDefault("options.idle", c.idle)
}

// SetReferences method are sets references to dependent components.
//   - references 	references to locate the component dependencies.
func (c *MemcachedLock) SetReferences(references cref.IReferences) {
	c.connectionResolver.SetReferences(references)
}

// IsOpen method are checks if the component is opened.
// Returns: true if the component has been opened and false otherwise.
func (c *MemcachedLock) IsOpen() bool {
	return c.client != nil
}

/// Open method are opens the component.
// Parameters:
//   - correlationId 	(optional) transaction id to trace execution through call chain.
// Retruns: error or nil no errors occured.
func (c *MemcachedLock) Open(correlationId string) error {
	connections, err := c.connectionResolver.ResolveAll(correlationId)

	if err == nil && len(connections) == 0 {
		err = cerr.NewConfigError(correlationId, "NO_CONNECTION", "Connection is not configured")
	}

	if err != nil {
		return err
	}

	var servers []string = make([]string, 0)
	for _, connection := range connections {
		host := connection.Host()
		port := connection.Port()
		if port == 0 {
			port = 11211
		}

		servers = append(servers, host+":"+strconv.FormatInt(int64(port), 10))
	}

	// options = {
	//     maxKeySize: c.maxKeySize,
	//     maxExpiration: c.maxExpiration,
	//     maxValue: c.maxValue,
	//     poolSize: c.poolSize,
	//     reconnect: c.reconnect,
	//     timeout: c.timeout,
	//     retries: c.retries,
	//     failures: c.failures,
	//     retry: c.retry,
	//     remove: c.remove,
	//     idle: c.idle
	// };

	c.client = memcache.New(servers...)
	c.client.Timeout = time.Duration(c.timeout) * time.Millisecond
	//c.client.MaxIdleConns = c.idle

	return nil
}

// Close method are closes component and frees used resources.
// Parameters:
//   - correlationId 	(optional) transaction id to trace execution through call chain.
//   - callback 			callback function that receives error or nil no errors occured.
func (c *MemcachedLock) Close(correlationId string) error {
	c.client = nil
	return nil
}

func (c *MemcachedLock) checkOpened(correlationId string) (state bool, err error) {
	if !c.IsOpen() {
		err = cerr.NewInvalidStateError(correlationId, "NOT_OPENED", "Connection is not opened")
		return false, err
	}

	return true, nil
}

// TryAcquireLock method are makes a single attempt to acquire a lock by its key.
// It returns immediately a positive or negative result.
// Parameters:
//    - correlationId     (optional) transaction id to trace execution through call chain.
//    - key               a unique lock key to acquire.
//    - ttl               a lock timeout (time to live) in milliseconds.
//  Returns: a lock result or error.
func (c *MemcachedLock) TryAcquireLock(correlationId string, key string, ttl int64) (result bool, err error) {

	state, err := c.checkOpened(correlationId)
	if !state {
		return false, err
	}

	lifetimeInSec := ttl / 1000
	item := memcache.Item{
		Key:        key,
		Value:      []byte("lock"),
		Expiration: int32(lifetimeInSec),
	}
	err = c.client.Add(&item)

	if err != nil && err == memcache.ErrNotStored {
		return false, nil
	}
	return err == nil, err

}

// ReleaseLock method are releases prevously acquired lock by its key.
//    - correlationId     (optional) transaction id to trace execution through call chain.
//    - key               a unique lock key to release.
//  Returns error or nil for success.
func (c *MemcachedLock) ReleaseLock(correlationId string, key string) error {
	state, err := c.checkOpened(correlationId)
	if !state {
		return err
	}
	err = c.client.Delete(key)
	if err != nil && err == memcache.ErrCacheMiss {
		err = nil
	}
	return err
}
