package build

import (
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cbuild "github.com/pip-services3-gox/pip-services3-components-gox/build"
	memcache "github.com/pip-services3-gox/pip-services3-memcached-gox/cache"
	memlock "github.com/pip-services3-gox/pip-services3-memcached-gox/lock"
)

// DefaultMemcachedFactory Creates Redis components by their descriptors.
// See MemcachedCache
// See MemcachedLock
type DefaultMemcachedFactory struct {
	*cbuild.Factory
	Descriptor               *cref.Descriptor
	MemcachedCacheDescriptor *cref.Descriptor
	MemcachedLockDescriptor  *cref.Descriptor
}

// NewDefaultMemcachedFactory Create a new instance of the factory.
func NewDefaultMemcachedFactory() *DefaultMemcachedFactory {
	c := DefaultMemcachedFactory{}
	c.Factory = cbuild.NewFactory()

	c.Descriptor = cref.NewDescriptor("pip-services", "factory", "memcached", "default", "1.0")
	c.MemcachedCacheDescriptor = cref.NewDescriptor("pip-services", "cache", "memcached", "*", "1.0")
	c.MemcachedLockDescriptor = cref.NewDescriptor("pip-services", "lock", "memcached", "*", "1.0")

	c.RegisterType(c.MemcachedCacheDescriptor, memcache.NewMemcachedCache[any])
	c.RegisterType(c.MemcachedLockDescriptor, memlock.NewMemcachedLock)
	return &c
}
