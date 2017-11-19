[![Build Status](https://travis-ci.org/frictionlessdata/datapackage-go.svg?branch=master)](https://travis-ci.org/frictionlessdata/datapackage-go) [![Coverage Status](https://coveralls.io/repos/github/frictionlessdata/datapackage-go/badge.svg?branch=master)](https://coveralls.io/github/frictionlessdata/datapackage-go?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/frictionlessdata/datapackage-go)](https://goreportcard.com/report/github.com/frictionlessdata/datapackage-go) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/frictionlessdata/chat) [![GoDoc](https://godoc.org/github.com/frictionlessdata/datapackage-go?status.svg)](https://godoc.org/github.com/frictionlessdata/datapackage-go/pkg)

# datapackage-go
A Go library for working with [Data Packages](http://specs.frictionlessdata.io/data-package/).

## Features

* [pkg.Package](https://godoc.org/github.com/frictionlessdata/datapackage-go/pkg#Package) class for working with data packages
* [Resource](https://godoc.org/github.com/frictionlessdata/datapackage-go/pkg#Resource) class for working with data resources
* [Valid](https://godoc.org/github.com/frictionlessdata/datapackage-go/pkg#Valid) function for validating data package descriptors

## Getting started

## Library Installation

This package uses [semantic versioning 2.0.0](http://semver.org/).

### Using dep

```sh
$ go get -u github.com/golang/dep/cmd/dep
$ dep init
$ dep ensure
```

### Using go get

```sh
$ go get -u github.com/frictionlessdata/datapackage-go/...
```
