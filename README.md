# g2-sdk-go-grpc

## :warning: WARNING: g2-sdk-go-grpc is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing `g2-sdk-go-grpc` packages provide a Go Software Development Kit
adhering to the
[g2-sdk-go](https://github.com/Senzing/g2-sdk-go) interfaces that
communicates with a
[Senzing gRPC server](https://github.com/Senzing/servegrpc).

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing/g2-sdk-go-grpc.svg)](https://pkg.go.dev/github.com/senzing/g2-sdk-go-grpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing/g2-sdk-go-grpc)](https://goreportcard.com/report/github.com/senzing/g2-sdk-go-grpc)
[![License](https://img.shields.io/badge/License-Apache2-brightgreen.svg)](https://github.com/Senzing/g2-sdk-go-grpc/blob/main/LICENSE)

[![gosec.yaml](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/gosec.yaml/badge.svg)](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/gosec.yaml)
[![go-test-linux.yaml](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/go-test-linux.yaml/badge.svg)](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/go-test-linux.yaml)
[![go-test-darwin.yaml](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/go-test-darwin.yaml/badge.svg)](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/go-test-darwin.yaml)
[![go-test-windows.yaml](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/go-test-windows.yaml/badge.svg)](https://github.com/Senzing/g2-sdk-go-grpc/actions/workflows/go-test-windows.yaml)

## Overview

The Senzing `g2-sdk-go-grpc` packages enable Go programs to call Senzing library functions
across a network to a
[Senzing gRPC server](https://github.com/Senzing/servegrpc).

`g2-sdk-go-grpc` packages implement the following
[g2-sdk-go](https://github.com/Senzing/g2-sdk-go)
interfaces:

1. [G2config](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2api#G2config)
1. [G2configmgr](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2api#G2configmgr)
1. [G2diagnostic](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2api#G2diagnostic)
1. [G2engine](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2api#G2engine)
1. [G2product](https://pkg.go.dev/github.com/senzing/g2-sdk-go/g2api#G2product)

Other implementations of the
[g2-sdk-go](https://github.com/Senzing/g2-sdk-go)
interface include:

- [g2-sdk-go-base](https://github.com/Senzing/g2-sdk-go-base) - for
  calling Senzing SDK APIs natively
- [g2-sdk-go-mock](https://github.com/Senzing/g2-sdk-go-mock) - for
  unit testing calls to the Senzing Go SDK
- [go-sdk-abstract-factory](https://github.com/Senzing/go-sdk-abstract-factory) - An
  [abstract factory pattern](https://en.wikipedia.org/wiki/Abstract_factory_pattern)
  for switching among implementations

## Use

(TODO:)

## References

1. [Development](docs/development.md)
1. [Errors](docs/errors.md)
1. [Examples](docs/examples.md)
1. [Package reference](https://pkg.go.dev/github.com/senzing/g2-sdk-go-grpc)
