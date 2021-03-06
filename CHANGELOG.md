# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



## [Unreleased]

## [0.6.0] - 2022-07-19

### Changed

- Add current and target release version to metrics.

## [0.5.0] - 2022-03-30

### Changed

- Upgrade cluster-api to v1.0.4 and use `v1beta1` types.
- Drop `apiextensions` and use `k8smetadata` for labels and annotations.

## [0.4.0] - 2022-03-21

### Added

- Add VerticalPodAutoscaler CR.

## [0.3.1] - 2021-11-24

### Changed

- Fix clearing metrics when cluster is deleted or invalid.

## [0.3.0] - 2021-10-04

### Added
- Push app to Azure app collection. 

## [0.2.1] - 2021-09-28

### Changed

- Updated out of hours email contact.

## [0.2.0] - 2021-09-15


### Added

- Add prometheus metrics about number of attempted, succeeded and failed upgrades.
- Add prometheus metrics about upgrade time of all clusters.
- Added documentation in the `README.md` file.

### Added

- Tests

### Fixed

- Fix out of office time.

## [0.1.1] - 2021-09-13

### Fixed

- Fix out of office check for weekends.

## [0.1.0] - 2021-09-10

### Added

- Add helm manifest.
- Add Upgrade-Schedule-Operator.


[Unreleased]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.6.0...HEAD
[0.6.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/giantswarm/upgrade-schedule-operator/releases/tag/v0.1.0
