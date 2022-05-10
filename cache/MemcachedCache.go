package cache

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	ccon "github.com/pip-services3-go/pip-services3-components-go/connect"
)

/*
MemcachedCache are distributed cache that stores values in Memcaches caching service.

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

- *:discovery:*:*:1.0    (optional) IDiscovery services to resolve connection

Example:

    cache := NewMemcachedCache();
    cache.Configure(cconf.NewConfigParamsFromTuples(
      "host", "localhost",
      "port", 11211,
    ));

    err := cache.Open("123")
      ...


    ret, err := cache.Store("123", "key1", []byte("ABC"))
    if err != nil {
    	...
    }

    res, err := cache.Retrive("123", "key1")
    value, _ := res.([]byte)
    fmt.Println(string(value))     // Result: "ABC"

*/
type MemcachedCache struct {
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

// NewMemcachedCache method are creates a new instance of this cache.
func NewMemcachedCache() *MemcachedCache {
	c := &MemcachedCache{
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
	return c
}

/*
   Configures component by passing configuration parameters.
   - config    configuration parameters to be set.
*/
func (c *MemcachedCache) Configure(config *cconf.ConfigParams) {
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

// SetReferences are sets references to dependent components.
//   - references 	references to locate the component dependencies.
func (c *MemcachedCache) SetReferences(references cref.IReferences) {
	c.connectionResolver.SetReferences(references)
}

// IsOpen Checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *MemcachedCache) IsOpen() bool {
	return c.client != nil
}

// Open method are opens the component.
// Parameters:
//   - correlationId 	(optional) transaction id to trace execution through call chain.
// Retruns: error or nil no errors occured.
func (c *MemcachedCache) Open(correlationId string) error {
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
// - correlationId 	(optional) transaction id to trace execution through call chain.
// - callback 			callback function that receives error or nil no errors occured.
func (c *MemcachedCache) Close(correlationId string) error {
	c.client = nil
	return nil
}

func (c *MemcachedCache) checkOpened(correlationId string) (state bool, err error) {
	if !c.IsOpen() {
		err = cerr.NewInvalidStateError(correlationId, "NOT_OPENED", "Connection is not opened")
		return false, err
	}

	return true, nil
}

// Retrieve method are retrieves cached value from the cache using its key.
// If value is missing in the cache or expired it returns nil.
// Parameters:
//   - correlationId     (optional) transaction id to trace execution through call chain.
//   - key               a unique value key.
// Retruns value *memcache.Item, err error
// cached value or error.
func (c *MemcachedCache) Retrieve(correlationId string, key string) (value interface{}, err error) {
	state, err := c.checkOpened(correlationId)
	if !state {
		return nil, err
	}
	item, err := c.client.Get(key)
	if err != nil && err == memcache.ErrCacheMiss {
		err = nil
	}
	if item != nil {
		var value interface{}
		err := json.Unmarshal(item.Value, &value)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
	return nil, err
}

// Retrieve method are retrieves cached value from the cache using its key.
// If value is missing in the cache or expired it returns nil.
// Parameters:
//   - correlationId     (optional) transaction id to trace execution through call chain.
//   - key               a unique value key.
// Retruns value *memcache.Item, err error
// cached value or error.
func (c *MemcachedCache) RetrieveAs(correlationId string, key string, result interface{}) (value interface{}, err error) {
	state, err := c.checkOpened(correlationId)
	if !state {
		return nil, err
	}
	item, err := c.client.Get(key)
	if err != nil && err == memcache.ErrCacheMiss {
		err = nil
	}
	if item != nil {
		err = json.Unmarshal(item.Value, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err
}

// Store method are stores value in the cache with expiration time.
// Parameters:
//   - correlationId     (optional) transaction id to trace execution through call chain.
//   - key               a unique value key.
//   - value             a value to store.
//   - timeout           expiration timeout in milliseconds.
// Returns: error or nil for success
func (c *MemcachedCache) Store(correlationId string, key string, value interface{}, timeout int64) (result interface{}, err error) {
	state, err := c.checkOpened(correlationId)
	if !state {
		return nil, err
	}
	timeoutInSec := int32(timeout) / 1000

	var val []byte

	switch v := value.(type) {
	case []byte:
		val = v
		break
	case string:
		val, err = json.Marshal(v)
		break
	default:
		val, err = json.Marshal(v)
		break
	}

	if err != nil {
		return nil, err
	}

	item := memcache.Item{
		Key:        key,
		Value:      val,
		Expiration: timeoutInSec,
	}
	return value, c.client.Set(&item)
}

// Remove method are removes a value from the cache by its key.
// Parameters:
//   - correlationId     (optional) transaction id to trace execution through call chain.
//   - key               a unique value key.
// Retruns: an error or nil for success
func (c *MemcachedCache) Remove(correlationId string, key string) error {
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
