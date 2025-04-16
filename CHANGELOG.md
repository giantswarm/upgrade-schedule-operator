# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [Unreleased]

### Fixed

- Disable logger development mode to avoid panicking, use zap as logger.
- Fix linting issues.

## [0.10.1] - 2024-08-01

### Fixed 

- Change label to check CAPI clusters.

## [0.10.0] - 2024-08-01

### Changed

- Support CAPI clusters.

## [0.9.0] - 2024-03-27

### Changed

- Configure `gsoci.azurecr.io` as the default container image registry.

### Added

- Add global.podSecurityStandards.enforced value for PSS migration.

## [0.8.0] - 2023-07-13

### Added 

- Added required values for pss policies.

## [0.7.4] - 2023-07-05

### Fixed

- Remove vintage monitoring labels.

## [0.7.3] - 2023-07-04

### Fixed

- Add missing team label on servicemonitor.

## [0.7.2] - 2023-04-03

### Fixed

- Add `installation` name into reconciler.

## [0.7.1] - 2023-03-31

### Fixed

- Allow required volume types in PSP so that pods can be admitted

## [0.7.0] - 2023-03-31

### Added

- Added the use of the runtime/default seccomp profile.
- Add installation name into slack notification.

## [0.6.1] - 2023-01-03

### Fixed

- Fixed scraping port in annotation for monitoring.

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


[Unreleased]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.10.1...HEAD
[0.10.1]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.10.0...v0.10.1
[0.10.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.7.4...v0.8.0
[0.7.4]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.7.3...v0.7.4
[0.7.3]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.7.2...v0.7.3
[0.7.2]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.7.1...v0.7.2
[0.7.1]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.6.1...v0.7.0
[0.6.1]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/giantswarm/upgrade-schedule-operator/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/giantswarm/upgrade-schedule-operator/releases/tag/v0.1.0
