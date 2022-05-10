# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Memcached components for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit.

The Memcached module contains the following components: MemcachedLock and MemcachedCache for working with locks and cache on the Memcached server.

The module contains the following packages:
- [**Build**](https://godoc.org/github.com/pip-services3-go/pip-services3-memcached-go/build) - a standard factory for constructing components
- [**Cache**](https://godoc.org/github.com/pip-services3-go/pip-services3-memcached-go/cache) - cache Components in Memcached
- [**Lock**](https://godoc.org/github.com/pip-services3-go/pip-services3-memcached-go/lock) - components of working with locks in Memcached

<a name="links"></a> Quick links:

* [Configuration](https://www.pipservices.org/recipies/configuration)
* [API Reference](https://godoc.org/github.com/pip-services3-go/pip-services3-memcached-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](https://www.pipservices.org/community/help)
* [Contribute](https://www.pipservices.org/community/contribute)


## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services3-go/pip-services3-memcached-go@latest
```

## Develop

For development you shall install the following prerequisites:
* Golang v1.12+
* Visual Studio Code or another IDE of your choice
* Docker
* Git

Run automated tests:
```bash
go test -v ./test/...
```

Generate API documentation:
```bash
./docgen.ps1
```

Before committing changes run dockerized test as:
```bash
./test.ps1
./clear.ps1
```

## Contacts

The library is created and maintained by **Sergey Seroukhov** and **Levichev Dmitry**.

The documentation is written by:
- **Levichev Dmitry**