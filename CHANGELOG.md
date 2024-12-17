# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], [markdownlint],
and this project adheres to [Semantic Versioning].

## [Unreleased]

-

## [0.8.6] - 2024-12-17

### Changed in 0.8.7

- Fixed `Reinitialize()`

## [0.8.5] - 2024-12-10

### Changed in 0.8.5

- Update dependencies

## [0.8.4] - 2024-10-30

### Changed in 0.8.4

- Add `Reinitialize()` and `Destroy` to `SzAbstractFactory`
- Update dependencies

## [0.8.3] - 2024-10-01

### Added in 0.8.3

- Method `PreprocessRecord()`

## [0.8.2] - 2024-09-11

### Changed in 0.8.2

- Update dependencies
- Added test cases.

## [0.8.1] - 2024-08-27

### Changed in 0.8.1

- Modify method calls to match Senzing API 4.0.0-24237

## [0.8.0] - 2024-08-23

### Changed in 0.8.0

- Change from `g2` to `sz`/`er`

## [0.7.3] - 2024-08-12

### Changed in 0.7.3

- Update to `template-go`
- Update tests

## [0.7.2] - 2024-06-26

### Changed in 0.7.2

- Synchronized with [sz-sdk-go-core](https://github.com/senzing-garage/sz-sdk-go-core) and [sz-sdk-go-mock](https://github.com/senzing-garage/sz-sdk-go-mock)
- Updated dependencies

## [0.7.1] - 2024-05-09

### Added in 0.7.1

- `SzDiagnostic.GetFeature`
- `SzEngine.FindInterestingEntitiesByEntityId`
- `SzEngine.FindInterestingEntitiesByRecordId`

### Deleted in 0.7.1

- `SzEngine.GetRepositoryLastModifiedTime`

### Changed in 0.7.1

- Migrated from `g2` to `sz`
- Updated dependencies

## [0.7.0] - 2024-03-01

### Changed in 0.7.0

- Updated dependencies
- Deleted methods not used in V4

## [0.6.0] - 2024-01-26

### Changed in 0.6.0

- Renamed module to `github.com/senzing-garage/g2-sdk-go-grpc`
- Refactor to [template-go](https://github.com/senzing-garage/template-go)
- Update dependencies
  - google.golang.org/grpc v1.61.0
  - github.com/senzing-garage/g2-sdk-go v0.9.0
  - github.com/senzing-garage/g2-sdk-proto/go v0.0.0-20240126210601-d02d3beb81d4

## [0.5.0] - 2024-01-02

### Changed in 0.5.0

- Refactor to [template-go](https://github.com/senzing-garage/template-go)
- Update dependencies
  - github.com/senzing-garage/go-common v0.4.0
  - github.com/senzing-garage/go-logging v1.4.0
  - github.com/senzing-garage/go-observing v0.3.0
  - github.com/senzing/g2-sdk-go v0.8.0
  - google.golang.org/grpc v1.60.1

## [0.4.4] - 2023-12-12

### Added in 0.4.4

- `ExportCSVEntityReportIterator` and `ExportJSONEntityReportIterator`

### Changed in 0.4.4

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.7.6
  - google.golang.org/grpc v1.60.0

## [0.4.3] - 2023-10-18

### Changed in 0.4.3

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.7.4
  - github.com/senzing/g2-sdk-proto/go v0.0.0-20231016131354-0d0fba649357
  - github.com/senzing-garage/go-common v0.3.1
  - github.com/senzing-garage/go-logging v1.3.3
  - github.com/senzing-garage/go-observing v0.2.8
  - google.golang.org/grpc v1.59.0

## [0.4.2] - 2023-10-13

### Changed in 0.4.2

- Changed from `int` to `int64` where required by the SenzingAPI
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.7.3
  - google.golang.org/grpc v1.58.3

### Deleted in 0.4.2

- `g2product.ValidateLicenseFile`
- `g2product.ValidateLicenseStringBase64`

## [0.4.1] - 2023-10-03

### Changed in 0.4.1

- Updated testing

## [0.4.0] - 2023-09-26

### Changed in 0.4.0

- Supports SenzingAPI 3.8.0
- Deprecated functions have been removed

## [0.3.3] - 2023-09-01

### Changed in 0.3.3

- Last version before SenzingAPI 3.8.0

## [0.3.2] - 2023-08-05

### Changed in 0.3.2

- Changed default port to 8261
- Moved to `go-logging`
- Refactor to `template-go`
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.8
  - github.com/senzing-garage/go-common v0.2.11
  - github.com/senzing-garage/go-logging v1.3.2
  - github.com/senzing-garage/go-observing v0.2.7
  - google.golang.org/grpc v1.57.0

## [0.3.1] - 2023-05-26

### Changed in 0.3.1

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.5
  - github.com/senzing/g2-sdk-proto/go v0.0.0-20230608182106-25c8cdc02e3c
  - github.com/senzing-garage/go-common v0.1.4
  - github.com/senzing-garage/go-logging v1.2.6
  - github.com/senzing-garage/go-observing v0.2.6
  - github.com/stretchr/testify v1.8.4
  - google.golang.org/grpc v1.56.0

## [0.3.0] - 2023-05-26

### Changed in 0.3.0

- Change `g2config.Load()` signature
- Added `gosec`
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.4
  - github.com/senzing/g2-sdk-proto/go v0.0.0-20230526140633-b44eb0f20e1b

## [0.2.7] - 2023-05-19

### Changed in 0.2.7

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.3

## [0.2.6] - 2023-05-11

### Changed in 0.2.6

- Update dependencies
  - github.com/senzing-garage/go-common v0.1.3
  - github.com/senzing-garage/go-logging v1.2.3

## [0.2.5] - 2023-05-10

### Changed in 0.2.5

- Added `GetObserverOrigin()` and `SetObserverOrigin()` to g2* packages
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.2
  - github.com/senzing-garage/go-observing v0.2.2
  - google.golang.org/grpc v1.55.0

## [0.2.4] - 2023-04-22

### Changed in 0.2.4

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.1

## [0.2.3] - 2023-04-21

### Changed in 0.2.3

- Changed `SetLogLevel(ctx context.Context, logLevel logger.Level)` to `SetLogLevel(ctx context.Context, logLevelName string)`

## [0.2.2] - 2023-03-29

### Changed in 0.2.2

- Added `helper.ConvertGrpcError()`
- Refactored documentation
- Updated dependencies

## [0.2.1] - 2023-02-21

### Changed in 0.2.1

- Change GetSdkId() signature.

## [0.2.0] - 2023-02-15

### Changed in 0.2.0

- Using refactored g2-sdk-go

## [0.1.0] - 2023-01-31

### Added to 0.1.0

- Initial implementation

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[markdownlint]: https://dlaa.me/markdownlint/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
