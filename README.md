# sz-sdk-go-grpc

If you are beginning your journey with [Senzing],
please start with [Senzing Quick Start guides].

You are in the [Senzing Garage] where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: sz-sdk-go-grpc is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing `sz-sdk-go-grpc` packages provide a [Go]
language Software Development Kit adhering to the
[sz-sdk-go] interfaces that communicates with a [Senzing gRPC server].

[![Go Reference Badge]][Package reference]
[![Go Report Card Badge]][Go Report Card]
[![License Badge]][License]
[![go-test-linux.yaml Badge]][go-test-linux.yaml]
[![go-test-darwin.yaml Badge]][go-test-darwin.yaml]
[![go-test-windows.yaml Badge]][go-test-windows.yaml]

[![golangci-lint.yaml Badge]][golangci-lint.yaml]

## Overview

The Senzing `g2-sdk-go-grpc` packages enable Go programs to call Senzing library functions
across a network to a
[Senzing gRPC server](https://github.com/senzing-garage/servegrpc).

Other implementations of the [sz-sdk-go]
interface include:

- [sz-sdk-go-core] - for calling Senzing SDK APIs natively
- [sz-sdk-go-mock] - for unit testing calls to the Senzing Go SDK
- [go-sdk-abstract-factory] - An [abstract factory pattern] for switching among implementations

## Use

(TODO:)

## References

1. [Development]
1. [Errors]
1. [Examples]
1. [Package reference]

[abstract factory pattern]: https://en.wikipedia.org/wiki/Abstract_factory_pattern
[Development]: docs/development.md
[Errors]: docs/errors.md
[Examples]: docs/examples.md
[Go Reference Badge]: https://pkg.go.dev/badge/github.com/senzing-garage/template-go.svg
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/senzing-garage/template-go
[Go Report Card]: https://goreportcard.com/report/github.com/senzing-garage/template-go
[go-sdk-abstract-factory]: https://github.com/senzing-garage/go-sdk-abstract-factory
[go-test-darwin.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-darwin.yaml/badge.svg
[go-test-darwin.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-darwin.yaml
[go-test-linux.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-linux.yaml/badge.svg
[go-test-linux.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-linux.yaml
[go-test-windows.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-windows.yaml/badge.svg
[go-test-windows.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-windows.yaml
[Go]: https://go.dev/
[golangci-lint.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/golangci-lint.yaml/badge.svg
[golangci-lint.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/golangci-lint.yaml
[License Badge]: https://img.shields.io/badge/License-Apache2-brightgreen.svg
[License]: https://github.com/senzing-garage/template-go/blob/main/LICENSE
[Package reference]: https://pkg.go.dev/github.com/senzing-garage/sz-sdk-go-core
[Senzing Garage]: https://github.com/senzing-garage
[Senzing gRPC server]: https://github.com/senzing-garage/servegrpc
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[Senzing]: https://senzing.com/
[sz-sdk-go-core]: https://github.com/senzing-garage/sz-sdk-go-core
[sz-sdk-go-mock]: https://github.com/senzing-garage/sz-sdk-go-mock
[sz-sdk-go]: https://github.com/senzing-garage/sz-sdk-go
