package build

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
	memcache "github.com/pip-services3-go/pip-services3-memcached-go/cache"
	memlock "github.com/pip-services3-go/pip-services3-memcached-go/lock"
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

	c.RegisterType(c.MemcachedCacheDescriptor, memcache.NewMemcachedCache)
	c.RegisterType(c.MemcachedLockDescriptor, memlock.NewMemcachedLock)
	return &c
}
